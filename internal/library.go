package internal

import (
	"fmt"
	"path"

	"github.com/BurntSushi/toml"
)

// Section represents a single part of a given source
// If we're splitting albums into multiple songs, each of these songs would have
// one of these sections
type Section struct {
	Name string
	// If this is not set, then implicitly this section goes to the end of the source
	HasEnd bool
	// The timestamp where this section starts
	StartTime string
	// The timestamp where this section ends
	EndTime string
}

// Source represents a single item in our library
// A source contains information about where to download some media, along with
// metadata about that media, and how to split up the media into multiple parts
type Source struct {
	// The name of this item
	Name string
	// The name of the artist that made this item
	Artist string
	// The nested path where we want this source to go to
	Path string
	// The URL we can use to download this source
	URL string
	// A list of sections that we need to split this source into
	Sections []Section
}

// DirectoryPath creats the path where we'd store the invididual sources
func (s *Source) DirectoryPath(base string) string {
	return path.Join(base, s.Path, s.Name)
}

// MP3Path creates the path where we'd save the entire source media
func (s *Source) MP3Path(base string) string {
	return path.Join(base, s.Path, s.Name+".mp3")
}

// M4APath creates the path where we'd save this full media with an m4a extension
// This is useful because we download using youtube-dl to this path, and then convert
// using another path
func (s *Source) M4APath(base string) string {
	return path.Join(base, s.Path, s.Name+".m4a")
}

// CoverArtPath returns the path to the cover art
func (s *Source) CoverArtPath(base string) string {
	return path.Join(base, s.Path, s.Name+".jpg")
}

// SectionDirectory returns the path where section files would be stored
func (s *Source) SectionDirectory(base string) string {
	return path.Join(base, s.Path, s.Name)
}

// SectionPath returns the path to a section file contained in this Source
// For example, if an album contains a specific song, this will give us the final path
// of that song, contained in the album directory
func (s *Source) SectionPath(base string, section *Section) string {
	return path.Join(s.SectionDirectory(base), section.Name+".mp3")
}

// Library represents a collection of sources we want to gather and download
type Library struct {
	// The individual sources that make up this library
	Sources []Source
}

// These types are used for automatic toml decoding
// So here we just have a timestamp, and then zip multiple sections together to
// get a start and end time
type rawSection struct {
	Name      string
	Timestamp string
}

type rawSource struct {
	Name       string
	Artist     string
	Path       string
	URL        string
	Namestamps []string
	Timestamps []string
}

type rawLibrary struct {
	Source []rawSource
}

// LoadLibrary parses a toml document to a library structure
func LoadLibrary(document string) (*Library, error) {
	library := Library{}
	rawLibrary := rawLibrary{}
	_, err := toml.Decode(document, &rawLibrary)
	if err != nil {
		return nil, err
	}
	library.Sources = make([]Source, 0, len(rawLibrary.Source))
	for _, rawSource := range rawLibrary.Source {
		namestampCount := len(rawSource.Namestamps)
		timestampCount := len(rawSource.Timestamps)
		if namestampCount != timestampCount {
			return nil, fmt.Errorf(
				"In the source named '%s': The number of names (%d) does not match the number of timestamps (%d)",
				rawSource.Name,
				namestampCount,
				timestampCount,
			)
		}
		newSections := make([]Section, 0, namestampCount)
		for i := 0; i < namestampCount; i++ {
			section := Section{
				Name:      rawSource.Namestamps[i],
				StartTime: rawSource.Timestamps[i],
			}
			if i+1 < namestampCount {
				section.HasEnd = true
				section.EndTime = rawSource.Timestamps[i+1]
			}
			newSections = append(newSections, section)
		}

		source := Source{
			Name:     rawSource.Name,
			Artist:   rawSource.Artist,
			Path:     rawSource.Path,
			URL:      rawSource.URL,
			Sections: newSections,
		}
		library.Sources = append(library.Sources, source)
	}
	return &library, nil
}
