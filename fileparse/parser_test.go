package fileparse_test

import (
	"bufio"
	"bytes"
	"fmt"
	"testing"

	"github.com/nickwells/fileparse.mod/fileparse"
	"github.com/nickwells/location.mod/location"
	"github.com/nickwells/testhelper.mod/v2/testhelper"
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

	testhelper.DiffString(t, "EchoParser.ParseLine(...)", "",
		buf.String(), expVal)
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

		testhelper.DiffInt(t, id, "error count", len(errs), tc.expErrs)
		testhelper.DiffInt(t, id, "files seen", s.FilesVisited(), tc.expFiles)
		testhelper.DiffInt(t, id, "lines read", s.LinesRead(), tc.expLines)
		testhelper.DiffInt(t, id, "parsings", s.LinesParsed(), tc.expParsings)
	}
}

func TestEmptyStats(t *testing.T) {
	var s fileparse.Stats

	id := "An empty Stats structure"

	testhelper.DiffInt(t, id, "FilesVisited()", s.FilesVisited(), 0)
	testhelper.DiffInt(t, id, "LinesRead()", s.LinesRead(), 0)
	testhelper.DiffInt(t, id, "LinesParsed()", s.LinesParsed(), 0)

	expStr := "files:   0   lines read:     0   parsed:     0"

	testhelper.DiffString(t, id, "String representation", s.String(), expStr)
}
