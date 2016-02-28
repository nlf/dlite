package plist

import (
	"bytes"
	"io/ioutil"
	"testing"
)

func BenchmarkXMLGenerate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		d := newXMLPlistGenerator(ioutil.Discard)
		d.generateDocument(plistValueTree)
	}
}

func BenchmarkXMLParse(b *testing.B) {
	buf := bytes.NewReader([]byte(plistValueTreeAsXML))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StartTimer()
		d := newXMLPlistParser(buf)
		d.parseDocument()
		b.StopTimer()
		buf.Seek(0, 0)
	}
}

func TestVariousIllegalXMLPlists(t *testing.T) {
	plists := []string{
		"<plist><doct><key>helo</key><string></string></doct></plist>",
		"<plist><dict><string>helo</string></dict></plist>",
		"<plist><dict><key>helo</key></dict></plist>",
		"<plist><integer>helo</integer></plist>",
		"<plist><real>helo</real></plist>",
		"<plist><data>*@&amp;%#helo</data></plist>",
		"<plist><date>*@&amp;%#helo</date></plist>",
		"<plist><date>*@&amp;%#helo</date></plist>",
		"<plist><integer>10</plist>",
		"<plist><real>10</plist>",
		"<plist><string>10</plist>",
		"<plist><dict>10</plist>",
		"<plist><dict><key>10</plist>",
		"<plist>",
		"<plist><data>",
		"<plist><date>",
		"<plist><array>",
		"<pl",
		"bplist00",
	}

	for _, plist := range plists {
		buf := bytes.NewReader([]byte(plist))
		d := newXMLPlistParser(buf)
		_, err := d.parseDocument()
		t.Logf("Error: %v", err)
		if err == nil {
			t.Error("Expected error, received nothing.")
		}
	}
}
