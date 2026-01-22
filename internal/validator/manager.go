package validator

import (
	"fmt"
	"sync"

	"github.com/LICODX/PoSSR-RNRCORE/internal/consensus/bft"
)

// Manager manages validator lifecycle (registration, activation, removal)
type Manager struct {
	mu sync.RWMutex

	// Active validator set (currently validating)
	activeSet *bft.ValidatorSet

	// Pending validators (will be activated next epoch)
	pendingSet map[[32]byte]*bft.Validator

	// Configuration
	minStake      uint64 // Minimum stake to become validator
	maxValidators int    // Maximum number of validators
	epochLength   uint64 // Blocks per epoch
}

// NewManager creates a new validator manager
func NewManager(minStake uint64, maxValidators int, epochLength uint64, initialValidators []*bft.Validator) *Manager {
	return &Manager{
		activeSet:     bft.NewValidatorSet(initialValidators),
		pendingSet:    make(map[[32]byte]*bft.Validator),
		minStake:      minStake,
		maxValidators: maxValidators,
		epochLength:   epochLength,
	}
}

// RegisterValidator registers a new validator (pending activation)
func (vm *Manager) RegisterValidator(validator *bft.Validator) error {
	vm.mu.Lock()
	defer vm.mu.Unlock()

	// Check minimum stake
	if validator.VotingPower < vm.minStake {
		return fmt.Errorf("insufficient stake: %d < %d", validator.VotingPower, vm.minStake)
	}

	// Check if already exists
	if vm.activeSet.HasAddress(validator.Address) {
		return fmt.Errorf("validator already active: %x", validator.Address[:4])
	}

	if _, exists := vm.pendingSet[validator.Address]; exists {
		return fmt.Errorf("validator already pending: %x", validator.Address[:4])
	}

	// Check max validators
	if vm.activeSet.Size()+len(vm.pendingSet) >= vm.maxValidators {
		return fmt.Errorf("max validators reached: %d", vm.maxValidators)
	}

	// Add to pending set
	vm.pendingSet[validator.Address] = validator

	fmt.Printf("[Validator] Registered pending validator %x with stake %d\n",
		validator.Address[:4], validator.VotingPower)

	return nil
}

// RotateEpoch rotates the validator set (activate pending, remove slashed)
func (vm *Manager) RotateEpoch(currentHeight uint64, slashedAddresses [][32]byte) {
	vm.mu.Lock()
	defer vm.mu.Unlock()

	if currentHeight%vm.epochLength != 0 {
		return
	}

	fmt.Printf("[Validator] Epoch rotation at height %d\n", currentHeight)

	// Remove slashed validators
	for _, addr := range slashedAddresses {
		if err := vm.activeSet.Remove(addr); err == nil {
			fmt.Printf("[Validator] Removed slashed validator %x\n", addr[:4])
		}
	}

	// Activate pending validators
	for addr, val := range vm.pendingSet {
		if err := vm.activeSet.Add(val); err == nil {
			fmt.Printf("[Validator] Activated validator %x\n", addr[:4])
			delete(vm.pendingSet, addr)
		}
	}
}

// GetActiveSet returns the current active validator set
func (vm *Manager) GetActiveSet() *bft.ValidatorSet {
	vm.mu.RLock()
	defer vm.mu.RUnlock()

	return vm.activeSet.Copy()
}

// GetPendingValidators returns pending validators
func (vm *Manager) GetPendingValidators() []*bft.Validator {
	vm.mu.RLock()
	defer vm.mu.RUnlock()

	pending := make([]*bft.Validator, 0, len(vm.pendingSet))
	for _, val := range vm.pendingSet {
		pending = append(pending, val)
	}
	return pending
}

// UpdateStake updates a validator's stake (voting power)
func (vm *Manager) UpdateStake(address [32]byte, newStake uint64) error {
	vm.mu.Lock()
	defer vm.mu.Unlock()

	// Check minimum stake
	if newStake < vm.minStake {
		return fmt.Errorf("stake below minimum: %d < %d", newStake, vm.minStake)
	}

	// Update active validator
	if err := vm.activeSet.UpdateVotingPower(address, newStake); err == nil {
		fmt.Printf("[Validator] Updated stake for %x: %d\n", address[:4], newStake)
		return nil
	}

	// Update pending validator
	if val, exists := vm.pendingSet[address]; exists {
		val.VotingPower = newStake
		fmt.Printf("[Validator] Updated pending stake for %x: %d\n", address[:4], newStake)
		return nil
	}

	return fmt.Errorf("validator not found: %x", address[:4])
}

// RemoveValidator removes a validator (e.g., voluntary exit)
func (vm *Manager) RemoveValidator(address [32]byte) error {
	vm.mu.Lock()
	defer vm.mu.Unlock()

	// Try removing from active set
	if err := vm.activeSet.Remove(address); err == nil {
		fmt.Printf("[Validator] Removed validator %x from active set\n", address[:4])
		return nil
	}

	// Try removing from pending set
	if _, exists := vm.pendingSet[address]; exists {
		delete(vm.pendingSet, address)
		fmt.Printf("[Validator] Removed validator %x from pending set\n", address[:4])
		return nil
	}

	return fmt.Errorf("validator not found: %x", address[:4])
}

// GetValidator returns a validator by address (active or pending)
func (vm *Manager) GetValidator(address [32]byte) (*bft.Validator, bool) {
	vm.mu.RLock()
	defer vm.mu.RUnlock()

	// Check active set
	if val := vm.activeSet.GetByAddress(address); val != nil {
		return val, true
	}

	// Check pending set
	if val, exists := vm.pendingSet[address]; exists {
		return val, true
	}

	return nil, false
}
