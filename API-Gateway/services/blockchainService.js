const axios = require('axios');
const dotenv = require('dotenv');

dotenv.config();

const getTransactionStatus = async (txHash) => {
    try {
        const response = await axios.post(process.env.BLOCKCHAIN_API_URL, {
            jsonrpc: '2.0',
            method: 'eth_getTransactionReceipt',
            params: [txHash],
            id: 1,
        });
        return response.data.result;
    } catch (err) {
        throw new Error('Error interacting with blockchain');
    }
};

module.exports = {getTransactionStatus};
