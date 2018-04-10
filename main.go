package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"strings"

	"github.com/logeable/shs/middleware"
)

func main() {
	shsfs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	port := shsfs.Int("port", 9929, "server port")
	dir := shsfs.String("dir", ".", "statics files directory")

	shsfs.Parse(os.Args[1:])

	p, err := getAbsPath(*dir)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(fmt.Sprintf("port: %d, dir: %s", *port, p))

	fs := http.FileServer(http.Dir(p))
	mux := http.NewServeMux()
	mux.Handle("/", middleware.LogMiddleware(fs))

	addr := fmt.Sprintf(":%d", *port)
	log.Fatal(http.ListenAndServe(addr, mux))
}

func getAbsPath(p string) (string, error) {
	if !path.IsAbs(p) {
		if strings.HasPrefix(p, "~/") {
			usr, err := user.Current()
			if err != nil {
				return "", err
			}
			return filepath.Clean(filepath.Join(usr.HomeDir, p[2:])), nil
		}
		wd, err := os.Getwd()
		if err != nil {
			return "", err
		}
		return filepath.Clean(filepath.Join(wd, p)), nil
	}
	return p, nil
}
