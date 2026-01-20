// RNRScan Explorer Logic

document.addEventListener('DOMContentLoaded', () => {
    fetchStats();
    fetchLatestBlocks();
    setInterval(fetchStats, 5000); // Live updates
});

async function fetchStats() {
    try {
        const response = await fetch('/api/stats');
        const data = await response.json();

        document.getElementById('latestBlock').innerText = `#${data.height}`;
        document.getElementById('txCount').innerText = `${data.mempoolSize} Pending`; // Using mempool as proxy for activity

    } catch (error) {
        console.error('Error fetching stats:', error);
    }
}

async function fetchLatestBlocks() {
    try {
        const response = await fetch('/api/blocks?limit=5');
        const blocks = await response.json();

        const container = document.getElementById('blocksList');
        container.innerHTML = '';

        blocks.forEach(block => {
            const timeAgo = getTimeAgo(block.timestamp);
            const html = `
                <div class="list-item fade-in">
                    <div class="item-icon">B</div>
                    <div class="item-content">
                        <div>
                            <a href="/block/${block.hash}" class="item-primary">${block.height}</a>
                            <div class="item-secondary">${timeAgo}</div>
                        </div>
                        <div class="text-right">
                            <div class="item-primary">Validator: ${block.miner.substring(0, 10)}...</div>
                            <div class="item-secondary">${block.txCount} txns</div>
                        </div>
                    </div>
                </div>
            `;
            container.innerHTML += html;
        });

        // Mock Transactions for Demo if API returns empty array (since network might be empty)
        populateMockTxs();

    } catch (error) {
        console.error('Error fetching blocks:', error);
        document.getElementById('blocksList').innerHTML = '<div class="text-error p-3">Failed to load blocks</div>';
    }
}

function populateMockTxs() {
    const container = document.getElementById('txList');
    container.innerHTML = '';

    // Generating some visual mock data to demonstrate UI style (since mainnet might be idle)
    // In production, fetch from /api/transactions
    const mocks = [
        { hash: '0x3f...e1', from: '0xAlice...', to: '0xBob...', amount: '100 RNR' },
        { hash: '0xa9...b2', from: '0xValidator...', to: '0xPool...', amount: '5000 RNR' },
        { hash: '0x7c...d9', from: '0xUser1...', to: '0xUser2...', amount: '12.5 RNR' },
    ];

    mocks.forEach(tx => {
        const html = `
            <div class="list-item fade-in">
                <div class="item-icon" style="background: rgba(0,201,167,0.1); color: var(--success);"><i class="fa-solid fa-file-contract"></i></div>
                <div class="item-content">
                    <div>
                        <a href="#" class="item-primary">${tx.hash}</a>
                        <div class="item-secondary">From ${tx.from}</div>
                    </div>
                    <div class="text-right">
                        <div class="badge" style="color:var(--success)">${tx.amount}</div>
                        <div class="item-secondary">To ${tx.to}</div>
                    </div>
                </div>
            </div>
        `;
        container.innerHTML += html;
    });
}

function getTimeAgo(timestamp) {
    const seconds = Math.floor((Date.now() / 1000) - timestamp);
    if (seconds < 60) return `${seconds} secs ago`;
    return `${Math.floor(seconds / 60)} mins ago`;
}

function performSearch() {
    const query = document.getElementById('searchInput').value;
    if (query) {
        window.location.href = `/api/search?q=${query}`; // Direct API for now, improved later
    }
}
