package fileparse

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/nickwells/location.mod/location"
)

// DefaultInclKeyword is the value which introduces the name of a file to be
// read and substituted into the current file
//
// DefaultCommentIntro is the default comment introducer - everything from
// this to the end of the line is ignored
const (
	DefaultInclKeyWord  string = "#include"
	DefaultCommentIntro string = "//"
)

// FP records the configuration of the file parser
type FP struct {
	fileType   string
	lineParser LineParser

	cmtIntro    string
	inclKeyWord string

	passBlanks bool
	dontTrim   bool

	stats Stats
}

// Stats returns the latest statistics for the FP
func (fp FP) Stats() Stats { return fp.stats }

// New initialises a file parser with the default comment characters and
// include keyword and the passed LineParser. The desc is used in error
// messages to identify the type of file being parsed
func New(desc string, lp LineParser) *FP {
	return &FP{
		fileType:    desc,
		lineParser:  lp,
		cmtIntro:    DefaultCommentIntro,
		inclKeyWord: DefaultInclKeyWord}
}

// PassBlankLines will change the behaviour of the Parser so that blank lines
// are passed to the LineParser
func (fp *FP) PassBlankLines() {
	fp.passBlanks = true
}

// DontTrimWhiteSpace will change the behaviour of the Parser so that
// white-space is no longer trimmed from the front and back of lines passed
// to the LineParser
func (fp *FP) DontTrimWhiteSpace() {
	fp.dontTrim = true
}

// SetCommentIntro changes the comment introducer from the default value. A
// comment is taken to run from the start of the comment introducer to the
// end of the line. Setting the comment introducer to the empty string will
// mean that comments are ignored, though whitespace will still be trimmed.
func (fp *FP) SetCommentIntro(cmtIntro string) {
	fp.cmtIntro = cmtIntro
}

// SetInclKeyWord changes the include keyword from the default value. Setting
// the include keyword to the empty string will turn off the include file
// mechanism
func (fp *FP) SetInclKeyWord(incl string) {
	fp.inclKeyWord = incl
}

// stripComment will remove comments. All text from the start of a comment as
// given by the comment intro to the end of the line is removed.
func (fp FP) stripComment(s string) string {
	if fp.cmtIntro == "" {
		return s
	}

	parts := strings.SplitN(s, fp.cmtIntro, 2)
	return parts[0]
}

// trimSpace will remove any white space from the beginning or end of the
// line
func (fp FP) trimSpace(s string) string {
	if fp.dontTrim {
		return s
	}
	return strings.TrimSpace(s)
}

// isAnInclLine returns the include file name and a bool indicating whether
// this is an include line or not.
//
// A line is an include line if the file parser's include keyword is
// non-empty and the line passed has the keyword as a prefix. If it is an
// include line then the include file name is everything after the keyword
// with any surrounding whitespace stripped off.
func (fp FP) isAnInclLine(line string) (inclFileName string, hasIncl bool) {
	inclFileName = ""
	hasIncl = fp.inclKeyWord != "" && strings.HasPrefix(line, fp.inclKeyWord)

	if hasIncl {
		inclFileName = strings.TrimPrefix(line, fp.inclKeyWord)
		inclFileName = strings.TrimSpace(inclFileName)
	}
	return inclFileName, hasIncl
}

// Parse will read the passed file, following include directives (and
// checking for loops). It will strip out any blank lines and comments, strip
// any white space from the front and back of the line and call the
// LineParser on any remaining text. It is the responsibility of the
// LineParser to perform any operations resulting from the parsed lines.
//
// This default behaviour can be changed by setting the include prefix, the
// comment introducer and the blank line and white-space trimming behaviour.
//
// Any errors detected will be returned. Note that more than one error is
// possible.
func (fp *FP) Parse(filename string) []error {
	fp.stats = Stats{} // reset the stats each time we parse
	inclChain := location.NewChain()
	return fp.parseFile(filename, inclChain)
}

// fixIncludeFileName returns the include file name with the directory of the
// current file prepended if it is not an absolute pathename (starts with a
// '/')
func fixIncludeFileName(inclFileName, currentFileName string) string {
	if filepath.IsAbs(inclFileName) {
		return inclFileName
	}
	return filepath.Join(filepath.Dir(currentFileName), inclFileName)
}

// noteStr returns a formatted note string for setting the location note
func (fp *FP) noteStr(inclChain location.LocChain) string {
	note := fp.fileType
	if s := inclChain.String(); s != "" {
		note += " : " + s
	}
	return note
}

// openFile normalises the file name and checks the include chain for
// loops. Then it creates a new location object and opens the file. Any
// errors are returned.
func (fp *FP) openFile(filename string, inclChain location.LocChain) (
	*os.File, *location.L, error) {
	fixedFileName, err := FixFileName(filename)
	if err != nil {
		return nil, nil, fmt.Errorf("%s: Couldn't expand: '%s' : %s",
			fp.noteStr(inclChain), filename, err.Error())
	}

	loopFound, loopMsg := inclChain.HasLoop(fixedFileName)
	if loopFound {
		return nil, nil,
			fmt.Errorf("loop found: '%s' has been visited before: %s",
				fixedFileName, loopMsg)
	}

	fd, err := os.Open(fixedFileName) // nolint: gosec
	if err != nil {
		return nil, nil, err
	}

	loc := location.New(fixedFileName)
	loc.SetNote(fp.noteStr(inclChain))

	return fd, loc, nil
}

// parseFile is called by Parse once the stats and the include chain have
// been established. It is then called recursively to parse the file and any
// included files.
func (fp *FP) parseFile(filename string, inclChain location.LocChain) []error {
	var errors = make([]error, 0)

	fd, loc, err := fp.openFile(filename, inclChain)
	if err != nil {
		return append(errors, err)
	}
	defer fd.Close() // nolint: errcheck

	fp.stats.filesVisited++
	scanner := bufio.NewScanner(fd)

	for scanner.Scan() {
		fp.stats.linesRead++
		originalLine := scanner.Text()
		loc.Incr()

		line := fp.stripComment(originalLine)

		line = fp.trimSpace(line)

		if line == "" && !fp.passBlanks {
			continue // ignore blank lines
		}

		inclFileName, hasIncl := fp.isAnInclLine(line)
		if hasIncl {
			if inclFileName == "" {
				loc.SetContent(originalLine)
				errors = append(errors, loc.Errorf("Missing include file name"))
				continue
			}

			inclFileName = fixIncludeFileName(inclFileName, loc.Source())

			errors = append(errors,
				fp.parseFile(inclFileName, append(inclChain, *loc))...)
			continue
		}

		fp.stats.linesParsed++
		if err = fp.lineParser.ParseLine(line, loc); err != nil {
			errors = append(errors, err)
		}
	}

	if err = scanner.Err(); err != nil {
		errors = append(errors, err)
	}

	return errors
}
