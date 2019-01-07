package fileparse_test

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/nickwells/fileparse.mod/fileparse"
	"github.com/nickwells/location.mod/location"
)

func TestEchoParser(t *testing.T) {
	var buf bytes.Buffer
	w := bufio.NewWriter(&buf)
	ep := fileparse.EchoParser{Writer: w}
	loc := location.New("testLoc")

	expVal := "line"
	ep.ParseLine(expVal, loc)
	w.Flush()
	expVal += "\n"

	if buf.String() != expVal {
		t.Errorf("EchoParser.ParseLine(...) failed: expected: '%s', got: '%s'",
			expVal, buf.String())
	}
}

func ExampleEchoParser() {
	var ep fileparse.EchoParser
	loc := location.New("testLoc")
	ep.ParseLine("line", loc)
	// Output: line
}

func TestParse(t *testing.T) {
	var np fileparse.NullParser
	fpNull := fileparse.New("intro", np)

	testCases := []struct {
		filename            string
		expectedErrCount    int
		expectedFileCount   int
		expectedLineCount   int
		expectedParsedCount int
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
		errs := fpNull.Parse(tc.filename)
		if ecount := len(errs); ecount != tc.expectedErrCount {
			t.Error("Parse(", tc.filename, ") failed - expected: ",
				tc.expectedErrCount, " errors, got: ", ecount, "\n",
				"errors: ", errs)
		}
		if fpNull.Stats().FilesVisited() != tc.expectedFileCount {
			t.Error("Parse(", tc.filename, ") failed - expected: ",
				tc.expectedFileCount, " files visited, got: ",
				fpNull.Stats().FilesVisited())
		}
		if fpNull.Stats().LinesRead() != tc.expectedLineCount {
			t.Error("Parse(", tc.filename, ") failed - expected: ",
				tc.expectedLineCount, " lines read, got: ",
				fpNull.Stats().LinesRead())
		}
		if fpNull.Stats().LinesParsed() != tc.expectedParsedCount {
			t.Error("Parse(", tc.filename, ") failed - expected: ",
				tc.expectedParsedCount, " lines parsed, got: ",
				fpNull.Stats().LinesParsed())
		}
	}
}

func TestStats(t *testing.T) {
	var s fileparse.Stats

	if fv := s.FilesVisited(); fv != 0 {
		t.Error("an empty Stats structure should have filesVisited: 0, has: ",
			fv)
	}
	if lr := s.LinesRead(); lr != 0 {
		t.Error("an empty Stats structure should have linesRead: 0, has: ",
			lr)
	}
	if lp := s.LinesParsed(); lp != 0 {
		t.Error("an empty Stats structure should have linesParsed: 0, has: ",
			lp)
	}
	expectedStr := "files:   0   lines read:     0   parsed:     0"
	if s := s.String(); s != expectedStr {
		t.Error(
			"an empty Stats structure should have a String representation of: ",
			expectedStr, " has: ", s)
	}
}
