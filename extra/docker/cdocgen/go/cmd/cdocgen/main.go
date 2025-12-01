package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/boschglobal/dse.doc/docker/cdocgen/internal/app/cdocgen"
)

// Usage, maximum width 80 character, indentation with space characters.
var usage = `
CDOCGEN (from Dynamic Simulation Environment - Documentation Project)

  Generate Markdown documentation from MD formatted C comments.
  Suitable for integration with Hugo.

Examples:
  cdocgen -input module.h -output module.md -title Module -linktitle Module

  cdocgen \
      -input module.h \
      -output module.md \
      -title Module \
      -linktitle Module \
      -frontmatter \"{\"name\": \"foo\", \"content\": \"bar\"}\" ")

Flags:
`

func main() {
	input := flag.String("input", "", "path of C header file (required)")
	output := flag.String("output", "", "path of generated Markdown documentation file (required)")
	cDir := flag.String("cdir", "", "search path (comma separated list) for C code files which contain documentation")
	title := flag.String("title", "", "title (in frontmatter) of the generated Markdown documentation (required)")
	linktitle := flag.String("linktitle", "", "linktitle (in frontmatter) of the generated Markdown documentation (required)")
	frontmatter := flag.String("frontmatter", "", "frontmatter (JSON string) for the generated Markdown documentation")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), usage)
		flag.PrintDefaults()
	}
	flag.Parse()
	if *input == "" {
		fmt.Fprintf(os.Stderr, "ERROR: --input must be specified\n")
		os.Exit(1)
	}
	if *output == "" {
		fmt.Fprintf(os.Stderr, "ERROR: --output must be specified\n")
		os.Exit(1)
	}
	if *title == "" {
		fmt.Fprintf(os.Stderr, "ERROR: --title must be specified\n")
		os.Exit(1)
	}
	if *linktitle == "" {
		fmt.Fprintf(os.Stderr, "ERROR: --linktitle must be specified\n")
		os.Exit(1)
	}
	err := cdocgen.Generate(*input, *output, *cDir, *title, *linktitle, *frontmatter)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: --Output file was not generated\n")
		os.Exit(1)
	}
}
