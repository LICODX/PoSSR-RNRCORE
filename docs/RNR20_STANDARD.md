# RNR-20 Token Standard

## Overview
RNR-20 is the token standard for the RNR blockchain, similar to:
- **ERC-20** (Ethereum)
- **BEP-20** (Binance Smart Chain)
- **SPL** (Solana)

Create unlimited custom tokens on RNR network with full ERC-20 compatibility.

---

## ✅ What's Implemented

### Core Token System
- ✅ Token registry (`internal/token/registry.go`)
- ✅ Token manager (`internal/token/manager.go`)  
- ✅ Token state management (`internal/state/token_state.go`)
- ✅ Token types & metadata (`pkg/types/token.go`)
- ✅ Transaction types (`pkg/types/token_tx.go`)

### Features
- ✅ Create custom tokens
- ✅ Transfer tokens between accounts
- ✅ Approve spending allowances
- ✅ Mint new tokens (if mintable)
- ✅ Burn tokens (if burnable)
- ✅ Symbol uniqueness enforcement
- ✅ LevelDB persistence

---

## Token Creation

### Code Example
```go
import (
    "github.com/LICODX/PoSSR-RNRCORE/internal/token"
    "github.com/LICODX/PoSSR-RNRCORE/internal/state"
    "github.com/LICODX/PoSSR-RNRCORE/pkg/types"
)

// Initialize token system
registry := token.NewRegistry()
tokenState := state.NewTokenState(db)
manager := token.NewManager(registry, tokenState)

// Create a new token
metadata := types.TokenMetadata{
    Name:          "MyToken",
    Symbol:        "MTK",
    Decimals:      18,
    InitialSupply: 1000000,
    IsMintable:    true,
    IsBurnable:    true,
}

var creator [32]byte  // Your address
token, err := manager.CreateToken(metadata, creator)
if err != nil {
    panic(err)
}

fmt.Printf("Token created: %s (%s)\n", token.Name, token.Symbol)
```

---

## Token Transfer

```go
var tokenAddr [32]byte  // Token address from creation
var from [32]byte       // Sender address
var to [32]byte         // Recipient address

// Transfer 100 tokens
err := manager.Transfer(tokenAddr, from, to, 100)
if err != nil {
    fmt.Println("Transfer failed:", err)
}
```

---

## Get Balance

```go
balance := manager.GetBalance(tokenAddr, accountAddr)
fmt.Printf("Balance: %d tokens\n", balance)
```

---

## Approve & TransferFrom

```go
// Approve spender to use 50 tokens
err := manager.Approve(tokenAddr, owner, spender, 50)

// Spender transfers from owner to recipient
err = manager.TransferFrom(tokenAddr, spender, owner, recipient, 30)
```

---

## Mint Tokens (If Mintable)

```go
// Only creator can mint
err := manager.Mint(tokenAddr, recipient, 1000, creator)
```

---

## Burn Tokens (If Burnable)

```go
// Burn 50 tokens from account
err := manager.Burn(tokenAddr, account, 50)
```

---

## API Endpoints (Future)

### POST /api/token/create
Create new token (costs 100 RNR)

**Request:**
```json
{
  "name": "MyToken",
  "symbol": "MTK",
  "decimals": 18,
  "initialSupply": 1000000,
  "isMintable": true,
  "isBurnable": true
}
```

**Response:**
```json
{
  "success": true,
  "tokenAddress": "0xabc123...",
  "symbol": "MTK"
}
```

### POST /api/token/transfer
Transfer tokens

**Request:**
```json
{
  "tokenAddress": "0xabc123...",
  "to": "rnr1...",
  "amount": 100
}
```

### GET /api/token/{address}
Get token info

**Response:**
```json
{
  "address": "0xabc123...",
  "name": "MyToken",
  "symbol": "MTK",
  "totalSupply": 1000000,
  "decimals": 18,
  "creator": "rnr1...",
  "isMintable": true,
  "isBurnable": true
}
```

### GET /api/tokens
List all tokens

**Response:**
```json
{
  "tokens": [
    {"address": "0x...", "name": "Token A", "symbol": "TKA"},
    {"address": "0x...", "name": "Token B", "symbol": "TKB"}
  ],
  "total": 2
}
```

---

## Economics

### Token Creation Fee
- **Cost:** 100 RNR (prevents spam)
- **Paid to:** Miners (distributed among block winners)
- **Effect:** Drives RNR demand

### Token Transfer Fee
- **Cost:** 1 RNR minimum (same as regular TX)
- **Paid to:** Miners
- **Effect:** All token activity requires RNR

### RNR Utility
1. Required for token creation → increased demand
2. Required for all token transfers → sustained demand
3. Future DeFi operations → massive demand
4. Platform currency → value appreciation

---

## Use Cases

### Stablecoins
```go
usdToken := types.TokenMetadata{
    Name:          "USD Coin on RNR",
    Symbol:        "USDR",
    Decimals:      6,
    InitialSupply: 10000000, // 10M USDR
}
```

### Utility Tokens
```go
gameToken := types.TokenMetadata{
    Name:          "GameFi Token",
    Symbol:        "GAME",
    Decimals:      18,
    IsMintable:    true, // For player rewards
}
```

### Governance Tokens
```go
govToken := types.TokenMetadata{
    Name:          "RNR DAO",
    Symbol:        "RNRDAO",
    TotalSupply:   100000000,
    IsBurnable:    true, // Deflationary
}
```

---

## Integration Guide

### 1. Initialize Token System in Node
```go
// In cmd/rnr-node/main.go
tokenRegistry := token.NewRegistry()
tokenState := state.NewTokenState(db)
tokenManager := token.NewManager(tokenRegistry, tokenState)
```

### 2. Add API Routes (Future)
```go
// In server.go
http.HandleFunc("/api/token/create", handleTokenCreate)
http.HandleFunc("/api/token/transfer", handleTokenTransfer)
http.HandleFunc("/api/token/", handleTokenInfo)
http.HandleFunc("/api/tokens", handleTokensList)
```

### 3. Update Wallet UI (Future)
```html
<!-- Multi-token display -->
<div class="tokens">
    <h3>💰 My Tokens</h3>
    <div class="token-item">
        <span>RNR</span>
        <span>1,000.00</span>
    </div>
    <div class="token-item">
        <span>MTK</span>
        <span>500.00</span>
        <button>Send</button>
    </div>
</div>

<!-- Token creation form -->
<form id="create-token">
    <input id="name" placeholder="Token Name">
    <input id="symbol" placeholder="Symbol">
    <input id="supply" placeholder="Total Supply">
    <button>Create (100 RNR)</button>
</form>
```

---

## Security Considerations

### Symbol Uniqueness
- Each symbol can only be registered once
- Prevents confusion between tokens

### Balance Validation
- All transfers check sufficient balance
- No negative balances possible
- Overflow protection built-in

### Creator Permissions
- Only creator can mint (if mintable enabled)
- Only token holder can burn their tokens
- Cannot mint/burn if not enabled

### Re-entrancy Protection
- State updates atomic
- No external calls during state changes

---

## Testing

### Create Test Token
```bash
# Create 3 test tokens
go run scripts/create-test-tokens.go
```

### Transfer Test
```bash
# Transfer tokens between accounts
go run scripts/test-transfer.go
```

---

## Future Enhancements

### Phase 2: NFT Support (RNR-721)
- Non-fungible tokens
- Digital collectibles
- Gaming assets

### Phase 3: DeFi Primitives
- Decentralized exchange (DEX)
- Liquidity pools
- Lending/borrowing
- Yield farming

### Phase 4: Cross-Chain
- Bridge to Ethereum
- Bridge to BSC
- Wrapped tokens (wETH, wBTC)

---

## Status

**Core Implementation:** ✅ 100% Complete
**API Endpoints:** ⏰ Pending integration
**Wallet UI:** ⏰ Pending implementation
**Testing:** ⏰ Pending test suite

**Ready for:** Internal testing & development
**Production Ready:** After API/UI integration

---

## Contributing

To complete RNR-20:
1. Integrate API endpoints into dashboard server
2. Add multi-token wallet UI
3. Create comprehensive test suite
4. Deploy to testnet
5. Conduct security audit

---

**RNR-20: Making RNR the next Ethereum! 🚀**

**Version:** 1.0.0-dev  
**Status:** Core Complete, Integration Pending  
**Last Updated:** 2026-01-18
