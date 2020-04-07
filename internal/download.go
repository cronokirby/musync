package internal

import (
	"os"
	"os/exec"
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
	if err := cmd.Run(); err != nil {
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
	if err := cmd.Run(); err != nil {
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

// Download will take a basePath, and download a single source into a big mp3
// This will also download the cover art as well
func Download(basePath string, source *Source) (bool, error) {
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
