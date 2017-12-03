package main

import (
	"github.com/gotk3/gotk3/gtk"

	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
)

type fileInfo struct {
	Name       string
	Icon       string
	ChildDirs  []fileInfo
	ChildFiles []fileInfo
}

type tcpKeepAliveListener struct {
	*net.TCPListener
}

func (ln tcpKeepAliveListener) Accept() (c net.Conn, err error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
}

// fileServer serves a file on a random port number. It shuts down if there
// is no download from the server within a period of time.
type fileServer struct {
	http.Server
	port     int
	listener net.Listener
}

// Map of selected files provided as command line args to the application.
var rootSelectedFiles map[string]bool

var sharedPath = "/"

func fileServerNew() (*fileServer, error) {
	fs := &fileServer{}
	fs.Server.Addr = ":"
	listener, err := net.Listen("tcp", fs.Server.Addr)
	fs.listener = listener
	if err != nil {
		return nil, err
	}
	fs.port = fs.listener.Addr().(*net.TCPAddr).Port
	return fs, nil
}

var absPath string

func (fs *fileServer) start(app *QrShare, qrWindow *gtk.ApplicationWindow) error {
	rootSelectedFiles = make(map[string]bool)
	for _, s := range app.files {
		rootSelectedFiles[path.Base(s)] = true
	}

	_, err := os.Stat(app.files[0])
	if err != nil {
		log.Println("Error starting server:", err)
		return err
	}

	absPath = getAbsPath(app.files)

	mux := http.NewServeMux()

	// Serve shared files under path sharedPath
	mux.Handle(sharedPath, http.StripPrefix(sharedPath,
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if path.Clean(r.URL.Path) == "/" {
				// Because there might be only one file to be served at the root.
				// Joining a URL path / at the end of a filename will make the
				// filepath a directory name.
				serve(w, r, absPath)
			} else {
				serve(w, r, path.Join(absPath, path.Clean(r.URL.Path)))
			}
		})))

	fs.Server.Handler = mux

	log.Println("Starting file sharing")
	return fs.Serve(tcpKeepAliveListener{fs.listener.(*net.TCPListener)})
}

// getParentDir returns the parent directory of the given file.
func getAbsPath(names []string) string {
	p := names[0]

	// Get the absolute path name for relative filename.
	if !path.IsAbs(p) {
		c, _ := os.Getwd()
		p = path.Join(c, p)
	}

	if len(names) == 1 {
		return p
	}

	return path.Dir(p)
}

// serve serves a file or directory request with the given file path.
func serve(w http.ResponseWriter, r *http.Request, filePath string) {
	f, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			log.Println("Requested non existent file:", err)
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, notFoundHTML)
			return
		}

		log.Println("Error opening file:", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, internalErrorHTML)
		return
	}

	defer f.Close()

	fStat, err := f.Stat()
	if err != nil {
		log.Println("Error getting file info:", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, internalErrorHTML)
		return
	}

	if fStat.IsDir() {
		serveDir(w, r, f)
		return
	}

	if (fStat.Mode() &^ 07777) == os.ModeSocket {
		log.Println("file is a socket: not serving it")
		w.WriteHeader(http.StatusNotFound) // maybe status forbidden??
		fmt.Fprint(w, notFoundHTML)
		return
	}

	serveFile(w, r, filePath)
}

// serveDir serves a directory content.
func serveDir(w http.ResponseWriter, r *http.Request, f *os.File) {
	fStats, err := f.Readdir(-1)
	if err != nil {
		log.Println("Error reading directory:", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, internalErrorHTML)
		return
	}

	fis := fileInfo{
		Name:       strings.TrimPrefix(f.Name(), absPath),
		ChildFiles: []fileInfo{},
		ChildDirs:  []fileInfo{},
	}

	// TODO: Sort the filenames
	// I think the slice return from os.Readdir is sorted already.

	for _, fStat := range fStats {
		// Skip hidden files.
		if fStat.Name()[0] == '.' {
			continue
		}

		// If path is root, filter files not in app.files.
		if _, ok := rootSelectedFiles[fStat.Name()]; !ok &&
			// http.StripPrefix removes / also. Need to check for that.
			(r.URL.Path == "/" || r.URL.Path == "") &&
			// It's okay if there is only one directory is to be served.
			len(rootSelectedFiles) != 1 {
			continue
		}

		if fStat.IsDir() {
			fis.ChildDirs = append(fis.ChildDirs, fileInfo{Name: fStat.Name(), Icon: iconFolder})
		} else {
			icon := getIcon(path.Join(f.Name(), fStat.Name()))
			fis.ChildFiles = append(fis.ChildFiles, fileInfo{Name: fStat.Name(), Icon: icon})
		}
	}

	tpl, err := template.New("t").Parse(listingHTML)
	if err != nil {
		log.Println("Error parsing template:", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, internalErrorHTML)
		return
	}

	err = tpl.Execute(w, fis)
	if err != nil {
		log.Println("Error executing template:", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, internalErrorHTML)
		return
	}
}

// serveFile serves a file using standard library http.ServeFile.
func serveFile(w http.ResponseWriter, r *http.Request, filePath string) {
	http.ServeFile(w, r, filePath)
}
