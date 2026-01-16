package txpool

import (
	"container/heap"
	"fmt"
	"sync"

	"github.com/LICODX/PoSSR-RNRCORE/pkg/types"
)

const MaxPoolSize = 10000

// Pool manages pending transactions
type Pool struct {
	transactions map[[32]byte]*types.Transaction
	queue        PriorityQueue
	mu           sync.RWMutex
}

// NewPool creates a new transaction pool
func NewPool() *Pool {
	pq := make(PriorityQueue, 0)
	heap.Init(&pq)

	return &Pool{
		transactions: make(map[[32]byte]*types.Transaction),
		queue:        pq,
	}
}

// Add adds a transaction to the pool
func (p *Pool) Add(tx types.Transaction) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Check pool size
	if len(p.transactions) >= MaxPoolSize {
		return fmt.Errorf("transaction pool full")
	}

	// Check for duplicate
	if _, exists := p.transactions[tx.ID]; exists {
		return fmt.Errorf("duplicate transaction")
	}

	// TODO: Validate transaction signature and nonce

	// Add to map
	p.transactions[tx.ID] = &tx

	// Add to priority queue (by nonce for now)
	heap.Push(&p.queue, &TxWithPriority{
		Tx:       &tx,
		Priority: int64(tx.Nonce),
	})

	return nil
}

// GetBest returns N highest priority transactions
func (p *Pool) GetBest(n int) []types.Transaction {
	p.mu.RLock()
	defer p.mu.RUnlock()

	var result []types.Transaction
	for i := 0; i < n && i < len(p.queue); i++ {
		result = append(result, *p.queue[i].Tx)
	}
	return result
}

// Remove removes a transaction from pool
func (p *Pool) Remove(txID [32]byte) {
	p.mu.Lock()
	defer p.mu.Unlock()

	delete(p.transactions, txID)
	// TODO: Remove from priority queue
}

// Size returns number of pending transactions
func (p *Pool) Size() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return len(p.transactions)
}

// TxWithPriority wraps transaction with priority
type TxWithPriority struct {
	Tx       *types.Transaction
	Priority int64 // Higher = more priority
	index    int
}

// PriorityQueue implements heap.Interface
type PriorityQueue []*TxWithPriority

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].Priority > pq[j].Priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*TxWithPriority)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.index = -1
	*pq = old[0 : n-1]
	return item
}
