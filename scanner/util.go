package scanner

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/Masterminds/semver"
)

func FindMatchingFiles(workingDir, pattern string) ([]string, error) {
	var matches []string
	err := filepath.Walk(workingDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if matched, err := filepath.Match(pattern, filepath.Base(path)); err != nil {
			return err
		} else if matched {
			matches = append(matches, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return matches, nil
}

func FindMatchingFolders(workingDir, pattern string) ([]string, error) {
	var matches []string
	err := filepath.Walk(workingDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			return nil
		}
		if matched, err := filepath.Match(pattern, filepath.Base(path)); err != nil {
			return err
		} else if matched {
			matches = append(matches, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return matches, nil
}

func ReadJSONFile(filePath string) (map[string]interface{}, error) {
	file, fileErr := os.Open(filePath)
	if fileErr != nil {
		return nil, fileErr
	}
	defer file.Close()
	myMap := map[string]interface{}{}
	decoder := json.NewDecoder(file)
	jsonErr := decoder.Decode(&myMap)
	if jsonErr != nil {
		return nil, jsonErr
	}
	return myMap, nil
}

func ReturnVersionObject(version string) (*semver.Version, error) {
	// strip weird stuff from the version slow for now
	version = strings.ReplaceAll(version, "^", "")
	v, err := semver.NewVersion(version)
	if err != nil {
		return nil, err
	}
	return v, nil
}
