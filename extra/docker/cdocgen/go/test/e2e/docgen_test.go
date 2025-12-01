//go:build e2e
// +build e2e

// to run this test: go test -tags e2e
package e2e

import (
	_ "github.com/boschglobal/dse.doc/docker/cdocgen/test/init/workingdir"

	"os"
	"os/exec"
	"regexp"
	"testing"
)

func TestCLI(t *testing.T) {
	tmpDir := t.TempDir()
	cmd := exec.Command("bin/cdocgen",
		"--input",
		"test/testdata/header.h",
		"--output",
		tmpDir+"/output_new.md",
		"--title",
		"header.h",
		"--linktitle",
		"linkheader.h",
		"--frontmatter",
		"{\"name\": \"foo\", \"content\": \"bar\"}",
		"--cdir",
		"/home/gok1abt/fsil/dse.doc/tools/cdocgen/test/testdata/cfiles",
	)
	_, err := cmd.Output()
	if err != nil {
		t.Errorf("Unexpected Error %s", err)
	}
	data, err := os.ReadFile(tmpDir + "/output_new.md")
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

func TestCLIDefault(t *testing.T) {
	tmpDir := t.TempDir()
	cmd := exec.Command("bin/cdocgen",
		"--input",
		"test/testdata/header.h",
		"--output",
		tmpDir+"/output_new.md",
		"--title",
		"header.h",
		"--linktitle",
		"linkheader.h",
		"--frontmatter",
		"{\"name\": \"foo\", \"content\": \"bar\"}",
	)
	_, err := cmd.Output()
	if err != nil {
		t.Errorf("Unexpected Error %s", err)
	}
	data, err := os.ReadFile(tmpDir + "/output_new.md")
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

func TestCLIWithoutFrontmatter(t *testing.T) {
	tmpDir := t.TempDir()
	cmd := exec.Command("bin/cdocgen",
		"--input",
		"test/testdata/header.h",
		"--output",
		tmpDir+"/output_new.md",
		"--title",
		"header.h",
		"--linktitle",
		"linkheader.h",
	)
	_, err := cmd.Output()
	if err != nil {
		t.Errorf("Unexpected Error %s", err)
	}
	data, err := os.ReadFile(tmpDir + "/output_new.md")
	if err != nil {
		t.Errorf("Error. Expected output.md . Not found. ")
	}
	testData := []struct {
		name string
	}{
		{"MyStruct"},
		{"## Module Level Doc"},
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

func TestCLIWithIncorrectFrontmatterFormat(t *testing.T) {
	tmpDir := t.TempDir()
	cmd := exec.Command("bin/cdocgen",
		"--input",
		"test/testdata/header.h",
		"--output",
		tmpDir+"/output_new.md",
		"--title",
		"header.h",
		"--linktitle",
		"linkheader.h",
		"--frontmatter",
		"{wrong format}",
	)
	_, err := cmd.Output()
	if err != nil {
		t.Errorf("Unexpected Error %s", err)
	}
	data, err := os.ReadFile(tmpDir + "/output_new.md")
	if err != nil {
		t.Errorf("Error. Expected output.md . Not found. ")
	}
	testData := []struct {
		name string
	}{
		{"wrong format"},
		{"name: foo"},
	}
	for _, td := range testData {
		regex, err := regexp.Compile(td.name)
		if err != nil {
			t.Errorf("Unexpected error")
		}
		if regex.MatchString(string(data)) {
			t.Errorf("Error. Expected '%s' not found. ", td)
		}
	}

}

func TestCLIWithIncorrectOutputFile(t *testing.T) {
	tmpDir := t.TempDir()
	cmd := exec.Command("bin/cdocgen",
		"--input",
		"test/testdata/header.h",
		"--output",
		tmpDir,
		"--title",
		"header.h",
		"--linktitle",
		"linkheader.h",
		"--frontmatter",
		"{wrong format}",
	)
	_, err := cmd.Output()
	if err == nil {
		t.Errorf("Error did not occur")
	}
}
