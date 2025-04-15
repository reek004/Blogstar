import { config as dotenvConfig } from 'dotenv';
import fs from 'fs';
import YAML from 'yaml';

export interface Config {
  geminiApiKey: string;
  projectId: string;
  location: string;
  defaultModels: string[];
  maxTokens: number;
  rateLimit: {
    requestsPerMinute: number;
    tiers: {
      free: number;
      basic: number;
      premium: number;
    };
    burst: boolean;
    delay_after: number;
    delay_ms: number;
  };
}

export function loadConfig(): Config {
  // Load .env file
  dotenvConfig();

  // Read config from YAML
  const configFile = fs.readFileSync('configs/config.yaml', 'utf-8');
  const config = YAML.parse(configFile) as Config;

  // Override with environment variables
  if (process.env.GEMINI_API_KEY) {
    config.geminiApiKey = process.env.GEMINI_API_KEY;
  }

  // Set defaults
  if (!config.defaultModels?.length) {
    config.defaultModels = ['gemini-1.5-pro'];
  }

  if (!config.rateLimit?.requestsPerMinute) {
    config.rateLimit = {
      requestsPerMinute: 5,
      tiers: {
        free: 5,
        basic: 10,
        premium: 15
      },
      burst: false,
      delay_after: 3,
      delay_ms: 1000
    };
  }

  return config;
}