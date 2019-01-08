/*

Package fileparse provides a standard mechanism for parsing a file. It
supports comments and include directives.

You first construct a file parser and set any options and a line parser
specific to the type of file you are working with and then call Parse to
process the file. Any errors found while parsing will be returned for further
processing.

*/
package fileparse
