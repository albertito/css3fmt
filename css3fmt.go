// css3fmt is an auto-indenter/formatter for CSS.
//
// See https://blitiri.com.ar/git/r/css3fmt for more details.
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/gorilla/css/scanner"
)

var (
	rewrite = flag.Bool("w", false,
		"Do not print reformatted sources to standard output."+
			" If a file's formatting is different from gofmt's,"+
			" overwrite it with the formatted version")
)

func main() {
	flag.Parse()

	if args := flag.Args(); len(args) > 0 {
		for _, fname := range args {
			f, err := os.Open(fname)
			if err != nil {
				fatalf("%s: %s\n", fname, err)
			}
			defer f.Close()

			s := indent(f)
			if *rewrite {
				err = ioutil.WriteFile(fname, []byte(s), 0660)
				if err != nil {
					fatalf("%s: %s\n", fname, err)
				}
			} else {
				os.Stdout.WriteString(s)
			}
		}
	} else {
		os.Stdout.WriteString(indent(os.Stdin))
	}
}

func fatalf(s string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, s, a...)
	os.Exit(1)
}

func indent(f *os.File) string {
	buf, err := ioutil.ReadAll(f)
	if err != nil {
		fatalf("%s: error reading: %s\n", f.Name(), err)
	}

	out := output{}

	s := scanner.New(string(buf))
scan:
	for {
		t := s.Next()
		switch t.Type {
		case scanner.TokenEOF:
			break scan
		case scanner.TokenError:
			fatalf("%s:%d:%d: error tokenizing: %s\n",
				f.Name(), t.Line, t.Column, t.Value)
		case scanner.TokenChar:
			switch t.Value {
			case "{":
				out.emit("{\n")
				out.indent++
			case "}":
				out.indent--
				out.emit("}\n")
				if out.indent == 0 {
					out.emit("\n")
				}
			case ";":
				if out.inFunc {
					out.inFunc = false
					out.indent--
				}
				out.emit(";\n")
			default:
				out.emit(t.Value)
			}
		case scanner.TokenS:
			if !strings.Contains(t.Value, "\n") {
				out.emit(" ")
			} else if out.inFunc {
				// Respect newline within functions, as users know best how to
				// break up arguments.
				out.emit("\n")
			}
			out.afterEmptyLine = strings.Contains(t.Value, "\n\n")
		case scanner.TokenComment:
			out.emitComment(t.Value)
		case scanner.TokenFunction:
			out.emit(t.Value)
			if !out.inFunc {
				out.inFunc = true
				out.indent++
			}
		default:
			out.emit(t.Value)
		}

		//fmt.Printf("\n«%s»\n", t)
	}

	return strings.Trim(out.buf.String(), "\n") + "\n"
}

type output struct {
	indent int

	afterN         bool
	afterEmptyLine bool
	inFunc         bool

	buf strings.Builder
}

func (o *output) emit(s string) {
	// Indent if we just came from a newline, UNLESS we're only printing a
	// newline, to avoid trailing spaces.
	if o.afterN && s != "\n" {
		for i := 0; i < o.indent; i++ {
			//o.buf.WriteString("‧‧‧‧")
			o.buf.WriteString("    ")
		}
	}

	o.buf.WriteString(s)
	o.afterN = strings.HasSuffix(s, "\n")
	o.afterEmptyLine = false
}

func (o *output) emitComment(s string) {
	// We preserve empty newlines before comments, so they can be used to
	// break long series of entries, or to group sections.
	if o.afterEmptyLine {
		o.emit("\n")
	}

	// Emit the lines adjusting indentation for "* ".
	lines := strings.Split(s, "\n")
	for _, l := range lines {
		if strings.HasPrefix(trimAllSp(l), "* ") {
			l = " " + trimAllSp(l)
		}
		o.emit(l + "\n")
	}
}

func trimAllSp(s string) string {
	return strings.Trim(s, " \t\r\n")
}
