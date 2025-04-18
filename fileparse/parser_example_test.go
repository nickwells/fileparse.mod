package fileparse_test

import (
	"fmt"

	"github.com/nickwells/fileparse.mod/fileparse"
	"github.com/nickwells/location.mod/location"
)

// ExampleFP_Parse demonstrates how the parsing will scan through various
// files, ignoring comments and blank lines and following includes
func ExampleFP_Parse() {
	p := fileparse.New("an example parser", fileparse.EchoParser{})
	p.Parse("testdata/example/startingFile")
	// Output: first content line
	// second content line
	// includefile line 1
	// includefile line 2
	// includefile line 3
	// includeFile3 line 1
	// includeFile3 line 2
	// includeFile2 line 1
	// last line
}

// ExampleFP_SetInclKeyWord demonstrates how the parsing will scan through
// various files, ignoring comments and blank lines. The include keyword has
// been set to the empty string and so no include directives are recognised;
// hence they are not followed and the lines are not subject to any special
// handling.
func ExampleFP_SetInclKeyWord() {
	p := fileparse.New("an example parser", fileparse.EchoParser{})
	p.SetInclKeyWord("")
	p.Parse("testdata/example/startingFile")
	// Output: first content line
	// second content line
	// @include includeFile1
	// @include includeFile2
	// last line
}

// ExampleFP_SetCommentIntro demonstrates how the parsing will scan through
// various files, ignoring blank lines. The include keyword and the comment
// intro have been set to the empty string and so no include directives are
// recognised and no text is treated as a comment.
func ExampleFP_SetCommentIntro() {
	p := fileparse.New("an example parser", fileparse.EchoParser{})
	p.SetInclKeyWord("")
	p.SetCommentIntro("")
	p.Parse("testdata/example/startingFile")
	// Output: # This is the first file read by the example
	// first content line
	// second content line
	// @include includeFile1
	// @include includeFile2
	// last line
}

// ExampleEchoParser demonstrates the behaviour of the EchoParser which just
// echos the lines it parses.
func ExampleEchoParser() {
	var ep fileparse.EchoParser

	loc := location.New("testLoc")

	err := ep.ParseLine("line", loc)
	if err != nil {
		fmt.Println("unexpected error :", err)
	}
	// Output: line
}
