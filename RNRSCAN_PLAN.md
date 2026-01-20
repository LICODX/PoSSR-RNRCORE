# RNRScan - Block Explorer Implementation Plan

## Architecture Overview

RNRScan will be a comprehensive blockchain explorer with the following components:

### Pages
1. **Homepage** (`/`) - Latest blocks & transactions
2. **Block Details** (`/block/:height`) - Block info, transactions
3. **Transaction Details** (`/tx/:hash`) - TX details, status
4. **Address Info** (`/address/:addr`) - Balance, TX history
5. **Search** - Universal search (block/tx/address)

### API Endpoints
1. `GET /api/blocks?page=1&limit=20` - List recent blocks
2. `GET /api/block/:height` - Block details
3. `GET /api/transactions?page=1&limit=20` - Recent TXs
4. `GET /api/tx/:hash` - Transaction details
5. `GET /api/address/:addr` - Address info
6. `GET /api/search?q=...` - Universal search

### Tech Stack
- **Backend:** Go (existing server)
- **Frontend:** Vanilla JS + Modern CSS
- **Real-time:** Polling every 3 seconds
- **Design:** Glassmorphism + Dark theme

## Implementation Steps
1. ✅ Create API endpoints in `dashboard/server.go`
2. ✅ Create explorer HTML pages in `dashboard/static/explorer/`
3. ✅ Add CSS styling (Etherscan-inspired but modern)
4. ✅ Implement search functionality
5. ✅ Add real-time updates

Let's start building!
