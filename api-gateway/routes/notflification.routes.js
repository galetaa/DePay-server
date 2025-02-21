const express = require('express');
const router = express.Router();

// Пример эндпоинта для отправки уведомлений
router.post('/send', (req, res) => {
  const { user_id, message } = req.body;

  if (!user_id || !message) {
    return res.status(400).json({ error: 'Invalid request data' });
  }

  res.json({ success: true, message: 'Notification sent' });
});

module.exports = router;