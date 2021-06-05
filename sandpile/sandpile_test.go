package main

import "testing"

func TestRunSerialParallel(t *testing.T) {
	var sandpileTest = []struct {
		size      int
		pile      int
		placement string
	}{
		{100, 10000, "random"},
		{100, 10000, "central"},
		{1000, 120000, "random"},
		{1000, 120000, "central"},
		{200, 23413, "random"},
		{200, 23413, "central"},
		{41, 100000, "random"},
		{41, 100000, "central"},
		{53, 1088000, "random"},
		{53, 1088000, "central"},
		{1234, 213419, "random"},
		{978, 1888888, "central"},
		{978, 1888888, "random"},
	}

	for _, tt := range sandpileTest {
		sBoard, pBoard := RunSerialParallel(tt.size, tt.pile, tt.placement)
		if !CheckBoard(sBoard, pBoard) {
			t.Errorf("serial and parallele board do not match")
		} else {
			t.Logf("OK, passed. size: %d, pile: %d, placement: %s", tt.size, tt.pile, tt.placement)
		}
	}
}
