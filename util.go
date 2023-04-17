package main

import (
	"archive/zip"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func Download(source, file string) error {
	log.Println("Downloading", source)

	out, err := os.Create(file)
	if err != nil {
		return err
	}

	resp, err := http.Get(source)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	resp.Body.Close()

	log.Println("Downloaded", file)

	return nil
}

func GetURLBody(url string) (string, error) {
	log.Println("Retrieving URL Body of", url)

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	resp.Body.Close()

	return string(body), nil
}

func UnzipFolder(source string, destDir string) error {
	log.Println("Extracting", source)

	zip, err := zip.OpenReader(source)
	if err != nil {
		return err
	}

	for _, file := range zip.File {
		// Roblox's Zip Files have windows paths in them
		filePath := filepath.Join(destDir, strings.ReplaceAll(file.Name, `\`, "/"))
		log.Println("Unzipping", filePath)

		if file.FileInfo().IsDir() {
			log.Println("Creating directory", filePath)

			if err := os.MkdirAll(filePath, file.Mode()); err != nil {
				return err
			}

			continue
		}

		destFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}

		fileZipped, err := file.Open()
		if err != nil {
			return err
		}

		if _, err := io.Copy(destFile, fileZipped); err != nil {
			return err
		}

		destFile.Close()
		fileZipped.Close()
	}

	zip.Close()

	return nil
}

// A MD5 file checksum verification function, used alongside
// Roblox's _provided_ MD5 checksums.
func VerifyFileMD5(filePath string, signature string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return err
	}

	if signature != hex.EncodeToString(hash.Sum(nil)) {
		return fmt.Errorf("file %s checksum mismatch: %x", filePath, hash.Sum(nil))
	}

	return nil
}

func NewestFile(glob string) string {
	// Since filepath.Glob sorts numerically, the 'newest' log files
	// will always be last (hence why retrieveing the last array element
	// is used), as they contain the date they were created at.
	// On-top of this, it also sorts alphabetically, so we only check for
	// log files that match the pattern.
	files, _ := filepath.Glob(glob)

	if len(files) < 1 {
		return ""
	}

	return files[len(files)-1]
}
