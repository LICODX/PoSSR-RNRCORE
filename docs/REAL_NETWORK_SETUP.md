# 🌍 Real Network Setup Guide (Public Internet/WAN)

This guide explains how to connect multiple PoSSR RNRCORE nodes across the **Public Internet**, allowing you to verify the blockchain with friends or colleagues in different physical locations.

---

## 📋 Prerequisites

1.  **Public IP Address**: At least one node (the "Genesis" node) needs a public IP address.
2.  **Open Ports**: You must allow **TCP Port 8001** through your firewall/router.
3.  **Go Installed**: All computers must have Go 1.20+ installed.

---

## 🛠️ Option A: Using a Cloud VPS (Recommended)

The easiest way to host a public node is using a Virtual Private Server (VPS) like DigitalOcean, AWS, Linode, or Vultr.

### 1. Rent a Linux VPS
- **OS**: Ubuntu 22.04 LTS (Recommended)
- **Specs**: 2 vCPU, 4GB RAM (Minimum for stable consensus)

### 2. Configure Firewall (UFW)
On your VPS terminal, allow the necessary ports:
```bash
sudo ufw allow 8001/tcp    # P2P Port
sudo ufw allow 9101/tcp    # Dashboard (Optional)
sudo ufw enable
```

### 3. Run the Genesis Node
Clone and build the project on your VPS:
```bash
git clone https://github.com/LICODX/PoSSR-RNRCORE.git
cd PoSSR-RNRCORE
go build -o rnr-node ./cmd/rnr-node
./rnr-node --port 8001
```

### 4. Get Your Address
Look at the startup logs for the "LibP2P identity":
```text
🌐 LibP2P GossipSub node started
   ID: QmYourNodeID123456789...
   Addresses:
     /ip4/127.0.0.1/tcp/8001/p2p/QmYourNodeID...
     /ip4/203.0.113.5/tcp/8001/p2p/QmYourNodeID...  <-- COPY THIS ONE (Public IP)
```

---

## 🏠 Option B: Home Network (Port Forwarding)

If you want to run a public node from your home computer:

### 1. Find Your Local IP
- **Windows**: Open Command Prompt, type `ipconfig`. Look for `IPv4 Address` (e.g., `192.168.1.5`).
- **Mac/Linux**: Type `ifconfig`.

### 2. Configure Router Port Forwarding
1.  Log in to your router admin panel (usually `192.168.1.1` or `192.168.0.1`).
2.  Find **"Port Forwarding"** or **"Virtual Servers"**.
3.  Add a new rule:
    - **Service Name**: RNR-P2P
    - **External Port**: 8001
    - **Internal Port**: 8001
    - **Internal IP**: Your computer's local IP (from step 1).
    - **Protocol**: TCP

### 3. Find Your Public IP
Go to [https://ifconfig.me](https://ifconfig.me) or Google "What is my IP".
*Example: 154.22.10.55*

---

## 🔗 Connecting a Peer Node (Computer B)

Now that **Computer A** (Genesis) is running and accessible externally, **Computer B** can connect to it.

1.  **On Computer B**, open a terminal.
2.  Run the node with the `--peer` flag pointing to Computer A's Public Address.

**Syntax:**
```bash
./rnr-node --port 8002 --peer /ip4/<COMPUTER_A_PUBLIC_IP>/tcp/8001/p2p/<COMPUTER_A_NODE_ID>
```

**Example:**
```bash
./rnr-node --port 8002 --peer /ip4/203.0.113.5/tcp/8001/p2p/QmHash123abc...
```

---

## ❓ Troubleshooting

### Connection Timed Out?
- **Check Firewall**: Ensure Windows Defender Firewall or UFW isn't blocking `rnr-node`.
- **Check Port Forwarding**: Use a tool like [CanYouSeeMe.org](https://canyouseeme.org/) to check if port 8001 is open on Computer A.

### "No Peers Found"?
- Ensure you copied the full Multiaddr correctly, including the `/p2p/Qm...` part.
- Try pinging Computer A's IP address from Computer B: `ping 203.0.113.5`.
