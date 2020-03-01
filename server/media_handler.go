package main

import (
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"
)

type file struct {
	Name      string    `json:"name"`
	IsDir     bool      `json:"isDir"`
	Files     []file    `json:"files"`
	Subtitles subtitles `json:"subtitles"`
}

type subtitles map[string]string

var mediaCache []file = nil

// MediaCacheHandler starts the caching of media periodically. Returns stop
func MediaCacheHandler(delayOptional ...time.Duration) chan bool {
	delay := 30 * time.Second
	if len(delayOptional) > 0 {
		delay = delayOptional[0]
	}

	stop := make(chan bool)

	go func() {
		for {
			var err error
			mediaCache, err = scanDirectory(GetEnv("MEDIA_PATH", "./"))
			if err != nil {
				mediaCache = nil
			}

			select {
			case <-time.After(delay):
			case <-stop:
				return
			}
		}
	}()

	return stop
}

// ListMediaAndSubtitlesHandler makes a JSON response of nested media and subtitles files
func ListMediaAndSubtitlesHandler(w http.ResponseWriter, req *http.Request) {
	if len(mediaCache) == 0 {
		movies, err := scanDirectory(GetEnv("MEDIA_PATH", "./"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		jsonResponse(w, movies, http.StatusOK)
	} else {
		jsonResponse(w, mediaCache, http.StatusOK)
	}
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

	// sort alphabetically
	sort.Slice(content, func(i, j int) bool { return content[i].Name() < content[j].Name() })

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
