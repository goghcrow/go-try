package test

import "testing"

func lit[T any](x T) (T, error) {
	return x, nil
}

func TestEval(t *testing.T) {
	_ = testSwitchStmt()
	_ = testForStmt()
	testRangeStmt()
}
