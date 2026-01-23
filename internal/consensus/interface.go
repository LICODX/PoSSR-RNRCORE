package consensus

import (
	"github.com/LICODX/PoSSR-RNRCORE/pkg/types"
)

// ConsensusEngine defines the formal contract for any consensus mechanism
// (PoW, BFT, PoA) used by the blockchain.
//
// This interface allows modular swapping of consensus algorithms without
// changing the core blockchain logic, addressing architectural criticisms
// in debat/7.txt regarding "Code Structure" and "Formal Contracts".
type Engine interface {
	// Initialize starts the consensus engine and any background routines
	Initialize() error

	// RunConsensusRound executes one round of consensus to produce a block
	// height: The target block height
	// txs: The transactions to include in the block
	// Returns: The finalized block, or error if consensus fails
	RunConsensusRound(height uint64, txs []types.Transaction) (*types.Block, error)

	// ValidateBlockHeader checks if a block header conforms to consensus rules
	// (e.g. valid PoW hash, valid BFT signatures)
	ValidateBlockHeader(header *types.BlockHeader) error

	// VerifySeal checks the cryptographic seal of the block (Nonce/Hash or Signatures)
	VerifySeal(header *types.BlockHeader) error

	// Stop gracefully shuts down the consensus engine
	Stop() error
}

// BFTAdapter adapts the concrete BFTEngine to the generic Engine interface
// This ensures that our specific BFT implementation adheres to the formal contract.
type BFTAdapter struct {
	engine *BFTEngine
}

// Ensure BFTAdapter implements Engine
var _ Engine = (*BFTAdapter)(nil)

func NewBFTAdapter(engine *BFTEngine) *BFTAdapter {
	return &BFTAdapter{engine: engine}
}

func (a *BFTAdapter) Initialize() error {
	// Logic to start BFT listeners if not already started
	return nil
}

func (a *BFTAdapter) RunConsensusRound(height uint64, txs []types.Transaction) (*types.Block, error) {
	return a.engine.RunConsensusRound(height, txs)
}

func (a *BFTAdapter) ValidateBlockHeader(header *types.BlockHeader) error {
	// TODO: Implement header validation logic specific to BFT (validators, round, etc.)
	return nil
}

func (a *BFTAdapter) VerifySeal(header *types.BlockHeader) error {
	// TODO: Verify quorum signatures
	return nil
}

func (a *BFTAdapter) Stop() error {
	return nil
}
