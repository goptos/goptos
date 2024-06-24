package io

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

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

func ListDirs(dir string) ([]string, error) {
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

func ListFiles(dirs []string) ([]string, error) {
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
			if v.Name() == strings.Split(dir, "/")[1]+".go" {
				names = append(names, dir+"/"+v.Name())
			}
		}
	}
	return names, nil
}

func CleanLine(line string) string {
	line = strings.Trim(line, " ")
	line = strings.Trim(line, "\t")
	return line
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
		fmt.Printf("%d\t%s\n", i, CleanLine(line))
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
