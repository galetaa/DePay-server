const { createLogger, transports, format } = require('winston');
const { ElasticsearchTransport } = require('winston-elasticsearch'); // Correct import

const esTransport = new ElasticsearchTransport({
  level: 'info',
  clientOpts: {
    node: process.env.ELASTICSEARCH_URL || 'http://localhost:9200',
  },
});

const logger = createLogger({
  level: 'info',
  format: format.combine(
    format.timestamp(),
    format.json()
  ),
  transports: [
    new transports.Console(),
    esTransport, // Correct instance
  ],
});

module.exports = logger;
