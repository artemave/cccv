package main

import (
	"log"
	"os"
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
	&Change{FileName: FileName("tmp_cccv.go"), Line: Line{Number: 22, Text: "\tFileName"}},
	&Change{FileName: FileName("tmp_cccv.go"), Line: Line{Number: 48, Text: "\to.O(results)"}},
}

func TestParsesDiff(t *testing.T) {
	RegisterTestingT(t)

	changes := getChanges(strings.NewReader(diff))
	Expect(*changes).To(Equal(expectedChanges))
}

func TestFindsDuplicates(t *testing.T) {
	RegisterTestingT(t)

	WriteFile("/tmp/some_file.go", func(f *os.File) {
		f.WriteString("writes\n")
		f.WriteString("o.O(results)\n")
		f.WriteString("writes\n")
		f.WriteString("type FileName string\n")
		f.WriteString("writes\n")
	})
	defer os.Remove("tmp_cccv.go")

	expectedResult := FileResult{
		FileName: FileName("/tmp/some_file.go"),
		Lines: []*Line{
			&Line{Number: 2, Text: "o.O(results)"},
			&Line{Number: 4, Text: "type FileName string"},
		},
	}
	result := GenResultForFile("/tmp/some_file.go", &expectedChanges)
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
	result := GenResultForFile("tmp_cccv.go", &expectedChanges)
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
