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
)

func Scan(f string) error {
	err := filepath.Walk(f, visit)
	if err != nil {
		fmt.Printf("walk error: %v\n", err)
		return err
	}
	return nil
}

func visit(p string, f os.FileInfo, e error) error {

	fmt.Printf("Visited: %s\n", p)
	if strings.EqualFold(path.Ext(f.Name()), ".cbz") || strings.EqualFold(path.Ext(f.Name()), ".cbr") {
		fmt.Printf("Found cbz file!\n")

		//TODO: parse comic info
		file, err := ioutil.ReadFile(p)
		if err != nil {
			log.Print("Error: ", err)
			return err
		}
		md5 := md5.Sum(file)
		fmt.Printf("MD5: %x\n", md5)

	}
	return nil
}

func watch(f []string) error {
	return nil
}
