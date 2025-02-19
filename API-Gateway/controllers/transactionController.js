const express = require('express');
const authenticateJWT = require('../middlewares/authMiddleware');
const router = express.Router();

// Подать подписанную транзакцию
router.post('/', authenticateJWT, (req, res) => {
  // Логика обработки транзакции
});

// Получить статус транзакции
router.post('/status', authenticateJWT, (req, res) => {
  // Логика получения статуса транзакции
});

module.exports = router;
