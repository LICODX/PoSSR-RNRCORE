# RNR Smart Contract Examples

This directory contains example smart contracts written in Rust that demonstrate the capabilities of the RNR blockchain smart contract system.

## 📋 Available Examples

### 1. Counter Contract (`counter.rs`)

A simple counter demonstrating basic contract functionality.

**Features**:
- Global state storage
- Multiple exported methods
- Safe arithmetic operations

**Methods**:
- `init(owner_ptr)` - Initialize contract
- `increment()` - Increment counter by 1
- `add(amount)` - Add custom amount
- `decrement()` - Decrement counter by 1
- `get()` - Get current value
- `reset()` - Reset to 0 (owner only)

**Compile**:
```bash
rustc --target wasm32-unknown-unknown -O --crate-type=cdylib counter.rs -o counter.wasm
```

### 2. Vesting Contract (`vesting.rs`)

Token vesting with time-lock functionality.

**Features**:
- Time-locked token releases
- Block height-based conditions
- Blockchain interaction
- State validation

**Methods**:
- `init(beneficiary_ptr, release_height, amount)` - Setup vesting
- `release()` - Release tokens if unlocked
- `is_released()` - Check release status
- `get_release_height()` - Get unlock height
- `get_amount()` - Get vesting amount
- `blocks_until_release()` - Remaining blocks

**Compile**:
```bash
rustc --target wasm32-unknown-unknown -O --crate-type=cdylib vesting.rs -o vesting.wasm
```

## 🛠️ Prerequisites

### Install Rust
```bash
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh
```

### Add WASM Target
```bash
rustup target add wasm32-unknown-unknown
```

## 🚀 Quick Start

### 1. Compile Contract
```bash
cd examples/contracts
rustc --target wasm32-unknown-unknown -O --crate-type=cdylib counter.rs
```

### 2. Deploy to RNR (Coming Soon)
```bash
rnr-contract deploy counter.wasm
```

### 3. Interact (Coming Soon)
```bash
rnr-contract call <address> increment
rnr-contract query <address> get
```

## 📖 Learning Path

1. **Start with Counter**: Understand basic structure
2. **Study Vesting**: Learn blockchain interaction
3. **Build Your Own**: Create custom logic

## 🔗 Host Functions Available

Your contracts can call these RNR blockchain functions:

| Function | Description | Gas |
|----------|-------------|-----|
| `get_balance(addr_ptr)` | Get address balance | 50 |
| `transfer(to_ptr, amount)` | Transfer RNR | 500 |
| `storage_read(key_ptr, key_len)` | Read storage | 100 |
| `storage_write(...)` | Write storage | 1,000 |
| `emit_event(...)` | Emit event | 200 |
| `get_caller()` | Get caller address | 50 |
| `get_block_height()` | Get current block | 50 |

## 💡 Best Practices

1. **No Standard Library**: Use `#![no_std]`
2. **No Main Function**: Use `#![no_main]`
3. **Export Functions**: Use `#[no_mangle]` and `pub extern "C"`
4. **Panic Handler**: Always include panic handler
5. **Static Variables**: Use `static mut` for state (careful!)
6. **Gas Awareness**: Optimize for low gas usage

## 🎯 Use Cases

- **Counter**: Basic state management, voting systems
- **Vesting**: Token locks, salary payments, ICO vesting
- **DeFi**: Build DEX, lending, staking (advanced)
- **NFT**: Minting, trading, royalties (advanced)
- **DAO**: Governance, proposals, voting (advanced)

## 🔐 Security Notes

⚠️ **These are EXAMPLES for learning**
- Not audited for production use
- Use `unsafe` carefully
- Test thoroughly before mainnet
- Consider professional audit for real applications

## 📚 Resources

- [RNR Smart Contracts Guide](../../docs/SMART_CONTRACTS.md)
- [Rust WASM Book](https://rustwasm.github.io/docs/book/)
- [WebAssembly Spec](https://webassembly.github.io/spec/)

## 🔜 Coming Soon

- AssemblyScript examples
- DeFi contract templates
- NFT contract examples
- Testing framework
- Interactive tutorial

---

**Ready to build?** Start with `counter.rs` and deploy your first smart contract! 🚀
