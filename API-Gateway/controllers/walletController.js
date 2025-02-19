const express = require('express');
const authenticateJWT = require('../middlewares/authMiddleware');
const router = express.Router();

// Добавить новый кошелек
router.post('/', authenticateJWT, (req, res) => {
  // Логика добавления кошелька
});

// Получить список кошельков
router.post('/list', authenticateJWT, (req, res) => {
  // Логика получения списка кошельков
});

module.exports = router;
