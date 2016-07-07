package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
	"git.chiefnoah.tech/chiefnoah/gocomics/models"
	"git.chiefnoah.tech/chiefnoah/gocomics/database"
	"strings"
	"fmt"
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

func comicListHandler(c *gin.Context) {

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
	fmt.Printf("Children folders: %+v", childrenFolders)
	childrenComicsCount := database.GetChildrenComicsCount(category.ID)
	childrenComics := models.CSComicCountResponse{
		Count: childrenComicsCount,
		URL_Path: "/comiclist?folder=" + category.Full,
	}
	folders := []models.CSFolder{}

	for _, v := range *childrenFolders {
		if(v.ID == 1) {
			continue
		}
		csfolder := models.CSFolder{}
		csfolder.URL_Path = "/folders/" + strings.Replace(v.Full, "\\", "/", -1)
		csfolder.Name = v.Name
		folders = append(folders, csfolder)
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