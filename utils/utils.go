package utils

import (
	"archive/zip"
	"fmt"
	"git.packetlostandfound.us/chiefnoah/gocomics/models"
	"io"
	"log"
	"os"
	"path/filepath"
)

const (
	IMAGES_DIR = ".images"
	CACHE_DIR  = ".temp"
)

//based off answer here https://stackoverflow.com/questions/20357223/easy-way-to-unzip-file-with-golang
//Extracts a cbz to the .temp directory inside a folder with the same name
func ExtractComic(comicfile *models.ComicFile) error {
	r, err := zip.OpenReader(comicfile.AbsolutePath)
	if err != nil {
		log.Print("Unable to extract cbz\n")
		return err
	}
	defer func() {
		if err := r.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	extractAndWrite := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer func() {
			if err := rc.Close(); err != nil {
				log.Fatal(err)
			}
		}()

		wd, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		dirname := filepath.Join(wd, CACHE_DIR, comicfile.Hash)
		os.MkdirAll(dirname, 0755)
		path := filepath.Join(dirname, f.Name)
		fmt.Println("Extracting to: ", path)
		//This probably isn't necessary because we're always dealing with .cbz/.zip files
		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer func() {
				if err := f.Close(); err != nil {
					log.Fatal(err)
				}
			}()
			_, err = io.Copy(f, rc)
			if err != nil {
				log.Fatal(err)
			}
		}
		return nil
	}

	for _, f := range r.File {
		//we only want the first file, so only extract that one
		err = extractAndWrite(f)
		if err != nil {
			return err
		}
	}
	return nil

}

func ExtractCoverImage(comicfile *models.ComicFile) error {
	r, err := zip.OpenReader(filepath.Join(comicfile.AbsolutePath, comicfile.FileName))
	if err != nil {
		log.Print("Unable to extract cbz ", err)
		return err
	}
	defer func() {
		if err := r.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	extractAndWrite := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer func() {
			if err := rc.Close(); err != nil {
				log.Fatal(err)
			}
		}()

		wd, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		dirname := filepath.Join(wd, IMAGES_DIR)
		path := filepath.Join(dirname, comicfile.Hash+".jpg") //EVERYTHING IS A JPG
		//This probably isn't necessary because we're always dealing with .cbz/.zip files
		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer func() {
				if err := f.Close(); err != nil {
					log.Fatal(err)
				}
			}()
			_, err = io.Copy(f, rc)
			if err != nil {
				log.Fatal(err)
			}
		}
		return nil
	}

	err = extractAndWrite(r.File[0])
	if err != nil {
		return err
	}
	return nil
}
