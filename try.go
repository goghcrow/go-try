package try

func Try0[E any](E)                                   {}
func Try[A, E any](A, E) (_ A)                        { return }
func Try2[A, B, E any](A, B, E) (_ A, _ B)            { return }
func Try3[A, B, C, E any](A, B, C, E) (_ A, _ B, _ C) { return }
