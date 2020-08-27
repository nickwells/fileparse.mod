package fileparse_test

import (
	"bufio"
	"bytes"
	"fmt"
	"testing"

	"github.com/nickwells/fileparse.mod/fileparse"
	"github.com/nickwells/location.mod/location"
	"github.com/nickwells/testhelper.mod/testhelper"
)

func TestEchoParser(t *testing.T) {
	var buf bytes.Buffer
	w := bufio.NewWriter(&buf)
	ep := fileparse.EchoParser{Writer: w}
	loc := location.New("testLoc")

	expVal := "line"
	err := ep.ParseLine(expVal, loc)
	if err != nil {
		t.Error("unexpected error :", err)
	}
	w.Flush()
	expVal += "\n"

	testhelper.CmpValString(t, "EchoParser.ParseLine(...)", "",
		buf.String(), expVal)
}

func ExampleEchoParser() {
	var ep fileparse.EchoParser
	loc := location.New("testLoc")
	err := ep.ParseLine("line", loc)
	if err != nil {
		fmt.Println("unexpected error :", err)
	}
	// Output: line
}

func TestParse(t *testing.T) {
	var np fileparse.NullParser
	fpNull := fileparse.New("intro", np)

	testCases := []struct {
		filename    string
		expErrs     int
		expFiles    int
		expLines    int
		expParsings int
	}{
		{"~NoSuchUser/NoSuchFile", 1, 0, 0, 0},
		{"./testdata/NoSuchFile", 1, 0, 0, 0},
		{"./testdata/Empty", 0, 1, 0, 0},
		{"./testdata/OneInclude", 0, 2, 1, 0},
		{"./testdata/IncludeLoopSelf", 1, 1, 1, 0},
		{"./testdata/IncludeLoopStart", 1, 3, 3, 0},
		{"./testdata/BadIncludeFormat", 1, 1, 1, 0},
		{"./testdata/BadIncludeFileNonexistent", 1, 1, 1, 0},
		{"./testdata/FileWithContent", 0, 1, 2, 1},
	}

	for _, tc := range testCases {
		id := fmt.Sprintf("Parse(%q)", tc.filename)
		errs := fpNull.Parse(tc.filename)
		s := fpNull.Stats()
		testhelper.CmpValInt(t, id, "error count", len(errs), tc.expErrs)
		testhelper.CmpValInt(t, id, "files seen", s.FilesVisited(), tc.expFiles)
		testhelper.CmpValInt(t, id, "lines read", s.LinesRead(), tc.expLines)
		testhelper.CmpValInt(t, id, "parsings", s.LinesParsed(), tc.expParsings)
	}
}

func TestEmptyStats(t *testing.T) {
	var s fileparse.Stats
	id := "An empty Stats structure"

	testhelper.CmpValInt(t, id, "FilesVisited()", s.FilesVisited(), 0)
	testhelper.CmpValInt(t, id, "LinesRead()", s.LinesRead(), 0)
	testhelper.CmpValInt(t, id, "LinesParsed()", s.LinesParsed(), 0)
	expStr := "files:   0   lines read:     0   parsed:     0"
	testhelper.CmpValString(t, id, "String representation", s.String(), expStr)
}
