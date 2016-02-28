# Ply
Property list pretty-printer powered by `howett.net/plist`.

_verb. work with (a tool, especially one requiring steady, rhythmic movements)._

## Installation

`go get github.com/DHowett/go-plist/ply`

## Usage

```
  ply [OPTIONS]

Application Options:
  -c, --convert=<format>    convert the property list to a new format (c=list for list) (pretty)
  -k, --key=<keypath>       A keypath! (/)
  -o, --out=<filename>      output filename
  -I, --indent              indent indentable output formats (xml, openstep, gnustep, json)

Help Options:
  -h, --help                Show this help message
```

## Features

### Keypath evaluation

```
$ ply file.plist
{
  x: {
       y: {
            z: 1024
          }
     }
}
$ ply -k x/y/z file.plist
1024
```

Keypaths are composed of a number of path expressions:

* `/name` - dictionary key access
* `[i]` - index array, string, or data
* `[i:j]` - silce array, string, or data in the range `[i, j)`
* `!` - parse the data value as a property list and use it as the base of evaluation for further path components
* `$(subexpression)` - evaluate `subexpression` and paste its value

#### Examples

Given the following property list:

```
{
	a = {
		b = {
			c = (1, 2, 3);
			d = hello;
		};
		data = <414243>;
	};
	sub = <7b0a0974 6869733d 22612064 69637469 6f6e6172 7920696e 73696465 20616e6f 74686572 20706c69 73742122 3b7d>;
	hello = subexpression;
}
```

##### pretty print
```
$ ply file.plist
{
  a: {
       b: {
            c: (
                 [0]: 1
                 [1]: 2
                 [2]: 3
               )
            d: hello
          }
       data: 00000000  41 42 43                                          |ABC.............|
     }
  hello: subexpression
  sub: 00000000  7b 0a 09 74 68 69 73 3d  22 61 20 64 69 63 74 69  |{..this="a dicti|
       00000010  6f 6e 61 72 79 20 69 6e  73 69 64 65 20 61 6e 6f  |onary inside ano|
       00000020  74 68 65 72 20 70 6c 69  73 74 21 22 3b 7d        |ther plist!";}..|
}
```

##### consecutive dictionary keys
```
$ ply file.plist -k 'a/b/d'
hello
```

##### array indexing
```
$ ply file.plist -k 'a/b/c[1]'
2
```

##### data hexdump
```
$ ply file.plist -k 'a/data'
00000000  41 42 43                                          |ABC.............|
```

##### data and array slicing
```
$ ply file.plist -k 'a/data[2:3]'
00000000  43                                                |C...............|
```

```
$ ply -k 'sub[0:10]' file.plist
00000000  7b 0a 09 74 68 69 73 3d  22 61                    |{..this="a......|
```

##### subplist parsing
```
$ ply -k 'sub!' file.plist
{
  this: a dictionary inside another plist!
}
```

##### subplist keypath evaluation
```
$ ply -k 'sub!/this' file.plist
a dictionary inside another plist!
```

##### subexpression evaluation
```
$ ply -k '/$(/a/b/d)' file.plist
subexpression
```

### Property list conversion

`-c <format>`, or `-c list` to list them all.

* Binary property list [`bplist`]
* XML [`xml`]
* GNUstep [`gnustep`, `gs`]
* OpenStep [`openstep`, `os`]
* JSON (for a subset of data types) [`json`]
* YAML [`yaml`]

#### Notes
By default, ply will emit the most compact representation it can for a given format. The `-I` flag influences the inclusion of whitespace.

Ply will overwrite the input file unless an output filename is specified with `-o <file>`.

### Property list subsetting

(and subset conversion)

```
$ ply -k '/a/b' -o file-a-b.plist -c openstep -I file.plist
$ cat file-a-b.plist
{
	c = (
		1,
		2,
		3,
	);
	d = hello;
}
```

#### Subplist extraction

```
$ ply -k '/sub!' -o file-sub.plist -c openstep -I file.plist
$ cat file-sub.plist
{
	this = "a dictionary inside another plist!";
}
```
