package plist

import (
	"bytes"
	"fmt"
	"testing"
)

func BenchmarkXMLEncode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewEncoder(&bytes.Buffer{}).Encode(plistValueTreeRawData)
	}
}

func BenchmarkBplistEncode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewBinaryEncoder(&bytes.Buffer{}).Encode(plistValueTreeRawData)
	}
}

func BenchmarkOpenStepEncode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewEncoderForFormat(&bytes.Buffer{}, OpenStepFormat).Encode(plistValueTreeRawData)
	}
}

func TestEncode(t *testing.T) {
	var failed bool
	for _, test := range tests {
		failed = false
		t.Logf("Testing Encode (%s)", test.Name)

		// A test that should render no output!
		errors := make(map[int]error)
		if test.ShouldFail && len(test.Expected) == 0 {
			_, err := Marshal(test.Data, XMLFormat)
			failed = failed || (test.ShouldFail && err == nil)
		}

		results := make(map[int][]byte)
		for fmt, dat := range test.Expected {
			results[fmt], errors[fmt] = Marshal(test.Data, fmt)
			failed = failed || (test.ShouldFail && errors[fmt] == nil)
			failed = failed || !bytes.Equal(dat, results[fmt])
		}

		if failed {
			t.Logf("Value: %#v", test.Data)
			if test.ShouldFail {
				t.Logf("Expected: Error")
			} else {
				printype := "%s"
				for fmt, dat := range test.Expected {
					if fmt == BinaryFormat {
						printype = "%+v"
					} else {
						printype = "%s"
					}
					t.Logf("Expected %s: "+printype+"\n", FormatNames[fmt], dat)
				}
			}

			printype := "%s"
			for fmt, dat := range results {
				if fmt == BinaryFormat {
					printype = "%+v"
				} else {
					printype = "%s"
				}
				t.Logf("Received %s: "+printype+"\n", FormatNames[fmt], dat)
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

func ExampleEncoder_Encode() {
	type sparseBundleHeader struct {
		InfoDictionaryVersion string `plist:"CFBundleInfoDictionaryVersion"`
		BandSize              uint64 `plist:"band-size"`
		BackingStoreVersion   int    `plist:"bundle-backingstore-version"`
		DiskImageBundleType   string `plist:"diskimage-bundle-type"`
		Size                  uint64 `plist:"size"`
	}
	data := &sparseBundleHeader{
		InfoDictionaryVersion: "6.0",
		BandSize:              8388608,
		Size:                  4 * 1048576 * 1024 * 1024,
		DiskImageBundleType:   "com.apple.diskimage.sparsebundle",
		BackingStoreVersion:   1,
	}

	buf := &bytes.Buffer{}
	encoder := NewEncoder(buf)
	err := encoder.Encode(data)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(buf.String())

	// Output: <?xml version="1.0" encoding="UTF-8"?>
	// <!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
	// <plist version="1.0"><dict><key>CFBundleInfoDictionaryVersion</key><string>6.0</string><key>band-size</key><integer>8388608</integer><key>bundle-backingstore-version</key><integer>1</integer><key>diskimage-bundle-type</key><string>com.apple.diskimage.sparsebundle</string><key>size</key><integer>4398046511104</integer></dict></plist>
}

func ExampleMarshal_xml() {
	type sparseBundleHeader struct {
		InfoDictionaryVersion string `plist:"CFBundleInfoDictionaryVersion"`
		BandSize              uint64 `plist:"band-size"`
		BackingStoreVersion   int    `plist:"bundle-backingstore-version"`
		DiskImageBundleType   string `plist:"diskimage-bundle-type"`
		Size                  uint64 `plist:"size"`
	}
	data := &sparseBundleHeader{
		InfoDictionaryVersion: "6.0",
		BandSize:              8388608,
		Size:                  4 * 1048576 * 1024 * 1024,
		DiskImageBundleType:   "com.apple.diskimage.sparsebundle",
		BackingStoreVersion:   1,
	}

	plist, err := MarshalIndent(data, XMLFormat, "\t")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(plist))

	// Output: <?xml version="1.0" encoding="UTF-8"?>
	// <!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
	// <plist version="1.0">
	// 	<dict>
	// 		<key>CFBundleInfoDictionaryVersion</key>
	// 		<string>6.0</string>
	// 		<key>band-size</key>
	// 		<integer>8388608</integer>
	// 		<key>bundle-backingstore-version</key>
	// 		<integer>1</integer>
	// 		<key>diskimage-bundle-type</key>
	// 		<string>com.apple.diskimage.sparsebundle</string>
	// 		<key>size</key>
	// 		<integer>4398046511104</integer>
	// 	</dict>
	// </plist>
}

func ExampleMarshal_gnustep() {
	type sparseBundleHeader struct {
		InfoDictionaryVersion string `plist:"CFBundleInfoDictionaryVersion"`
		BandSize              uint64 `plist:"band-size"`
		BackingStoreVersion   int    `plist:"bundle-backingstore-version"`
		DiskImageBundleType   string `plist:"diskimage-bundle-type"`
		Size                  uint64 `plist:"size"`
	}
	data := &sparseBundleHeader{
		InfoDictionaryVersion: "6.0",
		BandSize:              8388608,
		Size:                  4 * 1048576 * 1024 * 1024,
		DiskImageBundleType:   "com.apple.diskimage.sparsebundle",
		BackingStoreVersion:   1,
	}

	plist, err := MarshalIndent(data, GNUStepFormat, "\t")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(plist))

	// Output: {
	// 	CFBundleInfoDictionaryVersion = 6.0;
	// 	band-size = <*I8388608>;
	// 	bundle-backingstore-version = <*I1>;
	// 	diskimage-bundle-type = com.apple.diskimage.sparsebundle;
	// 	size = <*I4398046511104>;
	// }
}
