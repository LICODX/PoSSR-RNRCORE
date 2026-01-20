@echo off
echo ============================================
echo RNR-CORE PROJECT CLEANUP
echo ============================================
echo.
echo This script will remove temporary and build files
echo.

REM Remove log files
echo [1/5] Removing log files...
del /Q *.log 2>nul
echo   - Removed node logs and automation logs

REM Remove test executables (keep main rnr-node.exe)
echo [2/5] Removing test build files...
del /Q test-*.exe 2>nul
echo   - Removed test executables

REM Remove data directory
echo [3/5] Removing blockchain data...
if exist data (
    rmdir /S /Q data
    echo   - Removed data directory
)

REM Remove temporary text files
echo [4/5] Removing temporary files...
del /Q stress_report.txt 2>nul
del /Q tt.txt 2>nul
del /Q "uji shorting algorithm.txt" 2>nul
echo   - Removed temporary text files

REM Clean up build cache (optional)
echo [5/5] Cleaning Go build cache...
go clean -cache 2>nul
echo   - Cleaned Go build cache

echo.
echo ============================================
echo CLEANUP COMPLETE!
echo ============================================
echo.
echo Remaining important files:
echo   - Source code (cmd/, internal/, pkg/)
echo   - Documentation (README.md, WHITEPAPER.md, docs/)
echo   - Configuration (go.mod, .gitignore)
echo   - Main executable (rnr-node.exe)
echo   - Test scripts (RUN_*.bat)
echo.
pause
