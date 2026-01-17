@echo off
color 0A
title PoSSR RNRCORE - MAINNET NODE
echo ====================================================
echo   PoSSR RNRCORE MAINNET LAUNCHER
echo   Bootnode: 192.168.36.1 (THIS MACHINE)
echo   Genesis Address: rnr1pq03gqs8zg0...
echo ====================================================
echo.
echo [1] Initializing Node...
rmdir /S /Q data logs 2>nul
echo [2] Starting Node...
echo.
rnr-node.exe
pause
