package rt

func Ø()                                       {}
func I[A any](a A) A                           { return a }
func II[A, B any](a A, b B) (A, B)             { return a, b }
func III[A, B, C any](a A, b B, c C) (A, B, C) { return a, b, c }
