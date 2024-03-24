//go:build !try

// Code generated by github.com/goghcrow/go-try DO NOT EDIT.
package test

func tuple_assign_map_index() error {
	{
		𝘃𝗮𝗹𝟭, 𝗲𝗿𝗿𝟭 := ret1Err[int]()
		if 𝗲𝗿𝗿𝟭 != nil {
			return 𝗲𝗿𝗿𝟭
		}
		_, _ = map[int]int{}[𝘃𝗮𝗹𝟭]
	}
	{
		𝘃𝗮𝗹𝟮, 𝗲𝗿𝗿𝟮 := ret1Err[int]()
		if 𝗲𝗿𝗿𝟮 != nil {
			return 𝗲𝗿𝗿𝟮
		}
		v, ok := map[int]int{}[𝘃𝗮𝗹𝟮]
		_, _ = v, ok
	}
	{
		𝘃𝗮𝗹𝟯, 𝗲𝗿𝗿𝟯 := ret1Err[int]()
		if 𝗲𝗿𝗿𝟯 != nil {
			return 𝗲𝗿𝗿𝟯
		}
		var v, ok = map[int]int{}[𝘃𝗮𝗹𝟯]
		_, _ = v, ok
	}
	{
		var v int
		𝘃𝗮𝗹𝟰, 𝗲𝗿𝗿𝟰 := ret1Err[int]()
		if 𝗲𝗿𝗿𝟰 != nil {
			return 𝗲𝗿𝗿𝟰
		}
		v, _ = map[int]int{}[𝘃𝗮𝗹𝟰]
		_ = v
	}
	{
		var v int
		𝘃𝗮𝗹𝟱, 𝗲𝗿𝗿𝟱 := ret1Err[int]()
		if 𝗲𝗿𝗿𝟱 != nil {
			return 𝗲𝗿𝗿𝟱
		}
		v, ok := map[int]int{}[𝘃𝗮𝗹𝟱]
		_, _ = v, ok
	}
	{
		var v int
		var ok bool
		𝘃𝗮𝗹𝟲, 𝗲𝗿𝗿𝟲 := ret1Err[int]()
		if 𝗲𝗿𝗿𝟲 != nil {
			return 𝗲𝗿𝗿𝟲
		}
		v, ok = map[int]int{}[𝘃𝗮𝗹𝟲]
		_, _ = v, ok
	}
	return nil
}
func tuple_assign_any_map_index() error {
	{
		𝘃𝗮𝗹𝟭, 𝗲𝗿𝗿𝟭 := ret1Err[int]()
		if 𝗲𝗿𝗿𝟭 != nil {
			return 𝗲𝗿𝗿𝟭
		}
		_, _ = map[any]int{}[𝘃𝗮𝗹𝟭]
	}
	{
		𝘃𝗮𝗹𝟮, 𝗲𝗿𝗿𝟮 := ret1Err[int]()
		if 𝗲𝗿𝗿𝟮 != nil {
			return 𝗲𝗿𝗿𝟮
		}
		v, ok := map[any]int{}[𝘃𝗮𝗹𝟮]
		_, _ = v, ok
	}
	{
		𝘃𝗮𝗹𝟯, 𝗲𝗿𝗿𝟯 := ret1Err[int]()
		if 𝗲𝗿𝗿𝟯 != nil {
			return 𝗲𝗿𝗿𝟯
		}
		var v, ok = map[any]int{}[𝘃𝗮𝗹𝟯]
		_, _ = v, ok
	}
	{
		var v int
		𝘃𝗮𝗹𝟰, 𝗲𝗿𝗿𝟰 := ret1Err[int]()
		if 𝗲𝗿𝗿𝟰 != nil {
			return 𝗲𝗿𝗿𝟰
		}
		v, _ = map[any]int{}[𝘃𝗮𝗹𝟰]
		_ = v
	}
	{
		var v int
		𝘃𝗮𝗹𝟱, 𝗲𝗿𝗿𝟱 := ret1Err[int]()
		if 𝗲𝗿𝗿𝟱 != nil {
			return 𝗲𝗿𝗿𝟱
		}
		v, ok := map[any]int{}[𝘃𝗮𝗹𝟱]
		_, _ = v, ok
	}
	{
		var v int
		var ok bool
		𝘃𝗮𝗹𝟲, 𝗲𝗿𝗿𝟲 := ret1Err[int]()
		if 𝗲𝗿𝗿𝟲 != nil {
			return 𝗲𝗿𝗿𝟲
		}
		v, ok = map[any]int{}[𝘃𝗮𝗹𝟲]
		_, _ = v, ok
	}
	return nil
}
