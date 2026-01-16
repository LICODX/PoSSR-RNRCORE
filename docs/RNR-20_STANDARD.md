# RNR-20 Token Standard Specification

**Status:** Draft
**Type:** Layer 1 Protocol

## Abstract
The RNR-20 standard defines a common interface for fungible tokens on the PoSSR blockchain. It allows the creation of utility tokens, stablecoins (USDT-R), and wrapped assets (wBTC-R) that inherit the security and speed of the PoRS consensus.

## Specification

### Data Structure
Every RNR-20 token is a native object in the State Trie (not just a smart contract variable), ensuring performance.

```go
type Token struct {
    Symbol      string   // e.g. "USDT"
    Name        string   // "Tether USD (RNR)"
    TotalSupply uint64   
    Decimals    uint8    // 18
    Owner       [32]byte // Contract Address or EOA
}
```

### Methods

#### `transfer(to Address, value uint64)`
Moves `value` amount of tokens from caller to `to`.
*   **Event:** `Transfer(from, to, value)`

#### `approve(spender Address, value uint64)`
Approves `spender` to withdraw from your account multiple times, up to the `value` amount.

#### `transferFrom(from Address, to Address, value uint64)`
Moves `value` amount from `from` to `to` using the allowance mechanism.

## Wrapped Asset Architecture (Bridge)

To bring BTC/ETH/USDT to RNRCORE:

1.  **Lock & Mint:** User locks BTC in a multisig vault on the Bitcoin Network.
2.  **Oracle:** The RNR Bridge Oracle witnesses the lock.
3.  **Mint:** RNR-20 "wBTC" are minted to the user's RNR address.
4.  **Burn & Release:** User burns wBTC on RNR to release native BTC.

## Usage Example (Go SDK) (Planned)

```go
token := rnr.NewToken("USDT", "Tether", 1000000000)
tx := token.Transfer(recipient, 500)
client.SendTransaction(tx)
```
