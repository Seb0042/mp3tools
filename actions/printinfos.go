package actions

import (
	"fmt"
	"io/fs"
	"regexp"

	"github.com/bogem/id3v2"
)

func PrintInfos(s string, d fs.DirEntry, err error) error {
	libRegEx, err := regexp.Compile("^.+\\.(mp3)$")
	if err == nil && !d.IsDir() && libRegEx.MatchString(s) {
		tag, err := id3v2.Open(s, id3v2.Options{Parse: true})
		if err != nil {
			fmt.Println("Error while opening mp3 file: ", err)
		}
		fmt.Printf("File:%s .", s)
		// Read tags.
		fmt.Printf("Genre: %s, Title: %s, Artist: %s\n", tag.Genre(), tag.Title(), tag.Artist())
		fmt.Println("Comments")
		comments := tag.GetFrames(tag.CommonID("Comments"))
		for _, f := range comments {
			comment, ok := f.(id3v2.CommentFrame)
			if !ok {
				fmt.Println("Couldn't assert comment frame")
			} else {

				fmt.Printf("\t%s\n", comment.Text)
			}
		}

		tag.Close()

	}
	return nil
}
