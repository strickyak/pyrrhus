package main

import "encoding/json"
import "fmt"
import "go/ast"
import "go/parser"
import "go/token"
import "os"
import "strconv"
import "strings"

var _ = json.NewEncoder
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

func Str(x interface{}) string {
  return fmt.Sprintf("%v", x)
}

func FilterDotGo(info os.FileInfo) bool {
  s := info.Name()
  n := len(s)
  // return strings.Contains(s, ".go")
  return n>3 && s[n-3]=='.' && s[n-2]=='g' && s[n-1]=='o'
}

func GrokDir(dir string) {
  pkgs, err := parser.ParseDir(fset, dir, FilterDotGo, parser.Mode(0))
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
            if !ast.IsExported(x.Name.Name) {
	      continue
	    }
            fstr := funcDeclStr(x)
            fmt.Printf("@@ %s\n", fstr)

            if strings.Contains(fstr, "?") {
              fmt.Printf("FUNC #%d == %#v\n", i, dcl);
              fmt.Printf("   Recv: %#v\n", x.Recv)
	      if x.Recv != nil {
                for _, elem := range x.Recv.List {
                  fmt.Printf("      Elem: %#v\n", elem)
		  for _, rid := range elem.Names {
                    fmt.Printf("      Name: %#v\n", rid.Name)
		  }
                  fmt.Printf("      Type: %#v\n", elem.Type)
                  fmt.Printf("      ====: %s\n", typeStr(elem.Type))
                  fmt.Printf("      ====: %s\n", ast.Print(fset, elem.Type))
                }
              }

              fmt.Printf("  FUNC Type = %s\n", ast.Print(fset, x.Type))
            }
	  default:
            fmt.Printf("DECL #%d == %#v\n", i, dcl);
	}
      }
    }
  }
}

func typeStr(a interface{}) string {
  switch t := a.(type) {
  case (*ast.Ident):
    return "'" + t.Name + " "
  case (*ast.ArrayType):
    return "{ ARRAY " + typeStr(t.Elt) + " } "
  case (*ast.StarExpr):
    return "{ STAR " + typeStr(t.X) + " } "
  case (*ast.Ellipsis):
    return "{ ELLIPSIS " + typeStr(t.Elt) + " } "
  case (*ast.SelectorExpr):
    return "{ SEL " + typeStr(t.X) + " dot: " + typeStr(t.Sel) + " } "
  case (*ast.InterfaceType):
    return "{ INTERFACE todo: ...  } "
  case (*ast.FuncType):
    return "{ FN " + funcParamsResults(t) +  " } "
  case (*ast.MapType):
    return "{ MAP " + typeStr(t.Key) + " to: " + typeStr(t.Value) + " } "
  case (*ast.ChanType):
    return "{ CHAN " + Str(t.Dir) + " " + typeStr(t.Value) + " } "

  case (*ast.Object):
    panic("OBJECT")
    // return "{ OBJECT " + funcDeclStr(t.Decl.(*ast.FuncDecl)) + " } "
  }

  return fmt.Sprintf("{ ?WHAT? %#v } ", a)
}


func funcDeclStr(f *ast.FuncDecl) string {
  fstr := "{ "
  if f.Recv != nil {
    if len(f.Recv.List) != 1 { panic("f.Recv.List") }
    fstr += "METH  recv: " + typeStr(f.Recv.List[0].Type) + " "
  } else {
    fstr += "FUNC "
  }
  fstr += " name: " + f.Name.Name + " " + funcParamsResults(f.Type) + " "
  return fstr + " } "
}

func funcParamsResults(f *ast.FuncType) string {
  z := " args: { "

  // list of parameters
  params := f.Params
  i := 0
  if params == nil { panic("params is never nil but it is") }
  for _, lv := range params.List {
    pname := "_"
    if len(lv.Names) > 0 {
      pname = lv.Names[0].Name
    }
    tname := typeStr(lv.Type)
    if i > 0 { z += " , " }
    z +=  pname + " : " + tname
    i++
  }

  z += " } "
  z += " results: { "

  // list of parameters
  i = 0
  rr := f.Results
  if rr != nil {
    for _, lv := range rr.List {
      pname := "_"
      if len(lv.Names) > 0 {
        pname = lv.Names[0].Name
      }
      tname := typeStr(lv.Type)
      if i > 0 { z += " , " }
      z +=  pname + " : " + tname
    }
    i++
  }

  z += " } "
  return z
}

func main() {
  for _, dir := range os.Args[1:] {
    GrokDir(dir)
  }
}
