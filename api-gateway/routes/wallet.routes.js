const express = require('express');
const router = express.Router();

// Добавление нового кошелька
router.post('/', (req, res) => {
  const { user_id, wallet_name, coins, tokens } = req.body;

  if (!user_id || !wallet_name || !coins) {
    return res.status(400).json({ error: 'Invalid request data' });
  }

  const walletId = `wallet_${Math.random().toString(36).substr(2, 9)}`;
  res.json({ wallet_id: walletId, message: 'Wallet successfully added' });
});

// Получение списка кошельков
router.post('/list', (req, res) => {
  const { user_id } = req.body;

  if (!user_id) {
    return res.status(400).json({ error: 'User ID is required' });
  }

  const wallets = [
    { wallet_id: 'wallet123', wallet_name: 'Main Wallet', coins: { ETH: '0xabc123...', BTC: '1A1zP1eP...' }, tokens: { DAI: '0x12345678...' } },
  ];
  res.json({ wallets });
});

module.exports = router;