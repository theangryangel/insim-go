package geometry

type FixedPoint struct {
	X int32 `struct:"int32"`
	Y int32 `struct:"int32"`
	Z int32 `struct:"int32"`
}

type FloatingPoint struct {
	X float32 `struct:"float32"`
	Y float32 `struct:"float32"`
	Z float32 `struct:"float32"`
}
