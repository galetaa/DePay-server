const express = require('express');
const jwt = require('jsonwebtoken');
const router = express.Router();

// Регистрация пользователя
router.post('/register', (req, res) => {
    const {fname, lname, email, password} = req.body;

    if (!email || !password) {
        return res.status(400).json({error: 'Email and password are required'});
    }

    // Здесь должна быть проверка данных и запись в базу данных
    const user = {id: 'user123', fname, lname, email};

    const token = jwt.sign(user, process.env.JWT_SECRET, {expiresIn: '1h'});
    res.json({user_id: user.id, token});
});

// Аутентификация пользователя
router.post('/login', (req, res) => {
    const {email, password} = req.body;

    if (!email || !password) {
        return res.status(400).json({error: 'Email and password are required'});
    }

    // Здесь должна быть проверка учетных данных
    const user = {id: 'user123', email};

    const token = jwt.sign(user, process.env.JWT_SECRET, {expiresIn: '1h'});
    res.json({user_id: user.id, token});
});

module.exports = router;