package main

import (
	"testing"
	"time"
)

func TestNextScheduledTransmissionFrom(t *testing.T) {
	// set the offset for testing
	offset = "2"

	getNext := func(year int, month time.Month, day, hour, min, sec int) time.Time {
		now := time.Date(year, month, day, hour, min, sec, 0, time.UTC)
		return nextScheduledTransmissionFrom(now)
	}

	// Test: 12:01 should schedule for 12:02
	result := getNext(2023, 1, 1, 12, 1, 0)
	expected := time.Date(2023, 1, 1, 12, 2, 0, 0, time.UTC)
	if !result.Equal(expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}

	// Test: 12:09 should schedule for 12:12
	result = getNext(2023, 1, 1, 12, 9, 0)
	expected = time.Date(2023, 1, 1, 12, 12, 0, 0, time.UTC)
	if !result.Equal(expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}

	// Test: 12:12 should schedule for 12:22
	result = getNext(2023, 1, 1, 12, 12, 0)
	expected = time.Date(2023, 1, 1, 12, 22, 0, 0, time.UTC)
	if !result.Equal(expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}

	// Test: 12:59 should schedule for 13:02 (hour rollover)
	result = getNext(2023, 1, 1, 12, 59, 0)
	expected = time.Date(2023, 1, 1, 13, 2, 0, 0, time.UTC)
	if !result.Equal(expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}

	// Test: 23:59 should schedule for 00:02 next day (day rollover)
	result = getNext(2023, 1, 1, 23, 59, 0)
	expected = time.Date(2023, 1, 2, 0, 2, 0, 0, time.UTC)
	if !result.Equal(expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}
