package main

import (
	"encoding/hex"
	"fmt"
	"io"
	"reflect"
	"sort"
	"time"
)

func PrettyPrint(w io.Writer, val interface{}) {
	printValue(w, val, "")
}

func printMap(w io.Writer, tv reflect.Value, depth string) {
	fmt.Fprintf(w, "{\n")
	ss := make(sort.StringSlice, tv.Len())
	i := 0
	for _, kval := range tv.MapKeys() {
		if kval.Kind() == reflect.Interface {
			kval = kval.Elem()
		}

		if kval.Kind() != reflect.String {
			continue
		}

		ss[i] = kval.String()
		i++
	}
	sort.Sort(ss)
	for _, k := range ss {
		val := tv.MapIndex(reflect.ValueOf(k))
		v := val.Interface()
		nd := depth + "  "
		for i := 0; i < len(k)+2; i++ {
			nd += " "
		}
		fmt.Fprintf(w, "  %s%s: ", depth, k)
		printValue(w, v, nd)
	}
	fmt.Fprintf(w, "%s}\n", depth)
}

func printValue(w io.Writer, val interface{}, depth string) {
	switch tv := val.(type) {
	case map[interface{}]interface{}:
		printMap(w, reflect.ValueOf(tv), depth)
	case map[string]interface{}:
		printMap(w, reflect.ValueOf(tv), depth)
	case []interface{}:
		fmt.Fprintf(w, "(\n")
		for i, v := range tv {
			id := fmt.Sprintf("[%d]", i)
			nd := depth + "  "
			for i := 0; i < len(id)+2; i++ {
				nd += " "
			}
			fmt.Fprintf(w, "  %s%s: ", depth, id)
			printValue(w, v, nd)
		}
		fmt.Fprintf(w, "%s)\n", depth)
	case int64, uint64, string, float32, float64, bool, time.Time:
		fmt.Fprintf(w, "%+v\n", tv)
	case uint8:
		fmt.Fprintf(w, "0x%2.02x\n", tv)
	case []byte:
		l := len(tv)
		sxl := l / 16
		if l%16 > 0 {
			sxl++
		}
		sxl *= 16
		var buf [4]byte
		var off [8]byte
		var asc [16]byte
		var ol int
		for i := 0; i < sxl; i++ {
			if i%16 == 0 {
				if i > 0 {
					io.WriteString(w, depth)
				}
				buf[0] = byte(i >> 24)
				buf[1] = byte(i >> 16)
				buf[2] = byte(i >> 8)
				buf[3] = byte(i)
				hex.Encode(off[:], buf[:])
				io.WriteString(w, string(off[:])+"  ")
			}
			if i < l {
				hex.Encode(off[:], tv[i:i+1])
				if tv[i] < 32 || tv[i] > 126 {
					asc[i%16] = '.'
				} else {
					asc[i%16] = tv[i]
				}
			} else {
				off[0] = ' '
				off[1] = ' '
				asc[i%16] = '.'
			}
			off[2] = ' '
			ol = 3
			if i%16 == 7 || i%16 == 15 {
				off[3] = ' '
				ol = 4
			}
			io.WriteString(w, string(off[:ol]))
			if i%16 == 15 {
				io.WriteString(w, "|"+string(asc[:])+"|\n")
			}
		}
	default:
		fmt.Fprintf(w, "%#v\n", val)
	}
}
