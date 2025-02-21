const express = require('express');
const helmet = require('helmet');
const morgan = require('morgan');
const cors = require('cors');
const logger = require('./config/elasticsearch.config');

const app = express();

const corsMiddleware = require('./middlewares/cors.middleware');
const errorHandler = require('./middlewares/error-handling.middleware');

// Middleware
app.use(helmet());
app.use(cors());
app.use(express.json());
app.use(morgan('combined'));

app.use(corsMiddleware);
app.use(errorHandler);

// Logging Middleware
app.use((req, res, next) => {
  const start = Date.now();
  res.on('finish', () => {
    const responseTime = Date.now() - start;
    logger.info({
      ip: req.ip,
      method: req.method,
      url: req.url,
      status: res.statusCode,
      responseTime,
    });
  });
  next();
});

// Routes
app.use('/auth', require('./routes/auth.routes'));
app.use('/wallets', require('./routes/wallet.routes'));
app.use('/transactions', require('./routes/transaction.routes'));

// Error Handling Middleware
app.use((err, req, res, next) => {
  console.error(err.stack);
  res.status(500).send('Something broke!');
});

module.exports = app;
