package fileparse

import (
	"fmt"
	"testing"

	"github.com/nickwells/testhelper.mod/v2/testhelper"
)

func TestIsAnInclLine(t *testing.T) {
	var np NullParser
	fpNull := New("intro", np)

	type inclTest struct {
		line        string
		expFileName string
		expHasIncl  bool
	}

	testCases1 := []inclTest{
		{"has no include directive", "", false},
		{"@include ", "", true},
		{"@include xxx ", "xxx", true},
	}

	for _, tc := range testCases1 {
		id := fmt.Sprintf("isAnInclLine(%q)", tc.line)
		filename, hasIncl := fpNull.isAnInclLine(tc.line)

		testhelper.DiffString(t, id, "filename", filename, tc.expFileName)
		testhelper.DiffBool(t, id, "hasIncl", hasIncl, tc.expHasIncl)
	}

	fpNull.SetInclKeyWord("INCLUDE")

	testCases2 := []inclTest{
		{"has no include directive", "", false},
		{"INCLUDE ", "", true},
		{"INCLUDE xxx ", "xxx", true},
	}

	for _, tc := range testCases2 {
		id := fmt.Sprintf("isAnInclLine(%q) - include keyword: 'INCLUDE'",
			tc.line)
		filename, hasIncl := fpNull.isAnInclLine(tc.line)

		testhelper.DiffString(t, id, "filename", filename, tc.expFileName)
		testhelper.DiffBool(t, id, "hasIncl", hasIncl, tc.expHasIncl)
	}
}

func TestStripComment(t *testing.T) {
	var np NullParser
	fpNull := New("intro", np)

	type commentTest struct {
		line    string
		expLine string
	}

	testCases1 := []commentTest{
		{"abc # test", "abc "},
		{" # test", " "},
		{"   ", "   "},
		{"abc ", "abc "},
	}

	for _, tc := range testCases1 {
		id := fmt.Sprintf("stripComment(%q)", tc.line)
		stripped := fpNull.stripComment(tc.line)

		testhelper.DiffString(t, id, "stripped line", stripped, tc.expLine)
	}

	fpNull.SetCommentIntro("//")

	testCases2 := []commentTest{
		{"abc # test", "abc # test"},
		{" # test", " # test"},
		{"abc // test", "abc "},
		{" // test", " "},
		{"   ", "   "},
		{"abc ", "abc "},
	}

	for _, tc := range testCases2 {
		id := fmt.Sprintf("stripComment(%q) - comment intro: '#'", tc.line)
		stripped := fpNull.stripComment(tc.line)

		testhelper.DiffString(t, id, "stripped line", stripped, tc.expLine)
	}
}
