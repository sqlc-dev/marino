package main

import (
	"testing"

	"reflect"
)

func TestParseLine(t *testing.T) {
	add := parseLine("	add               \"ADD\"")
	if !reflect.DeepEqual(add, "ADD") {
		t.Fatalf("got %v, want %v", "ADD", add)
	}

	tso := parseLine("	tidbCurrentTSO    \"TIDB_CURRENT_TSO\"")
	if !reflect.DeepEqual(tso, "TIDB_CURRENT_TSO") {
		t.Fatalf("got %v, want %v", "TIDB_CURRENT_TSO", tso)
	}
}
