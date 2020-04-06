package main

import (
	"fmt"

	"github.com/alecthomas/kong"
)

var cli struct {
	Sync struct {
		Root string `name:"root" help:"The directory to sync files to" type:"existingdir" default:"."`
		Path string `arg:"" name:"path" help:"The file describing what to sync" type:"existingfile"`
	} `cmd:"" help:"Sync music files"`
	Add struct {
		Path string `arg:"" name:"path" help:"The file containing the files to sync" type:"existingfile"`
	} `cmd:"" help:"Add a new album or song to the syncing file"`
}

func sync(root string, path string) {
	fmt.Printf("Sync, root: %s, path: %s\n", root, path)
}

func add(path string) {
	fmt.Printf("Add, path: %s\n", path)
}

func main() {
	ctx := kong.Parse(&cli)
	switch ctx.Command() {
	case "sync <path>":
		sync(cli.Sync.Root, cli.Sync.Path)
	case "add <path>":
		add(cli.Add.Path)
	}
}
