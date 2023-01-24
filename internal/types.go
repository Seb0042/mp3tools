package actions

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"regexp"

	"github.com/dhowden/tag"
)

type FilesWalker struct {
	PrefixPath string
	Rootdir    string
	Genre      string
	Level      string
}

func (w *FilesWalker) GetFiles(s string, d fs.DirEntry, err error) error {

	libRegEx, err := regexp.Compile("^.+\\.(mp3|flac)$")
	if err == nil && !d.IsDir() && libRegEx.MatchString(s) {

		f, err := os.Open(s)
		if err != nil {
			fmt.Printf("error loading file: %v", err)
			return errors.New("Open")
		}
		defer f.Close()
		m, err := tag.ReadFrom(f)
		if err != nil {
			fmt.Printf("error reading file: %v\n", err)
			return errors.New("Read")
		}
		fmt.Printf("filepath: %s\n", s)
		printMetadata(m)
	}
	return nil
}
