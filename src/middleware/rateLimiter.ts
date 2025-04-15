import rateLimit from 'express-rate-limit';
import { Config } from '../config';
import { Request, Response } from 'express';

interface RateLimitConfig {
  windowMs: number;
  max: number;
  message: string;
  statusCode: number;
  headers: boolean;
}

export function createRateLimiter(config: Config) {
  const defaultConfig: RateLimitConfig = {
    windowMs: 60 * 1000, // 1 minute window
    max: config.rateLimit.requestsPerMinute,
    message: 'Too many requests, please try again later',
    statusCode: 429,
    headers: true,
  };

  // Create rate limiter middleware
  const limiter = rateLimit({
    windowMs: defaultConfig.windowMs,
    max: defaultConfig.max,
    message: { error: defaultConfig.message },
    statusCode: defaultConfig.statusCode,
    standardHeaders: defaultConfig.headers,
    legacyHeaders: false,
    keyGenerator: (req: Request): string => {
      // Use IP and user agent for rate limiting
      return `${req.ip}-${req.headers['user-agent']}`;
    },
    handler: (req: Request, res: Response): void => {
      const retryAfter = Math.ceil(defaultConfig.windowMs / 1000);
      res.setHeader('Retry-After', retryAfter);
      res.status(defaultConfig.statusCode).json({
        error: defaultConfig.message,
        retryAfter: retryAfter,
      });
    },
    skip: (req: Request): boolean => {
      // Skip rate limiting for health checks
      return req.path === '/health';
    }
  });

  return limiter;
}

// Additional rate limiter for specific endpoints or user tiers
export function createTieredRateLimiter(config: Config, tier: string) {
  const tierLimits: { [key: string]: number } = {
    'free': config.rateLimit.requestsPerMinute,
    'basic': config.rateLimit.requestsPerMinute * 2,
    'premium': config.rateLimit.requestsPerMinute * 5
  };

  return rateLimit({
    windowMs: 60 * 1000,
    max: tierLimits[tier] || tierLimits.free,
    message: { error: `Rate limit exceeded for ${tier} tier` }
  });
}