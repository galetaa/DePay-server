const express = require('express');
const router = express.Router();

// Отправка транзакции
router.post('/', (req, res) => {
  const { user_id, transaction } = req.body;

  if (!user_id || !transaction) {
    return res.status(400).json({ error: 'Invalid request data' });
  }

  const transactionId = `tx_${Math.random().toString(36).substr(2, 9)}`;
  res.json({ transaction_id: transactionId, status: 'accepted', message: 'Transaction accepted for processing' });
});

// Получение статуса транзакции
router.post('/status', (req, res) => {
  const { transaction_id } = req.body;

  if (!transaction_id) {
    return res.status(400).json({ error: 'Transaction ID is required' });
  }

  const status = 'confirmed'; // Заглушка
  res.json({ transaction_id, status, details: { blockNumber: 12345678, timestamp: '2025-02-04T12:45:00Z' } });
});

module.exports = router;