package main

import (
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
)

type file struct {
	Name      string    `json:"name"`
	IsDir     bool      `json:"isDir"`
	Files     []file    `json:"files"`
	Subtitles subtitles `json:"subtitles"`
}

type subtitles map[string]string

// ListMediaAndSubtitlesHandler makes a JSON response of nested media and subtitles files
func ListMediaAndSubtitlesHandler(w http.ResponseWriter, req *http.Request) {
	movies, err := scanDirectory(GetEnv("MEDIA_PATH", "./"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse(w, movies, http.StatusOK)
}

func scanDirectory(dirName string) (files []file, err error) {
	files = []file{}

	content, err := listDirectoryContent(dirName)
	if err != nil {
		return
	}

	for _, c := range content {
		f := file{Name: c.Name()}
		if c.IsDir() {
			subDirName := dirName + f.Name + "/"

			f.IsDir = true
			f.Files, err = scanDirectory(subDirName)
			if err != nil {
				return
			}
			f.Subtitles, err = scanForSubtitles(subDirName)
			if err != nil {
				return
			}

			files = append(files, f)
		}
	}

	return
}

func listDirectoryContent(dirName string) (content []os.FileInfo, err error) {
	// open directory
	directory, err := os.Open(dirName)
	defer directory.Close()
	if err != nil {
		return
	}

	// list all files and directories as slice
	content, err = directory.Readdir(-1)
	if err != nil {
		return
	}

	return
}

func scanForSubtitles(dirName string) (subs map[string]string, err error) {
	subs = make(subtitles)

	content, err := listDirectoryContent(dirName)
	if err != nil {
		return
	}

	// iterate over files and check if they match regex
	// if so add to subtitles
	r := regexp.MustCompile(`\.(\w+)\.srt$`)
	for _, f := range content {
		match := r.FindStringSubmatch(f.Name())
		if match != nil {
			subs[match[1]] = f.Name()
		}
	}

	return
}

// UpdateSubtitleHandler updates a given subtitles based on request.Body
func UpdateSubtitleHandler(w http.ResponseWriter, request *http.Request) {
	req := customRequest{request}
	params, missing := req.ValidateBody("subtitle", "offset")
	if len(missing) > 0 {
		jsonResponse(w, map[string]interface{}{
			"ok":      false,
			"message": "missing parameters",
			"missing": missing,
		}, http.StatusBadRequest)
		return
	}

	subtitle := GetEnv("MEDIA_PATH", "./") + params["subtitle"]
	_, err := os.Stat(subtitle)
	if os.IsNotExist(err) {
		jsonResponse(w, map[string]interface{}{
			"ok":      false,
			"message": err.Error(),
		}, http.StatusBadRequest)
		return
	}

	err = updateSubtitle(subtitle, params["offset"])
	if err != nil {
		return
	}

	jsonResponse(w, map[string]bool{"ok": true}, http.StatusOK)
}

func updateSubtitle(path string, offset string) (err error) {
	f, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}

	lines := strings.Split(string(f), "\n")

	r := regexp.MustCompile(`(\d+):(\d+):(\d+),(\d+) --> (\d+):(\d+):(\d+),(\d+)`)
	for i, line := range lines {
		matches := r.FindStringSubmatch(line)
		if matches != nil {
			t := timeframe{}
			t.SetFromMatches(matches).Offset(offset)
			lines[i] = t.String()
		}
	}
	output := strings.Join(lines, "\n")
	ioutil.WriteFile(path, []byte(output), 0644)

	return
}
