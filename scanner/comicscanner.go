package comicscanner

import (
	//"github.com/fsnotify/fsnotify"
	"path"
	"strings"
	"path/filepath"
	"os"
	"fmt"
	"crypto/md5"
	"io/ioutil"
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
		file, _ := ioutil.ReadFile(p)
		md5 := md5.Sum(file)
		fmt.Printf("MD5: %x\n", md5)

	}
	return nil
}

func watch(f []string) error {
	return nil
}