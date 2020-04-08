package internal

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/bogem/id3v2"
)

// Downloads the url as a big media file
func youtubeDL(url string, outPath string) error {
	cmd := exec.Command(
		"youtube-dl",
		url,
		"-x",
		"--write-thumbnail",
		"--audio-format",
		"m4a",
		"-o", outPath,
	)
	// We want to capture the stderr inside of the error, hence Output instead of Run
	if _, err := cmd.Output(); err != nil {
		return err
	}
	return nil
}

func convertToMP3(m4aPath string, mp3Path string) error {
	cmd := exec.Command(
		"ffmpeg",
		"-i", m4aPath,
		"-loglevel", "panic",
		"-y",
		"-vn",
		"-ar", "44100",
		"-ac", "2",
		"-ab", "192k",
		"-id3v2_version", "3",
		"-f", "mp3",
		mp3Path,
	)
	if _, err := cmd.Output(); err != nil {
		return err
	}
	return nil
}

func pathExists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	}
	return false
}

func fillSectionMetadata(base string, source *Source, section *Section, sectionIndex int) error {
	path := source.SectionPath(base, section)
	tags, err := id3v2.Open(path, id3v2.Options{Parse: true})
	if err != nil {
		return err
	}
	tags.SetTitle(section.Name)
	tags.SetAlbum(source.Name)
	tags.SetArtist(source.Artist)
	tags.AddTextFrame(tags.CommonID("Track number/Position in set"), tags.DefaultEncoding(), fmt.Sprintf("%d/%d", sectionIndex, len(source.Sections)))
	artwork, err := ioutil.ReadFile(source.CoverArtPath(base))
	if err != nil {
		return err
	}
	pic := id3v2.PictureFrame{
		Encoding:    id3v2.EncodingUTF8,
		MimeType:    "image/jpeg",
		PictureType: id3v2.PTFrontCover,
		Description: "Front cover",
		Picture:     artwork,
	}
	tags.AddAttachedPicture(pic)
	if err := tags.Save(); err != nil {
		return err
	}
	return nil
}

func createSection(base string, source *Source, section *Section, sectionIndex int) error {
	toSplit := source.MP3Path(base)
	outputPath := source.SectionPath(base, section)
	args := []string{
		"-i", toSplit,
		"-y",
		"-acodec", "copy",
		"-id3v2_version", "3",
		"-ss", section.StartTime,
	}
	if section.HasEnd {
		args = append(args, "-to", section.EndTime, outputPath)
	} else {
		args = append(args, outputPath)
	}
	cmd := exec.Command("ffmpeg", args...)
	if _, err := cmd.Output(); err != nil {
		return err
	}
	if err := fillSectionMetadata(base, source, section, sectionIndex); err != nil {
		return err
	}
	return nil
}

func createAlbumDirectory(base string, source *Source) error {
	toCreate := source.SectionDirectory(base)
	if err := os.MkdirAll(toCreate, os.ModePerm); err != nil {
		return err
	}
	return nil
}

func downloadSource(basePath string, source *Source) (bool, error) {
	m4aPath := source.M4APath(basePath)
	mp3Path := source.MP3Path(basePath)
	hasBeenDownloaded := pathExists(mp3Path)
	hasBeenSplit := pathExists(source.DirectoryPath(basePath))
	if hasBeenDownloaded || hasBeenSplit {
		return false, nil
	}
	// Now we actually download
	if err := youtubeDL(source.URL, m4aPath); err != nil {
		return false, err
	}
	if err := convertToMP3(m4aPath, mp3Path); err != nil {
		return false, err
	}
	if err := os.Remove(m4aPath); err != nil {
		return false, err
	}
	return true, nil
}

func splitSource(basePath string, source *Source) (bool, error) {
	hasBeenSplit := pathExists(source.DirectoryPath(basePath))
	if hasBeenSplit {
		return false, nil
	}
	if err := createAlbumDirectory(basePath, source); err != nil {
		return false, err
	}
	for i, v := range source.Sections {
		if err := createSection(basePath, source, &v, i); err != nil {
			return false, err
		}
	}
	return true, nil
}

// Download will take a basePath, and download a single source into a big mp3
// This will also download the cover art as well. We will then split the big mp3
// into the sections that compose it, annotating it with metadata
func Download(basePath string, source *Source) error {
	downloaded, err := downloadSource(basePath, source)
	if err != nil {
		return err
	}
	if !downloaded {
		fmt.Println("  Already Downloaded")
	}
	fmt.Printf("  Splitting '%s'\n", source.Name)
	split, err := splitSource(basePath, source)
	if err != nil {
		return err
	}
	if !split {
		fmt.Println("  Already Split")
	}
	// We can get rid of the downloaded mp3 now that we've split it
	downloadedPath := source.MP3Path(basePath)
	if err := os.Remove(downloadedPath); err != nil {
		return err
	}
	return nil
}
