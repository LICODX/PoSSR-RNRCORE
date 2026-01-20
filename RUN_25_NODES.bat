@echo off
REM RNR Blockchain - 25 Node Adversarial Network Test
REM 18 Honest Nodes + 7 Malicious Nodes
REM Byzantine Fault Tolerance Test

echo ========================================
echo RNR 25-Node Adversarial Network Test
echo ========================================
echo.
echo Network Composition:
echo - 18 Honest Nodes (Node 1-18)
echo - 7 Malicious Nodes (Node 19-25)
echo - Byzantine Ratio: 28%% (below 33%% threshold)
echo.
echo Test Features:
echo - Every 3 blocks: 5 random nodes send 2 RNR
echo - Every 25 blocks: 3 nodes create tokens
echo - Every 55 blocks: 7 nodes deploy contracts
echo.
echo Starting nodes...
echo [INFO] All output redirected to log files (node1.log - node25.log)
echo.

REM Clean previous data
if exist data rmdir /s /q data
mkdir data

REM Honest Nodes (1-18)
echo Starting Honest Nodes (1-18)...

start "RNR Node 1" cmd /c "rnr-node.exe --port 8001 --rpc-port 9001 --dashboard-port 9101 --data-dir data\node1 > node1.log 2>&1"
timeout /t 2 /nobreak >nul

start "RNR Node 2" cmd /c "rnr-node.exe --port 8002 --rpc-port 9002 --dashboard-port 9102 --data-dir data\node2 --peer /ip4/127.0.0.1/tcp/8001 > node2.log 2>&1"
start "RNR Node 3" cmd /c "rnr-node.exe --port 8003 --rpc-port 9003 --dashboard-port 9103 --data-dir data\node3 --peer /ip4/127.0.0.1/tcp/8001 > node3.log 2>&1"
start "RNR Node 4" cmd /c "rnr-node.exe --port 8004 --rpc-port 9004 --dashboard-port 9104 --data-dir data\node4 --peer /ip4/127.0.0.1/tcp/8001 > node4.log 2>&1"
start "RNR Node 5" cmd /c "rnr-node.exe --port 8005 --rpc-port 9005 --dashboard-port 9105 --data-dir data\node5 --peer /ip4/127.0.0.1/tcp/8001 > node5.log 2>&1"
timeout /t 1 /nobreak >nul

start "RNR Node 6" cmd /c "rnr-node.exe --port 8006 --rpc-port 9006 --dashboard-port 9106 --data-dir data\node6 --peer /ip4/127.0.0.1/tcp/8001 > node6.log 2>&1"
start "RNR Node 7" cmd /c "rnr-node.exe --port 8007 --rpc-port 9007 --dashboard-port 9107 --data-dir data\node7 --peer /ip4/127.0.0.1/tcp/8001 > node7.log 2>&1"
start "RNR Node 8" cmd /c "rnr-node.exe --port 8008 --rpc-port 9008 --dashboard-port 9108 --data-dir data\node8 --peer /ip4/127.0.0.1/tcp/8001 > node8.log 2>&1"
start "RNR Node 9" cmd /c "rnr-node.exe --port 8009 --rpc-port 9009 --dashboard-port 9109 --data-dir data\node9 --peer /ip4/127.0.0.1/tcp/8001 > node9.log 2>&1"
start "RNR Node 10" cmd /c "rnr-node.exe --port 8010 --rpc-port 9010 --dashboard-port 9110 --data-dir data\node10 --peer /ip4/127.0.0.1/tcp/8001 > node10.log 2>&1"
timeout /t 1 /nobreak >nul

start "RNR Node 11" cmd /c "rnr-node.exe --port 8011 --rpc-port 9011 --dashboard-port 9111 --data-dir data\node11 --peer /ip4/127.0.0.1/tcp/8001 > node11.log 2>&1"
start "RNR Node 12" cmd /c "rnr-node.exe --port 8012 --rpc-port 9012 --dashboard-port 9112 --data-dir data\node12 --peer /ip4/127.0.0.1/tcp/8001 > node12.log 2>&1"
start "RNR Node 13" cmd /c "rnr-node.exe --port 8013 --rpc-port 9013 --dashboard-port 9113 --data-dir data\node13 --peer /ip4/127.0.0.1/tcp/8001 > node13.log 2>&1"
start "RNR Node 14" cmd /c "rnr-node.exe --port 8014 --rpc-port 9014 --dashboard-port 9114 --data-dir data\node14 --peer /ip4/127.0.0.1/tcp/8001 > node14.log 2>&1"
start "RNR Node 15" cmd /c "rnr-node.exe --port 8015 --rpc-port 9015 --dashboard-port 9115 --data-dir data\node15 --peer /ip4/127.0.0.1/tcp/8001 > node15.log 2>&1"
timeout /t 1 /nobreak >nul

start "RNR Node 16" cmd /c "rnr-node.exe --port 8016 --rpc-port 9016 --dashboard-port 9116 --data-dir data\node16 --peer /ip4/127.0.0.1/tcp/8001 > node16.log 2>&1"
start "RNR Node 17" cmd /c "rnr-node.exe --port 8017 --rpc-port 9017 --dashboard-port 9117 --data-dir data\node17 --peer /ip4/127.0.0.1/tcp/8001 > node17.log 2>&1"
start "RNR Node 18" cmd /c "rnr-node.exe --port 8018 --rpc-port 9018 --dashboard-port 9118 --data-dir data\node18 --peer /ip4/127.0.0.1/tcp/8001 > node18.log 2>&1"
timeout /t 2 /nobreak >nul

echo.
echo Malicious Nodes (19-25)...
echo [WARNING] These nodes will exhibit Byzantine behavior!
echo.

REM Malicious Nodes (19-25)
start "RNR Node 19 [MALICIOUS-DOUBLESPEND]" cmd /c "rnr-node.exe --port 8019 --rpc-port 9019 --dashboard-port 9119 --data-dir data\node19 --peer /ip4/127.0.0.1/tcp/8001 > node19.log 2>&1"
start "RNR Node 20 [MALICIOUS-INVALIDTX]" cmd /c "rnr-node.exe --port 8020 --rpc-port 9020 --dashboard-port 9120 --data-dir data\node20 --peer /ip4/127.0.0.1/tcp/8001 > node20.log 2>&1"
start "RNR Node 21 [MALICIOUS-BLOCKSPAM]" cmd /c "rnr-node.exe --port 8021 --rpc-port 9021 --dashboard-port 9121 --data-dir data\node21 --peer /ip4/127.0.0.1/tcp/8001 > node21.log 2>&1"
start "RNR Node 22 [MALICIOUS-TXSPAM]" cmd /c "rnr-node.exe --port 8022 --rpc-port 9022 --dashboard-port 9122 --data-dir data\node22 --peer /ip4/127.0.0.1/tcp/8001 > node22.log 2>&1"
start "RNR Node 23 [MALICIOUS-SELFISH]" cmd /c "rnr-node.exe --port 8023 --rpc-port 9023 --dashboard-port 9123 --data-dir data\node23 --peer /ip4/127.0.0.1/tcp/8001 > node23.log 2>&1"
start "RNR Node 24 [MALICIOUS-FORK]" cmd /c "rnr-node.exe --port 8024 --rpc-port 9024 --dashboard-port 9124 --data-dir data\node24 --peer /ip4/127.0.0.1/tcp/8001 > node24.log 2>&1"
start "RNR Node 25 [MALICIOUS-SILENT]" cmd /c "rnr-node.exe --port 8025 --rpc-port 9025 --dashboard-port 9125 --data-dir data\node25 --peer /ip4/127.0.0.1/tcp/8001 > node25.log 2>&1"

echo.
echo ========================================
echo All 25 nodes started!
echo ========================================
echo.
echo Dashboards available at:
echo   Node 1:  http://localhost:9101
echo   Node 10: http://localhost:9110
echo   Node 18: http://localhost:9118
echo.
echo Logs being written to:
echo   node1.log - node25.log
echo   tx_automation.log
echo.
echo Waiting 30 seconds for network formation...
timeout /t 30 /nobreak

echo.
echo Starting automated transaction system...
start "TX Automation" cmd /c "go run simulation\automated_tx\main.go > tx_automation.log 2>&1"

echo.
echo ========================================
echo Test is running!
echo ========================================
echo.
echo Monitoring:
echo - Watch dashboards for network health
echo - Check log files: node1.log - node25.log
echo - Test will run for ~200 blocks (~33 minutes)
echo.
echo Useful commands:
echo   Get-Content node1.log -Tail 20 -Wait  (live tail)
echo   Select-String "MINING" node*.log      (search logs)
echo.
echo Press any key to view monitoring dashboard...
pause
start http://localhost:9101

echo.
echo Test in progress. All output in log files.
echo Press Ctrl+C to stop this script.
echo.
