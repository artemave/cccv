package main

import (
	"log"
	"os"
	"regexp"
	"strings"
	"testing"
	. "github.com/onsi/gomega"
)

const diff = `
diff --git a/README.md b/README.md
index c9caa23..8ad488f 100644
--- a/README.md
+++ b/README.md
@@ -7,4 +7,5 @@ Check if git diff (commit, pr) contains copy pasted code.

'''
 % go get github.com/artemave/cccv
-nonsense
+% go get github.com/artemave/cccv
 % git diff | cccv
diff --git a/tmp_cccv.go b/tmp_cccv.go
index af3d9dc..0e96324 100644
--- a/tmp_cccv.go
+++ b/tmp_cccv.go
@@ -15,9 +15,11 @@ import (
 )
 
 type FileName string
+type FileName string
 
 type Change struct {
 	FileName
+	FileName
 	Line
 }
 
@@ -43,6 +45,7 @@ func main() {
 	}
 
 	o.O(results)
+	o.O(results)
 }
 
 func GenResultForFile(fName string, changes *[]*Change) FileResult {
`

var expectedChanges = []*Change{
	&Change{FileName: FileName("README.md"), Line: Line{Number: 10, Text: "% go get github.com/artemave/cccv"}},
	&Change{FileName: FileName("tmp_cccv.go"), Line: Line{Number: 18, Text: "type FileName string"}},
	&Change{FileName: FileName("tmp_cccv.go"), Line: Line{Number: 48, Text: "\to.O(results)"}},
}

var config = Config{
	ExcludeLines:               []*regexp.Regexp{},
	ExcludeFiles:               []*regexp.Regexp{},
	MinLineLength:              10,
	IgnoreHunksOfLinesLessThan: 1,
}

func TestParsesDiff(t *testing.T) {
	RegisterTestingT(t)

	changes := getChanges(strings.NewReader(diff), config)
	Expect(*changes).To(Equal(expectedChanges))
}

func TestFindsDuplicates(t *testing.T) {
	RegisterTestingT(t)

	WriteFile("/tmp/some_file.go", func(f *os.File) {
		f.WriteString("writes\n")
		f.WriteString("writes\n")
		f.WriteString("writes\n")
		f.WriteString("type FileName string\n")
		f.WriteString("writes\n")
	})
	defer os.Remove("tmp_cccv.go")

	expectedResult := FileResult{
		FileName: FileName("/tmp/some_file.go"),
		Lines: []*Line{
			&Line{Number: 4, Text: "type FileName string"},
		},
	}
	result := GenResultForFile("/tmp/some_file.go", &expectedChanges, config)
	Expect(result).To(Equal(expectedResult))
}

func TestIgnoresLinesShorterThan10c(t *testing.T) {
	RegisterTestingT(t)

	WriteFile("/tmp/some_file.go", func(f *os.File) {
		f.WriteString("writes\n")
		f.WriteString("writes\n")
		f.WriteString("writes\n")
		f.WriteString("FileName\n")
		f.WriteString("writes\n")
	})
	defer os.Remove("tmp_cccv.go")

	expectedResult := FileResult{
		FileName: FileName("/tmp/some_file.go"),
		Lines:    []*Line{},
	}
	result := GenResultForFile("/tmp/some_file.go", &expectedChanges, config)
	Expect(result).To(Equal(expectedResult))
}

func TestConsidersIndentedButOtherwiseIdenticalLinesAsDuplicates(t *testing.T) {
	RegisterTestingT(t)

	WriteFile("/tmp/some_file.go", func(f *os.File) {
		f.WriteString("writes\n")
		f.WriteString("\t o.O(results) \n")
	})
	defer os.Remove("tmp_cccv.go")

	expectedResult := FileResult{
		FileName: FileName("/tmp/some_file.go"),
		Lines: []*Line{
			&Line{Number: 2, Text: "\t o.O(results) "},
		},
	}
	result := GenResultForFile("/tmp/some_file.go", &expectedChanges, config)
	Expect(result).To(Equal(expectedResult))
}

func TestDoesNotCountChangesThemselvesAsDuplicates(t *testing.T) {
	RegisterTestingT(t)

	WriteFile("tmp_cccv.go", func(f *os.File) {
		for i := 0; i < 17; i++ {
			f.WriteString("writes\n")
		}
		f.WriteString("type FileName string\n")
	})
	defer os.Remove("tmp_cccv.go")

	expectedResult := FileResult{
		FileName: FileName("tmp_cccv.go"),
		Lines:    []*Line{},
	}
	result := GenResultForFile("tmp_cccv.go", &expectedChanges, config)
	Expect(result).To(Equal(expectedResult))
}

func TestIncludesChangeOnlyOnce(t *testing.T) {
	RegisterTestingT(t)

	WriteFile("/tmp/some_file.go", func(f *os.File) {
		f.WriteString("writes\n")
		f.WriteString("type FileName string\n")
		f.WriteString("writes\n")
		f.WriteString("type FileName string\n")
	})
	defer os.Remove("tmp_cccv.go")

	expectedChanges := expectedChanges
	expectedChanges = append(expectedChanges, &Change{FileName: FileName("tmp_cccv.go"), Line: Line{Number: 28, Text: "type FileName string"}})

	expectedResult := FileResult{
		FileName: FileName("/tmp/some_file.go"),
		Lines: []*Line{
			&Line{Number: 2, Text: "type FileName string"},
			&Line{Number: 4, Text: "type FileName string"},
		},
	}
	result := GenResultForFile("/tmp/some_file.go", &expectedChanges, config)
	Expect(result).To(Equal(expectedResult))
}

func TestExcludesLinesYaml(t *testing.T) {
	RegisterTestingT(t)

	config := config
	config.ExcludeLines = []*regexp.Regexp{regexp.MustCompile("github")}

	WriteFile("/tmp/some_file.go", func(f *os.File) {
		f.WriteString("% go get github.com/artemave/cccv\n")
		f.WriteString("type FileName string\n")
	})
	defer os.Remove("tmp_cccv.go")

	expectedResult := FileResult{
		FileName: FileName("/tmp/some_file.go"),
		Lines: []*Line{
			&Line{Number: 2, Text: "type FileName string"},
		},
	}
	result := GenResultForFile("/tmp/some_file.go", &expectedChanges, config)
	Expect(result).To(Equal(expectedResult))
}

func TestIgnoresDuplicatesOfLessThanNLinesLong(t *testing.T) {
	RegisterTestingT(t)

	config := config
	config.IgnoreHunksOfLinesLessThan = 2

	WriteFile("/tmp/some_file.go", func(f *os.File) {
		f.WriteString("line 1\n")
		f.WriteString("added line 2\n")
		f.WriteString("added line 3\n")
		f.WriteString("line 4\n")
		f.WriteString("added line 5\n")
		f.WriteString("added line 6\n")
		f.WriteString("line 7\n")
	})
	defer os.Remove("tmp_cccv.go")

	changes := []*Change{
		&Change{FileName: FileName("README"), Line: Line{Number: 10, Text: "added line 2"}},
		&Change{FileName: FileName("README"), Line: Line{Number: 11, Text: "added line 3"}},
		&Change{FileName: FileName("README"), Line: Line{Number: 21, Text: "added line 5"}},
		&Change{FileName: FileName("README"), Line: Line{Number: 31, Text: "added line 6"}},
	}

	expectedResult := FileResult{
		FileName: FileName("/tmp/some_file.go"),
		Lines: []*Line{
			&Line{Number: 2, Text: "added line 2"},
			&Line{Number: 3, Text: "added line 3"},
		},
	}

	result := GenResultForFile("/tmp/some_file.go", &changes, config)
	Expect(result).To(Equal(expectedResult))
}

func WriteFile(fname string, callback func(f *os.File)) {
	f, err := os.Create(fname)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	callback(f)
	f.Sync()
}
