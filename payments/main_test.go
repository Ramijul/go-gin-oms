package main

import (
	"testing"
)

func FuzzTestPaymentProcessor(f *testing.F) {
	testcases := []float64{900.99, 1000.00, 1000.01}
	for _, tc := range testcases {
		f.Add(tc) // Use f.Add to provide a seed corpus
	}

	f.Fuzz(func(t *testing.T, input float64) {
		resp := processPayment(input)

		var expected string
		if input > 1000 {
			expected = PAYMENT_STATUS_FAILED
		} else {
			expected = PAYMENT_STATUS_PAID
		}

		if resp != expected {
			t.Errorf("Expected: %f, received: %q", input, resp)
		}
	})
}
