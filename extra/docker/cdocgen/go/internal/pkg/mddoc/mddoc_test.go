package mddoc

import (
	"reflect"
	"strings"

	_ "github.com/boschglobal/dse.doc/docker/cdocgen/test/init/workingdir"

	"encoding/json"
	"os"
	"regexp"
	"testing"

	"github.com/boschglobal/dse.clib/extra/go/ast"
)

//scan for all valid fragments
//process individual valid fragment
//In individual fragment, have a main key as heading and then other keys which are subheading (parametres, examples, etc )
//Match typedfs and functions from index, with individual fragment heading

func TestFileMissing(t *testing.T) {
	mddoc := Mddoc{}
	err := mddoc.Scan("test/testdata/missing_header.h")
	if err == nil {
		t.Errorf("Error did not occur.")
	}
}

type lenChecks struct {
	field    string
	expected int
}
type existChecks struct {
	field    string
	expected string
}

func TestScan(t *testing.T) {
	testTable := []struct {
		testCase    string
		headerPath  string
		scanPath    string
		lenChecks   []lenChecks
		existChecks []existChecks
	}{
		{
			testCase:   "default scan",
			headerPath: "test/testdata/header.h",
			lenChecks: []lenChecks{
				{"Fragments.ModuleLevelFragments", 3},
				{"Fragments.FunctionFragments", 0},
				{"Fragments.OtherFragments", 4},
			},
			existChecks: []existChecks{
				{"Fragments.OtherFragments", "myFunction"},
			},
		},
		{
			testCase:   "cfile scan",
			headerPath: "test/testdata/header.h",
			scanPath:   "test/testdata/cfiles",
			lenChecks: []lenChecks{
				{"Fragments.ModuleLevelFragments", 3},
				{"Fragments.FunctionFragments", 0},
				{"Fragments.OtherFragments", 4},
			},
			existChecks: []existChecks{
				{"Fragments.ModuleLevelFragments", "MyHeader.h"},
				{"Fragments.OtherFragments", "MyStruct"}, // TODO check for nested val["Example"]
				{"Fragments.OtherFragments", "myFunction"},
			},
		},
	}
	for _, ti := range testTable {
		mddoc := Mddoc{}
		if ti.scanPath != "" {
			if err := mddoc.Scan(ti.headerPath, ti.scanPath); err != nil {
				t.Errorf("(%s) Unexpected Error while scanning path: %s", ti.testCase, ti.headerPath)
			}
		} else {
			if err := mddoc.Scan(ti.headerPath); err != nil {
				t.Errorf("(%s) Unexpected Error while scanning path: %s", ti.testCase, ti.headerPath)
			}
		}
		for _, lc := range ti.lenChecks {
			r := reflect.ValueOf(mddoc)
			var f reflect.Value
			for _, fieldName := range strings.Split(lc.field, ".") {
				f = r.FieldByName(fieldName)
				r = reflect.ValueOf(f.Interface())
			}
			len := f.Len()
			if len != lc.expected {
				t.Errorf("(%s) Error: %s, expected %d (have %d)", ti.testCase, lc.field, lc.expected, len)
			}
		}
		for _, ec := range ti.existChecks {
			r := reflect.ValueOf(mddoc)
			var f reflect.Value
			for _, fieldName := range strings.Split(ec.field, ".") {
				f = r.FieldByName(fieldName)
				r = reflect.ValueOf(f.Interface())
			}
			ok := false
			for _, k := range f.MapKeys() {
				if k.String() == ec.expected {
					ok = true
					break
				}
			}
			if !ok {
				t.Errorf("Error: %s, did not contain expected %s", ec.field, ec.expected)
			}
		}

	}
}

func TestInvalidFragment(t *testing.T) {
	data := []struct {
		name string
	}{
		{"# Not Valid H1"},
		{"# Not Valid H2"},
		{"Not Valid H1 1 star"},
		{"Not Valid H1 3 star"},
		{"Not Valid H1"},
	}
	mddoc := Mddoc{}
	err := mddoc.Scan("test/testdata/header.h", "test/testdata/cfiles")
	if err != nil {
		t.Errorf("Unexpected Error while scanning md fragments")
	}
	for _, d := range data {
		_, ok := mddoc.Fragments.ModuleLevelFragments[d.name]
		if ok {
			t.Errorf("Unexpected error while scanning md fragments. %s found ", d.name)
		}
	}
	_, ok := mddoc.Fragments.OtherFragments["invalidFunction"]
	if ok {
		t.Errorf("Unexpected error while scanning md fragments. 'invalidFunction' found ")
	}
}

func TestMddocNoFunctionMatch(t *testing.T) {
	doc := Mddoc{}
	tmpDir := t.TempDir()
	ast := ast.Ast{
		Path: "test/testdata/header.h",
	}
	err := ast.Load()
	if err != nil {
		t.Errorf("Unexpected error.")
	}
	ast.Parse(&doc.Index)
	doc.Scan(ast.Path, "test/testdata/cfiles")
	err = doc.Generate(tmpDir + "/output.md")
	if err != nil {
		t.Errorf("Unexpected Error while Generate: %s", err)
	}
	data, err := os.ReadFile(tmpDir + "/output.md")
	if err != nil {
		t.Errorf("Error. Expected output.md . Not found. ")
	}
	function := "missingFunction"
	regex, err := regexp.Compile(function)
	if err != nil {
		t.Errorf("Unexpected error")
	}
	if regex.MatchString(string(data)) {
		t.Errorf("Error. unexpected 'missingFunction' found. ")
	}
}

func TestGenerate(t *testing.T) {
	doc := Mddoc{}
	tmpDir := t.TempDir()
	jsonStr := `{"name": "foo", "content": "bar"}`
	var frontmatter map[string]string
	err := json.Unmarshal([]byte(jsonStr), &frontmatter)
	if err != nil {
		t.Errorf("Unexpected error.")
	}
	doc.Frontmatter.SetTitle("header.h", "link_header.h")
	doc.Frontmatter.SetContent(frontmatter)
	ast := ast.Ast{
		Path: "test/testdata/header.h",
	}
	err = ast.Load()
	if err != nil {
		t.Errorf("Unexpected error.")
	}
	ast.Parse(&doc.Index)
	doc.Scan(ast.Path, "test/testdata/cfiles")
	err = doc.Generate(tmpDir + "/output.md")
	if err != nil {
		t.Errorf("Unexpected Error while Generate: %s", err)
	}
	data, err := os.ReadFile(tmpDir + "/output.md")
	if err != nil {
		t.Errorf("Error. Expected output.md . Not found. ")
	}
	testData := []struct {
		name string
	}{
		{"MyStruct"},
		{"## Module Level Doc"},
		{"name: foo"},
	}
	for _, td := range testData {
		regex, err := regexp.Compile(td.name)
		if err != nil {
			t.Errorf("Unexpected error")
		}
		if !regex.MatchString(string(data)) {
			t.Errorf("Error. Expected '%s' not found. ", td)
		}
	}
}

func TestMissingOutputFile(t *testing.T) {
	doc := Mddoc{}
	jsonStr := `{"name": "foo", "content": "bar"}`
	var frontmatter map[string]string
	err := json.Unmarshal([]byte(jsonStr), &frontmatter)
	if err != nil {
		t.Errorf("Unexpected error.")
	}
	doc.Frontmatter.SetTitle("header.h", "link_header.h")
	doc.Frontmatter.SetContent(frontmatter)
	ast := ast.Ast{
		Path: "test/testdata/header.h",
	}
	err = ast.Load()
	if err != nil {
		t.Errorf("Unexpected error.")
	}
	ast.Parse(&doc.Index)
	doc.Scan(ast.Path, "test/testdata/cfiles")
	err = doc.Generate("temporary/wrong/output.md")
	if err == nil {
		t.Errorf("Error did not occur")
	}
}
