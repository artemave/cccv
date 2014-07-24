package main

import (
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
+% go get github.com/artemave/cccv
 % git diff | cccv
diff --git a/cccv.go b/cccv.go
index af3d9dc..0e96324 100644
--- a/cccv.go
+++ b/cccv.go
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

func TestParsesDiff(t *testing.T) {
	RegisterTestingT(t)

	expectedChanges := []*Change{
		&Change{FileName: FileName("README.md"), Line: Line{Number: 10, Text: "% go get github.com/artemave/cccv"}},
		&Change{FileName: FileName("cccv.go"), Line: Line{Number: 18, Text: "type FileName string"}},
		&Change{FileName: FileName("cccv.go"), Line: Line{Number: 22, Text: "FileName"}},
		&Change{FileName: FileName("cccv.go"), Line: Line{Number: 48, Text: "o.O(results)"}},
	}

	changes := getChanges(strings.NewReader(diff))
	Expect(*changes).To(Equal(expectedChanges))
}
