package actions

import (
	"fmt"
	"io/fs"
	"regexp"

	"github.com/bogem/id3v2"
)

func CheckComments(s string, d fs.DirEntry, err error) error {

	commentsAllowed := map[string]bool{}

	commentsAllowed["None"] = true
	commentsAllowed["Mix"] = true
	commentsAllowed["Listen"] = true
	commentsAllowed["Radio"] = true
	commentsAllowed["Not"] = true
	libRegEx, err := regexp.Compile("^.+\\.(mp3)$")
	if err == nil && !d.IsDir() && libRegEx.MatchString(s) {
		tag, err := id3v2.Open(s, id3v2.Options{Parse: true})
		if err != nil {
			fmt.Println("Error while opening mp3 file: ", err)
		}
		comments := tag.GetFrames(tag.CommonID("Comments"))
		for _, f := range comments {
			comment, ok := f.(id3v2.CommentFrame)
			if !ok {
				fmt.Println("Couldn't assert comment frame")
			} else {

				_, cok := commentsAllowed[comment.Text]
				if !cok {
					fmt.Printf("Not allowed comment: %s, on file %s\n", comment.Text, s)
				}
			}
		}

		tag.Close()

	}
	return nil
}
