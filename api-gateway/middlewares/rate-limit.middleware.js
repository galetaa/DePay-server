const { RateLimiterRedis } = require('rate-limiter-flexible');
const redis = require('../config/redis.config');

const rateLimiter = new RateLimiterRedis({
  storeClient: redis,
  keyPrefix: 'ratelimit',
  points: 100, // 100 requests
  duration: 60, // per minute
});

async function rateLimitMiddleware(req, res, next) {
  try {
    await rateLimiter.consume(req.user?.id || req.ip);
    next();
  } catch (rejRes) {
    res.status(429).send('Too Many Requests');
  }
}

module.exports = rateLimitMiddleware;