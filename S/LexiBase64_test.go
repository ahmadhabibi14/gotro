package S

import (
	"fmt"
	"math"
	"testing"
	"time"
)

func TestInit(t *testing.T) {
	if len(i2c_cb63) != 64 {
		t.Errorf(`Character to integer array should have 64 items`)
	}
	if len(c2i_cb63) != 64 {
		t.Errorf(`Integer to character map should have 64 items`)
	}
}

func testOk[T int64 | uint64](t *testing.T, ok bool, oldv T, newv T, cb63 string, print bool) bool {
	if !ok {
		t.Errorf(`Error decoding [%s] back to integer, source: %d`, cb63, oldv)
		return false
	}
	if newv != oldv {
		t.Errorf(`Invalid decoding [%s] back to integer, should be: %d, got: %d`, cb63, oldv, newv)
		return false
	}
	if print {
		fmt.Printf("%20v to %12v back to %20v\n", oldv, cb63, newv)
	}
	return true
}

func testOk2(t *testing.T, ok bool, oldv string, newv string, cb63 int64, print bool) bool {
	if !ok {
		t.Errorf(`Error encoding %d back to string, source: [%s]`, cb63, oldv)
		return false
	}
	if newv != oldv {
		t.Errorf(`Invalid encoding %d back to string, should be: [%s], got: [%s]`, cb63, oldv, newv)
		return false
	}
	if print {
		fmt.Printf("%12v to %20v back to %12v\n", oldv, cb63, newv)
	}
	return true
}

func TestEncodeDecodeCB63(t *testing.T) {
	tested := map[int64]bool{}
	x := int64(1)
	globalCount := int64(1)
	for m := int64(2); m < math.MaxInt16*math.MaxInt8; m += x {
		if m == 3 {
			if testing.Short() {
				x = 14
			} else {
				x = 2
			}
		}
		last := int64(-1)
		testedCount := int64(0)
		for z := int64(0); z >= 0 && z < math.MaxInt64; z = z*m + 1 {
			if last == z {
				break
			}
			if tested[z] {
				testedCount += 1
				if testedCount > 4 {
					break
				}
				continue
			}
			cb := EncodeCB63(z, 0)
			v, ok := DecodeCB63[int64](cb)
			if !testOk(t, ok, z, v, cb, globalCount%100000 == 0) {
				return
			}
			tested[z] = true
			last = z
			globalCount += 1
		}
	}
	//fmt.Println(`Unique numbers tested:`, len(tested))
}

func TestEncodeDecodeCB63Uint(t *testing.T) {
	tested := map[uint64]bool{}
	x := uint64(1)
	globalCount := uint64(1)
	for m := uint64(2); m < math.MaxInt16*math.MaxInt8; m += x {
		if m == 3 {
			if testing.Short() {
				x = 14
			} else {
				x = 2
			}
		}
		first := uint64(0)
		last := uint64(math.MaxUint64)
		testedCount := uint64(0)
		for z := 0; z < 1000; z++ {
			if tested[first] {
				testedCount += 1
				if testedCount > 4 {
					break
				}
				continue
			}
			cb := EncodeCB63(first, 0)
			v, ok := DecodeCB63[uint64](cb)
			if !testOk(t, ok, first, v, cb, globalCount%100000 == 0) {
				return
			}
			tested[first] = true
			first *= 2 + 1

			cb = EncodeCB63(last, 0)
			v, ok = DecodeCB63[uint64](cb)
			if !testOk(t, ok, last, v, cb, globalCount%100000 == 0) {
				return
			}
			tested[last] = true
			last /= 2 - 1

			globalCount += 2
		}
	}
	//fmt.Println(`Unique numbers tested:`, len(tested))
}

func TestUnixNano(t *testing.T) {
	z := time.Now().UnixNano()
	cb := EncodeCB63(z, 0)
	v, ok := DecodeCB63[int64](cb)
	if !testOk(t, ok, z, v, cb, true) {
		return
	}
}

func TestMaxInt(t *testing.T) {
	arr := []int64{math.MaxInt8, math.MaxInt16, math.MaxInt32, math.MaxInt64, 0}
	z := int64(1)
	for {
		arr = append(arr, z)
		z *= 10
		if z < 0 {
			break
		}
		arr = append(arr, z-1)
	}
	for _, z := range arr {
		cb := EncodeCB63(z, MaxStrLenCB63)
		v, ok := DecodeCB63[int64](cb)
		if !testOk(t, ok, z, v, cb, true) {
			break
		}
	}
}

func TestMaxUInt(t *testing.T) {
	arr := []uint64{math.MaxUint8, math.MaxUint16, math.MaxUint32, math.MaxUint64, 0}
	z := uint64(1)
	for {
		arr = append(arr, z)
		z *= 10
		if z <= 0 {
			break
		}
		arr = append(arr, z-1)
	}
	for _, z := range arr {
		cb := EncodeCB63(z, MaxStrLenCB63)
		v, ok := DecodeCB63[uint64](cb)
		if !testOk(t, ok, z, v, cb, true) {
			break
		}
	}
}

func TestRepeated(t *testing.T) {
	arr := []string{}
	for i := 1; i < 12; i++ {
		for _, c := range i2c_cb63 {
			arr = append(arr, Repeat(string(c), i))
			if i == MaxStrLenCB63 && c == '6' {
				break // max is '6'+'z'*10
			}
		}
	}
	for _, z := range arr {
		v, ok := DecodeCB63[int64](z)
		cb := EncodeCB63(v, len(z))
		if !testOk2(t, ok, z, cb, v, true) {
			break
		}
	}
}

func TestRandomCB63(t *testing.T) {
	m := map[string]int{}
	for z := 0; z < 1000; z++ {
		v := RandomCB63(2)[:5]
		m[v]++
		if m[v] > 1 {
			t.Errorf(`RandomCB63 should return unique values: dup=%v iter=%v`, v, z+1)
			break
		}
	}
}
