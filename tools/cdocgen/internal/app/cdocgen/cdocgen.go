package cdocgen

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.boschdevcloud.com/fsil/fsil.go/ast"
	"github.com/boschglobal/dse.doc/docker/cdocgen/internal/pkg/mddoc"
)

func Generate(input string, output string, cDir string, title string, linktitle string, jsonString string) error {
	cDirList := strings.Split(cDir, ",")
	doc := mddoc.Mddoc{}
	var frontmatter map[string]string
	if len(jsonString) > 0 {
		err := json.Unmarshal([]byte(jsonString), &frontmatter)
		if err != nil {
			fmt.Printf("error: Frontmatter could not be parsed")
		}
	}
	doc.Frontmatter.SetTitle(title, linktitle)
	doc.Frontmatter.SetContent(frontmatter)
	ast := ast.Ast{
		Path: input,
	}
	err := ast.Load()
	if err != nil {
		fmt.Printf("Error: %s", err)
	}
	ast.Parse(&doc.Index)
	if cDirList[0] == "" {
		doc.Scan(ast.Path)
	} else {
		doc.Scan(ast.Path, cDirList...)
	}
	err = doc.Generate(output)
	return err
}
