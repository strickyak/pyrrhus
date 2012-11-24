package main

import "fmt"
import "go/ast"
// import "go/build"
import "go/parser"
import "go/token"
import "os"
import "strconv"
import "strings"

var _ = os.Open
var _ = ast.Walk
var _ = strconv.Atoi
var _ = strings.Contains

/*  [go/ast]
type Importer func(imports map[string]*Object, path string) (pkg *Object, err error)

An Importer resolves import paths to package Objects. The imports map records the packages already imported, indexed by package id (canonical import path). An Importer must determine the canonical import path and check the map to see if it is already present in the imports map. If so, the Importer can return the map entry. Otherwise, the Importer should load the package data for the given path into a new *Object (pkg), record pkg in the imports map, and then return pkg.

func NewPackage(fset *token.FileSet, files map[string]*File, importer Importer, universe *Scope) (*Package, error)

NewPackage creates a new Package node from a set of File nodes. It resolves unresolved identifiers across files and updates each file's Unresolved list accordingly. If a non-nil importer and universe scope are provided, they are used to resolve identifiers not declared in any of the package files. Any remaining unresolved identifiers are reported as undeclared. If the files belong to different packages, one package name is selected and files with different package names are reported and then ignored. The result is a package node and a scanner.ErrorList if there were errors.

[go/build]
type Package struct {
    Dir        string // directory containing package sources
    Name       string // package name
    Doc        string // documentation synopsis
    ImportPath string // import path of package ("" if unknown)
    Root       string // root of Go tree where this package lives
    SrcRoot    string // package source root directory ("" if unknown)
    PkgRoot    string // package install root directory ("" if unknown)
    BinDir     string // command install directory ("" if unknown)
    Goroot     bool   // package found in Go root
    PkgObj     string // installed .a file

    // Source files
    GoFiles   []string // .go source files (excluding CgoFiles, TestGoFiles, XTestGoFiles)
    CgoFiles  []string // .go source files that import "C"
    CFiles    []string // .c source files
    HFiles    []string // .h source files
    SFiles    []string // .s source files
    SysoFiles []string // .syso system object files to add to archive

    // Cgo directives
    CgoPkgConfig []string // Cgo pkg-config directives
    CgoCFLAGS    []string // Cgo CFLAGS directives
    CgoLDFLAGS   []string // Cgo LDFLAGS directives

    // Dependency information
    Imports   []string                    // imports from GoFiles, CgoFiles
    ImportPos map[string][]token.Position // line information for Imports

    // Test information
    TestGoFiles    []string                    // _test.go files in package
    TestImports    []string                    // imports from TestGoFiles
    TestImportPos  map[string][]token.Position // line information for TestImports
    XTestGoFiles   []string                    // _test.go files outside package
    XTestImports   []string                    // imports from XTestGoFiles
    XTestImportPos map[string][]token.Position // line information for XTestImports
}

[parser]

func ParseDir(fset *token.FileSet, path string, filter func(os.FileInfo) bool, mode Mode) (pkgs map[string]*ast.Package, first error)

ParseDir calls ParseFile for the files in the directory specified by path and returns a map of package name -> package AST with all the packages found. If filter != nil, only the files with os.FileInfo entries passing through the filter are considered. The mode bits are passed to ParseFile unchanged. Position information is recorded in the file set fset.

If the directory couldn't be read, a nil map and the respective error are returned. If a parse error occurred, a non-nil but incomplete map and the first error encountered are returned.


*/

//func DumpBuildInfo() {
//  fmt.Printf("Default Build Context: %#v\n", build.Default)
//
//  p, err := build.Import("fmt", "", build.AllowBinary)
//  fmt.Printf("err %#v\n", err);
//  fmt.Printf("p %#v\n", p);
//  
//  fset := token.NewFileSet()
//  //// universe := ast.NewScope(nil)
//  //// files := make(map[string]*ast.File)
//  //// // NewPackage(fset *token.FileSet, files map[string]*File, importer Importer, universe *Scope) (*Package, error)
//  //// pack, err := ast.NewPackage(fset, files, nil, universe)
//  //// fmt.Printf("err %#v\n", err);
//  //// fmt.Printf("pack %#v\n", pack);
//
//  pkgs, err := parser.ParseDir(fset, "/opt/go/src/pkg/fmt", nil, parser.Mode(0))
//  fmt.Printf("err %#v\n", err);
//  fmt.Printf("pkgs %#v\n", pkgs);
//}

var  fset = token.NewFileSet()

func GrokDir(dir string) {
  pkgs, err := parser.ParseDir(fset, dir, nil, parser.Mode(0))
  if err != nil {
    panic(fmt.Sprintf("ERROR <%q> IN DIR <%s>", err, dir))
  }
  for pk, pv := range pkgs {
    if strings.Contains(pk, "_test") { continue }
    fmt.Printf("pk %#v\n", pk);
    fmt.Printf("pv %#v\n", pv);

    for fk, fv := range pv.Files {
      fmt.Printf("fk %#v\n", fk);
      // fmt.Printf("fv %#v\n", fv);
      for i, dcl := range fv.Decls {

	switch x := dcl.(type) {
	  case (*ast.FuncDecl):
            fmt.Printf("FUNC #%d == %#v\n", i, x);
            fmt.Printf("         Type = %s\n", ast.Print(fset, x.Type))

            // retrieve the parameter's type
	    params := x.Type.Params
	    for lk, lv := range params.List {
	      tname := "?"
	      switch t := lv.Type.(type) {
	      	case (*ast.Ident):
		  tname = t.Name
	      }

	      fmt.Printf("       %d = %s (%s)\n", lk, lv.Names[0].Name, tname)
	    }
	  default:
            fmt.Printf("DECL #%d == %#v\n", i, dcl);
	}
      }
    }
  }
}

func typeStr(a interface{}) string {
  tname := "?"

  switch t := a.(type) {
  case (*ast.Ident):
    tname = t.Name
  case (*ast.ArrayType):
    tname = "[]" + typeStr(t.Elt)
  case (*ast.StarExpr):
    tname = "*" + typeStr(t.X)
  }

  return tname
}

func main() {
  for _, dir := range os.Args[1:] {
    GrokDir(dir)
  }
}
