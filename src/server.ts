import express, { Request, Response } from 'express';
import cors from 'cors';
import { ContentGenerator, GenerationOptions } from './content/generator';
import { GeminiClient } from './api/gemini';
import { loadConfig } from './config';
import { createRateLimiter, createTieredRateLimiter } from './middleware/rateLimiter';

interface GenerateRequest {
  content_type: string;
  topic: string;
  tone?: string;
  length?: number;
  additional_context?: string;
}

interface GenerateResponse {
  content: string;
  filename?: string;
  error?: string;
  rateLimit?: {
    remaining: any;
    reset: any;
  };
}

export class Server {
  private app = express();
  private generator: ContentGenerator;
  private config = loadConfig();

  constructor(geminiClient: GeminiClient) {
    this.generator = new ContentGenerator(geminiClient);
    this.setupMiddleware();
    this.setupRoutes();
  }

  private setupMiddleware() {
    this.app.use(cors());
    this.app.use(express.json());

    // Global rate limiter
    const globalLimiter = createRateLimiter(this.config);
    this.app.use(globalLimiter);

    // Tiered rate limiters for different endpoints
    const generateLimiter = createTieredRateLimiter(this.config, 'free');
    this.app.use('/api/generate', generateLimiter);
  }

  private setupRoutes() {
    this.app.post('/api/generate', this.handleGenerate.bind(this));
    this.app.get('/health', (_, res) => res.send('Server is running'));
  }

  private async handleGenerate(req: Request, res: Response) {
    try {
      const reqBody = req.body as GenerateRequest;

      // Validate request
      if (!reqBody.content_type || !reqBody.topic) {
        return res.status(400).json({ error: 'Content type and topic are required' });
      }

      // Convert request to generation options
      const opts: GenerationOptions = {
        contentType: reqBody.content_type,
        topic: reqBody.topic,
        tone: reqBody.tone,
        length: reqBody.length,
        additionalContext: reqBody.additional_context
      };

      // Generate content
      const content = await this.generator.generateContent(opts);
      const filename = await this.generator.saveContent(content, opts.contentType);

      // Track rate limit headers
      const remaining = res.getHeader('X-RateLimit-Remaining');
      const reset = res.getHeader('X-RateLimit-Reset');

      const response: GenerateResponse = { 
        content, 
        filename,
        rateLimit: {
          remaining,
          reset
        }
      };
      res.json(response);

    } catch (error) {
      console.error('Error generating content:', error);
      res.status(500).json({ 
        error: error instanceof Error ? error.message : 'Internal server error',
        rateLimit: {
          remaining: res.getHeader('X-RateLimit-Remaining'),
          reset: res.getHeader('X-RateLimit-Reset')
        }
      });
    }
  }

  start(port: number = 8080): void {
    this.app.listen(port, () => {
      console.log(`Server running on port ${port}`);
    });
  }
}

// Start server when run directly
if (require.main === module) {
  const config = loadConfig();
  const geminiClient = new GeminiClient(config);
  const server = new Server(geminiClient);
  
  const port = parseInt(process.env.PORT || '8080', 10);
  server.start(port);
}