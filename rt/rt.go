package rt

func Ø()                                                            {}
func Ƭ𝟭[A any](a A) A                                               { return a }
func Ƭ2[A, B any](a A, b B) (A, B)                                  { return a, b }
func Ƭ3[A, B, C any](a A, b B, c C) (A, B, C)                       { return a, b, c }
func Ƭ4[A, B, C, D any](a A, b B, c C, d D) (A, B, C, D)            { return a, b, c, d }
func Ƭ5[A, B, C, D, E any](a A, b B, c C, d D, e E) (A, B, C, D, E) { return a, b, c, d, e }
func Ƭ6[A, B, C, D, E, F any](a A, b B, c C, d D, e E, f F) (A, B, C, D, E, F) {
	return a, b, c, d, e, f
}
func Ƭ7[A, B, C, D, E, F, G any](a A, b B, c C, d D, e E, f F, g G) (A, B, C, D, E, F, G) {
	return a, b, c, d, e, f, g
}
func Ƭ8[A, B, C, D, E, F, G, H any](a A, b B, c C, d D, e E, f F, g G, h H) (A, B, C, D, E, F, G, H) {
	return a, b, c, d, e, f, g, h
}
func Ƭ9[A, B, C, D, E, F, G, H, I any](a A, b B, c C, d D, e E, f F, g G, h H, i I) (A, B, C, D, E, F, G, H, I) {
	return a, b, c, d, e, f, g, h, i
}

type E𝗿𝗿𝗼𝗿 = error
