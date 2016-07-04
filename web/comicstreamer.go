package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func dbInfoHandler(c *gin.Context) {
	c.String(http.StatusOK, `{"comic_count": 13398, "last_updated": "2015-08-31T20:16:58.035000", "id": "f03b53dbd5364377867227e23112d3c7", "created": "2015-06-18T19:13:35.030000"}`)
}

func versionHandler(c *gin.Context) {
	c.String(http.StatusOK, `{"last_build": "2016-07-03", "version": "0.0.7"}`)
}

func comicListHandler(c *gin.Context) {

}