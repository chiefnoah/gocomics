package web

import (

	"github.com/gin-gonic/gin"

	"net/http"

)

/*

Comic Streamer compatibility API endpoints and stuff goes here


*/
func dbInfoHandler(c *gin.Context) {
	c.String(http.StatusOK, `{"comic_count": 13398, "last_updated": "2015-08-31T20:16:58.035000", "id": "f03b53dbd5364377867227e23112d3c7", "created": "2015-06-18T19:13:35.030000"}`)
}

func versionHandler(c *gin.Context) {
	c.String(http.StatusOK, `{"last_build": "2016-07-03", "version": "0.0.7"}`)
}

/*
func comicListHandler(c *gin.Context) {

	result := models.CSComicResult{}

	c.JSON(http.StatusOK, result)

}

func foldersHandler(c *gin.Context) {
	path := c.Param("path")
	if path == "/" {
		path = "/0"
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

	c.JSON(http.StatusOK, result)
}
*/
