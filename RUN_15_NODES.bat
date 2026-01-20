@echo off
echo 🚀 Starting PoSSR 15-Node Testnet (Genesis + 14 Guests)...
echo.

:: 1. Force Kill Zombie Processes
echo 💀 Killing old rnr-node.exe processes...
taskkill /F /IM rnr-node.exe >nul 2>&1

:: 2. Clean previous data
echo 🧹 Cleaning data/ directory...
timeout /t 2 >nul
rmdir /s /q data 2>nul
mkdir data

:: 3. Start Node 1 (Bootnode - GENESIS AUTHORITY)
echo 🟢 Node 1 (Genesis) - Port 3000 / Dash 8080
start "Node 1 (Genesis)" cmd /c "rnr-node.exe -port 3000 -datadir ./data/node1 -genesis > node1.log 2>&1"
timeout /t 3

:: 4-17. Start Nodes 2-15 (Guests)
echo 🟢 Node 2 (Guest) - Port 3001 / Dash 8081
start "Node 2" cmd /c "rnr-node.exe -port 3001 -datadir ./data/node2 -peers /ip4/127.0.0.1/tcp/3000 > node2.log 2>&1"

echo 🟢 Node 3 (Guest) - Port 3002 / Dash 8082
start "Node 3" cmd /c "rnr-node.exe -port 3002 -datadir ./data/node3 -peers /ip4/127.0.0.1/tcp/3000 > node3.log 2>&1"

echo 🟢 Node 4 (Guest) - Port 3003 / Dash 8083
start "Node 4" cmd /c "rnr-node.exe -port 3003 -datadir ./data/node4 -peers /ip4/127.0.0.1/tcp/3000 > node4.log 2>&1"

echo 🟢 Node 5 (Guest) - Port 3004 / Dash 8084
start "Node 5" cmd /c "rnr-node.exe -port 3004 -datadir ./data/node5 -peers /ip4/127.0.0.1/tcp/3000 > node5.log 2>&1"

echo 🟢 Node 6 (Guest) - Port 3005 / Dash 8085
start "Node 6" cmd /c "rnr-node.exe -port 3005 -datadir ./data/node6 -peers /ip4/127.0.0.1/tcp/3000 > node6.log 2>&1"

echo 🟢 Node 7 (Guest) - Port 3006 / Dash 8086
start "Node 7" cmd /c "rnr-node.exe -port 3006 -datadir ./data/node7 -peers /ip4/127.0.0.1/tcp/3000 > node7.log 2>&1"

echo 🟢 Node 8 (Guest) - Port 3007 / Dash 8087
start "Node 8" cmd /c "rnr-node.exe -port 3007 -datadir ./data/node8 -peers /ip4/127.0.0.1/tcp/3000 > node8.log 2>&1"

echo 🟢 Node 9 (Guest) - Port 3008 / Dash 8088
start "Node 9" cmd /c "rnr-node.exe -port 3008 -datadir ./data/node9 -peers /ip4/127.0.0.1/tcp/3000 > node9.log 2>&1"

echo 🟢 Node 10 (Guest) - Port 3009 / Dash 8089
start "Node 10" cmd /c "rnr-node.exe -port 3009 -datadir ./data/node10 -peers /ip4/127.0.0.1/tcp/3000 > node10.log 2>&1"

echo 🟢 Node 11 (Guest) - Port 3010 / Dash 8090
start "Node 11" cmd /c "rnr-node.exe -port 3010 -datadir ./data/node11 -peers /ip4/127.0.0.1/tcp/3000 > node11.log 2>&1"

echo 🟢 Node 12 (Guest) - Port 3011 / Dash 8091
start "Node 12" cmd /c "rnr-node.exe -port 3011 -datadir ./data/node12 -peers /ip4/127.0.0.1/tcp/3000 > node12.log 2>&1"

echo 🟢 Node 13 (Guest) - Port 3012 / Dash 8092
start "Node 13" cmd /c "rnr-node.exe -port 3012 -datadir ./data/node13 -peers /ip4/127.0.0.1/tcp/3000 > node13.log 2>&1"

echo 🟢 Node 14 (Guest) - Port 3013 / Dash 8093
start "Node 14" cmd /c "rnr-node.exe -port 3013 -datadir ./data/node14 -peers /ip4/127.0.0.1/tcp/3000 > node14.log 2>&1"

echo 🟢 Node 15 (Guest) - Port 3014 / Dash 8094
start "Node 15" cmd /c "rnr-node.exe -port 3014 -datadir ./data/node15 -peers /ip4/127.0.0.1/tcp/3000 > node15.log 2>&1"

echo.
echo ✅ 15-Node Testnet Running!
echo.
echo 📜 Configuration:
echo    - Node 1: Genesis Authority (has RNR coins)
echo    - Nodes 2-15: Guests (will auto-register and get wallets)
echo.
echo 📁 LOG FILES: node1.log to node15.log
echo.
echo 🌐 Dashboards:
echo    Node 1:  http://localhost:8080
echo    Node 2:  http://localhost:8081
echo    Node 15: http://localhost:8094
echo.
echo ⏳ Wait 30 seconds for network to stabilize...
timeout /t 30
echo.
echo 📊 Check node1.log to see Genesis mining
echo 📊 Check node2.log to see registration status
echo.
pause
