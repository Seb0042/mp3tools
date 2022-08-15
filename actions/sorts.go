package actions

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/bogem/id3v2"
)

func dedup(s []rune) string {
	// iterate over all characters in slice except the last one
	for i := 0; i < len(s)-1; i++ {
		// iterate over all the remaining characeters
		for j := i + 1; j < len(s); j++ {
			if s[i] != s[j] {
				break // this wasn't a duplicate, move on to the next char
			}
			// we found a duplicate!
			s = append(s[:i], s[j:]...)
		}
	}
	return string(s)
}

type SortWalker struct {
	PrefixPath string
	Rootdir    string
}

func (w *SortWalker) SortByGenre(s string, d fs.DirEntry, err error) error {

	libRegEx, err := regexp.Compile("^.+\\.(mp3)$")
	li := strings.LastIndex(w.PrefixPath, "/")
	var target string

	if li == len(w.PrefixPath)-1 {
		w.PrefixPath = w.PrefixPath[:len(w.PrefixPath)-1]
	}

	if err == nil && !d.IsDir() && libRegEx.MatchString(s) {
		tag, err := id3v2.Open(s, id3v2.Options{Parse: true})
		if err != nil {
			fmt.Println("Error while opening mp3 file: ", err)
		}
		filename := filepath.Base(s)
		artist := tag.Artist()
		if artist == "" {
			tag.Close()
			fmt.Printf("File without artist: %s\n", s)
			return errors.New("File without artist: " + s)
		}
		target = filepath.Join(w.Rootdir, tag.Genre())
		target = filepath.Join(target, tag.Artist())
		target = filepath.Join(target, filename)

		if target != s {
			fmt.Printf("Moving (%s) to (%s)\n", s, target)

			err := os.MkdirAll(filepath.Dir(target), os.ModePerm)
			if err != nil {
				fmt.Println(err.Error())
			}
			os.Rename(s, target)
		}
	}
	return nil
}
