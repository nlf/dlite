package plist

import (
	"reflect"
	"testing"
	"time"
)

func BenchmarkStructUnmarshal(b *testing.B) {
	type Data struct {
		Intarray []uint64  `plist:"intarray"`
		Floats   []float64 `plist:"floats"`
		Booleans []bool    `plist:"booleans"`
		Strings  []string  `plist:"strings"`
		Dat      []byte    `plist:"data"`
		Date     time.Time `plist:"date"`
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var xval Data
		d := &Decoder{}
		d.unmarshal(plistValueTree, reflect.ValueOf(&xval))
	}
}

func BenchmarkInterfaceUnmarshal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var xval interface{}
		d := &Decoder{}
		d.unmarshal(plistValueTree, reflect.ValueOf(&xval))
	}
}
