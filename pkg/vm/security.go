package vm

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// SecurityConfig defines limits and protections for contract execution
// These limits prevent DoS attacks similar to Ethereum's Shanghai incident
type SecurityConfig struct {
	// Execution Limits
	MaxExecutionTime time.Duration // Max time per contract call
	MaxMemoryBytes   uint64        // Max memory per contract instance
	MaxCallDepth     uint32        // Max nested contract calls
	MaxStorageOps    uint32        // Max storage reads/writes per call

	// Rate Limiting (per block)
	MaxContractsPerBlock uint32 // Max contract calls in one block
	MaxGasPerBlock       uint64 // Max total gas per block
	MaxDeploysPerBlock   uint32 // Max contract deployments per block

	// Circuit Breaker
	EnableCircuitBreaker bool          // Enable emergency circuit breaker
	FailureThreshold     uint32        // Failures before circuit opens
	CircuitResetTime     time.Duration // Time before retry after circuit opens

	// Emergency Controls
	EmergencyPauseEnabled bool // Allow emergency pause
	AdminPauseEnabled     bool // Allow admin pause (testnet only)

	// Monitoring
	EnableMetrics  bool          // Enable prometheus metrics
	EnableAlerts   bool          // Enable alerting
	AlertThreshold time.Duration // Alert if execution > threshold
}

// DefaultSecurityConfig returns conservative production settings
func DefaultSecurityConfig() *SecurityConfig {
	return &SecurityConfig{
		// Execution Limits
		MaxExecutionTime: 5 * time.Second,
		MaxMemoryBytes:   64 * 1024 * 1024, // 64 MB
		MaxCallDepth:     128,
		MaxStorageOps:    1000,

		// Rate Limiting
		MaxContractsPerBlock: 1000,
		MaxGasPerBlock:       10_000_000, // 10M gas per block
		MaxDeploysPerBlock:   10,

		// Circuit Breaker
		EnableCircuitBreaker: true,
		FailureThreshold:     10,
		CircuitResetTime:     5 * time.Minute,

		// Emergency Controls
		EmergencyPauseEnabled: true,
		AdminPauseEnabled:     false, // Disabled for mainnet

		// Monitoring
		EnableMetrics:  true,
		EnableAlerts:   true,
		AlertThreshold: 3 * time.Second,
	}
}

// TestnetSecurityConfig returns relaxed settings for testing
func TestnetSecurityConfig() *SecurityConfig {
	config := DefaultSecurityConfig()
	config.MaxExecutionTime = 10 * time.Second
	config.MaxContractsPerBlock = 5000
	config.AdminPauseEnabled = true // Allow admin pause on testnet
	return config
}

// DevelopmentSecurityConfig returns minimal limits for development
func DevelopmentSecurityConfig() *SecurityConfig {
	config := DefaultSecurityConfig()
	config.MaxExecutionTime = 30 * time.Second
	config.MaxContractsPerBlock = 10000
	config.EnableCircuitBreaker = false
	config.AdminPauseEnabled = true
	return config
}

// ExecutionLimiter tracks and enforces execution limits
type ExecutionLimiter struct {
	config *SecurityConfig

	// Per-execution tracking
	startTime  time.Time
	memoryUsed uint64
	callDepth  uint32
	storageOps uint32

	mu sync.RWMutex
}

// NewExecutionLimiter creates a new execution limiter
func NewExecutionLimiter(config *SecurityConfig) *ExecutionLimiter {
	return &ExecutionLimiter{
		config:    config,
		startTime: time.Now(),
	}
}

// CheckExecutionTime verifies execution hasn't exceeded time limit
func (el *ExecutionLimiter) CheckExecutionTime() error {
	if time.Since(el.startTime) > el.config.MaxExecutionTime {
		return &SecurityViolation{
			Type:    "ExecutionTimeout",
			Message: fmt.Sprintf("Execution exceeded %v", el.config.MaxExecutionTime),
			Limit:   uint64(el.config.MaxExecutionTime),
			Actual:  uint64(time.Since(el.startTime)),
		}
	}
	return nil
}

// CheckMemory verifies memory usage is within limits
func (el *ExecutionLimiter) CheckMemory(additional uint64) error {
	el.mu.Lock()
	defer el.mu.Unlock()

	newTotal := el.memoryUsed + additional
	if newTotal > el.config.MaxMemoryBytes {
		return &SecurityViolation{
			Type:    "MemoryLimit",
			Message: "Memory limit exceeded",
			Limit:   el.config.MaxMemoryBytes,
			Actual:  newTotal,
		}
	}
	el.memoryUsed = newTotal
	return nil
}

// IncrementCallDepth increments and checks call depth
func (el *ExecutionLimiter) IncrementCallDepth() error {
	el.mu.Lock()
	defer el.mu.Unlock()

	el.callDepth++
	if el.callDepth > el.config.MaxCallDepth {
		return &SecurityViolation{
			Type:    "CallDepthExceeded",
			Message: "Maximum call depth exceeded",
			Limit:   uint64(el.config.MaxCallDepth),
			Actual:  uint64(el.callDepth),
		}
	}
	return nil
}

// DecrementCallDepth decrements call depth
func (el *ExecutionLimiter) DecrementCallDepth() {
	el.mu.Lock()
	defer el.mu.Unlock()
	if el.callDepth > 0 {
		el.callDepth--
	}
}

// IncrementStorageOps increments and checks storage operations
func (el *ExecutionLimiter) IncrementStorageOps() error {
	el.mu.Lock()
	defer el.mu.Unlock()

	el.storageOps++
	if el.storageOps > el.config.MaxStorageOps {
		return &SecurityViolation{
			Type:    "StorageOpsExceeded",
			Message: "Maximum storage operations exceeded",
			Limit:   uint64(el.config.MaxStorageOps),
			Actual:  uint64(el.storageOps),
		}
	}
	return nil
}

// BlockLimiter tracks and enforces per-block limits
type BlockLimiter struct {
	config *SecurityConfig

	// Per-block counters (reset on new block)
	contractCalls   uint32
	contractDeploys uint32
	totalGasUsed    uint64

	mu sync.RWMutex
}

// NewBlockLimiter creates a new block limiter
func NewBlockLimiter(config *SecurityConfig) *BlockLimiter {
	return &BlockLimiter{
		config: config,
	}
}

// CheckContractCall verifies contract call is allowed in this block
func (bl *BlockLimiter) CheckContractCall() error {
	bl.mu.Lock()
	defer bl.mu.Unlock()

	if bl.contractCalls >= bl.config.MaxContractsPerBlock {
		return &SecurityViolation{
			Type:    "MaxContractsPerBlock",
			Message: "Maximum contract calls per block reached",
			Limit:   uint64(bl.config.MaxContractsPerBlock),
			Actual:  uint64(bl.contractCalls),
		}
	}
	bl.contractCalls++
	return nil
}

// CheckDeploy verifies contract deployment is allowed in this block
func (bl *BlockLimiter) CheckDeploy() error {
	bl.mu.Lock()
	defer bl.mu.Unlock()

	if bl.contractDeploys >= bl.config.MaxDeploysPerBlock {
		return &SecurityViolation{
			Type:    "MaxDeploysPerBlock",
			Message: "Maximum contract deployments per block reached",
			Limit:   uint64(bl.config.MaxDeploysPerBlock),
			Actual:  uint64(bl.contractDeploys),
		}
	}
	bl.contractDeploys++
	return nil
}

// AddGas adds gas usage and checks block gas limit
func (bl *BlockLimiter) AddGas(gas uint64) error {
	bl.mu.Lock()
	defer bl.mu.Unlock()

	newTotal := bl.totalGasUsed + gas
	if newTotal > bl.config.MaxGasPerBlock {
		return &SecurityViolation{
			Type:    "MaxGasPerBlock",
			Message: "Maximum gas per block exceeded",
			Limit:   bl.config.MaxGasPerBlock,
			Actual:  newTotal,
		}
	}
	bl.totalGasUsed = newTotal
	return nil
}

// Reset resets counters for a new block
func (bl *BlockLimiter) Reset() {
	bl.mu.Lock()
	defer bl.mu.Unlock()

	bl.contractCalls = 0
	bl.contractDeploys = 0
	bl.totalGasUsed = 0
}

// CircuitBreaker implements circuit breaker pattern for emergency stops
type CircuitBreaker struct {
	config *SecurityConfig

	// Circuit state
	state        CircuitState
	failures     uint32
	lastFailTime time.Time

	mu sync.RWMutex
}

type CircuitState int

const (
	CircuitClosed   CircuitState = iota // Normal operation
	CircuitOpen                         // Blocking all operations
	CircuitHalfOpen                     // Testing if system recovered
)

// NewCircuitBreaker creates a new circuit breaker
func NewCircuitBreaker(config *SecurityConfig) *CircuitBreaker {
	return &CircuitBreaker{
		config: config,
		state:  CircuitClosed,
	}
}

// Allow checks if operation is allowed
func (cb *CircuitBreaker) Allow() error {
	if !cb.config.EnableCircuitBreaker {
		return nil
	}

	cb.mu.RLock()
	defer cb.mu.RUnlock()

	switch cb.state {
	case CircuitOpen:
		// Check if enough time has passed to try again
		if time.Since(cb.lastFailTime) > cb.config.CircuitResetTime {
			cb.mu.RUnlock()
			cb.mu.Lock()
			cb.state = CircuitHalfOpen
			cb.mu.Unlock()
			cb.mu.RLock()
			return nil
		}
		return &SecurityViolation{
			Type:    "CircuitBreakerOpen",
			Message: "Circuit breaker is OPEN - contract execution disabled",
		}
	case CircuitHalfOpen:
		return nil // Allow one request to test
	case CircuitClosed:
		return nil // Normal operation
	}
	return nil
}

// RecordSuccess records a successful execution
func (cb *CircuitBreaker) RecordSuccess() {
	if !cb.config.EnableCircuitBreaker {
		return
	}

	cb.mu.Lock()
	defer cb.mu.Unlock()

	if cb.state == CircuitHalfOpen {
		cb.state = CircuitClosed
		cb.failures = 0
	}
}

// RecordFailure records a failed execution
func (cb *CircuitBreaker) RecordFailure() {
	if !cb.config.EnableCircuitBreaker {
		return
	}

	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.failures++
	cb.lastFailTime = time.Now()

	if cb.failures >= cb.config.FailureThreshold {
		cb.state = CircuitOpen
		// TODO: Alert administrators
		fmt.Printf("⚠️ CIRCUIT BREAKER OPENED - %d failures detected\n", cb.failures)
	}
}

// ManualOpen manually opens the circuit (emergency use)
func (cb *CircuitBreaker) ManualOpen() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.state = CircuitOpen
	cb.lastFailTime = time.Now()
}

// ManualClose manually closes the circuit
func (cb *CircuitBreaker) ManualClose() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.state = CircuitClosed
	cb.failures = 0
}

// SecurityViolation represents a security limit violation
type SecurityViolation struct {
	Type    string
	Message string
	Limit   uint64
	Actual  uint64
}

func (sv *SecurityViolation) Error() string {
	if sv.Limit > 0 {
		return fmt.Sprintf("Security Violation [%s]: %s (limit: %d, actual: %d)",
			sv.Type, sv.Message, sv.Limit, sv.Actual)
	}
	return fmt.Sprintf("Security Violation [%s]: %s", sv.Type, sv.Message)
}

// SecurityMetrics tracks security-related metrics
type SecurityMetrics struct {
	TotalViolations     uint64
	TimeoutViolations   uint64
	MemoryViolations    uint64
	CallDepthViolations uint64
	StorageViolations   uint64
	CircuitBreakerTrips uint64

	mu sync.RWMutex
}

// GlobalSecurityMetrics is the global security metrics instance
var GlobalSecurityMetrics = &SecurityMetrics{}

// RecordViolation records a security violation
func (sm *SecurityMetrics) RecordViolation(violationType string) {
	atomic.AddUint64(&sm.TotalViolations, 1)

	switch violationType {
	case "ExecutionTimeout":
		atomic.AddUint64(&sm.TimeoutViolations, 1)
	case "MemoryLimit":
		atomic.AddUint64(&sm.MemoryViolations, 1)
	case "CallDepthExceeded":
		atomic.AddUint64(&sm.CallDepthViolations, 1)
	case "StorageOpsExceeded":
		atomic.AddUint64(&sm.StorageViolations, 1)
	case "CircuitBreakerOpen":
		atomic.AddUint64(&sm.CircuitBreakerTrips, 1)
	}
}

// GetMetrics returns current metrics snapshot
func (sm *SecurityMetrics) GetMetrics() map[string]uint64 {
	return map[string]uint64{
		"total_violations":      atomic.LoadUint64(&sm.TotalViolations),
		"timeout_violations":    atomic.LoadUint64(&sm.TimeoutViolations),
		"memory_violations":     atomic.LoadUint64(&sm.MemoryViolations),
		"call_depth_violations": atomic.LoadUint64(&sm.CallDepthViolations),
		"storage_violations":    atomic.LoadUint64(&sm.StorageViolations),
		"circuit_breaker_trips": atomic.LoadUint64(&sm.CircuitBreakerTrips),
	}
}
