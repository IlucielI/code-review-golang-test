package main

import (
	"fmt"
	"math"
)

// Vulnerable: Integer overflow in price calculation
func calculateTotalPrice(itemPrice int32, quantity int32) int32 {
	// No overflow check - can wrap to negative
	total := itemPrice * quantity
	return total
}

// Vulnerable: Integer overflow in allocation size
func allocateBuffer(size int) []byte {
	// If size is near MaxInt, size+1024 can overflow to small negative
	bufferSize := size + 1024
	return make([]byte, bufferSize)
}

// Vulnerable: Unsigned underflow
func calculateDiscount(price uint32, discount uint32) uint32 {
	// If discount > price, underflows to large positive number
	return price - discount
}

// Example exploitation
func exploitExample() {
	// Overflow to negative price
	price := int32(math.MaxInt32 / 2)
	quantity := int32(3)
	total := calculateTotalPrice(price, quantity) // Overflows
	fmt.Printf("Total (overflowed): %d\n", total)
}
