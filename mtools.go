package main

import (
	"database/sql"
	"flag"
	"fmt"
	actions "mtools/internal"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

func main() {

	dirPtr := flag.String("dir", ".", "directory to search files")
	actionPtr := flag.String("action", "printinfo", "what to do")
	flag.Parse()
	_, err := os.Stat(*dirPtr)
	if err != nil {
		fmt.Printf("Dir %s not found\n", *dirPtr)
		os.Exit(-1)
	}
	switch *actionPtr {
	case "dupes":

		os.Remove("./mtools_dupes.db")
		db, err := sql.Open("sqlite3", "./mtools_dupes.db")
		if err != nil {
			fmt.Println(err)
		}
		defer db.Close()
		sqlStmt := `create table files(id integer not null primary key, filepath text,filetype string, title text, album text, artist text,composer text,genre string, year int, filesize int64);`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			fmt.Println(err)
		}
		w := &actions.FilesWalker{DB: db}
		filepath.WalkDir(*dirPtr, w.GetFiles)
		sqlStmt = `select a.id, a.Artist,a.Title,a.Album,a.Year,a.Filetype,a.filesize,a.Filepath from files a,files b where a.title = b.title and b.filepath != a.filepath;`
		rows, err := db.Query(sqlStmt)
		if err != nil {
			fmt.Println(err, sqlStmt)
		}
		defer rows.Close()
		var id int
		var FilePath string
		var FileType string
		var Title string
		var Album string
		var Artist string
		var Year int
		var FileSize int64

		for rows.Next() {
			err = rows.Scan(&id, &Artist, &Title, &Album, &Year, &FileType, &FileSize, &FilePath)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Printf("|%6d|%-30s|%-30s|%30s|%4d|%6s|%12d|%s|\n", id, Artist, Title, Album, Year, FileType, FileSize, FilePath)
		}
		err = rows.Err()
		if err != nil {
			fmt.Println(err)
		}
	default:
		fmt.Printf("Action %s unknown\n", *actionPtr)
	}
}
