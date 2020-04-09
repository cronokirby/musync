package main

import (
	"fmt"
	"io/ioutil"

	"github.com/cronokirby/musync/internal"

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
	data, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Printf("Failed to read sync file: %v\n", err)
		return
	}
	lib, err := internal.LoadLibrary(string(data))
	if err != nil {
		fmt.Printf("Failed to parse library: %v\n", err)
		return
	}
	for _, v := range lib.Sources {
		fmt.Printf("Downloading '%s'\n", v.Name)
		if err := internal.Download(out, &v); err != nil {
			fmt.Printf("Error while downloading '%s': %v\n", v.Name, err)
			return
		}
	}
}

func add(path string) {
	if err := internal.PromptToAddSource(path); err != nil {
		fmt.Printf("Failed to add source: %v", err)
	}
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
