package plist

import (
	"reflect"
	"testing"
	"time"
)

func BenchmarkStructMarshal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		e := &Encoder{}
		e.marshal(reflect.ValueOf(plistValueTreeRawData))
	}
}

func BenchmarkMapMarshal(b *testing.B) {
	data := map[string]interface{}{
		"intarray": []interface{}{
			int(1),
			int8(8),
			int16(16),
			int32(32),
			int64(64),
			uint(2),
			uint8(9),
			uint16(17),
			uint32(33),
			uint64(65),
		},
		"floats": []interface{}{
			float32(32.0),
			float64(64.0),
		},
		"booleans": []bool{
			true,
			false,
		},
		"strings": []string{
			"Hello, ASCII",
			"Hello, 世界",
		},
		"data": []byte{1, 2, 3, 4},
		"date": time.Date(2013, 11, 27, 0, 34, 0, 0, time.UTC),
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		e := &Encoder{}
		e.marshal(reflect.ValueOf(data))
	}
}
