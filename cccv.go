package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/AndrewVos/o"
)

type FileName string

type Change struct {
	FileName
	Line
}

type Line struct {
	Number int
	Text   string
}

type FileResult struct {
	FileName
	Lines []*Line
}

func main() {
	results := []FileResult{}

	changes := getChanges(os.Stdin)
	allGitFiles := gitLsFiles()

	for _, fName := range allGitFiles {
		r := GenResultForFile(fName, changes)
		results = append(results, r)
	}

	o.O(results)
}

func GenResultForFile(fName string, changes *[]*Change) FileResult {
	file, _ := os.Open(fName)
	scanner := bufio.NewScanner(file)
	currentLineNumber := 0
	result := FileResult{FileName: FileName(fName), Lines: []*Line{}}

	for scanner.Scan() {
		line := scanner.Text()
		trimmedLine := strings.TrimSpace(line)
		currentLineNumber++

		for _, changedLine := range *changes {
			if changedLine.Text == trimmedLine {

				// exclude lines from the diff itself
				if string(changedLine.FileName) == fName && changedLine.Line.Number == currentLineNumber {
					continue
				}

				result.Lines = append(result.Lines, &Line{Number: currentLineNumber, Text: line})
			}
		}
	}
	return result
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

func getChanges(reader io.Reader) *[]*Change {
	scanner := bufio.NewScanner(reader)
	var currentFile string
	var currentLineNumber int
	changes := &[]*Change{}

	currentFileR := regexp.MustCompile(`^\+\+\+ ./(.*)$`)
	lineAddedR := regexp.MustCompile(`^\+{1}(.*\w+.*)`)
	lineRangeR := regexp.MustCompile(`^@@.*?\+(\d+?),`)

	for scanner.Scan() {
		currentLine := scanner.Text()

		if res := currentFileR.FindStringSubmatch(currentLine); res != nil {
			currentFile = res[1]

		} else if res := lineRangeR.FindStringSubmatch(currentLine); res != nil {
			r, err := strconv.Atoi(res[1])
			if err != nil {
				log.Fatal(err)
				os.Exit(1)
			}
			currentLineNumber = r

		} else if lineAddedR.MatchString(currentLine) == true {
			res := lineAddedR.FindStringSubmatch(currentLine)
			newChange := &Change{
				FileName: FileName(currentFile),
				Line:     Line{Text: strings.TrimSpace(res[1]), Number: currentLineNumber},
			}
			*changes = append(*changes, newChange)
			currentLineNumber++

		} else {
			currentLineNumber++
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
		os.Exit(1)
	}

	return changes
}
