package project

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/goptos/goptos/githubapi"
)

func untar(dst string, r io.Reader) error {
	gzr, err := gzip.NewReader(r)
	if err != nil {
		return err
	}
	defer gzr.Close()
	tr := tar.NewReader(gzr)
	for {
		header, err := tr.Next()
		switch {
		// if no more files are found return
		case err == io.EOF:
			return nil
		// return any other error
		case err != nil:
			return err
		// if the header is nil, just skip it (not sure how this happens)
		case header == nil:
			continue
		}
		// the target location where the dir/file should be created
		// target := filepath.Join(dst, header.Name)
		// strip the root folder away
		target := filepath.Join(dst, strings.Join(strings.Split(header.Name, "/")[1:], "/"))

		// the following switch could also be done using fi.Mode(), not sure if there
		// a benefit of using one vs. the other.
		// fi := header.FileInfo()

		// check the file type
		switch header.Typeflag {
		// if its a dir and it doesn't exist create it
		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, 0755); err != nil {
					return err
				}
			}
		// if it's a file create it
		case tar.TypeReg:
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}
			// copy over contents
			if _, err := io.Copy(f, tr); err != nil {
				return err
			}
			// deferring would cause each file close to wait until all operations have completed.
			f.Close()
		}
	}
}

func create(org string, repo string, version string) {
	var url = "https://github.com/" + org + "/" + repo + "/archive/refs/tags/" + version + ".tar.gz"
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		log.Fatalf("get(`%s`): %s", url, res.Status)
	}
	err = untar(".", res.Body)
	if err != nil {
		log.Fatal(err)
	}
}

func Init(org string, repo string, version string) {
	var tags = githubapi.GetRepoTags(org, repo)
	if version == "latest" {
		var latest = ""
		for _, tag := range tags {
			var ver = strings.Split(tag.Ref, "/")[2]
			if ver > latest {
				latest = ver
			}
		}
		if latest == "" {
			log.Fatal("unable to find latest version of github.com/" + org + "/" + repo)
		}
		create(org, repo, latest)
		return
	}
	var hit = false
	for _, tag := range tags {
		var ver = strings.Split(tag.Ref, "/")[2]
		if version == ver {
			hit = true
			break
		}
	}
	if !hit {
		log.Fatalf("unable to find version `%s` of github.com/"+org+"/"+repo, version)
	}
	create(org, repo, version)
}
