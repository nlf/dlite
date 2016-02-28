package plist

import (
	"bytes"
	"io/ioutil"
	"testing"
)

func BenchmarkOpenStepGenerate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		d := newTextPlistGenerator(ioutil.Discard, OpenStepFormat)
		d.generateDocument(plistValueTree)
	}
}

func BenchmarkOpenStepParse(b *testing.B) {
	buf := bytes.NewReader([]byte(plistValueTreeAsOpenStep))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StartTimer()
		d := newTextPlistParser(buf)
		d.parseDocument()
		b.StopTimer()
		buf.Seek(0, 0)
	}
}

func BenchmarkGNUStepParse(b *testing.B) {
	buf := bytes.NewReader([]byte(plistValueTreeAsGNUStep))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StartTimer()
		d := newTextPlistParser(buf)
		d.parseDocument()
		b.StopTimer()
		buf.Seek(0, 0)
	}
}

func TestTextCommentDecode(t *testing.T) {
	var testData = `{
		A=1 /* A is 1 because it is the first letter */;
		B=2; // B is 2 because comment-to-end-of-line.
		C=3;
		S = /not/a/comment/;
		S2 = /not*a/*comm*en/t;
	}`
	type D struct {
		A, B, C int
		S       string
		S2      string
	}
	actual := D{1, 2, 3, "/not/a/comment/", "/not*a/*comm*en/t"}
	var parsed D
	buf := bytes.NewReader([]byte(testData))
	decoder := NewDecoder(buf)
	err := decoder.Decode(&parsed)
	if err != nil {
		t.Error(err.Error())
	}

	if actual != parsed {
		t.Logf("Expected: %#v", actual)
		t.Logf("Received: %#v", parsed)
		t.Fail()
	}
}
