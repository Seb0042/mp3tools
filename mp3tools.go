package main

import (
	"flag"
	"fmt"
	mp3toolsactions "mp3tools/actions"

	"os"
	"path/filepath"
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
	case "printinfo":
		filepath.WalkDir(*dirPtr, mp3toolsactions.PrintInfos)
	case "checkcom":
		filepath.WalkDir(*dirPtr, mp3toolsactions.CheckComments)
	default:
		fmt.Printf("Action %s unknown\n", *actionPtr)
	}
}
