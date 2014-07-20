package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/AndrewVos/o"
)

type ChangedLine struct {
	FileName string
	Text     string
}

func main() {
	result := make(map[string][]string)
	changedLineItselfAlreadySkipped := make(map[ChangedLine]bool)

	changes := getChanges()
	allGitFiles := gitLsFiles()

	for _, fName := range allGitFiles {
		file, _ := os.Open(fName)
		scanner := bufio.NewScanner(file)
		uniqLines := make(map[string]bool)

		for scanner.Scan() {
			line := scanner.Text()
			trimmedLine := strings.TrimSpace(line)

			if uniqLines[trimmedLine] {
				continue
			}
			uniqLines[trimmedLine] = true

			for _, changedLine := range *changes {
				if changedLine.Text == trimmedLine {

					if changedLine.FileName == fName && !changedLineItselfAlreadySkipped[*changedLine] {
						changedLineItselfAlreadySkipped[*changedLine] = true
						continue
					}

					if result[fName] == nil {
						result[fName] = []string{}
					}
					result[fName] = append(result[fName], line)
				}
			}
		}
	}

	o.O(result)
}

func gitLsFiles() []string {
	files := []string{}
	cmd := exec.Command("git", "ls-files")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(stdout)

	for scanner.Scan() {
		files = append(files, scanner.Text())
	}

	return files
}

func getChanges() *[]*ChangedLine {
	scanner := bufio.NewScanner(os.Stdin)
	var currentFile string
	changes := &[]*ChangedLine{}

	currentFileR := regexp.MustCompile(`^\+\+\+ ./(.*)$`)
	lineAddedR := regexp.MustCompile(`^\+{1}(.*\w+.*)`)

	for scanner.Scan() {
		currentLine := scanner.Text()

		if currentFileR.MatchString(currentLine) {
			res := currentFileR.FindStringSubmatch(currentLine)
			currentFile = res[1]
		} else if lineAddedR.MatchString(currentLine) == true {
			res := lineAddedR.FindStringSubmatch(currentLine)
			newChange := &ChangedLine{FileName: currentFile, Text: strings.TrimSpace(res[1])}
			*changes = append(*changes, newChange)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
		os.Exit(1)
	}

	return changes
}
