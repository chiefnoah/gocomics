package web

import (
	"net/http"
	"net/url"
	"path/filepath"
	"git.chiefnoah.tech/chiefnoah/gocomics/models"
	"git.chiefnoah.tech/chiefnoah/gocomics/database"
	"strings"
	"fmt"
	"log"
	"encoding/json"
)


/*

Comic Streamer compatibility API endpoints and stuff goes here


 */
func dbInfoHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`{"comic_count": 13398, "last_updated": "2015-08-31T20:16:58.035000", "id": "f03b53dbd5364377867227e23112d3c7", "created": "2015-06-18T19:13:35.030000"}`))
}

func versionHandler(w http.ResponseWriter, r *http.Request) {
	//Dummy info lol
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"last_build": "2016-07-03", "version": "0.0.7"}`))
	w.WriteHeader(http.StatusOK)

}

func comicListHandler(w http.ResponseWriter, r *http.Request) {

	result := models.CSComicResult{}

	if err := json.NewEncoder(w).Encode(result); err != nil {
		log.Printf("JSON Ecode error: %s", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)

	}

}

func foldersHandler(w http.ResponseWriter, r *http.Request) {
	print("FUCK MUX")
	//pathParams := mux.Vars(r)
	requestUrl, err := url.Parse(r.RequestURI)
	if err != nil {
		log.Print("How did you even mess this up? How did it make it this far???")
		return
	}
	var path string
	if requestUrl.Path == "/" {
		path = "/0"
	} else {
		path = requestUrl.Path
	}

	base := filepath.Base(path)
	fmt.Println("Base: ", base)

	var query = models.Category{
		Name: base,
	}
	category := database.GetCategory(&query)
	childrenFolders := database.GetChildrenCategories(category.ID)
	childrenComicsCount := database.GetChildrenComicsCount(category.ID)
	childrenComics := models.CSComicCountResponse{
		Count: childrenComicsCount,
		URL_Path: "/comiclist?folder=" + url.QueryEscape(category.Full),
	}
	log.Print("URL_PATH: ", childrenComics.URL_Path)
	folders := []models.CSFolder{}

	for _, v := range *childrenFolders {
		if(v.ID == 1) {
			continue
		}
		csfolder := models.CSFolder{}
		//This is a stupid workaround to url escape all the folders in the category path without losing the
		//path separators.
		split_path := strings.Split(v.Full, "/")
		for i, folder := range split_path {
			//We can't just URL escape this because for some reason url.QueryEscape() will only replace
			// spaces with a + instead of %20, and there doesn't seem to be a way around it
			split_path[i] = strings.Replace(folder, " ", "%20", -1)
		}
		v.Full = strings.Join(split_path, "/")
		csfolder.URL_Path = "/folders/" + v.Full
		csfolder.Name = v.Name
		folders = append(folders, csfolder)
		log.Print("CS Folder URL_Path: ", csfolder.URL_Path)
	}

	result := models.CSFolderResponse{
		Current: category.Full,
		Folders: folders,
		Comics: childrenComics,
	}


	//for _, v := range categoryNames {
		//TODO: query database for category
	//}
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(result); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}