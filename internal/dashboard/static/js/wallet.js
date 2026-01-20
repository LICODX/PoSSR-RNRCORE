// RNR Wallet Logic

let currentAddress = "";
let currentBalance = 0;

document.addEventListener('DOMContentLoaded', () => {
    fetchWalletData();
    setInterval(fetchWalletData, 3000);
});

async function fetchWalletData() {
    try {
        const response = await fetch('/api/wallet');
        if (response.status === 404) {
            // No wallet found on node
            document.getElementById('balanceDisplay').innerText = "No Wallet";
            return;
        }

        const data = await response.json();
        currentAddress = data.address;
        currentBalance = data.balance;

        // Update UI
        document.getElementById('shortAddress').innerText = `${currentAddress.substring(0, 6)}...${currentAddress.substring(currentAddress.length - 4)}`;
        document.getElementById('balanceDisplay').innerText = `${currentBalance} RNR`;

        // Mock USD Price ($12.45)
        const usdValue = (currentBalance * 12.45).toFixed(2);
        document.getElementById('usdDisplay').innerText = `$${usdValue} USD`;

    } catch (error) {
        console.error("Wallet fetch error:", error);
    }
}

function copyAddress() {
    navigator.clipboard.writeText(currentAddress).then(() => {
        alert("Address copied!");
    });
}

function openSendModal() {
    const to = prompt("Enter Recipient Address:");
    if (!to) return;

    const amount = prompt("Enter Amount (RNR):");
    if (!amount) return;

    // Send TX via API
    fetch('/api/wallet/send', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
            to: to,
            amount: parseFloat(amount),
            fee: 0.01
        })
    })
        .then(res => res.json())
        .then(data => {
            if (data.success) {
                alert("Transaction Sent! Hash: " + data.txHash);
            } else {
                alert("Error: " + data.error);
            }
        });
}
