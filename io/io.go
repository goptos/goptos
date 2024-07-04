package io

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/goptos/utils"
)

var verbose = (*utils.Verbose).New(nil)

func ReadFile(fileName string) ([]string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var lines = []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text()) // add the new line to the stored lines
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return lines, nil
}

func WriteFile(fileName string, lines []string) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString(strings.Join(lines, "\n"))
	if err != nil {
		return err
	}
	return nil
}

func WritePath(path string) error {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func ListCompDirs(dir string) ([]string, error) {
	f, err := os.Open(dir)
	if err != nil {
		return nil, err
	}
	files, err := f.Readdir(0)
	if err != nil {
		return nil, err
	}
	var names = []string{}
	for _, v := range files {
		if v.IsDir() {
			names = append(names, dir+"/"+v.Name())
		}
	}
	return names, nil
}

func ListCompFiles(dirs []string) ([]string, error) {
	var names = []string{}
	for _, dir := range dirs {
		f, err := os.Open(dir)
		if err != nil {
			return nil, err
		}
		files, err := f.Readdir(0)
		if err != nil {
			return nil, err
		}
		for _, v := range files {
			if v.IsDir() {
				continue
			}
			var splitDirName = strings.Split(dir, "/")
			if v.Name() == splitDirName[len(splitDirName)-1]+".go" {
				names = append(names, dir+"/"+v.Name())
			}
		}
	}
	return names, nil
}

// Return the input string without any leading or trailing spaces " " or tabs "\t"
func CleanLine(line string) string {
	line = strings.Trim(line, " ")
	line = strings.Trim(line, "\t")
	return line
}

// Return the leading white space of the input string (if it has any spaces " " or tabs "\t")
func GetLeadingWhiteSpace(line string) string {
	var buffer = ""
	for _, char := range line {
		switch string(char) {
		case " ":
			buffer = buffer + " "
		case "\t":
			buffer = buffer + "\t"
		default:
			return buffer
		}
	}
	return buffer
}

func FindTag(tag string, lines []string) (int, error) {
	for i, line := range lines {
		var cleaned = CleanLine(line)
		if len(cleaned) < len(tag) {
			continue
		}
		if cleaned[0:len(tag)] == tag {
			return i, nil
		}
	}
	return -1, fmt.Errorf("no '%s' found", tag)
}

func FindSection(startTag string, endTag string, lines []string) (int, int, error) {
	var from = 0
	var hit = false
	for i, line := range lines {
		verbose.Printf(4, "%d\t%s\n", i, CleanLine(line))
		if CleanLine(line) == startTag {
			if from+1 < len(lines) {
				from = i + 1
				hit = true
			}
		}
		if hit {
			if CleanLine(line) == endTag {
				if i > from {
					return from, i, nil
				}
			}
		}
	}
	return -1, -1, fmt.Errorf("no '%s' found", startTag)
}
