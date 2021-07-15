// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Stringer is a tool to automate the creation of methods that satisfy the fmt.Stringer
// interface. Given the name of a (signed or unsigned) integer type T that has constants
// defined, stringer will create a new self-contained Go source file implementing
//	func (t T) String() string
// The file is created in the same package and directory as the package that defines T.
// It has helpful defaults designed for use with go generate.
//
// Stringer works best with constants that are consecutive values such as created using iota,
// but creates good code regardless. In the future it might also provide custom support for
// constant sets that are bit patterns.
//
// For example, given this snippet,
//
//	package painkiller
//
//	type Pill int
//
//	const (
//		Placebo Pill = iota
//		Aspirin
//		Ibuprofen
//		Paracetamol
//		Acetaminophen = Paracetamol
//	)
//
// running this command
//
//	stringer -type=Pill
//
// in the same directory will create the file pill_string.go, in package painkiller,
// containing a definition of
//
//	func (Pill) String() string
//
// That method will translate the value of a Pill constant to the string representation
// of the respective constant name, so that the call fmt.Print(painkiller.Aspirin) will
// print the string "Aspirin".
//
// Typically this process would be run using go generate, like this:
//
//	//go:generate stringer -type=Pill
//
// If multiple constants have the same value, the lexically first matching name will
// be used (in the example, Acetaminophen will print as "Paracetamol").
//
// With no arguments, it processes the package in the current directory.
// Otherwise, the arguments must name a single directory holding a Go package
// or a set of Go source files that represent a single Go package.
//
// The -type flag accepts a comma-separated list of types so a single run can
// generate methods for multiple types. The default output file is t_string.go,
// where t is the lower-cased name of the first type listed. It can be overridden
// with the -output flag.
//
// The -linecomment flag tells stringer to generate the text of any line comment, trimmed
// of leading spaces, instead of the constant name. For instance, if the constants above had a
// Pill prefix, one could write
//
//	PillAspirin // Aspirin
//
// to suppress it in the output.
package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/constant"
	"go/format"
	"go/token"
	"go/types"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"unsafe"

	"golang.org/x/tools/go/packages"
)

////////////////////////////////////////////////////////////////////////////////
//
// TODO:
// 	1. Use the same Marshal* methods for all Run methods
// 	2. rename "stringer" => "go-enum" (or something)
//
////////////////////////////////////////////////////////////////////////////////

const generateMarshalers = true
const generateTests = true

var (
	typeNames   = flag.String("type", "", "comma-separated list of type names; must be set")
	output      = flag.String("output", "", "output file name; default srcdir/<type>_string.go")
	trimprefix  = flag.String("trimprefix", "", "trim the `prefix` from the generated constant names")
	linecomment = flag.Bool("linecomment", false, "use line comment text as printed text when present")
	sql         = flag.Bool("sql", false, "generate database/sql.Scanner database/sql/driver.Valuer methods")
	buildTags   = flag.String("tags", "", "comma-separated list of build tags to apply")
)

// Usage is a replacement usage function for the flags package.
func Usage() {
	fmt.Fprintf(os.Stderr, "Usage of stringer:\n")
	fmt.Fprintf(os.Stderr, "\tstringer [flags] -type T [directory]\n")
	fmt.Fprintf(os.Stderr, "\tstringer [flags] -type T files... # Must be a single package\n")
	fmt.Fprintf(os.Stderr, "For more information, see:\n")
	fmt.Fprintf(os.Stderr, "\thttp://godoc.org/golang.org/x/tools/cmd/stringer\n")
	fmt.Fprintf(os.Stderr, "Flags:\n")
	flag.PrintDefaults()
}

func main() {
	log.SetFlags(0)
	log.SetPrefix("stringer: ")
	flag.Usage = Usage
	flag.Parse()
	if len(*typeNames) == 0 {
		flag.Usage()
		os.Exit(2)
	}
	types := strings.Split(*typeNames, ",")
	var tags []string
	if len(*buildTags) > 0 {
		tags = strings.Split(*buildTags, ",")
	}

	// We accept either one directory or a list of files. Which do we have?
	args := flag.Args()
	if len(args) == 0 {
		// Default: process whole package in current directory.
		args = []string{"."}
	}

	// Parse the package once.
	var dir string
	g := Generator{
		trimPrefix:  *trimprefix,
		lineComment: *linecomment,
		sql:         *sql,
	}
	if g.sql && !generateMarshalers {
		panic("cannot generate SQL without Marshalers")
	}
	// TODO(suzmue): accept other patterns for packages (directories, list of files, import paths, etc).
	if len(args) == 1 && isDirectory(args[0]) {
		dir = args[0]
	} else {
		if len(tags) != 0 {
			log.Fatal("-tags option applies only to directories, not when files are specified")
		}
		dir = filepath.Dir(args[0])
	}

	g.parsePackage(args, tags)

	// Print the header and package clause.
	g.Printf("// Code generated by \"go-enum %s\"; DO NOT EDIT.\n", strings.Join(os.Args[1:], " "))
	g.Printf("\n")
	g.Printf("package %s", g.pkg.name)
	g.Printf("\n")
	if g.sql {
		g.Printf("import \"database/sql/driver\"\n") // Return value for Value() methods
	}
	if generateMarshalers {
		g.Printf("import \"errors\"\n") // Used by marshal/unmarshal methods.
	}
	if g.sql {
		g.Printf("import \"fmt\"\n") // Used by sql methods for errors.
	}
	g.Printf("import \"strconv\"\n") // Used by all methods.

	// Print the header for the test file
	g.TPrintf(testFileHeader, strings.Join(os.Args[1:], " "), g.pkg.name)

	// Run generate for each type.
	for _, typeName := range types {
		g.generate(typeName)
	}

	// Format the output.
	src := g.format()

	// Write to file.
	outputName := *output
	if outputName == "" {
		baseName := fmt.Sprintf("%s_string.go", types[0])
		outputName = filepath.Join(dir, strings.ToLower(baseName))
	}
	err := ioutil.WriteFile(outputName, src, 0644)
	if err != nil {
		log.Fatalf("writing output: %s", err)
	}

	if generateTests {
		outputName := strings.Replace(*output, ".go", "_test.go", 1)
		if outputName == "" {
			baseName := fmt.Sprintf("%s_string_test.go", types[0])
			outputName = filepath.Join(dir, strings.ToLower(baseName))
		}
		src := g.formatTest()
		err := ioutil.WriteFile(outputName, src, 0644)
		if err != nil {
			log.Fatalf("writing test output: %s", err)
		}
	}
}

const testFileHeader = `
// Code generated by "stringer %s"; DO NOT EDIT.

package %s

import (
	"encoding"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

`

// isDirectory reports whether the named file is a directory.
func isDirectory(name string) bool {
	info, err := os.Stat(name)
	if err != nil {
		log.Fatal(err)
	}
	return info.IsDir()
}

// Generator holds the state of the analysis. Primarily used to buffer
// the output for format.Source.
type Generator struct {
	buf  bytes.Buffer // Accumulated output.
	tbuf bytes.Buffer // Accumulated test output.
	pkg  *Package     // Package we are scanning.

	trimPrefix  string
	lineComment bool
	sql         bool
}

func (g *Generator) Printf(format string, args ...interface{}) {
	fmt.Fprintf(&g.buf, format, args...)
}

func (g *Generator) TPrintf(format string, args ...interface{}) {
	fmt.Fprintf(&g.tbuf, format, args...)
}

// File holds a single parsed file and associated data.
type File struct {
	pkg  *Package  // Package to which this file belongs.
	file *ast.File // Parsed AST.
	// These fields are reset for each type being generated.
	typeName string  // Name of the constant type.
	values   []Value // Accumulator for constant values of that type.

	trimPrefix  string
	lineComment bool
	sql         bool
}

type Package struct {
	name  string
	defs  map[*ast.Ident]types.Object
	files []*File
}

// parsePackage analyzes the single package constructed from the patterns and tags.
// parsePackage exits if there is an error.
func (g *Generator) parsePackage(patterns []string, tags []string) {
	cfg := &packages.Config{
		Mode: packages.LoadSyntax,
		// TODO: Need to think about constants in test files. Maybe write type_string_test.go
		// in a separate pass? For later.
		Tests:      false,
		BuildFlags: []string{fmt.Sprintf("-tags=%s", strings.Join(tags, " "))},
	}
	pkgs, err := packages.Load(cfg, patterns...)
	if err != nil {
		log.Fatal(err)
	}
	if len(pkgs) != 1 {
		log.Fatalf("error: %d packages found", len(pkgs))
	}
	g.addPackage(pkgs[0])
}

// addPackage adds a type checked Package and its syntax files to the generator.
func (g *Generator) addPackage(pkg *packages.Package) {
	g.pkg = &Package{
		name:  pkg.Name,
		defs:  pkg.TypesInfo.Defs,
		files: make([]*File, len(pkg.Syntax)),
	}

	for i, file := range pkg.Syntax {
		g.pkg.files[i] = &File{
			file:        file,
			pkg:         g.pkg,
			trimPrefix:  g.trimPrefix,
			lineComment: g.lineComment,
		}
	}
}

// writeConstantChecks generates code that will fail if the constants change value.
func (g *Generator) writeConstantChecks(typeName string, values []Value) {
	// If testing is enabled write to these checks to the test buffer,
	// otherwise we won't be able to achieve 100% test coverage.
	w := &g.buf
	if generateTests {
		w = &g.tbuf
	}
	// Generate code that will fail if the constants change value.
	fmt.Fprintf(w, "func _() {\n")
	fmt.Fprintf(w, "\t// An \"invalid array index\" compiler error signifies that the constant values have changed.\n")
	fmt.Fprintf(w, "\t// Re-run the stringer command to generate them again.\n")
	fmt.Fprintf(w, "\tvar x [1]struct{}\n")
	for _, v := range values {
		fmt.Fprintf(w, "\t_ = x[%s - %s]\n", v.originalName, v.str)
	}
	fmt.Fprintf(w, "}\n")
}

// generate produces the String method for the named type.
func (g *Generator) generate(typeName string) {
	values := make([]Value, 0, 100)
	for _, file := range g.pkg.files {
		// Set the state for this run of the walker.
		file.typeName = typeName
		file.values = nil
		if file.file != nil {
			ast.Inspect(file.file, file.genDecl)
			values = append(values, file.values...)
		}
	}

	if len(values) == 0 {
		log.Fatalf("no values defined for type %s", typeName)
	}
	if generateMarshalers {
		checkForDuplicateValues(typeName, values)
		checkForDuplicateStrings(typeName, values)
	}
	// Generate code that will fail if the constants change value.
	g.writeConstantChecks(typeName, values)

	runs := splitIntoRuns(values)
	// The decision of which pattern to use depends on the number of
	// runs in the numbers. If there's only one, it's easy. For more than
	// one, there's a tradeoff between complexity and size of the data
	// and code vs. the simplicity of a map. A map takes more space,
	// but so does the code. The decision here (crossover at 10) is
	// arbitrary, but considers that for large numbers of runs the cost
	// of the linear scan in the switch might become important, and
	// rather than use yet another algorithm such as binary search,
	// we punt and use a map. In any case, the likelihood of a map
	// being necessary for any realistic example other than bitmasks
	// is very low. And bitmasks probably deserve their own analysis,
	// to be done some other day.
	multipleRuns := false
	switch {
	case len(runs) == 1:
		g.buildOneRun(runs, typeName)
	case len(runs) <= 10:
		multipleRuns = true
		g.buildMultipleRuns(runs, typeName)
	default:
		g.buildMap(runs, typeName)
	}
	if generateMarshalers {
		g.buildUnmarshalers(runs, typeName, multipleRuns)
	}
	if generateTests {
		g.buildTests(runs, typeName)
	}
}

// checkForDuplicateValues checks for duplicate values which make generating
// marshal/unmarshal methods impossible.
func checkForDuplicateValues(typeName string, values []Value) {
	dupes := false
	seen := make(map[uint64][]string, len(values))
	for _, v := range values {
		seen[v.value] = append(seen[v.value], v.originalName)
		dupes = dupes || len(seen[v.value]) > 1
	}
	if !dupes {
		return
	}
	var buf bytes.Buffer
	for val, names := range seen {
		if len(names) == 1 {
			continue
		}
		if buf.Len() != 0 {
			buf.WriteString("; ")
		}
		fmt.Fprintf(&buf, "%s == %d", names, val)
	}
	log.Fatalf("cannot generate marshal/unmarshal methods for type: %s found "+
		"duplicate values: %s", typeName, &buf)
}

// checkForDuplicateStrings checks for values that have duplicate string forms
// which is possible with the -linecomment flag and makes generating
// marshal/unmarshal methods impossible.
func checkForDuplicateStrings(typeName string, values []Value) {
	dupes := false
	seen := make(map[string][]string, len(values))
	for _, v := range values {
		seen[v.name] = append(seen[v.name], v.originalName)
		dupes = dupes || len(seen[v.name]) > 1
	}
	if !dupes {
		return
	}
	var buf bytes.Buffer
	for name, origNames := range seen {
		if len(origNames) == 1 {
			continue
		}
		if buf.Len() != 0 {
			buf.WriteString("; ")
		}
		fmt.Fprintf(&buf, "%s == %s", origNames, name)
	}
	log.Fatalf("cannot generate marshal/unmarshal methods for type: %s found "+
		"values with duplicate strings representations: %s",
		typeName, &buf)
}

// splitIntoRuns breaks the values into runs of contiguous sequences.
// For example, given 1,2,3,5,6,7 it returns {1,2,3},{5,6,7}.
// The input slice is known to be non-empty.
func splitIntoRuns(values []Value) [][]Value {
	// We use stable sort so the lexically first name is chosen for equal elements.
	sort.Stable(byValue(values))
	// Remove duplicates. Stable sort has put the one we want to print first,
	// so use that one. The String method won't care about which named constant
	// was the argument, so the first name for the given value is the only one to keep.
	// We need to do this because identical values would cause the switch or map
	// to fail to compile.
	j := 1
	for i := 1; i < len(values); i++ {
		if values[i].value != values[i-1].value {
			values[j] = values[i]
			j++
		}
	}
	values = values[:j]
	runs := make([][]Value, 0, 10)
	for len(values) > 0 {
		// One contiguous sequence per outer loop.
		i := 1
		for i < len(values) && values[i].value == values[i-1].value+1 {
			i++
		}
		runs = append(runs, values[:i])
		values = values[i:]
	}
	return runs
}

// format returns the gofmt-ed contents of the Generator's buffer.
func (g *Generator) format() []byte {
	src, err := format.Source(g.buf.Bytes())
	if err != nil {
		// Should never happen, but can arise when developing this code.
		// The user can compile the output to see the error.
		log.Printf("warning: internal error: invalid Go generated: %s", err)
		log.Printf("warning: compile the package to analyze the error")
		return g.buf.Bytes()
	}
	return src
}

// formatTest returns the gofmt-ed contents of the Generator's test buffer.
func (g *Generator) formatTest() []byte {
	src, err := format.Source(g.tbuf.Bytes())
	if err != nil {
		log.Printf("warning: internal error: invalid Go generated (test files): %s", err)
		log.Printf("warning: compile the package to analyze the error")
		return g.tbuf.Bytes()
	}
	return src
}

// Value represents a declared constant.
type Value struct {
	originalName string // The name of the constant.
	name         string // The name with trimmed prefix.
	// The value is stored as a bit pattern alone. The boolean tells us
	// whether to interpret it as an int64 or a uint64; the only place
	// this matters is when sorting.
	// Much of the time the str field is all we need; it is printed
	// by Value.String.
	value  uint64          // Will be converted to int64 when needed.
	signed bool            // Whether the constant is a signed type.
	str    string          // The string representation given by the "go/constant" package.
	kind   types.BasicKind // Underlying type, used when generating tests
}

func (v *Value) String() string {
	return v.str
}

// byValue lets us sort the constants into increasing order.
// We take care in the Less method to sort in signed or unsigned order,
// as appropriate.
type byValue []Value

func (b byValue) Len() int      { return len(b) }
func (b byValue) Swap(i, j int) { b[i], b[j] = b[j], b[i] }
func (b byValue) Less(i, j int) bool {
	if b[i].signed {
		return int64(b[i].value) < int64(b[j].value)
	}
	return b[i].value < b[j].value
}

// genDecl processes one declaration clause.
func (f *File) genDecl(node ast.Node) bool {
	decl, ok := node.(*ast.GenDecl)
	if !ok || decl.Tok != token.CONST {
		// We only care about const declarations.
		return true
	}
	// The name of the type of the constants we are declaring.
	// Can change if this is a multi-element declaration.
	typ := ""
	// Loop over the elements of the declaration. Each element is a ValueSpec:
	// a list of names possibly followed by a type, possibly followed by values.
	// If the type and value are both missing, we carry down the type (and value,
	// but the "go/types" package takes care of that).
	for _, spec := range decl.Specs {
		vspec := spec.(*ast.ValueSpec) // Guaranteed to succeed as this is CONST.
		if vspec.Type == nil && len(vspec.Values) > 0 {
			// "X = 1". With no type but a value. If the constant is untyped,
			// skip this vspec and reset the remembered type.
			typ = ""

			// If this is a simple type conversion, remember the type.
			// We don't mind if this is actually a call; a qualified call won't
			// be matched (that will be SelectorExpr, not Ident), and only unusual
			// situations will result in a function call that appears to be
			// a type conversion.
			ce, ok := vspec.Values[0].(*ast.CallExpr)
			if !ok {
				continue
			}
			id, ok := ce.Fun.(*ast.Ident)
			if !ok {
				continue
			}
			typ = id.Name
		}
		if vspec.Type != nil {
			// "X T". We have a type. Remember it.
			ident, ok := vspec.Type.(*ast.Ident)
			if !ok {
				continue
			}
			typ = ident.Name
		}
		if typ != f.typeName {
			// This is not the type we're looking for.
			continue
		}
		// We now have a list of names (from one line of source code) all being
		// declared with the desired type.
		// Grab their names and actual values and store them in f.values.
		for _, name := range vspec.Names {
			if name.Name == "_" {
				continue
			}
			// This dance lets the type checker find the values for us. It's a
			// bit tricky: look up the object declared by the name, find its
			// types.Const, and extract its value.
			obj, ok := f.pkg.defs[name]
			if !ok {
				log.Fatalf("no value for constant %s", name)
			}
			info := obj.Type().Underlying().(*types.Basic).Info()
			if info&types.IsInteger == 0 {
				log.Fatalf("can't handle non-integer constant type %s", typ)
			}
			kind := obj.Type().Underlying().(*types.Basic).Kind()
			value := obj.(*types.Const).Val() // Guaranteed to succeed as this is CONST.
			if value.Kind() != constant.Int {
				log.Fatalf("can't happen: constant is not an integer %s", name)
			}
			i64, isInt := constant.Int64Val(value)
			u64, isUint := constant.Uint64Val(value)
			if !isInt && !isUint {
				log.Fatalf("internal error: value of %s is not an integer: %s", name, value.String())
			}
			if !isInt {
				u64 = uint64(i64)
			}
			v := Value{
				originalName: name.Name,
				value:        u64,
				signed:       info&types.IsUnsigned == 0,
				str:          value.String(),
				kind:         kind,
			}
			if c := vspec.Comment; f.lineComment && c != nil && len(c.List) == 1 {
				v.name = strings.TrimSpace(c.Text())
			} else {
				v.name = strings.TrimPrefix(v.originalName, f.trimPrefix)
			}
			f.values = append(f.values, v)
		}
	}
	return false
}

// Helpers

// usize returns the number of bits of the smallest unsigned integer
// type that will hold n. Used to create the smallest possible slice of
// integers to use as indexes into the concatenated strings.
func usize(n int) int {
	switch {
	case n < 1<<8:
		return 8
	case n < 1<<16:
		return 16
	default:
		// 2^32 is enough constants for anyone.
		return 32
	}
}

// declareIndexAndNameVars declares the index slices and concatenated names
// strings representing the runs of values.
func (g *Generator) declareIndexAndNameVars(runs [][]Value, typeName string) {
	var indexes, names []string
	for i, run := range runs {
		index, name := g.createIndexAndNameDecl(run, typeName, fmt.Sprintf("_%d", i))
		if len(run) != 1 {
			indexes = append(indexes, index)
		}
		names = append(names, name)
	}
	g.Printf("const (\n")
	for _, name := range names {
		g.Printf("\t%s\n", name)
	}
	g.Printf(")\n\n")

	if len(indexes) > 0 {
		g.Printf("var (")
		for _, index := range indexes {
			g.Printf("\t%s\n", index)
		}
		g.Printf(")\n\n")
	}
}

// declareIndexAndNameVar is the single-run version of declareIndexAndNameVars
func (g *Generator) declareIndexAndNameVar(run []Value, typeName string) {
	index, name := g.createIndexAndNameDecl(run, typeName, "")
	g.Printf("const %s\n", name)
	g.Printf("var %s\n", index)
}

// createIndexAndNameDecl returns the pair of declarations for the run. The caller will add "const" and "var".
func (g *Generator) createIndexAndNameDecl(run []Value, typeName string, suffix string) (string, string) {
	b := new(bytes.Buffer)
	indexes := make([]int, len(run))
	for i := range run {
		b.WriteString(run[i].name)
		indexes[i] = b.Len()
	}
	nameConst := fmt.Sprintf("_%s_name%s = %q", typeName, suffix, b.String())
	nameLen := b.Len()
	b.Reset()
	fmt.Fprintf(b, "_%s_index%s = [...]uint%d{0, ", typeName, suffix, usize(nameLen))
	for i, v := range indexes {
		if i > 0 {
			fmt.Fprintf(b, ", ")
		}
		fmt.Fprintf(b, "%d", v)
	}
	fmt.Fprintf(b, "}")
	return b.String(), nameConst
}

// declareNameVars declares the concatenated names string representing all the values in the runs.
func (g *Generator) declareNameVars(runs [][]Value, typeName string, suffix string) {
	g.Printf("const _%s_name%s = ", typeName, suffix)
	b := new(bytes.Buffer)
	for _, run := range runs {
		for i := range run {
			fmt.Fprintf(b, "%s", run[i].name)
		}
	}
	g.Printf("%q\n", b)
}

// buildOneRun generates the variables and String method for a single run of contiguous values.
func (g *Generator) buildOneRun(runs [][]Value, typeName string) {
	values := runs[0]
	g.Printf("\n")
	g.declareIndexAndNameVar(values, typeName)
	// The generated code is simple enough to write as a Printf format.
	lessThanZero := ""
	if values[0].signed {
		lessThanZero = "i < 0 || "
	}
	if values[0].value == 0 { // Signed or unsigned, 0 is still 0.
		g.Printf(stringOneRun, typeName, usize(len(values)), lessThanZero)
		if generateMarshalers {
			g.Printf(stringOneRunMarshal, typeName, usize(len(values)), lessThanZero)
		}
		if g.sql {
			g.Printf(stringOneRunSQL, typeName, usize(len(values)), lessThanZero)
		}
	} else {
		g.Printf(stringOneRunWithOffset, typeName, values[0].String(), usize(len(values)), lessThanZero)
		if generateMarshalers {
			g.Printf(stringOneRunWithOffsetMarshal, typeName, values[0].String(), usize(len(values)), lessThanZero)
		}
		if g.sql {
			g.Printf(stringOneRunWithOffsetSQL, typeName, values[0].String(), usize(len(values)), lessThanZero)
		}
	}
}

// Arguments to format are:
//	[1]: type name
//	[2]: size of index element (8 for uint8 etc.)
//	[3]: less than zero check (for signed types)
const stringOneRun = `func (i %[1]s) String() string {
	if %[3]si >= %[1]s(len(_%[1]s_index)-1) {
		return "%[1]s(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _%[1]s_name[_%[1]s_index[i]:_%[1]s_index[i+1]]
}
`

const stringOneRunMarshal = `
func (i %[1]s) Valid() bool {
	return !(%[3]si >= %[1]s(len(_%[1]s_index)-1))
}

func (i %[1]s) MarshalText() ([]byte, error) {
	if %[3]si >= %[1]s(len(_%[1]s_index)-1) {
		return nil, errors.New("invalid %[1]s: " + strconv.FormatInt(int64(i), 10))
	}
	return []byte(_%[1]s_name[_%[1]s_index[i]:_%[1]s_index[i+1]]), nil
}
`

const stringOneRunSQL = `
func (i %[1]s) Value() (driver.Value, error) {
	if %[3]si >= %[1]s(len(_%[1]s_index)-1) {
		return nil, errors.New("invalid %[1]s: " + strconv.FormatInt(int64(i), 10))
	}
	return _%[1]s_name[_%[1]s_index[i]:_%[1]s_index[i+1]], nil
}
`

// Arguments to format are:
//	[1]: type name
//	[2]: lowest defined value for type, as a string
//	[3]: size of index element (8 for uint8 etc.)
//	[4]: less than zero check (for signed types)
/*
 */
const stringOneRunWithOffset = `func (i %[1]s) String() string {
	i -= %[2]s
	if %[4]si >= %[1]s(len(_%[1]s_index)-1) {
		return "%[1]s(" + strconv.FormatInt(int64(i + %[2]s), 10) + ")"
	}
	return _%[1]s_name[_%[1]s_index[i] : _%[1]s_index[i+1]]
}
`

const stringOneRunWithOffsetMarshal = `
func (i %[1]s) Valid() bool {
	i -= %[2]s
	return !(%[4]si >= %[1]s(len(_%[1]s_index)-1))
}

func (i %[1]s) MarshalText() ([]byte, error) {
	i -= %[2]s
	if %[4]si >= %[1]s(len(_%[1]s_index)-1) {
		return nil, errors.New("invalid %[1]s: " + strconv.FormatInt(int64(i + %[2]s), 10))
	}
	return []byte(_%[1]s_name[_%[1]s_index[i]:_%[1]s_index[i+1]]), nil
}
`

const stringOneRunWithOffsetSQL = `
func (i %[1]s) Value() (driver.Value, error) {
	i -= %[2]s
	if %[4]si >= %[1]s(len(_%[1]s_index)-1) {
		return nil, errors.New("invalid %[1]s: " + strconv.FormatInt(int64(i + %[2]s), 10))
	}
	return _%[1]s_name[_%[1]s_index[i]:_%[1]s_index[i+1]], nil
}
`

// buildMultipleRuns generates the variables and String method for multiple runs of contiguous values.
// For this pattern, a single Printf format won't do.
func (g *Generator) buildMultipleRuns(runs [][]Value, typeName string) {
	g.Printf("\n")
	g.declareIndexAndNameVars(runs, typeName)
	g.Printf("func (i %s) String() string {\n", typeName)
	g.Printf("\tswitch {\n")
	for i, values := range runs {
		if len(values) == 1 {
			g.Printf("\tcase i == %s:\n", &values[0])
			g.Printf("\t\treturn _%s_name_%d\n", typeName, i)
			continue
		}
		if values[0].value == 0 && !values[0].signed {
			// For an unsigned lower bound of 0, "0 <= i" would be redundant.
			g.Printf("\tcase i <= %s:\n", &values[len(values)-1])
		} else {
			g.Printf("\tcase %s <= i && i <= %s:\n", &values[0], &values[len(values)-1])
		}
		if values[0].value != 0 {
			g.Printf("\t\ti -= %s\n", &values[0])
		}
		g.Printf("\t\treturn _%s_name_%d[_%s_index_%d[i]:_%s_index_%d[i+1]]\n",
			typeName, i, typeName, i, typeName, i)
	}
	g.Printf("\tdefault:\n")
	g.Printf("\t\treturn \"%s(\" + strconv.FormatInt(int64(i), 10) + \")\"\n", typeName)
	g.Printf("\t}\n")
	g.Printf("}\n")

	if generateMarshalers {
		g.multipleRunsValid(runs, typeName)
	}
}

func (g *Generator) multipleRunsValid(runs [][]Value, typeName string) {
	g.Printf("\n")
	g.Printf("func (i %s) Valid() bool {\n", typeName)
	g.Printf("\tswitch {\n")
	for _, values := range runs {
		if len(values) == 1 {
			g.Printf("\tcase i == %s:\n", &values[0])
			continue
		}
		if values[0].value == 0 && !values[0].signed {
			// For an unsigned lower bound of 0, "0 <= i" would be redundant.
			g.Printf("\tcase i <= %s:\n", &values[len(values)-1])
		} else {
			g.Printf("\tcase %s <= i && i <= %s:\n", &values[0], &values[len(values)-1])
		}
	}
	g.Printf("\tdefault:\n")
	g.Printf("\t\treturn false\n")
	g.Printf("\t}\n")
	g.Printf("\treturn true\n")
	g.Printf("}\n")

	g.Printf(stringMultipleRunsMarshal, typeName)
	if g.sql {
		g.Printf(stringMultipleRunsSQL, typeName)
	}
}

const stringMultipleRunsMarshal = `
func (i %[1]s) MarshalText() ([]byte, error) {
	if i.Valid() {
		return []byte(i.String()), nil
	}
	return nil, errors.New("invalid %[1]s: " + strconv.FormatInt(int64(i), 10))
}
`

const stringMultipleRunsSQL = `
func (i %[1]s) Value() (driver.Value, error) {
	if i.Valid() {
		return i.String(), nil
	}
	return nil, errors.New("invalid %[1]s: " + strconv.FormatInt(int64(i), 10))
}
`

// buildMap handles the case where the space is so sparse a map is a reasonable fallback.
// It's a rare situation but has simple code.
func (g *Generator) buildMap(runs [][]Value, typeName string) {
	g.Printf("\n")
	g.declareNameVars(runs, typeName, "")
	g.Printf("\nvar _%s_map = map[%s]string{\n", typeName, typeName)
	n := 0
	for _, values := range runs {
		for _, value := range values {
			g.Printf("\t%s: _%s_name[%d:%d],\n", &value, typeName, n, n+len(value.name))
			n += len(value.name)
		}
	}
	g.Printf("}\n\n")
	g.Printf(stringMap, typeName)
	if generateMarshalers {
		g.Printf(stringMapMarhalers, typeName)
	}
	if g.sql {
		g.Printf(stringMapSQL, typeName)
	}
}

// Argument to format is the type name.
const stringMap = `func (i %[1]s) String() string {
	if str, ok := _%[1]s_map[i]; ok {
		return str
	}
	return "%[1]s(" + strconv.FormatInt(int64(i), 10) + ")"
}
`

const stringMapMarhalers = `
func (i %[1]s) Valid() bool {
	_, ok := _%[1]s_map[i]
	return ok
}

func (i %[1]s) MarshalText() ([]byte, error) {
	if str, ok := _%[1]s_map[i]; ok {
		return []byte(str), nil
	}
	return nil, errors.New("invalid %[1]s: " + strconv.FormatInt(int64(i), 10))
}
`

const stringMapSQL = `
func (i %[1]s) Value() (driver.Value, error) {
	if str, ok := _%[1]s_map[i]; ok {
		return str, nil
	}
	return nil, errors.New("invalid %[1]s: " + strconv.FormatInt(int64(i), 10))
}
`

const genericScanSQL = `
func (i *%[1]s) Scan(src interface{}) error {
	switch s := src.(type) {
	case string:
		return i.Set(s)
	case []byte:
		return i.UnmarshalText(s)
	default:
		return fmt.Errorf("cannot scan type %%T into %[1]s", src)
	}
}
`

func countValues(runs [][]Value) int {
	n := 0
	for _, values := range runs {
		n += len(values)
	}
	return n
}

func (g *Generator) buildUnmarshalers(runs [][]Value, typeName string, multipleRuns bool) {
	count := countValues(runs)
	if count == 0 {
		log.Fatalf("no values defined for type %s", typeName)
	}
	// Use a map when there are more than 32 values. A switch is slightly faster
	// with 64 values but adds a lot of code for a marginal gain.
	//
	// See: internal/bench_lookup/bench_lookup_test.go for benchmark results.
	if count <= 32 {
		g.buildUnmarshalersSwitch(runs, typeName, multipleRuns)
	} else {
		g.buildUnmarshalersMap(runs, typeName, multipleRuns)
	}
}

func (g *Generator) buildUnmarshalersSwitch(runs [][]Value, typeName string, multipleRuns bool) {
	const errFormat = `
		if len(s) <= 32 {
			err = errors.New("malformed %[1]s: " + string(s))
		} else {
			err = errors.New("malformed %[1]s: " + string(s[0:29]) + "...")
		}
`

	marshalers := []struct {
		funcName, switchVal string
	}{
		{"Set(s string)", "s"},
		{"UnmarshalText(s []byte)", "string(s)"},
	}
	for _, m := range marshalers {
		g.Printf("\nfunc (i *%s) %s (err error) {\n", typeName, m.funcName)
		g.Printf("\tswitch %s {\n", m.switchVal)

		if multipleRuns {
			for i, values := range runs {
				if len(values) == 1 {
					g.Printf("\tcase _%s_name_%d:\n", typeName, i)
					g.Printf("\t\t*i = %s\n", values[0].originalName)
					continue
				}
				n := 0
				for _, value := range values {
					g.Printf("\tcase _%s_name_%d[%d:%d]:\n", typeName, i, n, n+len(value.name))
					g.Printf("\t\t*i = %s\n", value.originalName)
					n += len(value.name)
				}
			}
		} else {
			n := 0
			for _, values := range runs {
				// TODO: avoid index on single values (use Prime test)
				for _, value := range values {
					g.Printf("\tcase _%s_name[%d:%d]:\n", typeName, n, n+len(value.name))
					g.Printf("\t\t*i = %s\n", value.originalName)
					n += len(value.name)
				}
			}
		}

		g.Printf("\tdefault :\n")
		g.Printf(errFormat[1:] /* remove leading newline */, typeName)
		g.Printf("\t}\n")
		g.Printf("\treturn err\n")
		g.Printf("}\n\n")
	}
	g.Printf("\n")
	if g.sql {
		g.Printf(genericScanSQL, typeName)
		g.Printf("\n")
	}
}

func (g *Generator) buildUnmarshalersMap(runs [][]Value, typeName string, multipleRuns bool) {
	g.Printf("\nvar _%s_lookup_map = map[string]%s{\n", typeName, typeName)
	if multipleRuns {
		for i, values := range runs {
			n := 0
			for _, value := range values {
				g.Printf("\t_%s_name_%d[%d:%d]: %s,\n", typeName, i, n, n+len(value.name), &value)
				n += len(value.name)
			}
		}
	} else {
		n := 0
		for _, values := range runs {
			for _, value := range values {
				g.Printf("\t_%s_name[%d:%d]: %s,\n", typeName, n, n+len(value.name), &value)
				n += len(value.name)
			}
		}
	}
	g.Printf("}\n\n")
	g.Printf(stringMapUnmarshalers, typeName)
	g.Printf("\n")
	if g.sql {
		g.Printf(genericScanSQL, typeName)
		g.Printf("\n")
	}
}

// TODO: consider renaming
const stringMapUnmarshalers = `
func (i *%[1]s) Set(s string) error {
	if v, ok := _%[1]s_lookup_map[s]; ok {
		*i = v
		return nil
	}
	if len(s) <= 32 {
		return errors.New("malformed %[1]s: " + s)
	}
	return errors.New("malformed %[1]s: " + s[0:29] + "...")
}

func (i *%[1]s) UnmarshalText(s []byte) error {
	if v, ok := _%[1]s_lookup_map[string(s)]; ok {
		*i = v
		return nil
	}
	if len(s) <= 32 {
		return errors.New("malformed %[1]s: " + string(s))
	}
	return errors.New("malformed %[1]s: " + string(s[0:29]) + "...")
}
`

func typeMinMax(typeName string, kind types.BasicKind) (min, max uint64) {
	// use u to defeat the compiler's overflow check
	u := func(i int64) uint64 {
		return uint64(i)
	}
	switch kind {
	case types.Int:
		if unsafe.Sizeof(int(0)) == 8 {
			return u(math.MinInt64), math.MaxInt64
		} else {
			return u(math.MinInt32), math.MaxInt32
		}
	case types.Int8:
		return u(math.MinInt8), math.MaxInt8
	case types.Int16:
		return u(math.MinInt16), math.MaxInt16
	case types.Int32:
		return u(math.MinInt32), math.MaxInt32
	case types.Int64:
		return u(math.MinInt64), math.MaxInt64
	case types.Uint8:
		return 0, math.MaxUint8
	case types.Uint16:
		return 0, math.MaxUint16
	case types.Uint32:
		return 0, math.MaxUint32
	case types.Uint64:
		return 0, math.MaxUint64
	case types.Uint:
		fallthrough
	case types.Uintptr:
		if unsafe.Sizeof(uint(0)) == 8 {
			return 0, math.MaxUint64
		} else {
			return 0, math.MaxUint32
		}
	default:
		log.Fatalf("invalid kind: %d for type: %s", kind, typeName)
		panic("unreachable")
	}
}

func (g *Generator) buildInvalidValues(runs [][]Value, typeName string) map[uint64]Value {
	if len(runs) == 0 {
		log.Fatalf("no values defined for type %s", typeName)
	}

	values := make(map[uint64]bool)
	for _, run := range runs {
		for _, v := range run {
			values[v.value] = true
		}
	}

	invalid := make(map[uint64]Value)
	first := runs[0][0]
	signed := first.signed
	min, max := typeMinMax(typeName, first.kind)
	if !values[min] {
		invalid[min] = Value{signed: signed, value: min}
	}
	if !values[0] {
		invalid[0] = Value{signed: signed, value: 0}
	}
	negOne := -1 // work around the compilers overflow check
	if signed && !values[uint64(negOne)] {
		invalid[uint64(negOne)] = Value{signed: signed, value: uint64(negOne)}
	}
	if !values[max] {
		invalid[max] = Value{signed: signed, value: max}
	}

	// TODO: a lot of the uint64 conversions can probably be removed

	for _, run := range runs {
		if len(runs) == 0 {
			continue // can this happen?
		}
		first := run[0]
		last := run[len(run)-1]
		if signed {
			if int64(first.value) > int64(min) {
				u := uint64(int64(first.value) - 1)
				if !values[u] {
					invalid[u] = Value{signed: signed, value: u}
				}
			}
			// TODO: this is probably redundant
			if int64(last.value) < int64(max) {
				u := uint64(int64(last.value) + 1)
				if !values[u] {
					invalid[u] = Value{signed: signed, value: u}
				}
			}
		} else {
			if first.value > 0 {
				u := first.value - 1
				if !values[u] {
					invalid[u] = Value{signed: signed, value: u}
				}
			}
			if last.value < max {
				u := last.value + 1
				if !values[u] {
					invalid[u] = Value{signed: signed, value: u}
				}
			}
		}
	}

	for u, v := range invalid {
		if signed {
			v.str = strconv.FormatInt(int64(v.value), 10)
		} else {
			v.str = strconv.FormatUint(v.value, 10)
		}
		invalid[u] = v
	}
	return invalid
}

func (g *Generator) buildTests(runs [][]Value, typeName string) {
	invalid := g.buildInvalidValues(runs, typeName)

	values := make([]Value, 0, 100+len(invalid))
	for _, run := range runs {
		values = append(values, run...)
	}
	for _, v := range invalid {
		values = append(values, v)
	}

	sort.Stable(byValue(values))

	var buf bytes.Buffer
	for _, v := range values {
		if _, ok := invalid[v.value]; ok {
			fmt.Fprintf(&buf, "\t\t{%[1]s(%[2]s), \"%[1]s(%[3]d)\", false},\n",
				typeName, v.str, int64(v.value))
		} else {
			fmt.Fprintf(&buf, "\t\t{%s, %q, true},\n", v.originalName, v.name)
		}
	}

	if g.sql {
		g.TPrintf(testTemplate, typeName, buf.String(), testTemplateSQL)
	} else {
		g.TPrintf(testTemplate, typeName, buf.String(), "", "")
	}
	g.TPrintf("\n")

	buf.Reset()
	for _, run := range runs {
		for _, v := range run {
			fmt.Fprintf(&buf, "\t\t{%[1]s, %[2]q, []byte(%[2]q)},\n", v.originalName, v.name)
		}
	}
	if g.sql {
		g.TPrintf(benchmarkTemplate, typeName, buf.String(), benchmarkTemplateSQL)
	} else {
		g.TPrintf(benchmarkTemplate, typeName, buf.String(), "")
	}
	g.TPrintf("\n")
}

// Arguments to format are:
//	[1]: type name
//	[2]: values to test
const testTemplate = `
var (
	_ fmt.Stringer             = %[1]s(0)
	_ encoding.TextMarshaler   = %[1]s(0)
	_ encoding.TextUnmarshaler = (*%[1]s)(nil)
	_ func() bool              = %[1]s(0).Valid    // Valid()
	_ func(string) error       = (*%[1]s)(nil).Set // Set()
)

func TestGeneratedEnum_%[1]s(t *testing.T) {
	const _TypeName = "%[1]s"
	var tests = []struct {
		Val   %[1]s
		Str   string
		Valid bool
	}{
%[2]s
	}

	testUnmarshalError := func(t *testing.T, err error, s string) {
		t.Helper()
		if err == nil {
			t.Error("expected a non-nil error")
			return
		}
		var exp string
		if len(s) <= 32 {
			exp = fmt.Sprintf("malformed %%s: %%s", _TypeName, s)
		} else {
			exp = fmt.Sprintf("malformed %%s: %%s...", _TypeName, string(s[0:29]))
		}
		if err.Error() != exp {
			t.Errorf("unmarshal error: got: %%s want: %%s", err.Error(), exp)
		}
	}

	t.Run("Valid", func(t *testing.T) {
		for _, x := range tests {
			if x.Val.Valid() != x.Valid {
				t.Errorf("%%+v: got: %%t want: %%t", x, x.Val.Valid(), x.Valid)
			}
		}
	})

	t.Run("String", func(t *testing.T) {
		for _, x := range tests {
			str := x.Val.String()
			if str != x.Str {
				t.Errorf("%%+v: got: %%q want: %%q", x, str, x.Str)
			}
		}
	})

	t.Run("Set", func(t *testing.T) {
		var zeroValue %[1]s
		for _, x := range tests {
			var v %[1]s
			err := v.Set(x.Str)
			if (err == nil) != x.Valid {
				t.Errorf("%%+v: %%v", x, err)
				continue
			}
			exp := x.Val
			if !x.Valid {
				testUnmarshalError(t, err, x.Str)
				exp = zeroValue
			}
			if v != exp {
				t.Errorf("%%+v: got: %%s want: %%s", x, v, exp)
			}
		}

		// Test that we don't include long strings in the error message
		var v %[1]s
		invalid := strings.Repeat("a", 256) + "\x00" // this should not collide
		testUnmarshalError(t, v.Set(invalid), invalid)
	})

	t.Run("MarshalJSON", func(t *testing.T) {
		for _, x := range tests {
			data, err := json.Marshal(x.Val)
			if (err == nil) != x.Valid {
				t.Errorf("%%+v: %%v", x, err)
				continue
			}
			if !x.Valid {
				if data != nil {
					t.Errorf("%%+v: expected []byte(nil) on error got: %%v", x, data)
				}
				merr, ok := err.(*json.MarshalerError)
				if !ok {
					t.Errorf("%%+v: invalid error type: %%T", x, err)
				}
				exp := fmt.Sprintf("invalid %%s: %%d", _TypeName, int64(x.Val))
				if merr.Err.Error() != exp {
					t.Errorf("%%+v: got: %%s want: %%s", x, merr.Err.Error(), exp)
				}
				continue
			}

			exp, err := json.Marshal(x.Val.String())
			if err != nil {
				t.Errorf("%%+v: %%v", x, err)
				continue
			}
			if string(data) != string(exp) {
				t.Errorf("%%+v: got: '%%s' want: '%%s'", x, data, exp)
			}
			var v %[1]s
			if err := json.Unmarshal(data, &v); err != nil {
				t.Errorf("%%+v: %%v", x, err)
				continue
			}
			if v != x.Val {
				t.Errorf("%%+v: got: %%s want: %%s", x, v, x.Val)
			}
		}
	})

	t.Run("UnmarshalJSON", func(t *testing.T) {
		var zeroValue %[1]s
		for _, x := range tests {
			data, err := json.Marshal(x.Val)
			if (err == nil) != x.Valid {
				t.Errorf("%%+v: %%v", x, err)
				continue
			}
			if !x.Valid {
				// Set the data to the string value, which is also invalid.
				data, err = json.Marshal(x.Str)
				if err != nil {
					t.Fatalf("%%+v: %%v", x, err)
				}
			}

			var v %[1]s
			err = json.Unmarshal(data, &v)
			if (err == nil) != x.Valid {
				t.Errorf("%%+v: %%v", x, err)
				continue
			}
			if !x.Valid {
				b := data
				if len(b) > 1 && b[0] == '"' {
					b = b[1 : len(b)-1]
				}
				testUnmarshalError(t, err, string(b))
			}
			exp := x.Val
			if !x.Valid {
				exp = zeroValue
			}
			if v != exp {
				t.Errorf("%%+v: got: %%s want: %%s", x, v, exp)
			}
		}

		// Test that we don't include long strings in the error message
		var v %[1]s
		invalid := strings.Repeat("a", 256) + "\\u0000" // this should not collide
		err := json.Unmarshal([]byte("\""+invalid+"\""), &v)
		testUnmarshalError(t, err, invalid)
	})

	t.Run("MarshalText", func(t *testing.T) {
		for _, x := range tests {
			data, err := x.Val.MarshalText()
			if (err == nil) != x.Valid {
				t.Errorf("%%+v: %%v", x, err)
				continue
			}
			if !x.Valid {
				if data != nil {
					t.Errorf("%%+v: expected []byte(nil) on error got: %%v", x, data)
				}
				exp := fmt.Sprintf("invalid %%s: %%d", _TypeName, int64(x.Val))
				if err.Error() != exp {
					t.Errorf("%%+v: got: %%s want: %%s", x, err.Error(), exp)
				}
				continue
			}

			if string(data) != x.Val.String() {
				t.Errorf("%%+v: got: '%%s' want: '%%s'", x, data, x.Val.String())
			}
			var v %[1]s
			if err := v.UnmarshalText(data); err != nil {
				t.Errorf("%%+v: %%v", x, err)
				continue
			}
			if v != x.Val {
				t.Errorf("%%+v: got: %%s want: %%s", x, v, x.Val)
			}
		}
	})
	t.Run("UnmarshalText", func(t *testing.T) {
		var zeroValue %[1]s
		for _, x := range tests {
			data, err := x.Val.MarshalText()
			if (err == nil) != x.Valid {
				t.Errorf("%%+v: %%v", x, err)
				continue
			}
			if !x.Valid {
				// Set the data to the string value, which is also invalid.
				data = []byte(x.Str)
			}

			var v %[1]s
			err = v.UnmarshalText(data)
			if (err == nil) != x.Valid {
				t.Errorf("%%+v: %%v", x, err)
				continue
			}
			if !x.Valid {
				testUnmarshalError(t, err, string(data))
			}
			exp := x.Val
			if !x.Valid {
				exp = zeroValue
			}
			if v != exp {
				t.Errorf("%%+v: got: %%s want: %%s", x, v, exp)
			}
		}

		// invalid values
		for _, data := range [][]byte{nil, {}} {
			var v %[1]s
			if err := v.UnmarshalText(data); err == nil {
				t.Errorf("expected an error unmarshaling: %%v: %%v", data, err)
			}
		}

		// Test that we don't include long strings in the error message
		var v %[1]s
		invalid := strings.Repeat("a", 256) + "\x00" // this should not collide
		testUnmarshalError(t, v.UnmarshalText([]byte(invalid)), invalid)
	})
%[3]s
}
`

const testTemplateSQL = `
	t.Run("Value", func(t *testing.T) {
		for _, x := range tests {
			value, err := x.Val.Value()
			if (err == nil) != x.Valid {
				t.Errorf("%%+v: %%v", x, err)
				continue
			}
			if !x.Valid {
				if value != nil {
					t.Errorf("%%+v: expected nil on error got: %%v", x, value)
				}
				exp := fmt.Sprintf("invalid %%s: %%d", _TypeName, int64(x.Val))
				if err.Error() != exp {
					t.Errorf("%%+v: got: %%s want: %%s", x, err.Error(), exp)
				}
				continue
			}

			if value.(string) != x.Val.String() {
				t.Errorf("%%+v: got: '%%s' want: '%%s'", x, value, x.Val.String())
			}
			var v %[1]s
			if err := v.Scan(value); err != nil {
				t.Errorf("%%+v: %%v", x, err)
				continue
			}
			if v != x.Val {
				t.Errorf("%%+v: got: %%s want: %%s", x, v, x.Val)
			}
		}
	})
	t.Run("Scan", func(t *testing.T) {
		var zeroValue %[1]s
		for _, x := range tests {
			value, err := x.Val.Value()
			if (err == nil) != x.Valid {
				t.Errorf("%%+v: %%v", x, err)
				continue
			}
			if !x.Valid {
				// Set the value to the string value, which is also invalid.
				value = x.Str
			}

			var v %[1]s
			err = v.Scan(value)
			if (err == nil) != x.Valid {
				t.Errorf("%%+v: %%v", x, err)
				continue
			}
			if !x.Valid {
				testUnmarshalError(t, err, value.(string))
			}
			exp := x.Val
			if !x.Valid {
				exp = zeroValue
			}
			if v != exp {
				t.Errorf("%%+v: got: %%s want: %%s", x, v, exp)
			}
		}

		// invalid values
		for _, data := range []interface{}{nil, []byte{}, 123} {
			var v %[1]s
			if err := v.Scan(data); err == nil {
				t.Errorf("expected an error unmarshaling: %%v: %%v", data, err)
			}
		}

		// Test that we don't include long strings in the error message
		var v %[1]s
		invalid := strings.Repeat("a", 256) + "\x00" // this should not collide
		testUnmarshalError(t, v.Scan([]byte(invalid)), invalid)
	})
`

// Arguments to format are:
//	[1]: type name
//	[2]: valid values to benchmark with
const benchmarkTemplate = `
func BenchmarkGeneratedEnum_%[1]s(b *testing.B) {
	var tests = [...]struct {
		Val   %[1]s
		Str   string
		Bytes []byte
	}{
%[2]s
	}
	b.ResetTimer()
	b.Run("Valid", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			t := tests[i%%len(tests)]
			t.Val.Valid()
		}
	})
	b.Run("String", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			t := tests[i%%len(tests)]
			t.Val.String()
		}
	})
	b.Run("Set", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			t := tests[i%%len(tests)]
			t.Val.Set(t.Str)
		}
	})
	b.Run("MarshalText", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			t := tests[i%%len(tests)]
			t.Val.MarshalText()
		}
	})
	b.Run("UnmarshalText", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			t := tests[i%%len(tests)]
			t.Val.UnmarshalText(t.Bytes)
		}
	})
%[3]s
}
`

const benchmarkTemplateSQL = `
	b.Run("Value", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			t := tests[i%%len(tests)]
			t.Val.Value()
		}
	})
	b.Run("Scan", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			t := tests[i%%len(tests)]
			t.Val.Scan(t.Bytes)
		}
	})
`
