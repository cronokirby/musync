package internal

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"unicode"

	"github.com/BurntSushi/toml"
)

func readLine() string {
	r := bufio.NewReader(os.Stdin)
	line, err := r.ReadString('\n')
	if err != nil {
		panic(err)
	}
	if len(line) < 2 {
		return line
	}
	return line[:len(line)-1]
}

func promptForURL() string {
	fmt.Println("URL:")
	return readLine()
}

func promptForName() string {
	fmt.Println("Name:")
	return readLine()
}

func promptForArtist() string {
	fmt.Println("Artist:")
	return readLine()
}

func promptForPath() string {
	fmt.Println("Path:")
	return readLine()
}

func readYesNo() bool {
	answer := readLine()
	answer = strings.ToLower(answer)
	if answer == "y" || answer == "yes" {
		return true
	}
	return false
}

func readLineHandlingEOF(r *bufio.Reader) (string, error) {
	var bytes []byte
	for {
		line, isPrefix, err := r.ReadLine()
		if err != nil {
			return "", err
		}
		bytes = append(bytes, line...)
		if !isPrefix {
			break
		}
	}
	return string(bytes), nil
}

func readLinesUntilEOF() ([]string, error) {
	r := bufio.NewReader(os.Stdin)
	var lines []string
	for {
		line, err := readLineHandlingEOF(r)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return nil, err
			}
		}
		lines = append(lines, line)
	}
	return lines, nil
}

var nameRegexes = []string{
	`^\d+.\s([^\(]+)\s\(\d+:\d+\)`,
	`^\d+.\s([^\(]+)\s\d+:\d+`,
	`^\d+(?::\d+)+.+-\s+(.*)$`,
}
var timeRegexes = []string{
	`\((\d+:\d+)\)`,
	`(\d+:\d+)`,
	`(\d+(:\d+)+)`,
}

func tryRegexNumber(i int, lines []string) ([]string, []string) {
	nameRegex := regexp.MustCompile(nameRegexes[i])
	timeRegex := regexp.MustCompile(timeRegexes[i])
	var names []string
	var times []string
	for _, line := range lines {

		line = strings.TrimFunc(line, func(r rune) bool {
			return !unicode.IsGraphic(r)
		})

		nameMatches := nameRegex.FindStringSubmatch(line)
		if len(nameMatches) < 1 || nameMatches[1] == "" {
			return nil, nil
		}
		timeMatches := timeRegex.FindStringSubmatch(line)
		if len(timeMatches) < 1 || timeMatches[1] == "" {
			return nil, nil
		}
		names = append(names, nameMatches[1])
		times = append(times, timeMatches[1])
	}
	return names, times
}

// This will try and extract names and times from a pasting of the description
func getNamesAndTimesFromPaste(count int, lines []string) ([]string, []string) {
	if len(lines) != count {
		fmt.Printf("Expecting %d songs, but only found %d lines\n", count, len(lines))
		return nil, nil
	}
	for i := 0; i < len(nameRegexes); i++ {
		names, times := tryRegexNumber(i, lines)
		if names != nil && times != nil {
			fmt.Println("Extracted the following:")
			for j := 0; j < len(names); j++ {
				fmt.Printf("  %s (%s)\n", names[j], times[j])
			}
			fmt.Println("Does this look right to you? (Y/N)")
			looksRight := readYesNo()
			if looksRight {
				return names, times
			}
		}
	}
	return nil, nil
}

// Strategy: try and extract the names and times using regexes, and then just ask for
// each time and name if it's necessary
func promptForNamesAndTimes() ([]string, []string, error) {
	fmt.Println("How many songs are there?")
	var count int
	fmt.Scanf("%d", &count)
	var names []string
	var times []string
	fmt.Println("Would you like to try and extract names / times from a description? (Y/N)")
	shouldExtract := readYesNo()
	if shouldExtract {
		fmt.Println("Please paste that description, end with CTRL+D:")
		lines, err := readLinesUntilEOF()
		if err != nil {
			return nil, nil, err
		}
		names, times = getNamesAndTimesFromPaste(count, lines)
		if names != nil && times != nil {

			return names, times, nil
		}
		fmt.Println("We couldn't extract the timestamps using our builtin regexes")
	}
	fmt.Println("Let's enter each timestamp and song name then.")
	fmt.Println("Names:")
	names = make([]string, 0, count)
	for i := 0; i < count; i++ {
		name := readLine()
		names = append(names, name)
	}
	fmt.Println("Timestamps:")
	times = make([]string, 0, count)
	for i := 0; i < count; i++ {
		time := readLine()
		times = append(times, time)
	}
	return names, times, nil
}

// This will use the command line to ask the user for information in order to construct
// a new source
func promptForRawSource() (*rawSource, error) {
	url := promptForURL()
	name := promptForName()
	artist := promptForArtist()
	path := promptForPath()
	names, times, err := promptForNamesAndTimes()
	if err != nil {
		return nil, err
	}
	return &rawSource{
		Name:       name,
		Artist:     artist,
		Path:       path,
		URL:        url,
		Namestamps: names,
		Timestamps: times,
	}, nil
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
	newSource, err := promptForRawSource()
	if err != nil {
		return err
	}
	return writeRawSourceTo(path, newSource)
}
