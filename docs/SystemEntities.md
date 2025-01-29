### **Users**  
System users are cryptocurrency wallet owners who use the application for NFC-based payments  

| Name          | Data Type        | Description                        | Constraints                     |  
|--------------|-----------------|------------------------------------|---------------------------------|  
| `user_id`    | UUID             | Unique identifier of the user      | PRIMARY KEY                     |  
| `email`      | VARCHAR(255)     | User's email address               | UNIQUE, NOT NULL                |  
| `phone`      | VARCHAR(15)      | User's phone number                | UNIQUE, NOT NULL                |  
| `created_at` | TIMESTAMP        | Date and time of user registration | DEFAULT CURRENT_TIMESTAMP       |  

```json  
{  
  "user_id": "550e8400-e29b-41d4-a716-446655440000",  
  "email": "user@example.com",  
  "phone": "+12345678900",  
  "created_at": "2024-03-20T12:00:00Z"  
}  
```  

---

### **Wallets**  
Cryptocurrency wallets linked to users.  

| Name         | Data Type        | Description                                     | Constraints                     |  
|-------------|-----------------|-------------------------------------------------|---------------------------------|  
| `wallet_id`  | UUID             | Unique identifier of the wallet                 | PRIMARY KEY                     |  
| `user_id`    | UUID             | User identifier                                 | FOREIGN KEY (Users.user_id)     |  
| `address`    | VARCHAR(42)      | Wallet address (e.g., 0x...)                    | NOT NULL                        |  
| `network`    | VARCHAR(50)      | Blockchain network (Ethereum, Solana, TRC, etc) | NOT NULL                        |  
| `connected_at` | TIMESTAMP        | Date and time of wallet linking                 | DEFAULT CURRENT_TIMESTAMP       |  

```json  
{  
  "wallet_id": "550e8400-e29b-41d4-a716-446655440001",  
  "user_id": "550e8400-e29b-41d4-a716-446655440000",  
  "address": "0x742d35Cc6634C0532925a3b844Bc454e4438f44e",  
  "network": "ethereum",  
  "connected_at": "2024-03-20T12:05:00Z"  
}  
```  

---

### **Transactions**  
**Description:** Records of payments initiated via NFC.  

| Name         | Data Type        | Description                                   | Constraints                     |  
|-------------|-----------------|-----------------------------------------------|---------------------------------|  
| `tx_id`     | UUID             | Unique identifier of the transaction          | PRIMARY KEY                     |  
| `user_id`   | UUID             | User identifier                               | FOREIGN KEY (Users.user_id)     |  
| `session_id` | VARCHAR(50)      | NFC session identifier                        | NOT NULL                        |  
| `amount`    | DECIMAL(18, 8)   | Transaction amount                            | NOT NULL                        |  
| `currency`  | VARCHAR(10)      | Currency (USDC, USDT, ETH, etc)               | NOT NULL                        |  
| `status`    | VARCHAR(20)      | Transaction status (pending/confirmed/failed) | NOT NULL                        |  
| `tx_hash`   | VARCHAR(66)      | Blockchain transaction hash                   | UNIQUE                          |  
| `created_at` | TIMESTAMP        | Date and time of transaction creation         | DEFAULT CURRENT_TIMESTAMP       |  

```json  
{  
  "tx_id": "550e8400-e29b-41d4-a716-446655440002",  
  "user_id": "550e8400-e29b-41d4-a716-446655440000",  
  "session_id": "abc123",  
  "amount": 50.0,  
  "currency": "USDC",  
  "status": "confirmed",  
  "tx_hash": "0x88df016429689c079f3b2f6ad39fa052532c56795b733da78a91ebe6a713944b",  
  "created_at": "2024-03-20T12:10:00Z"  
}  
```  

---

### **NFC Sessions**  
**Description:** Temporary sessions created by the terminal for payment processing.  

| Name         | Data Type        | Description                       | Constraints                     |  
|-------------|-----------------|-----------------------------------|---------------------------------|  
| `session_id` | VARCHAR(50)      | Unique session identifier         | PRIMARY KEY                     |  
| `merchant_id` | VARCHAR(50)     | Merchant identifier               | NOT NULL                        |  
| `amount`    | DECIMAL(18, 8)   | Payment amount                    | NOT NULL                        |  
| `expires_at` | TIMESTAMP        | Session expiration date and time  | NOT NULL                        |  

```json  
{  
  "session_id": "abc123",  
  "merchant_id": "store_456",  
  "amount": 50.0,  
  "expires_at": "2024-03-20T12:15:00Z"  
}  
```  

---

### **Gas Pool**  
A reserve fund for covering gas fees for transactions.  

| Name        | Data Type        | Description                            | Constraints                     |  
|------------|-----------------|----------------------------------------|---------------------------------|  
| `pool_id`  | UUID             | Unique identifier of the gas pool      | PRIMARY KEY                     |  
| `merchant_id` | VARCHAR(50)   | Merchant identifier                    | NOT NULL                        |  
| `balance`  | DECIMAL(18, 8)   | Pool balance in the specified currency | NOT NULL                        |  
| `currency` | VARCHAR(10)      | Pool currency (ETH, SOL, etc)          | NOT NULL                        |  

```json  
{  
  "pool_id": "550e8400-e29b-41d4-a716-446655440003",  
  "merchant_id": "store_456",  
  "balance": 10.0,  
  "currency": "ETH"  
}  
```  

---

### **Entity Relationships**  
- **User ↔ Wallet:**   
  `Users.user_id` → `Wallets.user_id`.  

- **User ↔ Transaction:**  
  `Users.user_id` → `Transactions.user_id`.  

- **Session ↔ Transaction:**  
   `Sessions.session_id` → `Transactions.session_id`.  

- **Merchant ↔ Gas Pool:**  
   `Sessions.merchant_id` → `GasPool.merchant_id`.  
