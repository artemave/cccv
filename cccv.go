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

	"github.com/kr/pretty"
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

func (fr *FileResult) HasDuplicates() bool {
	return len(fr.Lines) > 0
}

func main() {
	results := []FileResult{}

	changes := getChanges(os.Stdin)
	allGitFiles := gitLsFiles()

	for _, fName := range allGitFiles {
		r := GenResultForFile(fName, changes)
		results = append(results, r)
	}

	thereAreDuplicates := false
	for _, r := range results {
		if r.HasDuplicates() {
			thereAreDuplicates = true
		}
	}

	if thereAreDuplicates {
		pretty.Print(results)
		os.Exit(1)
	}
}

func GenResultForFile(fName string, changes *[]*Change) FileResult {
	file, _ := os.Open(fName)
	scanner := bufio.NewScanner(file)
	currentLineNumber := 0
	result := FileResult{FileName: FileName(fName), Lines: []*Line{}}

	for scanner.Scan() {
		line := scanner.Text()
		currentLineNumber++

		for _, change := range *changes {
			if strings.TrimFunc(change.Text, TrimF) == strings.TrimFunc(line, TrimF) {

				// exclude lines from the diff itself
				if string(change.FileName) == fName && change.Line.Number == currentLineNumber {
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
	lineRemovedR := regexp.MustCompile(`^\-{1}`)
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

		} else if lineAddedR.MatchString(currentLine) {
			res := lineAddedR.FindStringSubmatch(currentLine)

			if len(strings.TrimFunc(res[1], TrimF)) <= 10 {
				currentLineNumber++
				continue
			}

			newChange := &Change{
				FileName: FileName(currentFile),
				Line:     Line{Text: res[1], Number: currentLineNumber},
			}
			*changes = append(*changes, newChange)
			currentLineNumber++

		} else if !lineRemovedR.MatchString(currentLine) {
			currentLineNumber++
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
		os.Exit(1)
	}

	return changes
}

func TrimF(c rune) bool { return (c == 32 || c == 9) }
