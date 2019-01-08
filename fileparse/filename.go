package fileparse

import (
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

// FixFileName first cleans the name (using filepath.Clean) and then, if it
// starts with a '~', it will expand the '~' into the home directory of the
// appropriate user
func FixFileName(fileName string) (string, error) {
	var err error
	fixedFileName := filepath.Clean(fileName)

	if fixedFileName[0] == '~' {
		pathParts := strings.Split(fixedFileName, string(os.PathSeparator))

		pathParts[0], err = ExpandTilde(pathParts[0])
		fixedFileName = filepath.Join(pathParts...)
	}
	return fixedFileName, err
}

// ExpandTilde will take the tilde-prefixed name and convert it into the
// user's home directory. As in the Unix bash shell a '~' on its own
// represents the home directory of the user running the code and with a
// following name it represents the home directory of the named user.
func ExpandTilde(tildeStr string) (string, error) {
	var usr *user.User
	var err error

	if tildeStr == "~" {
		usr, err = user.Current()
	} else {
		userName := tildeStr[1:]
		usr, err = user.Lookup(userName)
	}
	if err != nil {
		return "", err
	}
	return usr.HomeDir, err
}
