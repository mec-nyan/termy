package byteme

import "unsafe"

// Memory safe language my a**
// We can use []byte(s) but it seems to be, lets say, not ideal...
func UnsafeStrToBytes(s string) []byte {
	if s == "" {
		return nil
	}
	return unsafe.Slice(unsafe.StringData(s), len(s))
}
