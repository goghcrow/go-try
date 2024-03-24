package rt

func Ø()                                                           {}
func I[A any](a A) A                                               { return a }
func II[A, B any](a A, b B) (A, B)                                 { return a, b }
func III[A, B, C any](a A, b B, c C) (A, B, C)                     { return a, b, c }
func IV[A, B, C, D any](a A, b B, c C, d D) (A, B, C, D)           { return a, b, c, d }
func V[A, B, C, D, E any](a A, b B, c C, d D, e E) (A, B, C, D, E) { return a, b, c, d, e }
func VI[A, B, C, D, E, F any](a A, b B, c C, d D, e E, f F) (A, B, C, D, E, F) {
	return a, b, c, d, e, f
}
func VII[A, B, C, D, E, F, G any](a A, b B, c C, d D, e E, f F, g G) (A, B, C, D, E, F, G) {
	return a, b, c, d, e, f, g
}
func VIII[A, B, C, D, E, F, G, H any](a A, b B, c C, d D, e E, f F, g G, h H) (A, B, C, D, E, F, G, H) {
	return a, b, c, d, e, f, g, h
}
func IX[A, B, C, D, E, F, G, H, I any](a A, b B, c C, d D, e E, f F, g G, h H, i I) (A, B, C, D, E, F, G, H, I) {
	return a, b, c, d, e, f, g, h, i
}
func X[A, B, C, D, E, F, G, H, I, J any](a A, b B, c C, d D, e E, f F, g G, h H, i I, j J) (A, B, C, D, E, F, G, H, I, J) {
	return a, b, c, d, e, f, g, h, i, j
}
