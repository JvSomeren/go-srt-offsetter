package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type customRequest struct {
	*http.Request
}

// ValidateBody check if all passed parameter names are passen with the request
func (req customRequest) ValidateBody(params ...string) (values map[string]string, missing []string) {
	values = make(map[string]string)

	decoder := json.NewDecoder(req.Body)
	var vars map[string]interface{}
	_ = decoder.Decode(&vars)

	for _, param := range params {
		value := vars[param]
		if value == nil {
			missing = append(missing, param)
		} else {
			values[param] = fmt.Sprint(value)
		}
	}

	return
}

type spaHandler struct {
	staticPath string
	indexPath  string
	basePath   string
}

func jsonResponse(w http.ResponseWriter, res interface{}, statusCode int) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(res)
}

// ServeHTTP inspects the URL path to locate a file within the static dir
// on the SPA handler. If a file is found, it will be served. If not, the
// file located at the index path on the SPA handler will be served. This
// is suitable behavior for serving an SPA (single page application).
func (h spaHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// get the absolute path to prevent directory traversal
	path, err := filepath.Abs(req.URL.Path)
	if err != nil {
		// if we failed to get the absolute path respond with a 400 bad request
		// and stop
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// remove basePath from path
	path = strings.Replace(path, h.basePath, "", 1)

	// prepend the path with the path to the static directory
	path = filepath.Join(h.staticPath, path)
	fmt.Println(path)

	// check whether a file exists at the given path
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		// file does not exist, serve index.html
		http.ServeFile(w, req, filepath.Join(h.staticPath, h.indexPath))
		return
	} else if err != nil {
		// if we got an error (that wasn't that the file doesn't exist) stating the
		// file, return a 500 internal server error and stop
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// otherwise, use http.FileServer to serve the static dir
	http.StripPrefix(h.basePath, http.FileServer(http.Dir(h.staticPath))).ServeHTTP(w, req)
}

func healthCheck(w http.ResponseWriter, req *http.Request) {
	jsonResponse(w, map[string]bool{"ok": true}, http.StatusOK)
}

func registerAPIRoutes(router *mux.Router) {
	router.HandleFunc("/health", healthCheck).Methods("GET")
	router.HandleFunc("/media", ListMediaAndSubtitlesHandler).Methods("GET")
	router.HandleFunc("/subtitle", UpdateSubtitleHandler).Methods("PUT")
}

// CreateServer returns a pointer to a newly created http.Server
func CreateServer() *http.Server {
	router := mux.NewRouter()

	api := router.PathPrefix("/api/").Subrouter()
	registerAPIRoutes(api)

	staticPath := GetEnv("STATIC_PATH", "../webapp/public")
	basePath := GetEnv("BASE_PATH", "/")
	spa := spaHandler{staticPath: staticPath, indexPath: "index.html", basePath: basePath}
	router.PathPrefix(basePath).Handler(spa)

	// redirect all root requests to the given basePath
	if basePath != "/" {
		router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, basePath, http.StatusPermanentRedirect)
		})
	}

	originsOk := handlers.AllowedOrigins([]string{`*`})
	headersOk := handlers.AllowedHeaders([]string{"Content-Type", "X-Requested-With"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	srv := &http.Server{
		Handler:      handlers.CORS(originsOk, headersOk, methodsOk)(router),
		Addr:         ":" + GetEnv("PORT", "8080"),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	return srv
}
