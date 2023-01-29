package actions

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"regexp"

	"database/sql"

	"github.com/dhowden/tag"
	_ "github.com/mattn/go-sqlite3"
)

type FileInfos struct {
	FilePath string
	FileType tag.FileType
	Title    string
	Album    string
	Artist   string
	Composer string
	Genre    string
	Year     int
	FileSize int64
}

type FilesWalker struct {
	DB *sql.DB
}

func (w *FilesWalker) GetFiles(s string, d fs.DirEntry, err error) error {

	libRegEx, err := regexp.Compile("^.+\\.(mp3|flac)$")
	if err == nil && !d.IsDir() && libRegEx.MatchString(s) {

		f, err := os.Open(s)
		if err != nil {
			fmt.Printf("error loading file: %v =>%s\n", err, s)
			return errors.New("Open")
		}
		defer f.Close()
		m, err := tag.ReadFrom(f)
		if err != nil {
			fmt.Printf("error reading file: %v =>%s\n", err, s)
			//return errors.New("Read")
			return nil
		}
		var size int64
		fi, err := f.Stat()
		if err != nil {
			size = -1
		} else {
			size = fi.Size()
		}
		tx, err := w.DB.Begin()
		if err != nil {
			fmt.Println(err)
		}
		stmt, err := tx.Prepare("insert into files (filepath,filetype, title, album, artist,composer,genre, year, filesize)values (?,?,?,?,?,?,?,?,?)")
		if err != nil {
			fmt.Println(err)
		}
		defer stmt.Close()
		_, err = stmt.Exec(s, m.FileType(), m.Title(), m.Album(), m.Artist(), m.Composer(), m.Genre(), m.Year(), size)
		if err != nil {
			fmt.Println(err)
		}

		err = tx.Commit()
		if err != nil {
			fmt.Println(err, stmt)
		}

		/*	*w = append(*w, FileInfos{FilePath: s,
			FileType: m.FileType(),
			Title:    m.Title(),
			Album:    m.Album(),
			Artist:   m.Artist(),
			Composer: m.Composer(),
			Genre:    m.Genre(),
			Year:     m.Year(),
			FileSize: size})

		//printMetadata(m)*/
	}
	return nil
}
