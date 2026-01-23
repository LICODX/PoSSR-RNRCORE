# RnR Core: JSON-RPC API Reference

> **Status**: Experimental / Partial Implementation
> The current API is designed to be **Ethereum-Compatible** (`eth_` namespace) to support existing tooling (Metamask, etc.) in the future.

**Base URL**: `http://localhost:9001` (Default)
**Format**: JSON-RPC 2.0

---

## 📚 Endpoints

### `eth_blockNumber`
Returns the height of the most recent block.

**Parameters**: None

**Response**:
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": 1024
}
```

---

### `eth_getBalance`
Returns the balance of the account of given address.

**Parameters**:
1. `ADDRESS`: [Required] 20-byte address to check for balance.
2. `BLOCK`: [Optional] "latest", "earliest", or "pending".

**Response**:
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": 1000000000  // Wei/RnR-gwei
}
```
*Note: Current implementation returns a mock value for testing.*

---

### `eth_getBlockByNumber`
Returns information about a block by block number.

**Parameters**:
1. `BLOCK`: [Required] Hex encoded block number.
2. `FULL_TX`: [Required] Boolean. If true, returns full transaction objects.

**Response**:
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "number": "0x1b4",
    "hash": "0x...",
    "transactions": [...]
  }
}
```

---

### `eth_sendRawTransaction`
Creates new message call transaction or a contract creation for signed transactions.

**Parameters**:
1. `DATA`: [Required] The signed transaction data.

**Response**:
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": "0xe670ec64341771606e55d6b4ca35a1a6b75ee3d5145a99d05921026d1527331"
}
```
*Note: Returns transaction hash on success.*

---

## 🔮 Future Endpoints (Planned)

The following endpoints are planned for Phase 1:

- `rnr_getShardAssignment`: View validator shard responsibility.
- `rnr_getValidatorSet`: View current BFT committee.
- `rnr_getSlashingStats`: View slashed validators.
