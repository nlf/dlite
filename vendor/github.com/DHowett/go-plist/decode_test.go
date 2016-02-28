package plist

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"
)

func BenchmarkXMLDecode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		var bval interface{}
		buf := bytes.NewReader([]byte(plistValueTreeAsXML))
		b.StartTimer()
		decoder := NewDecoder(buf)
		decoder.Decode(bval)
		b.StopTimer()
	}
}

func BenchmarkBplistDecode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		var bval interface{}
		buf := bytes.NewReader(plistValueTreeAsBplist)
		b.StartTimer()
		decoder := NewDecoder(buf)
		decoder.Decode(bval)
		b.StopTimer()
	}
}

func TestLaxDecode(t *testing.T) {
	var laxTestDataStringsOnlyAsXML = `{B=1;D="2013-11-27 00:34:00 +0000";I64=1;F64="3.0";U64=2;}`
	d := LaxTestData{}
	buf := bytes.NewReader([]byte(laxTestDataStringsOnlyAsXML))
	decoder := NewDecoder(buf)
	decoder.lax = true
	err := decoder.Decode(&d)
	if err != nil {
		t.Error(err.Error())
	}

	if d != laxTestData {
		t.Logf("Expected: %#v", laxTestData)
		t.Logf("Received: %#v", d)
		t.Fail()
	}
}

func TestIllegalLaxDecode(t *testing.T) {
	i := int64(0)
	u := uint64(0)
	f := float64(0)
	b := false
	plists := []struct {
		pl string
		d  interface{}
	}{
		{"<string>abc</string>", &i},
		{"<string>abc</string>", &u},
		{"<string>def</string>", &f},
		{"<string>ghi</string>", &b},
		{"<string>jkl</string>", []byte{0x00}},
	}

	for _, plist := range plists {
		buf := bytes.NewReader([]byte(plist.pl))
		decoder := NewDecoder(buf)
		decoder.lax = true
		err := decoder.Decode(plist.d)
		t.Logf("Error: %v", err)
		if err == nil {
			t.Error("Expected error, received nothing.")
		}
	}
}

func TestIllegalDecode(t *testing.T) {
	i := int64(0)
	b := false
	plists := []struct {
		pl string
		d  interface{}
	}{
		{"<string>abc</string>", &i},
		{"<data>ABC=</data>", &i},
		{"<real>34.1</real>", &i},
		{"<true>def</true>", &i},
		{"<date>2010-01-01T00:00:00Z</date>", &i},
		{"<integer>0</integer>", &b},
		{"<array><integer>0</integer></array>", &b},
		{"<dict><key>a</key><integer>0</integer></dict>", &b},
		{"<array><true/><true/><true/></array>", &[1]int{1}},
	}

	for _, plist := range plists {
		buf := bytes.NewReader([]byte(plist.pl))
		decoder := NewDecoder(buf)
		err := decoder.Decode(plist.d)
		t.Logf("Error: %v", err)
		if err == nil {
			t.Error("Expected error, received nothing.")
		}
	}
}

func TestDecode(t *testing.T) {
	var failed bool
	for _, test := range tests {
		failed = false

		t.Logf("Testing Decode (%s)", test.Name)

		d := test.DecodeData
		if d == nil {
			d = test.Data
		}

		testData := reflect.ValueOf(d)
		if !testData.IsValid() || isEmptyInterface(testData) {
			continue
		}
		if testData.Kind() == reflect.Ptr || testData.Kind() == reflect.Interface {
			testData = testData.Elem()
		}
		d = testData.Interface()

		results := make(map[int]interface{})
		errors := make(map[int]error)
		for fmt, dat := range test.Expected {
			if test.SkipDecode[fmt] {
				continue
			}
			val := reflect.New(testData.Type()).Interface()
			_, errors[fmt] = Unmarshal(dat, val)

			vt := reflect.ValueOf(val)
			if vt.Kind() == reflect.Ptr || vt.Kind() == reflect.Interface {
				vt = vt.Elem()
				val = vt.Interface()
			}

			results[fmt] = val

			if !reflect.DeepEqual(d, val) {
				failed = true
			}
		}

		if results[BinaryFormat] != nil && results[XMLFormat] != nil {
			if !reflect.DeepEqual(results[BinaryFormat], results[XMLFormat]) {
				t.Log("Binary and XML decoding yielded different values.")
				t.Log("Binary:", results[BinaryFormat])
				t.Log("XML   :", results[XMLFormat])
				failed = true
			}
		}

		if failed {
			t.Logf("Expected: %#v\n", d)

			for fmt, dat := range results {
				t.Logf("Received %s: %#v\n", FormatNames[fmt], dat)
			}
			for fmt, err := range errors {
				if err != nil {
					t.Logf("Error %s: %v\n", FormatNames[fmt], err)
				}
			}
			t.Log("FAILED")
			t.Fail()
		}
	}
}

func TestInterfaceDecode(t *testing.T) {
	var xval interface{}
	buf := bytes.NewReader([]byte{98, 112, 108, 105, 115, 116, 48, 48, 214, 1, 13, 17, 21, 25, 27, 2, 14, 18, 22, 26, 28, 88, 105, 110, 116, 97, 114, 114, 97, 121, 170, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 16, 1, 16, 8, 16, 16, 16, 32, 16, 64, 16, 2, 16, 9, 16, 17, 16, 33, 16, 65, 86, 102, 108, 111, 97, 116, 115, 162, 15, 16, 34, 66, 0, 0, 0, 35, 64, 80, 0, 0, 0, 0, 0, 0, 88, 98, 111, 111, 108, 101, 97, 110, 115, 162, 19, 20, 9, 8, 87, 115, 116, 114, 105, 110, 103, 115, 162, 23, 24, 92, 72, 101, 108, 108, 111, 44, 32, 65, 83, 67, 73, 73, 105, 0, 72, 0, 101, 0, 108, 0, 108, 0, 111, 0, 44, 0, 32, 78, 22, 117, 76, 84, 100, 97, 116, 97, 68, 1, 2, 3, 4, 84, 100, 97, 116, 101, 51, 65, 184, 69, 117, 120, 0, 0, 0, 8, 21, 30, 41, 43, 45, 47, 49, 51, 53, 55, 57, 59, 61, 68, 71, 76, 85, 94, 97, 98, 99, 107, 110, 123, 142, 147, 152, 157, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 29, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 166})
	decoder := NewDecoder(buf)
	err := decoder.Decode(&xval)
	if err != nil {
		t.Log("Error:", err)
		t.Fail()
	}
}

func TestFormatDetection(t *testing.T) {
	type formatTest struct {
		expectedFormat int
		data           []byte
	}
	plists := []formatTest{
		{BinaryFormat, []byte{98, 112, 108, 105, 115, 116, 48, 48, 85, 72, 101, 108, 108, 111, 8, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 14}},
		{XMLFormat, []byte(`<string>&lt;*I3&gt;</string>`)},
		{InvalidFormat, []byte(`bplist00`)}, // Looks like a binary property list, and bplist does not have fallbacks(!)
		{OpenStepFormat, []byte(`(1,2,3,4,5)`)},
		{OpenStepFormat, []byte(`<abab>`)},
		{GNUStepFormat, []byte(`(1,2,<*I3>)`)},
		{OpenStepFormat, []byte{0x00}},
	}

	for _, fmttest := range plists {
		fmt, err := Unmarshal(fmttest.data, nil)
		if fmt != fmttest.expectedFormat {
			t.Errorf("Wanted %s, received %s.", FormatNames[fmttest.expectedFormat], FormatNames[fmt])
		}
		if err != nil {
			t.Logf("Error: %v", err)
		}
	}
}

func ExampleDecoder_Decode() {
	type sparseBundleHeader struct {
		InfoDictionaryVersion string `plist:"CFBundleInfoDictionaryVersion"`
		BandSize              uint64 `plist:"band-size"`
		BackingStoreVersion   int    `plist:"bundle-backingstore-version"`
		DiskImageBundleType   string `plist:"diskimage-bundle-type"`
		Size                  uint64 `plist:"size"`
	}

	buf := bytes.NewReader([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
	<dict>
		<key>CFBundleInfoDictionaryVersion</key>
		<string>6.0</string>
		<key>band-size</key>
		<integer>8388608</integer>
		<key>bundle-backingstore-version</key>
		<integer>1</integer>
		<key>diskimage-bundle-type</key>
		<string>com.apple.diskimage.sparsebundle</string>
		<key>size</key>
		<integer>4398046511104</integer>
	</dict>
</plist>`))

	var data sparseBundleHeader
	decoder := NewDecoder(buf)
	err := decoder.Decode(&data)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(data)

	// Output: {6.0 8388608 1 com.apple.diskimage.sparsebundle 4398046511104}
}
