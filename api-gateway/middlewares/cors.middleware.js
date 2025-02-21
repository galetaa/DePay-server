const cors = require('cors');

const corsMiddleware = cors({
  origin: process.env.CLIENT_URL || '*', // Разрешаем доступ только с определенного домена
  credentials: true,
});

module.exports = corsMiddleware;