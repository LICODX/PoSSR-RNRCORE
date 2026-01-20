// Simple Counter Contract in Rust
// Compile: rustc --target wasm32-unknown-unknown -O --crate-type=cdylib counter.rs

#![no_std]
#![no_main]

use core::panic::PanicInfo;

#[panic_handler]
fn panic(_info: &PanicInfo) -> ! {
    loop {}
}

// Global counter storage
static mut COUNTER: u64 = 0;
static mut OWNER: [u8; 32] = [0u8; 32];

// Host function imports from RNR blockchain
extern "C" {
    fn get_caller() -> i32;  // Returns pointer to caller address
    fn emit_event(event_ptr: i32, event_len: i32);
}

// Initialize contract with owner
#[no_mangle]
pub extern "C" fn init(owner_ptr: i32) {
    unsafe {
        // In production, copy owner address from WASM memory
        // For now, we'll just mark it as initialized
        COUNTER = 0;
    }
}

// Increment counter by 1
#[no_mangle]
pub extern "C" fn increment() -> u64 {
    unsafe {
        COUNTER += 1;
        COUNTER
    }
}

// Increment counter by custom amount
#[no_mangle]
pub extern "C" fn add(amount: u64) -> u64 {
    unsafe {
        COUNTER += amount;
        COUNTER
    }
}

// Decrement counter by 1
#[no_mangle]
pub extern "C" fn decrement() -> u64 {
    unsafe {
        if COUNTER > 0 {
            COUNTER -= 1;
        }
        COUNTER
    }
}

// Get current counter value
#[no_mangle]
pub extern "C" fn get() -> u64 {
    unsafe { COUNTER }
}

// Reset counter (only owner can call)
#[no_mangle]
pub extern "C" fn reset() -> i32 {
    unsafe {
        // In production, verify caller is owner
        COUNTER = 0;
        1 // Success
    }
}
