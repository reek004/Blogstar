import { GoogleGenerativeAI, GenerativeModel, GenerationConfig } from '@google/generative-ai';
import { Config } from '../config';

export class GeminiClient {
  private client: GoogleGenerativeAI;
  private model: GenerativeModel;

  constructor(config: Config) {
    this.client = new GoogleGenerativeAI(config.geminiApiKey);
    this.model = this.client.getGenerativeModel({
      model: config.defaultModels[0],
      generationConfig: {
        maxOutputTokens: config.maxTokens,
        temperature: 0.7,
        topP: 0.95,
      },
    });
  }

  async generateContent(prompt: string): Promise<string> {
    try {
      const result = await this.model.generateContent(prompt);
      const response = await result.response;
      const text = response.text();
      return text || '';

    } catch (error) {
      console.error('Error generating content:', error);
      if (error instanceof Error) {
        throw new Error(`Gemini API Error: ${error.message}`);
      }
      throw new Error('Failed to generate content');
    }
  }
}