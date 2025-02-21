const { createLogger, transports } = require('winston');
const Elasticsearch = require('winston-elasticsearch');

const esTransport = new Elasticsearch({
  level: 'info',
  clientOpts: {
    node: process.env.ELASTICSEARCH_URL || 'http://localhost:9200',
  },
});

const logger = createLogger({
  level: 'info',
  format: require('winston').format.combine(
    require('winston').format.timestamp(),
    require('winston').format.json()
  ),
  transports: [
    new transports.Console(),
    esTransport,
  ],
});

module.exports = logger;