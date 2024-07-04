package codegen

import (
	"fmt"
	"log"
	"strings"

	"github.com/goptos/goptos/io"
	"github.com/goptos/stateparser"
)

const (
	xHtmlElem     int = 0
	HtmlEventAttr int = 1
	xHtmlAttr     int = 2
	HtmlText      int = 3
	HtmlDynText   int = 100
)

func View(src string) {
	const varTag = "var view *Elem"
	const viewStartTag = "/* View"
	const viewEndTag = "*/"
	const codeStartTag = "/* macro:generated:view:start */"
	const codeEndTag = "/* macro:generated:view:end */"

	var srcDir = src + "/comp"
	var parser = stateparser.New()

	// List all dirs in comp directory (we start in src/)
	dirs, err := io.ListCompDirs(srcDir)
	if err != nil {
		log.Printf("ListDirs() %s\n", err)
		return
	}
	if len(dirs) < 1 {
		log.Printf("ListDirs() no directories found in %s\n", srcDir)
		return
	}

	// List all files in each dir found
	files, err := io.ListCompFiles(dirs)
	if err != nil {
		log.Printf("ListFiles() %s\n", err)
		return
	}
	if len(files) < 1 {
		log.Printf("ListFiles() no components found in %s\n", srcDir)
		return
	}

	// Process each component file
	for _, file := range files {
		log.Printf("[%s]\n", file)
		lines, err := io.ReadFile(file)
		if err != nil {
			log.Printf("  ReadFile() %s\n", err)
			return
		}

		// find var (to determine it's leading white spaces)
		varLine, err := io.FindTag(varTag, lines)
		if err != nil {
			log.Printf("  FindTag() %s\n", err)
			return
		}
		var varBuff = io.GetLeadingWhiteSpace(lines[varLine])

		// remove previous generated code
		from, to, err := io.FindSection(codeStartTag, codeEndTag, lines)
		if err == nil {
			lines = append(lines[:from-1], lines[to:]...)
			lines[from-1] = varTag
		}

		// find var again (to receive generated code)
		varLine, err = io.FindTag(varTag, lines)
		if err != nil {
			log.Printf("  FindTag() %s\n", err)
			return
		}

		// find view template
		from, to, err = io.FindSection(viewStartTag, viewEndTag, lines)
		if err != nil {
			log.Printf("  FindSection() %s\n", err)
			return
		}

		// enable imports
		for i, line := range lines {
			var tmp = strings.Split(io.CleanLine(line), " ")
			if len(tmp) != 5 {
				continue
			}
			if tmp[0] != "_" {
				continue
			}
			if strings.Join(tmp[2:5], " ") != "/* macro:import */" {
				continue
			}
			if i-1 >= 0 {
				if io.CleanLine(lines[i-1]) != tmp[1] {
					lines[i] = io.GetLeadingWhiteSpace(lines[i]) + tmp[1] + "\r\n" + lines[i]
				}
			}
		}

		// generate go code from template
		err = parser.ParseView(strings.Join(lines[from:to], "\n"))
		if err != nil {
			log.Printf("  View() %s\n", err)
			return
		}

		// replace var with generated code
		lines[varLine] = fmt.Sprintf("%s\n%s\n%s",
			varBuff+codeStartTag,
			varBuff+varTag+" = "+parser.Result,
			varBuff+codeEndTag)

		// correct file ending line
		lines[len(lines)-1] = lines[len(lines)-1] + "\r\n"

		err = io.WriteFile(file, lines)
		if err != nil {
			fmt.Printf("writeFile() %s\n", err)
			return
		}
	}
}
