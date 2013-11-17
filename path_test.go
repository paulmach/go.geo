package geo

import (
	"math"
	"math/rand"
	"testing"
)

func TestEncode(t *testing.T) {
	for loop := 0; loop < 100; loop++ {
		p := NewPath()
		for i := 0; i < 100; i++ {
			p.Push(&Point{rand.Float64(), rand.Float64()})
		}

		encoded := p.Encode()
		for _, c := range encoded {
			if c < 63 || c > 127 {
				t.Errorf("Encode, result out of range: %d", c)
			}
		}
	}
}

func TestEncodeDecode(t *testing.T) {
	for loop := 0; loop < 100; loop++ {

		p := NewPath()
		for i := 0; i < 100; i++ {
			p.Push(&Point{rand.Float64(), rand.Float64()})
		}

		encoded := p.Encode()
		path := Decode(encoded)

		if path.Length() != 100 {
			t.Fatalf("EncodeDecode, length mismatch: %d != 100", path.Length())
		}

		for i := 0; i < 100; i++ {
			a := p.GetAt(i)
			b := path.GetAt(i)

			if e := math.Abs(a[0] - b[0]); e > 1e-5 {
				t.Errorf("EncodeDecode, X error too big: %f", e)
			}

			if e := math.Abs(a[1] - b[1]); e > 1e-5 {
				t.Errorf("EncodeDecode, Y error too big: %f", e)
			}
		}
	}
}
