# RNR Smart Contracts Guide

## Overview

RNR Smart Contracts bring programmable logic to the RNR blockchain using WebAssembly (WASM) technology. Contracts can interact with both native RNR tokens and RNR-20 tokens.

## 🔄 Execution Model: All-Nodes Validation

### Who Runs the VM?

**ALL validator nodes execute smart contracts.** This ensures:

1. **Block Creator (Miner/Validator)**:
   - Executes contracts FIRST when creating a new block
   - Includes execution results in the block (state changes, gas used, events)
   - Adds contract transactions to the block

2. **Other Validator Nodes**:
   - When receiving a new block, RE-EXECUTE all contracts in that block
   - Verify that results EXACTLY match what the block creator claimed
   - If results differ → block is REJECTED

3. **Why All Nodes?**:
   - **Consensus**: All nodes must agree on execution results
   - **Security**: Prevents malicious nodes from faking execution
   - **State Synchronization**: All nodes maintain identical state
   - **Trustless**: No need to trust any single node

This is the same model used by Bitcoin (script validation) and Ethereum 1.0 (smart contract execution).

## 🏗️ Architecture

### WASM Virtual Machine

RNR uses **Wasmer** as the WASM runtime:
- High performance execution
- Secure sandboxing
- Memory-safe
- Supports Rust, AssemblyScript, C, and more

### Gas Metering

Every operation costs gas to prevent infinite loops and DoS attacks:

| Operation | Gas Cost |
|-----------|----------|
| Contract Deploy | 100,000 |
| Contract Call | 10,000 |
| Storage Read | 100 |
| Storage Write | 1,000 |
| Transfer | 500 |
| Get Balance | 50 |
| Emit Event | 200 |

**Gas is 100-1000x cheaper than Ethereum** to encourage adoption.

## 📝 Host Functions

Contracts can call these blockchain functions:

### Core Functions

```rust
// Get RNR balance of an address
extern "C" fn get_balance(address_ptr: i32) -> i64;

// Transfer RNR to another address
extern "C" fn transfer(to_addr_ptr: i32, amount: i64) -> i32;

// Get the caller's address
extern "C" fn get_caller() -> i32;

// Get current block height
extern "C" fn get_block_height() -> i64;
```

### Storage Functions

```rust
// Read from contract storage
extern "C" fn storage_read(key_ptr: i32, key_len: i32) -> i32;

// Write to contract storage
extern "C" fn storage_write(key_ptr: i32, key_len: i32, val_ptr: i32, val_len: i32);
```

### Event System

```rust
// Emit an event for indexing
extern "C" fn emit_event(event_ptr: i32, event_len: i32);
```

## 🚀 Quick Start

### 1. Write a Contract (Rust Example)

```rust
// counter.rs
#![no_std]
#![no_main]

use core::panic::PanicInfo;

#[panic_handler]
fn panic(_info: &PanicInfo) -> ! {
    loop {}
}

static mut COUNTER: u64 = 0;

#[no_mangle]
pub extern "C" fn increment() -> u64 {
    unsafe {
        COUNTER += 1;
        COUNTER
    }
}

#[no_mangle]
pub extern "C" fn get_count() -> u64 {
    unsafe { COUNTER }
}
```

### 2. Compile to WASM

```bash
rustc --target wasm32-unknown-unknown -O --crate-type=cdylib counter.rs -o counter.wasm
```

### 3. Deploy Contract

```go
// Deploy via Go
bytecode, _ := os.ReadFile("counter.wasm")

payload := types.ContractDeployPayload{
    Bytecode: bytecode,
    InitArgs: []byte{},
}

tx := types.Transaction{
    Type:    types.TxTypeContractDeploy,
    From:    myAddress,
    Payload: encodePayload(payload),
    Gas:     200000,
}

contractAddr := deployTransaction(tx)
```

### 4. Call Contract

```go
payload := types.ContractCallPayload{
    ContractAddr: contractAddr,
    Method:       "increment",
    Args:         []byte{},
}

tx := types.Transaction{
    Type:    types.TxTypeContractCall,
    From:    myAddress,
    Payload: encodePayload(payload),
    Gas:     50000,
}

result := callTransaction(tx)
```

## 💡 Example: Token Vesting Contract

```rust
// vesting.rs
#![no_std]
#![no_main]

extern "C" {
    fn get_block_height() -> i64;
    fn transfer(to_addr_ptr: i32, amount: i64) -> i32;
    fn get_caller() -> i32;
}

static mut BENEFICIARY: [u8; 32] = [0u8; 32];
static mut RELEASE_HEIGHT: i64 = 0;
static mut AMOUNT: i64 = 0;
static mut RELEASED: bool = false;

#[no_mangle]
pub extern "C" fn init(beneficiary_ptr: i32, release_height: i64, amount: i64) {
    unsafe {
        // Store vesting parameters
        RELEASE_HEIGHT = release_height;
        AMOUNT = amount;
        // Copy beneficiary address
        // In production, read from WASM memory at beneficiary_ptr
    }
}

#[no_mangle]
pub extern "C" fn release() -> i32 {
    unsafe {
        if RELEASED {
            return 0; // Already released
        }

        let current_height = get_block_height();
        if current_height < RELEASE_HEIGHT {
            return 0; // Too early
        }

        // Transfer tokens to beneficiary
        let result = transfer(0, AMOUNT); // beneficiary address at ptr 0
        if result == 1 {
            RELEASED = true;
            return 1; // Success
        }
        return 0; // Transfer failed
    }
}

#[no_mangle]
pub extern "C" fn is_released() -> i32 {
    unsafe {
        if RELEASED { 1 } else { 0 }
    }
}

#[panic_handler]
fn panic(_info: &core::panic::PanicInfo) -> ! {
    loop {}
}
```

## 🔐 Security Best Practices

1. **Reentrancy Protection**: Use checks-effects-interactions pattern
2. **Integer Overflow**: Use safe math operations
3. **Access Control**: Validate caller addresses
4. **Gas Limits**: Set appropriate gas limits
5. **Audit**: Have contracts audited before mainnet deployment

## 🎯 Use Cases

### DeFi Applications
- Decentralized exchanges (DEX)
- Lending protocols
- Yield farming
- Liquidity pools

### Token Features
- Token vesting schedules
- Multi-signature wallets
- DAO governance
- Staking mechanisms

### NFTs & Gaming
- NFT minting and trading
- In-game asset management
- Play-to-earn mechanics
- Digital collectibles

### Enterprise
- Supply chain tracking
- Identity verification
- Document notarization
- Multi-party agreements

## 📊 Gas Estimation

Example gas costs for common operations:

```go
// Deploy a simple counter contract
DeployGas: 150,000

// Call increment() function
CallGas: 20,000

// Complex DeFi swap operation
SwapGas: 100,000

// NFT mint
MintGas: 80,000
```

**Total cost example** (at 1 gas = 0.00001 RNR):
- Deploy contract: 1.5 RNR
- 10 function calls: 0.2 RNR
- **Much cheaper than Ethereum** (100-1000x cheaper)

## 🛠️ Development Tools (Coming Soon)

- **RNR Contract CLI**: Deploy, test, and interact with contracts
- **Testing Framework**: Unit test your contracts
- **IDE Extensions**: VSCode syntax highlighting and debugging
- **Block Explorer**: View contract state and transactions
- **Remix-like IDE**: Web-based contract development

## 🔄 Hybrid Model: Native Tokens + Smart Contracts

RNR offers the best of both worlds:

| Feature | Native RNR-20 | Smart Contracts |
|---------|---------------|-----------------|
| Speed | ⚡ Instant | 🚀 Fast |
| Cost | 💰 Minimal | 💵 Low |
| Simplicity | ✅ Very Easy | 🔧 More Complex |
| Flexibility | 📋 Standard | 🎨 Unlimited |
| Security | 🛡️ Built-in | 🔒 Auditable |

**Recommendation**:
- Use **RNR-20** for standard tokens (payments, rewards, simple transfers)
- Use **Smart Contracts** for complex logic (DeFi, NFTs, governance)

## 🌐 Language Support

Write contracts in your preferred language:

- ✅ **Rust** (Recommended - best performance)
- ✅ **AssemblyScript** (JavaScript-like, easy to learn)
- ✅ **C/C++** (Maximum control)
- 🔜 **Solidity** (via EVM compatibility layer - planned)

## 📚 Next Steps

1. **Read the Technical Docs**: `docs/SMART_CONTRACTS_TECHNICAL.md`
2. **Try Examples**: `examples/contracts/`
3. **Join Community**: Discord/Telegram for support
4. **Start Building**: Deploy your first contract!

## ⚠️ Current Status

**Phase 2: WASM Integration** - IN PROGRESS 🔧

✅ Completed:
- Core contract types
- VM interface
- Gas metering
- Contract state management
- WASM runtime integration
- Host functions

🔧 In Progress:
- Full WASM memory management
- Advanced host functions
- Contract testing framework

🔜 Coming Next:
- Contract CLI tools
- Developer documentation
- Example contracts library
- Testnet contract deployment

---

**RNR Smart Contracts**: Bringing programmable money to the blockchain, powered by WASM ⚡
