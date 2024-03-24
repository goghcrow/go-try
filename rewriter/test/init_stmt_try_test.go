//go:build try

package test

import (
	"go/ast"

	. "github.com/goghcrow/go-try"
)

// IfStmt | SwitchStmt | TypeSwitchStmt | ForStmt

func init_expr_stmt() error {
	{
		if id(42); Try(ret1Err[bool]()) {
		}
	}
	{
		switch id(42); Try(ret1Err[bool]()) {
		}
	}
	{
		switch id(42); n := Try(ret1Err[ast.Node]()).(type) {
		default:
			_ = n
		}
	}
	{
		for id(42); Try(ret1Err[bool]()); {
		}
	}
	return nil
}

func init_send_stmt() error {
	var ch chan<- int
	{
		if ch <- 42; Try(ret1Err[bool]()) {
		}
	}
	{
		switch ch <- 42; Try(ret1Err[bool]()) {
		}
	}
	{
		switch ch <- 42; n := Try(ret1Err[ast.Node]()).(type) {
		default:
			_ = n
		}
	}
	{
		for ch <- 42; Try(ret1Err[bool]()); {
		}
	}
	return nil
}

func init_incdec_stmt(i int) error {
	{
		if i++; Try(ret1Err[bool]()) {
		}
	}
	{
		switch i++; Try(ret1Err[bool]()) {
		}
	}
	{
		switch i++; n := Try(ret1Err[ast.Node]()).(type) {
		default:
			_ = n
		}
	}
	{
		for i++; Try(ret1Err[bool]()); {
		}
	}
	return nil
}

func init_assign_stmt() error {
	{
		if i := 0; Try(ret1Err[bool]()) {
			_ = i
		}
	}
	{
		switch i := 0; Try(ret1Err[bool]()) {
		default:
			_ = i
		}
	}
	{
		switch i := 0; n := Try(ret1Err[ast.Node]()).(type) {
		default:
			_ = i
			_ = n
		}
	}
	{
		for i := 0; Try(ret1Err[bool]()); {
			_ = i
		}
	}
	return nil
}
