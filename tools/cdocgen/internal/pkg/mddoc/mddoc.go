package mddoc

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.boschdevcloud.com/fsil/fsil.go/ast"
)

type Mddoc struct {
	Frontmatter Frontmatter
	Index       ast.Index
	Fragments   FragmentMap
}

type Frontmatter struct {
	title     string
	linkTitle string
	content   map[string]string
}

func (f *Frontmatter) SetTitle(title string, linkTitle string) {
	f.title = title
	f.linkTitle = linkTitle
}

func (f *Frontmatter) SetContent(content map[string]string) {
	f.content = make(map[string]string)
	f.content = content
}

type FragmentMap struct {
	ModuleLevelFragments map[string]string
	FunctionFragments    map[string]string
	OtherFragments       map[string]map[string]string
}

func checkEquals(line string) bool {
	length := 3
	pattern := fmt.Sprintf(`^([=]{%d,}|[-]{%d,})\n*$`, length, length)
	re := regexp.MustCompile(pattern)
	if re.MatchString(line) {
		return true
	} else {
		return false
	}
}

func checkMdUnderlineL1(line string) bool {
	pattern := fmt.Sprintf(`^([=]{%d,})\n*$`, 3)
	re := regexp.MustCompile(pattern)
	if re.MatchString(line) {
		return true
	} else {
		return false
	}
}

func checkMdUnderlineL2(line string) bool {
	pattern := fmt.Sprintf(`^([-]{%d,})\n*$`, 3)
	re := regexp.MustCompile(pattern)
	if re.MatchString(line) {
		return true
	} else {
		return false
	}
}

func removeLeadingSpacesAndAsterisk(line string, index int) string {
	res := line
	for i := 0; i < index && i < len(line); i++ {
		if line[i] == 32 || line[i] == 42 || line[i] == 57 {
			res = line[i+1:]
		} else {
			break
		}
	}
	return res
}

func findCFiles(directory string) []string {
	cFiles := []string{}

	filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return nil
		}

		if !info.IsDir() && strings.HasSuffix(info.Name(), ".c") {
			cFiles = append(cFiles, path)
		}

		return nil
	})

	return cFiles
}

func scanOtherFragments(fragMap *FragmentMap, sections []string, index int) {
	key := sections[0]
	innerMap := make(map[string]string)
	innerKey := key
	valid := false
	for num, line := range sections[1:] {
		if checkEquals(line) {
			// Previous line
			if len(sections[num]) > 0 {
				// Remove this line from the current map.
				innerMap[innerKey] = strings.Replace(innerMap[innerKey], sections[num], "", -1)
				// Create a new inner map entry
				innerKey = sections[num]
				valid = true
			}
		} else {
			innerMap[innerKey] = innerMap[innerKey] + line + "\n"
		}
	}
	if valid {
		fragMap.OtherFragments[key] = innerMap
	}
}

func mapMdUnderline(lines []string, adjust int) []string {
	var result []string
	for idx, line := range lines {
		if (idx + 1) < len(lines) {
			// H1 (===)
			if checkMdUnderlineL1((lines[idx+1])) {
				h := strings.Repeat("#", 1+adjust) + " " + line
				result = append(result, h)
				result = append(result, "")
				continue
			}
			if checkMdUnderlineL2((lines[idx+1])) {
				h := strings.Repeat("#", 2+adjust) + " " + line
				result = append(result, h)
				result = append(result, "")
				continue
			}
		}
		if checkEquals(line) {
			// Part of a header, don't emit.
			continue
		}
		result = append(result, line)
	}
	return result
}

func scan(file string, mddoc *Mddoc) error {
	code, err := os.ReadFile(file)
	if err != nil {
		fmt.Println("Error while reading header: ", err)
		return err
	}
	pattern := `/\*\*([^*][\s\S]*?)\*/`
	re := regexp.MustCompile(pattern)
	matches := re.FindAllStringSubmatch(string(code), -1)
	for _, match := range matches {
		commentBlock := strings.ReplaceAll(match[1], "\r", "")
		lines := strings.Split(commentBlock, "\n")
		if lines[0] == "" && len(lines) > 0 {
			lines = lines[1:]
		}
		// Detect and adjust for the comment style, based on first line.
		index := strings.IndexFunc(lines[0], func(r rune) bool {
			return r != ' ' && r != '*'
		})
		for i, v := range lines {
			lines[i] = removeLeadingSpacesAndAsterisk(v, index)
		}
		if len(lines) < 2 {
			continue
		}

		// A section has been identified.
		// FIXME move to 1.21 for entire repo and then use slices.Contains().
		contains := func(s []string, v string) bool {
			for i := range s {
				if v == s[i] {
					return true
				}
			}
			return false
		}
		if contains(mddoc.Index.Functions, lines[0]) {
			funcName := lines[0]
			startIndex := 1
			if checkEquals(lines[1]) {
				startIndex = 2
			}
			mdLines := mapMdUnderline(lines[startIndex:], 2)
			for _, line := range mdLines {
				mddoc.Fragments.FunctionFragments[funcName] += line + "\n"
			}
		} else if checkEquals(lines[1]) {
			fragName := lines[0]
			mdLines := mapMdUnderline(lines[2:], 1)
			for _, line := range mdLines {
				mddoc.Fragments.ModuleLevelFragments[fragName] += line + "\n"
			}
		} else {
			// Uncertain, therefore add to Other fragments.
			scanOtherFragments(&mddoc.Fragments, lines, index)
		}
	}
	return nil
}

func (mddoc *Mddoc) Scan(header string, scanDir ...string) error {
	mddoc.Fragments.ModuleLevelFragments = make(map[string]string)
	mddoc.Fragments.FunctionFragments = make(map[string]string)
	mddoc.Fragments.OtherFragments = make(map[string]map[string]string)

	// Scan the header file.
	err := scan(header, mddoc)
	if err != nil {
		return err
	}

	// Scan the C files.
	if len(scanDir) == 0 {
		folder, _ := os.Getwd()
		scanDir = append(scanDir, folder)
	}
	for _, folder := range scanDir {
		cFiles := findCFiles(folder)
		for _, cFile := range cFiles {
			err := scan(cFile, mddoc)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (mddoc *Mddoc) Generate(path string) error {
	doc := "---\n"
	doc += fmt.Sprintf("title: %s\n", mddoc.Frontmatter.title)
	doc += fmt.Sprintf("linkTitle: %s\n", mddoc.Frontmatter.linkTitle)
	for key, value := range mddoc.Frontmatter.content {
		doc += fmt.Sprintf("%s: %s\n", key, value)
	}
	doc += "---\n"
	for header, text := range mddoc.Fragments.ModuleLevelFragments {
		doc += fmt.Sprintf("## %s\n\n", header)
		doc += fmt.Sprintf("%s\n\n", text)
	}
	doc += "## Typedefs\n\n"
	// Sort the typedef.
	typedefNames := make([]string, 0, len(mddoc.Index.Typedefs))
	for k := range mddoc.Index.Typedefs {
		typedefNames = append(typedefNames, k)
	}
	sort.Strings(typedefNames)
	// Emit the sorted typedefs.
	for _, typedef := range typedefNames {
		variables := mddoc.Index.Typedefs[typedef]
		doc += fmt.Sprintf("### %s\n\n", typedef)
		if len(variables) > 0 {
			doc += "```c\n"
			doc += fmt.Sprintf("typedef struct %s {\n", typedef)
			for _, name := range variables {
				name = strings.Replace(name, " * ", "* ", 1)
				name = strings.Replace(name, " ** ", "** ", 1)
				doc += fmt.Sprintf("    %s;\n", name)
			}
			doc += "}\n"
			doc += "```\n\n"
		}
		if fragment, exists := mddoc.Fragments.OtherFragments[typedef]; exists {
			if value, exists := fragment[typedef]; exists {
				doc += value
			}
			for key, value := range fragment {
				if key != typedef {
					doc += fmt.Sprintf("#### %s\n\n", key)
					doc += value
				}
			}
		}

	}
	// Emit the sorted functions.
	doc += "## Functions\n\n"
	sort.Strings(mddoc.Index.Functions)
	for _, function := range mddoc.Index.Functions {
		if md, exists := mddoc.Fragments.FunctionFragments[function]; exists {
			doc += fmt.Sprintf("### %s\n", function)
			doc += fmt.Sprintf("%s\n\n", md)
		}
	}

	file, err := os.Create(path)
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		return err // FIXME need to raise an error, call needs to print and os.Exit(1)
	}
	defer file.Close()
	_, err = file.WriteString(doc)
	if err != nil {
		fmt.Printf("Error writing to file: %v\n", err)
		return err
	}
	fmt.Printf("Content saved to %s\n", path)
	return nil
}
