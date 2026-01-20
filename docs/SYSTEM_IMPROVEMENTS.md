# System Improvements - 17 Jan 2026

## ✅ IMPLEMENTED CHANGES

### 1. Registration System Overhaul

**Old Behavior:**
- Guest nodes received **10 RNR** initial funding
- No limit on number of registrations
- No tracking of registered nodes

**New Behavior:**
- Guest nodes receive **1 RNR** initial funding (minimal balance for transactions)
- **Maximum 10 guest nodes** allowed
- Registration counter tracks usage (X/10)
- HTTP 429 error when limit reached

**Code Changes:**
- `internal/dashboard/server.go`:
  - Added `registeredNodesCount` field to Server struct
  - Added limit check before registration
  - Reduced funding transaction: `10 RNR → 1 RNR`
  - Enhanced logging: `"Sent 1 RNR to ... (Total: X/10)"`

---

### 2. CLI Transparency Enhancement

**Old Behavior:**
```
🔨 Mining on top of Block #5 [Diff: 100] with 1 txs (inc. Coinbase)
```
No information about WHO receives the reward.

**New Behavior:**
```
🔨 Mining Block #6 [Diff: 100] | TXs: 1 (inc. Coinbase) | Reward → rnr1pq03gqs8zg0sgqg...
```

Shows:
- ✅ Which block being mined
- ✅ Transaction count
- ✅ **Reward recipient address** (first 20 chars + "...")

**Code Changes:**
- `cmd/rnr-node/main.go`:
  - Enhanced mining log with reward recipient
  - Shows `nodeWallet.Address[:20]+"..."` for clarity

---

### 3. Genesis Wallet Clarification

**Confirmed Design:**
- Genesis wallet is **normal wallet** (like any other node)
- Special privilege: **Only creates Genesis Block** (block #0)
- After Genesis: Must mine like regular node to earn coins
- Gets **100 RNR reward per mined block** (same as others)
- **Cost of registration:** 1 RNR × 10 guests = 10 RNR total
- **ROI:** Earns 100 RNR per block mined (10 RNR cost recovered in < 1 block)

---

## 🎨 DASHBOARD/GUI IMPROVEMENTS NEEDED

User feedback: *"Wallet GUI kurang baik tampilan dan fungsinya"*

### Current State Analysis

**Dashboard Location:** `internal/dashboard/static/`

**Current Features:**
- Basic blockchain stats (height, difficulty)
- Mempool transaction count
- Mock TPS calculation

**Missing Features (User Complaints):**

1. **Wallet Balance Display**
   - ❌ Cannot see current RNR balance
   - ❌ No transaction history
   - ❌ No address display

2. **Mining Status**
   - ❌ No indication if node is mining
   - ❌ No hash rate display
   - ❌ No reward history

3. **Network Visibility**
   - ❌ Peer count shows "0" (bug)
   - ❌ No network topology view
   - ❌ No registration status for guests

4. **Transaction Creation**
   - ❌ Cannot send RNR to other addresses
   - ❌ No transaction form

5. **Visual Design**
   - Current: Basic HTML/CSS
   - Needed: Modern, responsive UI

---

## 📋 RECOMMENDED GUI ENHANCEMENTS

### Priority 1: Wallet Information Panel

**Add to Dashboard:**
```
┌─────────────────────────────────┐
│ 💰 Wallet Balance               │
│ 150.00 RNR                      │
│                                 │
│ 📍 Address                      │
│ rnr1pq03gqs8zg0sgqg...         │
│ [Copy] [QR Code]                │
└─────────────────────────────────┘
```

**API Endpoint Needed:** `GET /api/wallet`
```json
{
  "address": "rnr1...",
  "balance": 150.00,
  "nonce": 5
}
```

### Priority 2: Mining Status Widget

```
┌─────────────────────────────────┐
│ ⛏️ Mining Status                │
│ ● Active (Mining Block #127)   │
│                                 │
│ Rewards Earned: 5 blocks        │
│ Total Mined: 500 RNR            │
└─────────────────────────────────┘
```

### Priority 3: Transaction History

```
┌─────────────────────────────────┐
│ 📜 Recent Transactions          │
│                                 │
│ ✅ +100 RNR (Coinbase)          │
│    Block #126 - 2 min ago       │
│                                 │
│ ✅ +100 RNR (Coinbase)          │
│    Block #125 - 3 min ago       │
│                                 │
│ ❌ -1 RNR (Sent)                │
│    To: rnr1abc... - 5 min ago   │
└─────────────────────────────────┘
```

### Priority 4: Send Transaction Form

```
┌─────────────────────────────────┐
│ 💸 Send RNR                     │
│                                 │
│ To Address:                     │
│ [rnr1_________________]         │
│                                 │
│ Amount (RNR):                   │
│ [________] RNR                  │
│                                 │
│       [Send Transaction]        │
└─────────────────────────────────┘
```

**API Endpoint Needed:** `POST /api/send`
```json
{
  "to": "rnr1...",
  "amount": 5.0
}
```

### Priority 5: Network Peer Display

```
┌─────────────────────────────────┐
│ 🌐 Network Peers (5)            │
│                                 │
│ ● Node 2 (127.0.0.1:3001)       │
│ ● Node 3 (127.0.0.1:3002)       │
│ ● Node 4 (127.0.0.1:3003)       │
│ ● Node 5 (127.0.0.1:3004)       │
│ ● Node 6 (127.0.0.1:3005)       │
└─────────────────────────────────┘
```

**Fix:** Get actual peer count from LibP2P

---

## 🎯 NEXT ACTIONS

**Untuk user:**
1. Restart network dengan executable baru:
   ```
   .\RUN_15_NODES.bat
   ```

2. Verify new behavior:
   - Check logs for: `"Sent 1 RNR to ... (Total: X/10)"`
   - Check logs for: `"Reward → rnr1..."`
   - Try registering 11th node (should fail with "limit reached")

**Untuk developer (next session):**
1. Implement Priority 1-2 dashboard features (wallet info, mining status)
2. Fix peer count display
3. Add transaction history API + UI
4. Create send transaction form
5. Modernize CSS (optional: use TailwindCSS or Bootstrap)

---

**Status:** ✅ Core improvements implemented
**Build:** ✅ Successful  
**Ready for:** Testing with new executable
