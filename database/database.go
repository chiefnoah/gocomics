package database

import (
	"git.packetlostandfound.us/chiefnoah/gocomics/config"
	"git.packetlostandfound.us/chiefnoah/gocomics/models"
	"github.com/HouzuoGuo/tiedot/db"
	"github.com/fatih/structs"
	"go.uber.org/zap"
)

//Global database connection. Shared with all threads, cause thread safe :3
var globalDB *ComicDB

type ComicDB struct {
	td *db.DB
	//idk what else should go in here lol
	Log *zap.SugaredLogger
}

type Database interface {
	//getDbConnection() (*db, error)
	//Comic Handling
	GetComicInfo(*models.ComicInfo) *[]map[string]interface{}
	AddComicInfo(*models.ComicInfo) error
	UpdateComicInfo(*models.ComicInfo) error
	DeleteComicInfo(file *models.ComicInfo) error

	//Folder Handling
	GetFolders(*models.Folder) *[]models.Folder
	AddFolder(*models.Folder) error
	UpdateFolder(*models.Folder) error
	DeleteFolder(*models.Folder) error

	CloseDbConnection() error
}

func GetComicDatabase() (*ComicDB, error) {
	if globalDB != nil {
		return globalDB, nil
	}
	logger, err := zap.NewDevelopment()
	if err != nil {
		//the logger failed. Idk what to even do here
		return nil, err
	}
	sugar := logger.Sugar()
	config := config.LoadConfigFile()
	db, err := db.OpenDB(config.DatabaseFolder)
	if err != nil {
		sugar.Errorf("Unable to initialize database: %s", err)
		return nil, err
	}
	//TODO: "create" all tables here
	err = db.Create("comicinfo")
	if err != nil {
		sugar.Warnf("Unable to create comicinfo table: %s", err)
	}
	cdb := ComicDB{
		td:  db,
		Log: sugar,
	}
	globalDB = &cdb
	return globalDB, nil
}

func (database *ComicDB) GetComicInfo(comic *models.ComicInfo) *[]map[string]interface{} {
	ci := database.td.Use("comicinfo")

	query := "all"/*map[string]interface{}{
		"eq":    comic.Hash,
		"in":    []interface{}{"Hash"},
		"limit": 1,
	}*/
	result := make(map[int]struct{})
	//You apparently have to index anything you want to query against
	//err := ci.Index([]string{"Hash"})
	//database.Log.Debugf("All index: %+x: ", ci.AllIndexes())
	//if err != nil {
	//	database.Log.Infof("Unable to index: %s", err)
	//}
	if err := db.EvalQuery(query, ci, &result); nil != err {
		database.Log.Errorf("Unable to retrieve comic from database: %s", err)
	}
	output := make([]map[string]interface{}, len(result))

	for id := range result {
		readBack, err := ci.Read(id)
		if err != nil {
			database.Log.Errorf("Unable to read document: %s", err)
		}
		database.Log.Debugf("ID: %s\nResulting doc: %v\n", id, readBack)
		readBack["ID"] = id
		output = append(output, readBack)
	}

	return &output
}

func (db *ComicDB) AddComicInfo(comic *models.ComicInfo) error {
	ci := db.td.Use("comicinfo")
	c := structs.Map(comic)
	//TODO: decide if I want to handle the document ID after inserting
	_, err := ci.Insert(c)
	return err
}
