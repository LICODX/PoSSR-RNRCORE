# RNR Security Best Practices

## Critical Security Principles

**Defense in Depth:** Multiple layers of security  
**Principle of Least Privilege:** Minimum necessary access  
**Zero Trust:** Verify everything, trust nothing  

---

## Wallet Security

### Private Key Protection

**CRITICAL: Your private key = Your funds!**

1. **Password Protection** ✅
   ```bash
   # ALWAYS use wallet encryption
   rnr-node -wallet-password "Strong_Password_16+_Chars!"
   
   # NEVER use cleartext wallets in production
   ```

2. **Password Requirements**
   - Minimum 16 characters
   - Mix: uppercase, lowercase, numbers, symbols
   - Unique (never reuse)
   - Not in dictionary
   - Changed quarterly

3. **Secure Storage**
   ```bash
   # Encrypted wallet file
   chmod 600 ~/.rnr/data/node_wallet.json
   
   # Store password in vault (NOT in code!)
   # - HashiCorp Vault
   # - AWS Secrets Manager
   # - Azure Key Vault
   ```

4. **Backup Strategy**
   ```bash
   # Daily encrypted backups
   tar -czf wallet-backup-$(date +%Y%m%d).tar.gz.enc \
       --encrypt-with-passphrase ~/.rnr/data/node_wallet.json
   
   # Store in 3+ locations:
   # 1. Local encrypted drive
   # 2. Cloud storage (encrypted)
   # 3. Offline USB (encrypted)
   ```

5. **Cold Storage**
   - Generate wallet on air-gapped machine
   - Store mnemonic offline (paper wallet)
   - Use hardware wallet for large amounts
   - Hot wallet: < 10% of holdings

### Genesis Wallet Security (CRITICAL!)

```bash
# NEVER hardcode Genesis mnemonic!

# Generate unique mnemonic on air-gapped machine
go run cmd/genesis-wallet/main.go > genesis.secret

# Store securely:
# 1. Paper backup (fireproof safe)
# 2. Encrypted USB (bank vault)
# 3. Multi-sig recovery (3-of-5 founders)

# Use environment variable
export GENESIS_MNEMONIC="<from secure vault>"

# OR secure file (NOT in git!)
echo "mnemonic..." > ~/.rnr/genesis.secret
chmod 400 ~/.rnr/genesis.secret
```

---

## Node Security

### Firewall Configuration

```bash
# Ubuntu/Debian
sudo ufw default deny incoming
sudo ufw default allow outgoing

# Allow P2P port ONLY
sudo ufw allow 9900/tcp comment 'RNR P2P'

# Dashboard should use reverse proxy (NOT direct!)
# sudo ufw allow 8080/tcp  # DON'T DO THIS!

sudo ufw enable
sudo ufw status
```

### SSH Hardening

```bash
# Disable password authentication
sudo nano /etc/ssh/sshd_config
# Set: PasswordAuthentication no
# Set: PubkeyAuthentication yes
# Set: PermitRootLogin no

# Use SSH keys only
ssh-keygen -t ed25519 -C "rnr-node-admin"

# Restrict SSH to specific IPs
sudo ufw allow from YOUR_IP to any port 22

sudo systemctl restart ssh
```

### System Hardening

```bash
# 1. Disable unused services
sudo systemctl disable bluetooth
sudo systemctl disable cups

# 2. Enable automatic security updates
sudo apt install unattended-upgrades
sudo dpkg-reconfigure -plow unattended-upgrades

# 3. Install fail2ban (brute force protection)
sudo apt install fail2ban
sudo systemctl enable fail2ban

# 4. Enable AppArmor/SELinux
sudo systemctl enable apparmor

# 5. Disable IPv6 (if not used)
echo "net.ipv6.conf.all.disable_ipv6 = 1" | sudo tee -a /etc/sysctl.conf
sudo sysctl -p
```

---

## Network Security

### Reverse Proxy (REQUIRED for Production)

**nginx Configuration:**
```nginx
# /etc/nginx/sites-available/rnr
server {
    listen 443 ssl http2;
    server_name dashboard.yournode.com;
    
    # SSL Certificate (Let's Encrypt)
    ssl_certificate /etc/letsencrypt/live/yournode.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/yournode.com/privkey.pem;
    
    # Security headers
    add_header Strict-Transport-Security "max-age=31536000" always;
    add_header X-Frame-Options "DENY" always;
    add_header X-Content-Type-Options "nosniff" always;
    
    # Rate limiting
    limit_req_zone $binary_remote_addr zone=api:10m rate=10r/s;
    limit_req zone=api burst=20 nodelay;
    
    # Proxy to RNR dashboard
    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

### DDoS Protection

```bash
# 1. CloudFlare (Recommended)
# - Enable "Under Attack" mode during DDoS
# - Use firewall rules to block bad actors
# - Rate limit API endpoints

# 2. iptables rate limiting
sudo iptables -A INPUT -p tcp --dport 9900 -m state --state NEW -m recent --set
sudo iptables -A INPUT -p tcp --dport 9900 -m state --state NEW -m recent --update --seconds 60 --hitcount 10 -j DROP
```

### VPN/Private Network

```bash
# Use WireGuard for node-to-node communication
sudo apt install wireguard

# Configure private network
# Only expose P2P port to trusted peers
```

---

## Operational Security

### Access Control

```bash
# Create dedicated user (NOT root!)
sudo useradd -m -s /bin/bash rnr-node
sudo usermod -aG sudo rnr-node

# Restrict file permissions
sudo chown -R rnr-node:rnr-node /blockchain/rnr/
sudo chmod 700 /blockchain/rnr/data
```

### Logging & Monitoring

```bash
# 1. Centralized logging
# Ship logs to ELK/Splunk/DataDog

# 2. Security event monitoring
sudo auditctl -w /blockchain/rnr/data -p wa -k rnr_wallet

# 3. File integrity monitoring
sudo apt install aide
sudo aideinit
sudo aide --check

# 4. Intrusion detection
sudo apt install ossec-hids
```

### Alerts

```yaml
# alerts.yml
alerts:
  - name: UnauthorizedAccess
    condition: failed_ssh_attempts > 5
    action: email_admin + block_ip
  
  - name: WalletFileModified
    condition: wallet_file_changed == true
    action: email_admin + snapshot_system
  
  - name: SuddenPeerDrop
    condition: peer_count < 3
    action: email_admin + restart_node
  
  - name: HighMemoryUsage
    condition: memory_usage > 95%
    action: email_admin + restart_node
```

---

## Incident Response

### Response Plan

**Tier 1: Suspected Breach**
1. Disconnect node from network
2. Take snapshot of system
3. Analyze logs for evidence
4. Rotate all credentials

**Tier 2: Confirmed Breach**
1. Shut down node immediately
2. Move funds to cold wallet
3. Wipe and rebuild system
4. Conduct forensic analysis
5. Report to community (if network-wide)

**Tier 3: Network Attack**
1. Coordinate with core devs
2. Deploy emergency patch
3. Hard fork if necessary
4. Post-mortem analysis

### Emergency Contacts

```
Security Team: security@rnr.network
PGP Key: 0x1234ABCD
Discord: @security-team
Phone: +1-XXX-XXX-XXXX (24/7)
```

---

## Compliance & Legal

### Data Protection

```bash
# GDPR Compliance (if in EU)
# - Encrypt all user data
# - Implement right to deletion
# - Maintain audit logs

# Log retention policy
find /var/log/rnr -name "*.log" -mtime +90 -delete
```

### Regulatory Compliance

- **KYC/AML:** If running exchange/custodial service
- **Tax Reporting:** Track mining income
- **License Requirements:** Check local regulations

---

## Security Checklist

**Daily:**
- [ ] Check logs for anomalies
- [ ] Verify node is online
- [ ] Monitor system resources

**Weekly:**
- [ ] Review security alerts
- [ ] Check for software updates
- [ ] Backup wallet
- [ ] Test restore procedure

**Monthly:**
- [ ] Rotate passwords
- [ ] Update firewall rules
- [ ] Review access logs
- [ ] Conduct security scan

**Quarterly:**
- [ ] Penetration testing
- [ ] Security audit
- [ ] Disaster recovery drill
- [ ] Update documentation

---

## Common Attacks & Mitigations

### 1. Private Key Theft
**Attack:** Malware steals wallet file  
**Mitigation:** Encrypted wallet + strong password

### 2. Man-in-the-Middle
**Attack:** Intercept P2P communication  
**Mitigation:** Verify peer certificates, use VPN

### 3. DDoS
**Attack:** Flood node with requests  
**Mitigation:** CloudFlare, rate limiting, firewall

### 4. Social Engineering
**Attack:** Trick operator into revealing credentials  
**Mitigation:** Security training, verify all requests

### 5. Supply Chain Attack
**Attack:** Compromised dependencies  
**Mitigation:** Verify checksums, use trusted sources

---

## Security Resources

- **Security Advisories:** https://security.rnr.network
- **Bug Bounty:** https://bounty.rnr.network
- **Security Forum:** https://forum.rnr.network/security
- **Emergency Hotline:** security@rnr.network

---

## Security Audit

**Professional Audits:**
- CertiK: https://certik.com
- Trail of Bits: https://trailofbits.com
- OpenZeppelin: https://openzeppelin.com

**Self-Audit Tools:**
- `lynis` - Security auditing
- `rkhunter` - Rootkit detection
- `ossec` - Intrusion detection

---

**Remember: Security is a process, not a product!** 🔒

**Last Updated:** 2026-01-18  
**Version:** 1.0.0
