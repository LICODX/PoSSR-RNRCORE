# RNR Node Operator Guide

## Overview
This guide helps you set up and run an RNR validator node for the mainnet.

---

## System Requirements

### Minimum Specifications
- **CPU:** 4 cores (2.5 GHz+)
- **RAM:** 8 GB
- **Storage:** 100 GB SSD
- **Network:** 100 Mbps up/down
- **OS:** Windows 10+, Ubuntu 20.04+, macOS 11+

### Recommended Specifications
- **CPU:** 8 cores (3.0 GHz+)
- **RAM:** 16 GB
- **Storage:** 500 GB NVMe SSD
- **Network:** 1 Gbps up/down
- **OS:** Ubuntu 22.04 LTS

---

## Installation

### Method 1: Pre-built Binary

**Windows:**
```powershell
# Download latest release
Invoke-WebRequest -Uri "https://github.com/LICODX/PoSSR-RNRCORE/releases/latest/download/rnr-node-windows.exe" -OutFile "rnr-node.exe"

# Verify checksum
Get-FileHash rnr-node.exe -Algorithm SHA256
```

**Linux:**
```bash
# Download latest release
wget https://github.com/LICODX/PoSSR-RNRCORE/releases/latest/download/rnr-node-linux

# Make executable
chmod +x rnr-node-linux

# Verify checksum
sha256sum rnr-node-linux
```

### Method 2: Build from Source

```bash
# Install Go 1.21+
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz

# Clone repository
git clone https://github.com/LICODX/PoSSR-RNRCORE.git
cd PoSSR-RNRCORE

# Build
go build -o rnr-node ./cmd/rnr-node

# Verify
./rnr-node -version
```

---

## Configuration

### Create Data Directory
```bash
mkdir -p /blockchain/rnr/data
```

### Set Wallet Password
```bash
# Option 1: Environment variable (recommended)
export WALLET_PASSWORD="YourSecurePassword123!"

# Option 2: Store in secure file
echo "YourSecurePassword123!" > /secure/wallet-password.txt
chmod 600 /secure/wallet-password.txt
```

### Configure Firewall
```bash
# Allow P2P port
sudo ufw allow 9900/tcp

# Allow dashboard (optional, use reverse proxy in production)
sudo ufw allow 8080/tcp

# Enable firewall
sudo ufw enable
```

---

## Running a Node

### Start as Validator Node
```bash
./rnr-node \
  -datadir /blockchain/rnr/data \
  -p2pport 9900 \
  -dashboard 8080 \
  -wallet-password "$WALLET_PASSWORD" \
  -peers "mainnet-seed1.rnr.network:9900,mainnet-seed2.rnr.network:9900"
```

### Start as Genesis Node (Mainnet Authority)
```bash
# Only for initial Genesis node operator!
export GENESIS_MNEMONIC="your twelve word mnemonic phrase here..."

./rnr-node \
  -genesis \
  -datadir /blockchain/rnr/data \
  -p2pport 9900 \
  -dashboard 8080 \
  -wallet-password "$WALLET_PASSWORD"
```

---

## Running as System Service

### Linux (systemd)

Create service file: `/etc/systemd/system/rnr-node.service`

```ini
[Unit]
Description=RNR Validator Node
After=network.target

[Service]
Type=simple
User=rnr
Group=rnr
WorkingDirectory=/blockchain/rnr
Environment="WALLET_PASSWORD=YourSecurePassword"
ExecStart=/usr/local/bin/rnr-node \
  -datadir /blockchain/rnr/data \
  -p2pport 9900 \
  -dashboard 8080 \
  -wallet-password "$WALLET_PASSWORD" \
  -peers "mainnet-seed1.rnr.network:9900"

Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```

Enable and start:
```bash
sudo systemctl daemon-reload
sudo systemctl enable rnr-node
sudo systemctl start rnr-node

# Check status
sudo systemctl status rnr-node

# View logs
sudo journalctl -u rnr-node -f
```

### Windows (NSSM)

```powershell
# Download NSSM
Invoke-WebRequest -Uri "https://nssm.cc/release/nssm-2.24.zip" -OutFile "nssm.zip"
Expand-Archive nssm.zip

# Install service
.\nssm\win64\nssm.exe install RNR-Node "C:\rnr\rnr-node.exe"
.\nssm\win64\nssm.exe set RNR-Node AppDirectory "C:\rnr"
.\nssm\win64\nssm.exe set RNR-Node AppParameters "-datadir C:\rnr\data -p2pport 9900"

# Start service
.\nssm\win64\nssm.exe start RNR-Node
```

---

## Monitoring Your Node

### Dashboard
Access at: `http://localhost:8080`

**Metrics to Monitor:**
- Block Height (should increase every ~10s)
- Peer Count (should be 10+)
- Mempool Size (varies with network activity)
- TPS (transactions per second)

### Prometheus Metrics
Endpoint: `http://localhost:8080/metrics`

**Key Metrics:**
```
rnr_block_height          # Current height
rnr_peer_count            # Connected peers
rnr_tps                   # Transactions/sec
rnr_mempool_size          # Pending TXs
```

### Health Check
```bash
curl http://localhost:8080/api/stats
```

Expected response:
```json
{
  "blockHeight": 12345,
  "tps": 150,
  "peers": 25,
  "mempoolSize": 42
}
```

---

## Maintenance

### Backup Wallet
```bash
# Encrypted wallet (IMPORTANT!)
cp /blockchain/rnr/data/node_wallet.json /backup/wallet-$(date +%Y%m%d).json

# Store password separately and securely
```

### Update Node
```bash
# Stop node
sudo systemctl stop rnr-node

# Backup data
tar -czf /backup/rnr-data-$(date +%Y%m%d).tar.gz /blockchain/rnr/data

# Download new version
wget https://github.com/LICODX/PoSSR-RNRCORE/releases/latest/download/rnr-node-linux
chmod +x rnr-node-linux
sudo mv rnr-node-linux /usr/local/bin/rnr-node

# Restart
sudo systemctl start rnr-node
```

### Database Cleanup
```bash
# Pruning happens automatically every 25 blocks
# To manually clean old data:
./rnr-node -datadir /blockchain/rnr/data -prune

# Check disk usage
du -sh /blockchain/rnr/data
```

---

## Security Best Practices

1. **Wallet Password**
   - Use 16+ character password
   - Mix letters, numbers, symbols
   - Never commit to version control
   - Store in secure vault (HashiCorp Vault, AWS Secrets Manager)

2. **Firewall**
   - Only open necessary ports (9900 for P2P)
   - Use reverse proxy for dashboard (nginx + HTTPS)
   - Enable DDoS protection (CloudFlare)

3. **System Hardening**
   - Keep OS updated
   - Disable unused services
   - Use SSH key authentication
   - Enable fail2ban

4. **Backups**
   - Daily automated backups
   - Store in multiple locations
   - Test recovery process monthly

5. **Monitoring**
   - Set up alerts (PagerDuty, email)
   - Monitor disk space
   - Track peer count
   - Watch for errors in logs

---

## Troubleshooting

See [TROUBLESHOOTING.md](TROUBLESHOOTING.md) for common issues and solutions.

---

## Support

- **Documentation:** https://docs.rnr.network
- **Discord:** https://discord.gg/rnr
- **GitHub Issues:** https://github.com/LICODX/PoSSR-RNRCORE/issues
- **Email:** support@rnr.network

---

**Last Updated:** 2026-01-18  
**Version:** 1.0.0
