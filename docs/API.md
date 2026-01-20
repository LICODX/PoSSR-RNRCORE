# RNR Core API Documentation

## Base URL
```
http://localhost:8080/api
```

## Authentication
Currently no authentication required for local node access. Production deployments should use reverse proxy with API keys.

---

## Wallet Endpoints

### GET /api/wallet
Get current wallet information.

**Response:**
```json
{
  "address": "rnr1pq03gqs8zg0sgqg7zsw3u8sgqqdp7rsrzuy3wxg7pyyqxrcspsr3cqq7qvqs78c2zyrpqzqdqvfq7xs8pcgq2m9d04",
  "balance": 1000,
  "publicKey": "0x1234...",
  "nonce": 5
}
```

**Status Codes:**
- `200 OK` - Success
- `500 Internal Server Error` - Wallet not initialized

---

### POST /api/wallet/send
Send RNR to another address.

**Request Body:**
```json
{
  "to": "rnr1...",
  "amount": 100,
  "fee": 1
}
```

**Response:**
```json
{
  "success": true,
  "txHash": "abc123def456...",
  "message": "Transaction submitted to mempool"
}
```

**Validation:**
- `to`: Must be valid Bech32 RNR address
- `amount`: Must be > 0
- `fee`: Must be >= 1 (MinTxFee)

**Status Codes:**
- `200 OK` - Transaction created
- `400 Bad Request` - Invalid parameters
- `500 Internal Server Error` - Transaction creation failed

---

## Blockchain Endpoints

### GET /api/stats
Get current blockchain statistics.

**Response:**
```json
{
  "blockHeight": 12345,
  "tps": 150,
  "peers": 25,
  "mempoolSize": 42,
  "networkHashrate": "1.2 TH/s"
}
```

---

### GET /api/blocks
Get recent blocks (paginated).

**Query Parameters:**
- `limit` (optional): Number of blocks to return (default: 20, max: 100)
- `offset` (optional): Offset for pagination (default: 0)

**Response:**
```json
{
  "blocks": [
    {
      "height": 12345,
      "hash": "abc123...",
      "timestamp": 1705552800,
      "txCount": 10,
      "difficulty": 1000000
    }
  ],
  "total": 12345
}
```

---

### GET /api/block/{height}
Get specific block by height.

**Path Parameters:**
- `height`: Block height (0 to current tip)

**Response:**
```json
{
  "height": 12345,
  "hash": "abc123...",
  "prevHash": "def456...",
  "merkleRoot": "ghi789...",
  "timestamp": 1705552800,
  "difficulty": 1000000,
  "nonce": 987654321,
  "txCount": 10,
  "vrfSeed": "jkl012..."
}
```

**Status Codes:**
- `200 OK` - Block found
- `404 Not Found` - Block doesn't exist

---

### GET /api/mining
Get current mining status.

**Response:**
```json
{
  "status": "active",
  "currentBlock": 12345,
  "difficulty": 1000000,
  "lastBlockTime": "14:32:10"
}
```

---

## Transaction Endpoints

### GET /api/transactions
Get recent transactions.

**Query Parameters:**
- `limit` (optional): Number of transactions (default: 20)
- `offset` (optional): Pagination offset

**Response:**
```json
{
  "transactions": [
    {
      "hash": "tx123...",
      "from": "rnr1...",
      "to": "rnr1...",
      "amount": 100,
      "fee": 1,
      "status": "confirmed",
      "blockHeight": 12345,
      "timestamp": 1705552800
    }
  ],
  "total": 5678
}
```

---

### GET /api/tx/{hash}
Get transaction by hash.

**Path Parameters:**
- `hash`: Transaction hash

**Response:**
```json
{
  "hash": "tx123...",
  "from": "rnr1...",
  "to": "rnr1...",
  "amount": 100,
  "fee": 1,
  "nonce": 5,
  "status": "confirmed",
  "blockHeight": 12345,
  "timestamp": 1705552800,
  "signature": "sig789..."
}
```

---

### GET /api/address/{address}
Get address information and transaction history.

**Path Parameters:**
- `address`: RNR Bech32 address

**Response:**
```json
{
  "address": "rnr1...",
  "balance": 1000,
  "nonce": 5,
  "txCount": 42,
  "transactions": [
    {
      "hash": "tx123...",
      "type": "sent",
      "amount": 100,
      "timestamp": 1705552800
    }
  ]
}
```

---

## Search Endpoint

### GET /api/search?q={query}
Universal search for blocks, transactions, or addresses.

**Query Parameters:**
- `q`: Search query (block height, tx hash, or address)

**Response:**
```json
{
  "type": "block|transaction|address",
  "result": { /* type-specific data */ }
}
```

---

## Metrics Endpoint

### GET /metrics
Prometheus-compatible metrics endpoint.

**Response Format:** Prometheus text format

**Metrics Exposed:**
- `rnr_block_height` - Current blockchain height
- `rnr_tps` - Transactions per second
- `rnr_peer_count` - Connected peers
- `rnr_mempool_size` - Pending transactions
- `rnr_block_time_seconds` - Block time histogram
- `rnr_mining_difficulty` - Current difficulty
- `rnr_blocks_processed_total` - Total blocks processed
- `rnr_transactions_processed_total` - Total transactions

---

## Error Responses

All error responses follow this format:
```json
{
  "success": false,
  "error": "Error message description"
}
```

**Common Error Codes:**
- `400` - Bad Request (invalid parameters)
- `404` - Not Found (resource doesn't exist)
- `500` - Internal Server Error (server-side issue)

---

## Rate Limiting

**Current:** No rate limiting (development)

**Production Recommendations:**
- 100 requests/minute per IP
- 1000 requests/hour per IP
- Use reverse proxy (nginx/CloudFlare) for rate limiting

---

## WebSocket API (Future)

**Planned endpoints:**
- `ws://localhost:8080/ws/blocks` - Real-time block updates
- `ws://localhost:8080/ws/transactions` - Real-time transaction stream
- `ws://localhost:8080/ws/mempool` - Mempool updates

---

## SDK Examples

### JavaScript/TypeScript
```javascript
// Get wallet info
const response = await fetch('http://localhost:8080/api/wallet');
const wallet = await response.json();
console.log(wallet.balance);

// Send transaction
const tx = await fetch('http://localhost:8080/api/wallet/send', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    to: 'rnr1...',
    amount: 100,
    fee: 1
  })
});
const result = await tx.json();
console.log(result.txHash);
```

### Python
```python
import requests

# Get stats
response = requests.get('http://localhost:8080/api/stats')
stats = response.json()
print(f"Block Height: {stats['blockHeight']}")

# Send transaction
tx_data = {
    'to': 'rnr1...',
    'amount': 100,
    'fee': 1
}
result = requests.post('http://localhost:8080/api/wallet/send', json=tx_data)
print(result.json())
```

### Go
```go
type WalletInfo struct {
    Address string  `json:"address"`
    Balance float64 `json:"balance"`
    Nonce   int     `json:"nonce"`
}

resp, _ := http.Get("http://localhost:8080/api/wallet")
var wallet WalletInfo
json.NewDecoder(resp.Body).Decode(&wallet)
fmt.Printf("Balance: %.2f RNR\n", wallet.Balance)
```

---

**Version:** 1.0.0  
**Last Updated:** 2026-01-18  
**Maintainer:** RNR Core Team
