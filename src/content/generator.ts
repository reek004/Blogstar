import fs from 'fs/promises';
import path from 'path';
import { GeminiClient } from '../api/gemini';

export type ContentType = 'blog_post' | 'article' | 'social_media' | 'script';

export interface GenerationOptions {
  contentType: ContentType | string;
  topic: string;
  tone?: string;
  length?: number;
  additionalContext?: string;
}

export class ContentGenerator {
  constructor(private geminiClient: GeminiClient) {}

  async generateContent(opts: GenerationOptions): Promise<string> {
    const prompt = this.constructPrompt(opts);
    return await this.geminiClient.generateContent(prompt);
  }

  private constructPrompt(opts: GenerationOptions): string {
    let prompt = `Write a ${opts.contentType} about '${opts.topic}'`;
    
    if (opts.tone) {
      prompt += ` in a ${opts.tone} tone`;
    }
    
    if (opts.length) {
      prompt += `. Aim for approximately ${opts.length} words`;
    }
    
    if (opts.additionalContext) {
      prompt += `. Additional context: ${opts.additionalContext}`;
    }
    
    return prompt;
  }

  async saveContent(content: string, contentType: string): Promise<string> {
    const outputDir = 'generated_content';
    await fs.mkdir(outputDir, { recursive: true });

    const timestamp = new Date().toISOString().replace(/[:\.]/g, '').slice(0, 15);
    const filename = path.join(outputDir, `${contentType}_${timestamp}.txt`);

    await fs.writeFile(filename, content, 'utf-8');
    return filename;
  }
}