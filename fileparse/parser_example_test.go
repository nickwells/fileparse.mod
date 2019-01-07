package fileparse_test

import "github.com/nickwells/fileparse.mod/fileparse"

// ExampleParser1 this demonstrates how the parsing will scan through various
// files, ignoring comments and blank lines and following includes
func ExampleParser1() {
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

// ExampleParser2 this demonstrates how the parsing will scan through various
// files, ignoring comments and blank lines. The include keyword has been set
// to the empty string and so no include directives are recognised; hence
// they are not followed and the lines are not subject to any special
// handling.
func ExampleParser2() {
	p := fileparse.New("an example parser", fileparse.EchoParser{})
	p.SetInclKeyWord("")
	p.Parse("testdata/example/startingFile")
	// Output: first content line
	// second content line
	// #include includeFile1
	// #include includeFile2
	// last line
}

// ExampleParser3 this demonstrates how the parsing will scan through various
// files, ignoring blank lines. The include keyword and the comment intro
// have been set to the empty string and so no include directives are
// recognised and no text is treated as a comment.
func ExampleParser3() {
	p := fileparse.New("an example parser", fileparse.EchoParser{})
	p.SetInclKeyWord("")
	p.SetCommentIntro("")
	p.Parse("testdata/example/startingFile")
	// Output: // This is the first file read by the example
	// first content line
	// second content line
	// #include includeFile1
	// #include includeFile2
	// last line
}
