@echo off
echo 🚀 Starting PoSSR 3-Node Testnet (Bootstrapped via Genesis)...
echo.

:: 1. Force Kill Zombie Processes
echo 💀 Killing old rnr-node.exe processes...
taskkill /F /IM rnr-node.exe >nul 2>&1

:: 2. Clean previous data
echo 🧹 Cleaning data/ directory...
timeout /t 2 >nul
rmdir /s /q data 2>nul
mkdir data

:: 3. Start Node 1 (Bootnode - GENESIS AUTHORITY) - LOG TO FILE
echo 🟢 Starting Node 1 (Genesis Authority) - Port 3000 / Dash 8080
start "Node 1 (Genesis Authority)" cmd /c "rnr-node.exe -port 3000 -datadir ./data/node1 -genesis > node1.log 2>&1"
timeout /t 5

:: 4. Start Node 2 (Guest) - LOG TO FILE
echo 🟢 Starting Node 2 (Guest) - Port 3001 / Dash 8081
start "Node 2 (Guest)" cmd /c "rnr-node.exe -port 3001 -datadir ./data/node2 -peers /ip4/127.0.0.1/tcp/3000 > node2.log 2>&1"
timeout /t 2

:: 5. Start Node 3 (Guest) - LOG TO FILE
echo 🟢 Starting Node 3 (Guest) - Port 3002 / Dash 8082
start "Node 3 (Guest)" cmd /c "rnr-node.exe -port 3002 -datadir ./data/node3 -peers /ip4/127.0.0.1/tcp/3000 > node3.log 2>&1"

echo.
echo ✅ Testnet Running!
echo.
echo 📁 LOG FILES CREATED:
echo    - node1.log (Genesis Node)
echo    - node2.log (Guest Node 2)
echo    - node3.log (Guest Node 3)
echo.
echo 👉 Tunggu 10 detik, lalu buka file node1.log, node2.log, node3.log dengan editor.
echo.
pause
