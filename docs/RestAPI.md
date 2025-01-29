
## **Designing the REST API**

### **1. General Principles**  
The API is built on the following key principles:  

- Protocol: All requests must go through HTTPS for security
- Data Format: JSON is used for all data exchanges
- Authentication: Access is secured using JWT (JSON Web Token)
- Response Codes:  
  - `200 OK` — Successful request
  - `400 Bad Request` — Invalid input data
  - `401 Unauthorized` — Missing access token
  - `403 Forbidden` — Insufficient permissions
  - `404 Not Found` — Requested resource does not exist
  - `500 Internal Server Error` — Unexpected server error

---

### **2. API Endpoints**  

#### **2.1. Authentication & Wallet Management**  

##### **2.1.1. Connecting a Wallet**  
- Endpoint: `POST /api/wallets/connect`  
- Purpose: Connect a cryptocurrency wallet using WalletConnect
- Headers:  
  - `Content-Type: application/json`  
- Request Body:  
  ```json  
  {  
    "wallet_address": "0x742d35Cc6634C0532925a3b844Bc454e4438f44e",  
    "network": "ethereum",  
    "signature": "0x88df016429689c079f3b2f6ad39fa052532c56795b733da78a91ebe6a713944b"  
  }  
  ```  
- Success Response:  
  ```json  
  {  
    "status": "success",  
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",  
    "expires_in": 3600  
  }  
  ```  
- Error Response:  
  ```json  
  {  
    "status": "error",  
    "message": "Invalid signature"  
  }  
  ```  

##### **2.1.2. Retrieving Wallets**  
- Endpoint: `GET /api/wallets`  
- Purpose: Get a list of connected wallets
- Headers:  
  - `Authorization: Bearer <access_token>`  
- Success Response:  
  ```json  
  [  
    {  
      "address": "0x742d35Cc6634C0532925a3b844Bc454e4438f44e",  
      "network": "ethereum",  
      "connected_at": "2024-03-20T12:05:00Z"  
    }  
  ]  
  ```  
- Error Response:  
  ```json  
  {  
    "status": "error",  
    "message": "Unauthorized"  
  }  
  ```  

---

#### **2.2. Managing NFC Sessions**  

##### **2.2.1. Creating a Session**  
- Endpoint: `POST /api/sessions/init`  
- Purpose: Start an NFC payment session (initiated by the terminal)
- Request Body:  
  ```json  
  {  
    "merchant_id": "store_123",  
    "amount": 50.0,  
    "currency": "USDC"  
  }  
  ```  
- Success Response:  
  ```json  
  {  
    "session_id": "abc123",  
    "expires_at": "2024-03-20T12:15:00Z"  
  }  
  ```  
- Error Response:  
  ```json  
  {  
    "status": "error",  
    "message": "Invalid merchant ID"  
  }  
  ```  

---

#### **2.3. Initiating Transactions**  

##### **2.3.1. Processing a Payment**  
- Endpoint: `POST /api/transactions/init`  
- Purpose: Process a payment using NFC
- Headers:  
  - `Authorization: Bearer <access_token>`  
  - `Content-Type: application/json`  
- Request Body:  
  ```json  
  {  
    "session_id": "abc123",  
    "wallet_address": "0x742d35Cc6634C0532925a3b844Bc454e4438f44e",  
    "signed_tx_data": "0x88df016429689c079f3b2f6ad39fa052532c56795b733da78a91ebe6a713944b"  
  }  
  ```  
- Success Response:  
  ```json  
  {  
    "tx_id": "550e8400-e29b-41d4-a716-446655440002",  
    "status": "pending",  
    "tx_hash": "0x88df016429689c079f3b2f6ad39fa052532c56795b733da78a91ebe6a713944b"  
  }  
  ```  
- Error Response:  
  ```json  
  {  
    "status": "error",  
    "message": "Session expired"  
  }  
  ```  

---

#### **2.4. Checking Transaction Status**  

##### **2.4.1. Retrieving Transaction Status**  
- Endpoint: `GET /api/transactions/{tx_id}`  
- Purpose: Get the status of a transaction
- Headers:  
  - `Authorization: Bearer <access_token>`  
- Success Response:  
  ```json  
  {  
    "tx_id": "550e8400-e29b-41d4-a716-446655440002",  
    "status": "confirmed",  
    "block_number": 19234323  
  }  
  ```  
- Error Response:  
  ```json  
  {  
    "status": "error",  
    "message": "Transaction not found"  
  }  
  ```