package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	actions "mtools/internal"
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
		w := actions.FilesWalker{PrefixPath: *dirPtr, Rootdir: ""}
		filepath.WalkDir(*dirPtr, w.GetFiles)
	default:
		fmt.Printf("Action %s unknown\n", *actionPtr)
	}
}
