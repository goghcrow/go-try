package test

func func1[A, R any](a A) (_ R, _ error) { return }

func return1[A any]() (_ A)                          { return }
func ret0Err() (_ error)                             { return }
func ret1Err[A any]() (_ A, _ error)                 { return }
func ret2Err[A, B any]() (_ A, _ B, _ error)         { return }
func ret3Err[A, B, C any]() (_ A, _ B, _ C, _ error) { return }

func consume1[A any](A)             {}
func consume2[A, B any](A, B)       {}
func consume3[A, B, C any](A, B, C) {}

func id[A any](a A) A { return a }
