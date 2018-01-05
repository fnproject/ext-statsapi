package stats

import (
	"strconv"
	"testing"
)

// This file contains test utilities only (no tests)

func assertNoError(t *testing.T, assertionText string, err error) {
	if err != nil {
		t.Fatal(assertionText + " FAILED due to error: " + err.Error())
	}
}

func assertIntsEqual(t *testing.T, assertionText string, expected int, actual int) {
	if actual != expected {
		t.Fatal(assertionText + " FAILED: expected " + strconv.Itoa(expected) + ", actual " + strconv.Itoa(actual))
	}
}

func assertStringsEqual(t *testing.T, assertionText string, expected string, actual string) {
	if actual != expected {
		t.Fatal(assertionText + " FAILED: expected " + expected + ", actual " + actual)
	}
}

func assertNotNil(t *testing.T, assertionText string, actual interface{}) {
	if actual == nil {
		t.Fatal(assertionText + " FAILED: expected a non-nil value")
	}
}
