const express = require('express');
const bodyParser = require('body-parser');
const cors = require('cors');
const dotenv = require('dotenv');
const winston = require('winston');
const socketIo = require('socket.io');
const mongoose = require('mongoose');
dotenv.config();

// Подключение к MongoDB
mongoose.connect(process.env.MONGODB_URI, {useNewUrlParser: true, useUnifiedTopology: true})
    .then(() => {
        console.log('Connected to MongoDB');
    })
    .catch((err) => {
        console.error('Error connecting to MongoDB:', err);
    });


// Загружаем переменные окружения
dotenv.config();

const app = express();
const port = process.env.PORT || 3000;

// Настройка логирования
const logger = winston.createLogger({
    transports: [
        new winston.transports.Console({
            format: winston.format.combine(
                winston.format.colorize(),
                winston.format.simple()
            ),
        }),
    ],
});

// Middleware
app.use(cors());
app.use(bodyParser.json());

// Подключаем маршруты
const authRoutes = require('./controllers/authController');
const walletRoutes = require('./controllers/walletController');
const transactionRoutes = require('./controllers/transactionController');

// Используем роутеры
app.use('/auth', authRoutes);
app.use('/wallets', walletRoutes);
app.use('/transactions', transactionRoutes);

// Запуск сервера
const server = app.listen(port, () => {
    logger.info(`Server is running on http://localhost:${port}`);
});

// WebSocket для уведомлений
const io = socketIo(server);
io.on('connection', (socket) => {
    logger.info('A user connected via WebSocket');
    socket.on('disconnect', () => {
        logger.info('A user disconnected');
    });
});

module.exports = {app, io};


