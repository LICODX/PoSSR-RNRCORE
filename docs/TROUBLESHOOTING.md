# RNR Troubleshooting Guide

## Quick Diagnosis

**Node won't start?** → [#Node Startup Issues](#node-startup-issues)  
**No peers connecting?** → [#Network Issues](#network-issues)  
**Not mining blocks?** → [#Mining Issues](#mining-issues)  
**Wallet problems?** → [#Wallet Issues](#wallet-issues)  
**Performance slow?** → [#Performance Issues](#performance-issues)

---

## Node Startup Issues

### Error: "Failed to open database"

**Symptoms:**
```
Error: leveldb: corrupted database
Failed to open database: corruption detected
```

**Causes:**
- Database corruption
- Unclean shutdown
- Disk full

**Solutions:**
```bash
# 1. Check disk space
df -h /blockchain/rnr/data

# 2. Try repair
./rnr-node -datadir /blockchain/rnr/data -repair-db

# 3. If repair fails, restore from backup
rm -rf /blockchain/rnr/data/chaindata
tar -xzf /backup/rnr-data-latest.tar.gz -C /blockchain/rnr/data

# 4. Last resort: resync from genesis
rm -rf /blockchain/rnr/data
./rnr-node -datadir /blockchain/rnr/data
```

---

### Error: "Wallet not initialized"

**Symptoms:**
```
[ERROR] Wallet not initialized
No wallet found at path
```

**Causes:**
- First run (no wallet created yet)
- Wallet file deleted
- Wrong datadir path

**Solutions:**
```bash
# 1. Check wallet exists
ls -la /blockchain/rnr/data/node_wallet.json

# 2. If exists but not loading, check password
./rnr-node -wallet-password "YourPassword"

# 3. If deleted, restore from backup
cp /backup/node_wallet.json /blockchain/rnr/data/

# 4. If never created, let node create one
./rnr-node -wallet-password "NewPassword123"
```

---

### Error: "Port already in use"

**Symptoms:**
```
Error: bind: address already in use
Failed to start P2P server on port 9900
```

**Causes:**
- Another instance running
- Port not released from previous crash

**Solutions:**
```bash
# 1. Check what's using the port
netstat -tulpn | grep 9900
lsof -i :9900

# 2. Kill previous instance
pkill rnr-node

# 3. Use different port
./rnr-node -p2pport 9901

# 4. Wait for port release (30-60 seconds)
sleep 60 && ./rnr-node
```

---

## Network Issues

### No Peers Connecting

**Symptoms:**
- Peer count stays at 0
- Dashboard shows "No peers"
- Not receiving blocks

**Diagnosis:**
```bash
# 1. Check firewall
sudo ufw status | grep 9900

# 2. Test port accessibility
nc -zv YOUR_PUBLIC_IP 9900

# 3. Check logs
grep "peer" ~/.rnr/logs/node.log | tail -20
```

**Solutions:**
```bash
# 1. Open firewall port
sudo ufw allow 9900/tcp

# 2. Check NAT/router port forwarding
# Forward external 9900 → internal 9900

# 3. Use explicit peer list
./rnr-node -peers "seed1.rnr.network:9900,seed2.rnr.network:9900"

# 4. Check network connectivity
ping seed1.rnr.network
```

---

### Frequent Disconnections

**Symptoms:**
- Peers connect then disconnect
- "Connection reset by peer" errors
- Unstable peer count

**Causes:**
- Network instability
- NAT timeout
- Firewall issues

**Solutions:**
```bash
# 1. Increase connection limits
ulimit -n 65536

# 2. Adjust TCP keepalive
sudo sysctl -w net.ipv4.tcp_keepalive_time=60

# 3. Use stable peers only
./rnr-node -peers "stable-peer1:9900,stable-peer2:9900"

# 4. Check network quality
mtr seed1.rnr.network
```

---

### Slow Block Propagation

**Symptoms:**
- Receiving blocks late
- Often on wrong fork
- Many orphaned blocks

**Diagnosis:**
```bash
# Check network latency
ping -c 100 seed1.rnr.network | tail -1

# Expected: < 100ms average
```

**Solutions:**
```bash
# 1. Connect to geographically closer peers
./rnr-node -peers "asia-seed:9900,eu-seed:9900"

# 2. Increase bandwidth allocation
# Upgrade internet plan

# 3. Use VPN to closer region
# WireGuard to AWS region near other nodes

# 4. Check system time sync
timedatectl status
# If not synced: sudo timedatectl set-ntp on
```

---

## Mining Issues

### Not Finding Blocks

**Symptoms:**
- Zero blocks found after 24+ hours
- Mining shows "active" but no rewards

**Diagnosis:**
```bash
# 1. Check if mining is actually running
curl http://localhost:8080/api/mining

# 2. Check peer count (need 5+)
curl http://localhost:8080/api/stats | jq '.peers'

# 3. Check wallet has funds for fees
curl http://localhost:8080/api/wallet | jq '.balance'
```

**Solutions:**
```bash
# 1. Ensure wallet has RNR for tx fees
# Get from faucet or transfer

# 2. Verify good network connection
# Need < 100ms latency to peers

# 3. Check CPU is not throttled
cat /sys/devices/system/cpu/cpu*/cpufreq/scaling_governor
# Should be "performance" not "powersave"

# 4. Increase hardware resources
# More CPU cores = higher chance
```

---

### Low Hashrate

**Symptoms:**
- Dashboard shows low hashrate
- Slower than similar hardware

**Solutions:**
```bash
# 1. Close other applications
killall chrome firefox

# 2. Set CPU to performance mode
sudo cpupower frequency-set -g performance

# 3. Disable CPU throttling
echo performance | sudo tee /sys/devices/system/cpu/cpu*/cpufreq/scaling_governor

# 4. Check thermal throttling
sensors | grep -i temp
# If > 80°C, improve cooling

# 5. Allocate more cores
taskset -c 0-7 ./rnr-node  # Use cores 0-7
```

---

### Orphaned Blocks

**Symptoms:**
- Finding blocks but they're not in main chain
- "Block rejected: already exists" errors

**Causes:**
- Network latency too high
- Clock drift
- Fork conflicts

**Solutions:**
```bash
# 1. Sync system time
sudo ntpdate -u pool.ntp.org
sudo timedatectl set-ntp on

# 2. Reduce network latency
# Connect to closer peers

# 3. Check clock offset
ntpq -p
# Offset should be < 100ms

# 4. Increase network bandwidth
# Upgrade to faster connection
```

---

## Wallet Issues

### Wrong Password Error

**Symptoms:**
```
Error: invalid password
Failed to decrypt wallet
```

**Solutions:**
```bash
# 1. Verify password
# Check for typos, caps lock

# 2. Try environment variable
export WALLET_PASSWORD="YourPassword"
./rnr-node -wallet-password "$WALLET_PASSWORD"

# 3. If forgotten, restore from backup
cp /backup/node_wallet.json /blockchain/rnr/data/

# 4. If backup also encrypted with unknown password
# Create new wallet, lose old funds :(
rm /blockchain/rnr/data/node_wallet.json
./rnr-node -wallet-password "NewPassword"
```

---

### Balance Not Updating

**Symptoms:**
- Mining blocks but balance stays 0
- Transactions sent but balance unchanged

**Diagnosis:**
```bash
# 1. Check blockchain height
curl http://localhost:8080/api/stats | jq '.blockHeight'

# 2. Check if syncing
# Height should increase every ~10s

# 3. Check transaction pool
curl http://localhost:8080/api/transactions
```

**Solutions:**
```bash
# 1. Wait for full sync
# Can take hours if far behind

# 2. Restart node
sudo systemctl restart rnr-node

# 3. Clear cache and resync
rm -rf /blockchain/rnr/data/cache
./rnr-node

# 4. Check logs for state errors
grep "state" ~/.rnr/logs/node.log | tail -50
```

---

### Transaction Stuck in Mempool

**Symptoms:**
- Sent transaction but never confirmed
- TX shows "pending" for hours

**Causes:**
- Fee too low
- Nonce conflict
- Network congestion

**Solutions:**
```bash
# 1. Check transaction status
curl http://localhost:8080/api/tx/YOUR_TX_HASH

# 2. If fee too low, send replacement TX with higher fee
# (with same nonce)

# 3. Wait for mempool to clear
# Typically 1-2 blocks

# 4. Restart node to refresh mempool
sudo systemctl restart rnr-node
```

---

## Performance Issues

### High Memory Usage

**Symptoms:**
- Node using 10+ GB RAM
- OOM (Out of Memory) kills
- System becomes sluggish

**Diagnosis:**
```bash
# Check memory usage
free -h
ps aux | grep rnr-node | awk '{print $4}'
```

**Solutions:**
```bash
# 1. Increase swap space
sudo dd if=/dev/zero of=/swapfile bs=1G count=8
sudo chmod 600 /swapfile
sudo mkswap /swapfile
sudo swapon /swapfile

# 2. Enable pruning (if disabled)
./rnr-node -prune

# 3. Restart node periodically
sudo systemctl restart rnr-node

# 4. Upgrade RAM
# Minimum 8GB, recommended 16GB
```

---

### High Disk Usage

**Symptoms:**
- Disk over 80% full
- "No space left on device" errors

**Solutions:**
```bash
# 1. Check disk usage
df -h /blockchain/rnr/data
du -sh /blockchain/rnr/data/*

# 2. Enable automatic pruning
# Pruning happens every 25 blocks automatically

# 3. Clean old logs
find ~/.rnr/logs -name "*.log" -mtime +30 -delete

# 4. Move to larger disk
# Stop node, rsync data, update path

# 5. Upgrade storage
# Minimum 100GB, recommended 500GB SSD
```

---

### Slow Query Performance

**Symptoms:**
- Dashboard loads slowly
- API requests timeout
- Block explorer laggy

**Solutions:**
```bash
# 1. Optimize database
./rnr-node -datadir /blockchain/rnr/data -compact-db

# 2. Use SSD instead of HDD
# 10-100x faster!

# 3. Increase database cache
# (if option available in config)

# 4. Use read replicas
# Deploy multiple nodes, query different ones
```

---

## Database Issues

### Corrupted Database

**Symptoms:**
```
leveldb: corruption detected
Invalid block stored at height X
```

**Solutions:**
```bash
# 1. Try repair
./rnr-node -repair-db

# 2. Restore from backup
sudo systemctl stop rnr-node
rm -rf /blockchain/rnr/data/chaindata
tar -xzf /backup/latest-backup.tar.gz -C /blockchain/rnr/data
sudo systemctl start rnr-node

# 3. Resync from network (last resort)
rm -rf /blockchain/rnr/data/chaindata
./rnr-node  # Will resync from genesis
```

---

## System Issues

### Clock Drift

**Symptoms:**
- "Block timestamp too far in future" errors
- Blocks constantly rejected

**Diagnosis:**
```bash
# Check time sync
timedatectl status

# Check NTP offset
ntpq -p
```

**Solutions:**
```bash
# 1. Enable NTP
sudo timedatectl set-ntp on

# 2. Force immediate sync
sudo ntpdate -u pool.ntp.org

# 3. Install chrony (better than NTP)
sudo apt install chrony
sudo systemctl enable chronyd
```

---

### File Descriptor Limits

**Symptoms:**
```
Error: too many open files
Cannot accept connection
```

**Solutions:**
```bash
# 1. Check current limit
ulimit -n

# 2. Increase for session
ulimit -n 65536

# 3. Increase permanently
echo "* soft nofile 65536" | sudo tee -a /etc/security/limits.conf
echo "* hard nofile 65536" | sudo tee -a /etc/security/limits.conf

# 4. For systemd service
# Add to [Service] section:
# LimitNOFILE=65536
```

---

## Getting Help

### Collect Diagnostic Info

```bash
# Run diagnostic script
cat > diagnose.sh << 'EOF'
#!/bin/bash
echo "=== System Info ==="
uname -a
df -h
free -h

echo "=== Node Status ==="
curl -s http://localhost:8080/api/stats | jq .

echo "=== Recent Logs ==="
tail -100 ~/.rnr/logs/node.log

echo "=== Network ==="
netstat -tulpn | grep 9900
EOF

chmod +x diagnose.sh
./diagnose.sh > diagnostic-report.txt
```

### Support Channels

- **Discord:** https://discord.gg/rnr (#support channel)
- **GitHub Issues:** https://github.com/LICODX/PoSSR-RNRCORE/issues
- **Forum:** https://forum.rnr.network
- **Email:** support@rnr.network

**When reporting issues, include:**
1. Operating system & version
2. Node version (`./rnr-node -version`)
3. Diagnostic report (above)
4. Steps to reproduce
5. Expected vs actual behavior

---

## Emergency Recovery

### Complete System Recovery

```bash
# 1. Backup wallet (CRITICAL!)
cp /blockchain/rnr/data/node_wallet.json ~/wallet-backup.json

# 2. Stop node
sudo systemctl stop rnr-node

# 3. Purge all data
rm -rf /blockchain/rnr/data/*

# 4. Restore wallet
cp ~/wallet-backup.json /blockchain/rnr/data/node_wallet.json

# 5. Restart node (will resync)
sudo systemctl start rnr-node

# 6. Monitor progress
journalctl -u rnr-node -f
```

---

**Still stuck? Ask in Discord #support!** 💬

**Last Updated:** 2026-01-18  
**Version:** 1.0.0
