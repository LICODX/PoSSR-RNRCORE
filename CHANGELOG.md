# PoSSR RNR-CORE - Changelog

## Version 2.1 - 17 January 2026

### 🎯 User-Requested Improvements

#### Registration System
- **CHANGED:** Guest funding reduced from 10 RNR → **1 RNR**
- **ADDED:** Registration limit of **maximum 10 guest nodes**
- **ADDED:** Registration counter displaying `(Total: X/10)`
- **SECURITY:** HTTP 429 error when registration limit exceeded

#### CLI Transparency
- **ENHANCED:** Mining logs now show Coinbase reward recipient
- **FORMAT:** `🔨 Mining Block #X | Reward → rnr1...`
- **BENEFIT:** Users can verify which wallet receives block rewards

#### Genesis Wallet Behavior
- **CLARIFIED:** Genesis is normal wallet (must mine for rewards like others)
- **ROLE:** Only creates Genesis Block (#0), then participates equally
- **ECONOMICS:** 1 RNR × 10 guests = 10 RNR cost, recovered in < 1 mined block

### 🐛 Bug Fixes (Previous Session)
- **FIXED:** Block time updated to 60 seconds (was 6 seconds)
- **FIXED:** Removed mock transaction causing nonce validation errors
- **FIXED:** Invalid nonce: expected 1, got UnixNano() issue resolved

### 📊 Technical Details

**Files Modified:**
1. `internal/dashboard/server.go`
   - Added `registeredNodesCount` tracking
   - Implemented 10-node registration limit
   - Reduced funding: `CreateTransaction(receiverHex, 1, ...)`

2. `cmd/rnr-node/main.go`
   - Enhanced mining log with reward recipient address
   - Format: `nodeWallet.Address[:20]+"..."`

3. `internal/params/constants.go`
   - Updated `BlockTime = 60` seconds

### 🔜 Known Issues / Future Improvements

#### Dashboard/GUI (User Feedback: "kurang baik")
- [ ] Add wallet balance display
- [ ] Add transaction history
- [ ] Show mining status (active/inactive)
- [ ] Fix peer count display (currently shows 0)
- [ ] Add send transaction form
- [ ] Modernize UI design

#### Network
- [ ] Peer count logging issue (cosmetic - network functions correctly)

### 📦 How to Test New Features

1. **Test Registration Limit:**
   ```bash
   # Start 15 nodes (will hit limit at 10 guests)
   .\RUN_15_NODES.bat
   
   # Check node1.log for:
   # "Sent 1 RNR to ... (Total: X/10)"
   ```

2. **Verify Reward Transparency:**
   ```bash
   # Check any node log for:
   # "Reward → rnr1pq03gqs8zg0sgqg..."
   ```

3. **Confirm 60s Block Time:**
   ```bash
   # Blocks should appear every ~60 seconds
   # Check timestamp delta in logs
   ```

---

**Previous Version:** 2.0 (Mainnet-ready with 15-node support)  
**Current Version:** 2.1 (User improvements + transparency)  
**Next Version:** 2.2 (Dashboard enhancements planned)
