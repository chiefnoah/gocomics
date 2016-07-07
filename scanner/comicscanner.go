package comicscanner

import (
	//"github.com/fsnotify/fsnotify"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"git.chiefnoah.tech/chiefnoah/gocomics/database"
	"git.chiefnoah.tech/chiefnoah/gocomics/models"
	"git.chiefnoah.tech/chiefnoah/gocomics/utils"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

var root string
var dbhandler *database.Dbhandler

func Scan(f string) error {
	root = f
	//base := filepath.Base(f)
	//models.Category{ID: 1, Name: base, Parent: 1, IsRoot:true, Full: base}
	//generates the temp and image directories
	setupDirs()
	dbhandler = database.BeginTransaction()
	err := dbhandler.ExecuteSql(`INSERT OR IGNORE INTO Category(ID, Name, Parent, IsRoot, Full) VALUES(?, ?, ?, ?, ?)`, 1, "0", 1, true, "0")
	if err != nil {
		log.Println("Error creating start category dir: ", err)
	}
	defer dbhandler.FinishTransaction()
	err = filepath.Walk(f, visit)
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
		checksum := md5.Sum(file)
		n := len(checksum)
		comicfile.Hash = hex.EncodeToString(checksum[:n])
		comicfile.FileSize = int64(f.Size())
		rel, _ := filepath.Rel(root, p)
		comicfile.RelativePath = filepath.Dir(filepath.ToSlash(filepath.Join(root, rel)))
		comicfile.FileName = f.Name()
		if !path.IsAbs(root) {
			ab, err := filepath.Abs(p)
			if err != nil {
				log.Print("Couldn't get relative path: ", err)
			}
			comicfile.AbsolutePath = filepath.Dir(filepath.ToSlash(ab))
		}

		//fmt.Printf("MD5: %s\n", comicfile.Hash)
		dbhandler.AddComic(&models.ComicInfo{}, &comicfile)
		go generateCoverImage(&comicfile)

	} else {

		dir := filepath.Base(filepath.Dir(p))
		name := filepath.Base(p)
		category := models.Category{Name: name, Parent: dir, IsRoot: false}
		if dir == root {
			category.IsRoot = true
		}
		dbhandler.AddCategory(&category)
	}
	return nil
}

func watch(f []string) error {
	return nil
}

//TODO: read metadata from file or filename
func setupDirs() {
	os.MkdirAll(utils.IMAGES_DIR, 0755)
	os.MkdirAll(utils.CACHE_DIR, 0755)
}

func generateCoverImage(comicfile *models.ComicFile) {

	if _, f := os.Stat(filepath.Join(utils.IMAGES_DIR, comicfile.Hash)); os.IsNotExist(f) {
		fmt.Println("No cover image found, generating")
		err := utils.ExtractCoverImage(comicfile)
		if err != nil {
			log.Println("Extraction error: ", err)
		}
	}
}
