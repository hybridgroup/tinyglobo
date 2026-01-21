package main

import "testing"

func TestTransmitFrequency(t *testing.T) {
	tests := []struct {
		band     int
		expected uint64
	}{
		{1, 14097020},
		{2, 14097060},
		{3, 14097140},
		{4, 14097180},
	}

	for _, tt := range tests {
		got := transmitFrequency(tt.band)
		if got != tt.expected {
			t.Errorf("transmitFrequency(%d) = %d; want %d", tt.band, got, tt.expected)
		}
	}
}
