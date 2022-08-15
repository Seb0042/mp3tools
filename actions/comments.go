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
					fmt.Printf("On file: %s, comment not allowed: %s\n", s, comment.Text)
				}
			}
		}

		tag.Close()

	}
	return nil
}

func CleanComments(s string, d fs.DirEntry, err error) error {

	commentsAllowed := map[string]bool{}
	// var savedFrame id3v2.CommentFrame
	savedFrame := id3v2.CommentFrame{}
	nbErrors := 0
	bck := false

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
		i := 0
		for _, f := range comments {
			comment, ok := f.(id3v2.CommentFrame)
			if !ok {
				fmt.Println("Couldn't assert comment frame")
			} else {
				_, cok := commentsAllowed[comment.Text]
				if !cok {
					nbErrors++
				} else {
					bck = true
					savedFrame = comment
				}
			}
			i++
		}
		if nbErrors != 0 {
			fmt.Printf("Wrong values for comment. ")
			tag.DeleteFrames(tag.CommonID("Comments"))
			if bck {
				tag.AddCommentFrame(savedFrame)
			} else {
				i = 0
			}
		}

		if i == 0 {
			fmt.Printf("No comment. Adding None. ")
			savedFrame.Encoding = id3v2.EncodingUTF8
			savedFrame.Language = "eng"
			savedFrame.Text = "None"
			tag.AddCommentFrame(savedFrame)
		}
		if !(nbErrors == 0 && i == 1) {
			fmt.Printf("Saving modifications to: %s\n", s)
			tag.Save()
		}

		tag.Close()

	}
	return nil
}
