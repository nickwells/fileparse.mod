package fileparse

import "testing"

func TestIsAnInclLine(t *testing.T) {
	var np NullParser
	fpNull := New("intro", np)

	type inclTest struct {
		line             string
		expectedFileName string
		expectedHasIncl  bool
	}
	testCases1 := []inclTest{
		{"has no include directive", "", false},
		{"#include ", "", true},
		{"#include xxx ", "xxx", true}}
	for _, it := range testCases1 {
		filename, hasIncl := fpNull.isAnInclLine(it.line)

		if filename != it.expectedFileName {
			t.Error("isAnInclLine(", it.line, ") failed\n",
				"expected filename: ", it.expectedFileName, "\n",
				"got: ", filename)
		}
		if hasIncl != it.expectedHasIncl {
			t.Error("isAnInclLine(", it.line, ") failed\n",
				"expected hasIncl: ", it.expectedHasIncl, "\n",
				"got: ", hasIncl)
		}
	}
	fpNull.SetInclKeyWord("INCLUDE")
	testCases2 := []inclTest{
		{"has no include directive", "", false},
		{"INCLUDE ", "", true},
		{"INCLUDE xxx ", "xxx", true}}
	for _, it := range testCases2 {
		filename, hasIncl := fpNull.isAnInclLine(it.line)

		if filename != it.expectedFileName {
			t.Error("isAnInclLine(", it.line, ") failed\n",
				"expected filename: ", it.expectedFileName, "\n",
				"got: ", filename)
		}
		if hasIncl != it.expectedHasIncl {
			t.Error("isAnInclLine(", it.line, ") failed\n",
				"expected hasIncl: ", it.expectedHasIncl, "\n",
				"got: ", hasIncl)
		}
	}
}

func TestStripComment(t *testing.T) {
	var np NullParser
	fpNull := New("intro", np)

	type commentTest struct {
		line         string
		expectedLine string
	}

	testCases1 := []commentTest{
		{"abc // test", "abc"},
		{" // test", ""},
		{"   ", ""},
		{"abc ", "abc"},
	}

	for _, ct := range testCases1 {
		strippedLine := fpNull.stripComment(ct.line)

		if strippedLine != ct.expectedLine {
			t.Error("stripComment(", ct.line, ") failed\n",
				"expected the stripped line to be: '", ct.expectedLine, "'\n",
				"got: '", strippedLine, "'")
		}
	}

	fpNull.SetCommentIntro("#")
	testCases2 := []commentTest{
		{"abc // test", "abc // test"},
		{" // test", "// test"},
		{"abc # test", "abc"},
		{" # test", ""},
		{"   ", ""},
		{"abc ", "abc"},
	}

	for _, ct := range testCases2 {
		strippedLine := fpNull.stripComment(ct.line)

		if strippedLine != ct.expectedLine {
			t.Error("stripComment(", ct.line, ") failed\n",
				"expected the stripped line to be: '", ct.expectedLine, "'\n",
				"got: '", strippedLine, "'")
		}
	}
}
