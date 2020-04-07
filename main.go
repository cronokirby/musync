package main

import (
	"fmt"

	"github.com/alecthomas/kong"
)

var cli struct {
	Sync struct {
		Out  string `help:"The directory to sync files to" type:"existingdir" default:"."`
		Path string `arg:"" name:"path" help:"The file describing what to sync" type:"existingfile"`
	} `cmd:"" help:"Sync music files"`
	Add struct {
		Path string `arg:"" name:"path" help:"The file containing the files to sync" type:"existingfile"`
	} `cmd:"" help:"Add a new album or song to the syncing file"`
}

func sync(out string, path string) {
	fmt.Printf("Sync, out: %s, path: %s\n", out, path)
}

func add(path string) {
	fmt.Printf("Add, path: %s\n", path)
}

func main() {
	ctx := kong.Parse(&cli)
	switch ctx.Command() {
	case "sync <path>":
		sync(cli.Sync.Out, cli.Sync.Path)
	case "add <path>":
		add(cli.Add.Path)
	}
}
