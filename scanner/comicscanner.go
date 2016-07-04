package comicscanner

import (
	//"github.com/fsnotify/fsnotify"
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"log"
	"git.chiefnoah.tech/chiefnoah/gocomics/models"
	"git.chiefnoah.tech/chiefnoah/gocomics/database"
	"encoding/hex"
)

var root string
var dbhandler *database.Dbhandler

func Scan(f string) error {
	root = f
	dbhandler = database.BeginTransaction()
	defer dbhandler.FinishTransaction()
	err := filepath.Walk(f, visit)
	if err != nil {
		fmt.Printf("walk error: %v\n", err)
		return err
	}
	log.Print("DONE WALKING!")
	return nil
}

func visit(p string, f os.FileInfo, e error) error {

	fmt.Printf("Visited: %s\n", p)
	if strings.EqualFold(path.Ext(f.Name()), ".cbz") || strings.EqualFold(path.Ext(f.Name()), ".cbr") {
		//fmt.Printf("Found cbz file!\n")

		//TODO: parse comic info
		file, err := ioutil.ReadFile(p)
		if err != nil {
			log.Print("Error: ", err)
			return err
		}

		var comicfile models.ComicFile
		//TODO: somehow get comic info based on filename/directory structure
		//TODO: generate cover images using hash
		checksum := md5.Sum(file)
		n := len(checksum)
		comicfile.Hash = hex.EncodeToString(checksum[:n])
		comicfile.FileSize = int64(f.Size())
		rel, _ := filepath.Rel(root, p)
		comicfile.RelativePath = filepath.ToSlash(rel)
		comicfile.FileName = f.Name()
		if !path.IsAbs(root) {
			ab, err := filepath.Abs(p)
			if err != nil {
				log.Print("Couldn't get relative path: ", err)
			}
			comicfile.AbsolutePath = filepath.ToSlash(ab)
		}

		//fmt.Printf("MD5: %s\n", comicfile.Hash)
		dbhandler.AddComic(models.ComicInfo{}, comicfile)

	}
	return nil
}

func watch(f []string) error {
	return nil
}
