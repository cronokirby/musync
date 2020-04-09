package internal

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

func promptForURL() string {
	fmt.Println("URL:")
	var url string
	fmt.Scanln(&url)
	return url
}

func promptForName() string {
	fmt.Println("Name:")
	var name string
	fmt.Scanln(&name)
	return name
}

func promptForArtist() string {
	fmt.Println("Artist:")
	var artist string
	fmt.Scanln(&artist)
	return artist
}

func promptForPath() string {
	fmt.Println("Path:")
	var path string
	fmt.Scanln(&path)
	return path
}

// Strategy: just ask for each time and name, make sure they're the same length
func promptForNamesAndTimes() ([]string, []string) {
	fmt.Println("Let's enter each timestamp and song name then.")
	fmt.Println("How many songs are there?")
	var count int
	fmt.Scanf("%d", &count)
	fmt.Println("Names:")
	names := make([]string, 0, count)
	for i := 0; i < count; i++ {
		var name string
		fmt.Scanf("%s", &name)
		names = append(names, name)
	}
	fmt.Println("Timestamps:")
	times := make([]string, 0, count)
	for i := 0; i < count; i++ {
		var time string
		fmt.Scanf("%s", &time)
		times = append(times, time)
	}
	return names, times
}

// This will use the command line to ask the user for information in order to construct
// a new source
func promptForRawSource() *rawSource {
	url := promptForURL()
	name := promptForName()
	artist := promptForArtist()
	path := promptForPath()
	names, times := promptForNamesAndTimes()
	return &rawSource{
		Name:       name,
		Artist:     artist,
		Path:       path,
		URL:        url,
		Namestamps: names,
		Timestamps: times,
	}
}

// This will append a raw source to a given file.
// If the file doesn't exist, it will create it with just this source.
func writeRawSourceTo(path string, source *rawSource) error {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("Couldn't open %s: %w", path, err)
	}
	if _, err := file.WriteString("\n"); err != nil {
		return fmt.Errorf("Couldn't write to %s: %w", path, err)
	}
	library := rawLibrary{Source: []rawSource{*source}}
	encoder := toml.NewEncoder(file)
	if err := encoder.Encode(library); err != nil {
		return fmt.Errorf("Couldn't encode or write new source: %w", err)
	}
	return nil
}

// PromptToAddSource will ask the user for information about a source, before adding it to path
func PromptToAddSource(path string) error {
	newSource := promptForRawSource()
	return writeRawSourceTo(path, newSource)
}
