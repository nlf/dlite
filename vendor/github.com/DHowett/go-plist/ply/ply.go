package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	"github.com/jessevdk/go-flags"
	"gopkg.in/yaml.v1"
	"howett.net/plist"
)

//import "github.com/mgutz/ansi"

const (
	PrettyFormat = 100 + iota
	JSONFormat
	YAMLFormat
	RawFormat
)

var nameFormatMap = map[string]int{
	"x":        plist.XMLFormat,
	"xml":      plist.XMLFormat,
	"xml1":     plist.XMLFormat,
	"b":        plist.BinaryFormat,
	"bin":      plist.BinaryFormat,
	"binary":   plist.BinaryFormat,
	"binary1":  plist.BinaryFormat,
	"o":        plist.OpenStepFormat,
	"os":       plist.OpenStepFormat,
	"openstep": plist.OpenStepFormat,
	"step":     plist.OpenStepFormat,
	"g":        plist.GNUStepFormat,
	"gs":       plist.GNUStepFormat,
	"gnustep":  plist.GNUStepFormat,
	"pretty":   PrettyFormat,
	"json":     JSONFormat,
	"yaml":     YAMLFormat,
	"r":        RawFormat,
	"raw":      RawFormat,
}

var opts struct {
	Convert string `short:"c" long:"convert" description:"convert the property list to a new format (c=list for list)" default:"pretty" value-name:"<format>"`
	Keypath string `short:"k" long:"key" description:"A keypath!" default:"/" value-name:"<keypath>"`
	Output  string `short:"o" long:"out" description:"output filename" default:"" value-name:"<filename>"`
	Indent  bool   `short:"I" long:"indent" description:"indent indentable output formats (xml, openstep, gnustep, json)"`
}

func main() {
	parser := flags.NewParser(&opts, flags.Default)
	args, err := parser.Parse()
	if err != nil {
		parser.WriteHelp(os.Stderr)
		fmt.Fprintln(os.Stderr, err)
		return
	}

	if opts.Convert == "list" {
		formats := make([]string, len(nameFormatMap))
		i := 0
		for k, _ := range nameFormatMap {
			formats[i] = k
			i++
		}

		fmt.Fprintln(os.Stderr, "Supported output formats:")
		fmt.Fprintln(os.Stderr, strings.Join(formats, ", "))
		return
	}

	if len(args) < 1 {
		parser.WriteHelp(os.Stderr)
		return
	}

	filename := args[0]

	keypath := opts.Keypath
	if len(keypath) == 0 {
		c := strings.Index(filename, ":")
		if c > -1 {
			keypath = filename[c+1:]
			filename = filename[:c]
		}
	}

	file, err := os.Open(filename)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}

	var val interface{}
	switch strings.ToLower(filepath.Ext(filename)) {
	case ".json", ".yaml", ".yml":
		buf := &bytes.Buffer{}
		io.Copy(buf, file)
		err = yaml.Unmarshal(buf.Bytes(), &val)
	default:
		dec := plist.NewDecoder(file)
		err = dec.Decode(&val)
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}
	file.Close()

	convert := strings.ToLower(opts.Convert)
	format, ok := nameFormatMap[convert]
	if !ok {
		fmt.Fprintf(os.Stderr, "unknown output format %s\n", convert)
		return
	}

	output := opts.Output
	newline := false
	var outputStream io.WriteCloser
	if format < PrettyFormat && output == "" {
		// Writing a plist, but no output filename. Save to original.
		output = filename
	} else if format >= PrettyFormat && output == "" {
		// Writing a non-plist, but no output filename: Stdout
		outputStream = os.Stdout
		newline = true
	} else if output == "-" {
		// - means stdout.
		outputStream = os.Stdout
		newline = true
	}

	if outputStream == nil {
		outfile, err := os.Create(output)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			return
		}
		outputStream = outfile
	}

	keypathContext := &KeypathWalker{}
	rval, err := keypathContext.WalkKeypath(reflect.ValueOf(val), keypath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}
	val = rval.Interface()

	switch {
	case format >= 0 && format < PrettyFormat:
		enc := plist.NewEncoderForFormat(outputStream, format)
		if opts.Indent {
			enc.Indent("\t")
		}
		err := enc.Encode(val)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			return
		}
	case format == PrettyFormat:
		PrettyPrint(outputStream, rval.Interface())
	case format == JSONFormat:
		var out []byte
		var err error
		if opts.Indent {
			out, err = json.MarshalIndent(val, "", "\t")
		} else {
			out, err = json.Marshal(val)
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			return
		}
		outputStream.Write(out)
	case format == YAMLFormat:
		out, err := yaml.Marshal(val)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			return
		}
		outputStream.Write(out)
	case format == RawFormat:
		newline = false
		switch rval.Kind() {
		case reflect.String:
			outputStream.Write([]byte(val.(string)))
		case reflect.Slice:
			if rval.Elem().Kind() == reflect.Uint8 {
				outputStream.Write(val.([]byte))
			}
		default:
			binary.Write(outputStream, binary.LittleEndian, val)
		}
	}
	if newline {
		fmt.Fprintf(outputStream, "\n")
	}
	outputStream.Close()
}

type KeypathWalker struct {
	rootVal *reflect.Value
	curVal  reflect.Value
}

func (ctx *KeypathWalker) Split(data []byte, atEOF bool) (advance int, token []byte, err error) {
	mode, oldmode := 0, 0
	depth := 0
	tok, subexpr := "", ""
	// modes:
	// 0: normal string, separated by /
	// 1: array index (reading between [])
	// 2: found $, looking for ( or nothing
	// 3: found $(, reading subkey, looking for )
	// 4: "escape"? unused as yet.
	if len(data) == 0 && atEOF {
		return 0, nil, io.EOF
	}
each:
	for _, v := range data {
		advance++
		switch {
		case mode == 4:
			// Completing an escape sequence.
			tok += string(v)
			mode = 0
			continue each
		case mode == 0 && v == '/':
			if tok != "" {
				break each
			} else {
				continue each
			}
		case mode == 0 && v == '[':
			if tok != "" {
				// We have encountered a [ after text, we want only the text
				advance-- // We don't want to consume this character.
				break each
			} else {
				tok += string(v)
				mode = 1
			}
		case mode == 1 && v == ']':
			mode = 0
			tok += string(v)
			break each
		case mode == 0 && v == '!':
			if tok == "" {
				tok = "!"
				break each
			} else {
				// We have encountered a ! after text, we want the text
				advance-- // We don't want to consume this character.
				break each
			}
		case (mode == 0 || mode == 1) && v == '$':
			oldmode = mode
			mode = 2
		case mode == 2:
			if v == '(' {
				mode = 3
				depth++
				subexpr = ""
			} else {
				// We didn't emit the $ to begin with, so we have to do it here.
				tok += "$" + string(v)
				mode = 0
			}
		case mode == 3 && v == '(':
			subexpr += string(v)
			depth++
		case mode == 3 && v == ')':
			depth--
			if depth == 0 {
				newCtx := &KeypathWalker{rootVal: ctx.rootVal}
				subexprVal, e := newCtx.WalkKeypath(*ctx.rootVal, subexpr)
				if e != nil {
					return 0, nil, errors.New("Dynamic subexpression " + subexpr + " failed: " + e.Error())
				}
				if subexprVal.Kind() == reflect.Interface {
					subexprVal = subexprVal.Elem()
				}
				s := ""
				if subexprVal.Kind() == reflect.String {
					s = subexprVal.String()
				} else if subexprVal.Kind() == reflect.Uint64 {
					s = strconv.Itoa(int(subexprVal.Uint()))
				} else {
					return 0, nil, errors.New("Dynamic subexpression " + subexpr + " evaluated to non-string/non-int.")
				}
				tok += s
				mode = oldmode
			} else {
				subexpr += string(v)
			}
		case mode == 3:
			subexpr += string(v)
		default:
			tok += string(v)
		}

	}
	return advance, []byte(tok), nil
}

func (ctx *KeypathWalker) WalkKeypath(val reflect.Value, keypath string) (reflect.Value, error) {
	if keypath == "" {
		return val, nil
	}

	if ctx.rootVal == nil {
		ctx.rootVal = &val
	}

	ctx.curVal = val

	scanner := bufio.NewScanner(strings.NewReader(keypath))
	scanner.Split(ctx.Split)
	for scanner.Scan() {
		token := scanner.Text()
		if ctx.curVal.Kind() == reflect.Interface {
			ctx.curVal = ctx.curVal.Elem()
		}

		switch {
		case len(token) == 0:
			continue
		case token[0] == '[': // array
			s := token[1 : len(token)-1]
			if ctx.curVal.Kind() != reflect.Slice && ctx.curVal.Kind() != reflect.String {
				return reflect.ValueOf(nil), errors.New("keypath attempted to index non-indexable with " + s)
			}

			colon := strings.Index(s, ":")
			if colon > -1 {
				var err error
				var si, sj int
				is := s[:colon]
				js := s[colon+1:]
				if is != "" {
					si, err = strconv.Atoi(is)
					if err != nil {
						return reflect.ValueOf(nil), err
					}
				}
				if js != "" {
					sj, err = strconv.Atoi(js)
					if err != nil {
						return reflect.ValueOf(nil), err
					}
				}
				if si < 0 || sj > ctx.curVal.Len() {
					return reflect.ValueOf(nil), errors.New("keypath attempted to index outside of indexable with " + s)
				}
				ctx.curVal = ctx.curVal.Slice(si, sj)
			} else {
				idx, _ := strconv.Atoi(s)
				ctx.curVal = ctx.curVal.Index(idx)
			}
		case token[0] == '!': // subplist!
			if ctx.curVal.Kind() != reflect.Slice || ctx.curVal.Type().Elem().Kind() != reflect.Uint8 {
				return reflect.Value{}, errors.New("Attempted to subplist non-data.")
			}
			byt := ctx.curVal.Interface().([]uint8)
			buf := bytes.NewReader(byt)
			dec := plist.NewDecoder(buf)
			var subval interface{}
			dec.Decode(&subval)
			ctx.curVal = reflect.ValueOf(subval)
		default: // just a string
			if ctx.curVal.Kind() != reflect.Map {
				return reflect.ValueOf(nil), errors.New("keypath attempted to descend into non-map using key " + token)
			}
			if token != "" {
				ctx.curVal = ctx.curVal.MapIndex(reflect.ValueOf(token))
			}
		}
	}
	err := scanner.Err()
	if err != nil {
		return reflect.ValueOf(nil), err
	}
	return ctx.curVal, nil
}
