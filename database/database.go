package database

import (
	"database/sql"
	"log"
	"strings"

	"git.chiefnoah.tech/chiefnoah/gocomics/models"
	_ "github.com/mattn/go-sqlite3"
	sq "github.com/Masterminds/squirrel"

)

//holds a reference to a database connection and a transaction used for large database processes
//like adding a shitton of comics ;P
type Dbhandler struct {
	Transaction *sql.Tx
	Db          *sql.DB
}

func Init() {

	var CREATE_USER_PROGRESS string = `CREATE TABLE IF NOT EXISTS "UserProgress" (
	'ID'	INTEGER PRIMARY KEY AUTOINCREMENT,
	'ComicID'	INTEGER,
	'Read'	INTEGER DEFAULT 0,
	'Completed'	INTEGER DEFAULT 0,
	'DateLastRead'	INTEGER DEFAULT 0,
	'DateCompleted'	INTEGER DEFAULT 0,
	'LastReadPage'	INTEGER DEFAULT 0,
	'UserID'	INTEGER,
	FOREIGN KEY('ComicID') REFERENCES 'Comic'('ID'),
	FOREIGN KEY('UserID') REFERENCES User(ID)
);`
	var CREATE_USER string = `CREATE TABLE IF NOT EXISTS 'User' (
	'ID'	INTEGER PRIMARY KEY AUTOINCREMENT,
	'Name'	TEXT NOT NULL UNIQUE,
	'Password'	TEXT NOT NULL,
	'APIKey'	TEXT
);`
	var CREATE_GENRES_BRIDGE string = `CREATE TABLE IF NOT EXISTS 'GenresBridge' (
	'ComicID'	INTEGER,
	'GenreID'	INTEGER,
	FOREIGN KEY('ComicID') REFERENCES Comic(ID),
	FOREIGN KEY('GenreID') REFERENCES Genres(ID)
);`
	var CREATE_GENRES string = `CREATE TABLE IF NOT EXISTS 'Genres' (
	'ID'	INTEGER PRIMARY KEY AUTOINCREMENT,
	'Genre'	TEXT NOT NULL UNIQUE
);`
	var CREATE_CREDIT string = `CREATE TABLE IF NOT EXISTS 'Credit' (
	'ID'	INTEGER PRIMARY KEY AUTOINCREMENT,
	'Author'	TEXT,
	'Publisher'	TEXT,
	'Other'		TEXT,
	'ComicID'	INTEGER,
	FOREIGN KEY('ComicID') REFERENCES Comic(ID)
);`
	var CREATE_CATEGORY string = `CREATE TABLE IF NOT EXISTS 'Category' (
	'ID'		INTEGER PRIMARY KEY AUTOINCREMENT,
	'Name'		TEXT NOT NULL UNIQUE,
	'Parent'	INTEGER NOT NULL,
	'IsRoot'	INTEGER DEFAULT 0,
	'Full'		TEXT NOT NULL UNIQUE,
	FOREIGN KEY('Parent') REFERENCES Category(ID)
);`
	var CREATE_COMIC_FILE string = `CREATE TABLE IF NOT EXISTS 'ComicFile' (
	'ID'	INTEGER PRIMARY KEY AUTOINCREMENT,
	'RelativePath'	TEXT NOT NULL,
	'AbsolutePath'	TEXT NOT NULL,
	'FileName' TEXT NOT NULL,
	'Hash'	TEXT NOT NULL UNIQUE,
	'Filesize'	INTEGER DEFAULT 0
);`
	var CREATE_COMIC string = `CREATE TABLE IF NOT EXISTS "Comic" (
	'ID'	INTEGER PRIMARY KEY AUTOINCREMENT,
	'Title'	TEXT NOT NULL,
	'Series'	TEXT,
	'IssueNumber'	REAL DEFAULT 0.0,
	'PageCount'	INTEGER,
	'ComicFileID'	INTEGER,
	'Hash'		TEXT NOT NULL UNIQUE,
	'Volume'	TEXT,
	'DateAdded'	INTEGER DEFAULT 0,
	'PublishDate'	INTEGER,
	'Synopsis'	TEXT,
	'Rating'	REAL DEFAULT 0.0,
	'Status'	TEXT,
	'CategoryID'	INTEGER DEFAULT 1,
	FOREIGN KEY('CategoryID') REFERENCES 'Category'('ID')
	FOREIGN KEY('ComicFileID') REFERENCES 'ComicFile'('ID')
);`
	var CREATE_CHARACTERS_BRIDGE string = `CREATE TABLE IF NOT EXISTS 'CharactersBridge' (
	'ComicID'	INTEGER NOT NULL,
	'CharacterID'	INTEGER NOT NULL,
	FOREIGN KEY('ComicID') REFERENCES Comic(ID),
	FOREIGN KEY('CharacterID') REFERENCES Character(ID)
);`
	var CREATE_CHARACTER string = `CREATE TABLE IF NOT EXISTS "Character" (
	'ID'	INTEGER PRIMARY KEY AUTOINCREMENT,
	'Name'	TEXT NOT NULL UNIQUE
);`
	var CREATE_BOOKMARKS string = `CREATE TABLE IF NOT EXISTS 'Bookmarks' (
	'ComicID'	INTEGER,
	'PageNumber'	INTEGER NOT NULL,
	FOREIGN KEY('ComicID') REFERENCES Comic(ID)
);`
	var PRAGMAS string = `PRAGMA foreign_keys = ON; VACUUM`

	db, err := sql.Open("sqlite3", "./library.mdb")
	if err != nil {
		log.Fatal("Unable to open database: ", err)
	}

	defer db.Close()

	_, err = db.Exec(CREATE_USER_PROGRESS)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(CREATE_USER)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(CREATE_GENRES_BRIDGE)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(CREATE_GENRES)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(CREATE_CREDIT)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(CREATE_CATEGORY)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(CREATE_COMIC_FILE)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(CREATE_COMIC)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(CREATE_CHARACTERS_BRIDGE)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(CREATE_CHARACTER)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(CREATE_BOOKMARKS)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(PRAGMAS)
	if err != nil {
		log.Fatal(err)
	}

	log.Print("Done initializing database")

}

func (h *Dbhandler) AddCategory(category *models.Category) error {

	INSERT_CATEGORY := `INSERT OR REPLACE INTO Category(ID, Name, Parent, IsRoot, Full) VALUES(IFNULL((SELECT ID FROM Category WHERE Name = ?), (SELECT SEQ from sqlite_sequence WHERE name='Category') + 1) ,?, IFNULL((SELECT ID FROM Category WHERE Name = ?), 1), ?, IFNULL((SELECT Full FROM Category WHERE Name = ?), "0") || ? || ?)`

	stmt, err := h.Db.Prepare(INSERT_CATEGORY)
	if err != nil {
		log.Println("Error preparing statement: ", err)
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(category.Name, category.Name, category.Parent, category.IsRoot, category.Parent, "/", category.Name)
	if err != nil {
		log.Println("Error adding category: ", err)
	}
	return err
}

func (h *Dbhandler) GetCategory(query *models.Category) *models.Category {


	//var value *models.Category
	var row *sql.Row
	var value models.Category

	sql := `SELECT ID, Name, Parent, IsRoot, Full FROM Category WHERE `
	if query.ID > 0 {
		sql += `ID = ?`
		row = h.Db.QueryRow(sql, query.ID)
	} else if query.Name != "" {
		sql += `Name = ?`
		row = h.Db.QueryRow(sql, query.Name)
	} else if query.ParentId > 0 {
		sql += `Parent = ?`
		row = h.Db.QueryRow(sql, query.ParentId)
	} else if query.Full != "" {
		sql += `Full LIKE ?`
		row = h.Db.QueryRow(sql, query.Full)
	} else if query.Parent != "" {
		sql += `Full LIKE ?`
		row = h.Db.QueryRow(sql, query.Parent)
	}

	err := row.Scan(&value.ID, &value.Name, &value.ParentId, &value.IsRoot, &value.Full)
	if err != nil {
		log.Println("Unable to get Category: ", err)
		return nil
	}

	return &value

}

func (h *Dbhandler) GetChildrenCategories(ID int) *[]models.Category {

	sql := `SELECT ID, Name, Parent, IsRoot, Full FROM Category WHERE Parent = ?`

	rows, err := h.Db.Query(sql, ID)
	if err != nil {
		log.Println("Unable to get children categories: ", err)
		return &[]models.Category{}
	}
	children := []models.Category{}
	for rows.Next() {
		category := models.Category{}
		rows.Scan(&category.ID, &category.Name, &category.ParentId, &category.IsRoot, &category.Full)
		children = append(children, category)
	}

	return &children
}

func (h *Dbhandler) GetChildrenComicsCount(ID int) int {

	sql := `SELECT COUNT(*) FROM Comic WHERE CategoryID = ?`
	row := h.Db.QueryRow(sql, ID)
	var count int = 0
	err := row.Scan(&count)
	if err != nil {
		log.Println("Couldn't get children comics: ", err)
		count = 0
	}

	return count
}

func GetChildrenComicsById(ID int) *[]models.ComicWrapper {

	return nil
}

func GetChildrenComicsByFolder(folder string) *[]models.ComicWrapper {

	return nil
}

func CleanCategory() error {
	//TODO: ONLY IF CATEGORY IS DIRECTORY STRUCTURE: remove rows that have no comics or children categories in them
	return nil
}

func (h *Dbhandler) AddComic(comic *models.ComicInfo, file *models.ComicFile) error {

	//The Hash field is shared because I felt like it. It doesn't really need to be shared but it helps to have a
	//uniquely identifying field besides the ID generated by the database
	INSERT_COMIC_FILE_INFO := `INSERT INTO ComicFile(RelativePath, AbsolutePath, FileName, Hash, Filesize)
	VALUES(?, ?, ?, ?, ?)`

	INSERT_COMIC_INFO := `INSERT INTO Comic(Title, Series, IssueNumber, PageCount, ComicFileID, Hash, Volume,
	DateAdded, PublishDate, Synopsis, Rating, Status, CategoryID)
	VALUES (?, ?, ?, ?, (SELECT ID FROM ComicFile WHERE ComicFile.Hash = ?), ?, ?, ?, ?, ?, ?, ?,
	IFNULL((SELECT ID FROM Category WHERE Full = (SELECT '0\' || RelativePath FROM ComicFile WHERE ComicFile.Hash = ?)), 1))`

	stmt, err := h.Db.Prepare(INSERT_COMIC_FILE_INFO)
	defer stmt.Close()
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = stmt.Exec(file.RelativePath, file.AbsolutePath, file.FileName, file.Hash, file.FileSize)
	if err != nil {
		//If the unique contraint on the HASH field fails, we just update the data
		if strings.Contains(err.Error(), "UNIQUE") {
			log.Print("File already in database, updating to latest info...")
			sql := `UPDATE ComicFile SET RelativePath = ?, AbsolutePath = ?, FileName = ? WHERE Hash = ?`
			st, err := h.Db.Prepare(sql)
			defer st.Close()
			if err != nil {
				log.Println("Unable to prepare statement: ", err)
				return err
			}
			//log.Printf("Updating with:%+v\n", file)
			_, err = st.Exec(file.RelativePath, file.AbsolutePath, file.FileName, file.Hash)
			if err != nil {
				log.Println("error updating: ", err)
				return err
			}

		} else {
			log.Println("Error inserting: ", err)
			return err
		}
	}

	stmt2, err := h.Db.Prepare(INSERT_COMIC_INFO)
	defer stmt2.Close()
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = stmt2.Exec(comic.Title, comic.Series, comic.IssueNumber, comic.PageCount, file.Hash, file.Hash, comic.Volume,
		comic.DateAdded, comic.PublishDate, comic.Synopsis, comic.Rating, comic.Status, file.Hash)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE") {
			//this means the record already exists. We don't want to overwrite tags here so we do nothing
			log.Println("Already have metadata on: ", comic.Title, " ", comic.Series)
		} else {
			log.Println("Database error: ", err)
			return err
		}
	}

	//log.Printf("Results:%+v\n", res)

	return nil
}

func (h *Dbhandler) GetComics(query *models.ComicWrapper) *[]models.ComicWrapper {

	sql := sq.Select("Comic.ID, Title, Series, IssueNumber, PageCount, Hash, Volume, DateAdded, PublishDate, Synopsis, Rating, Status, RelativePath, AbsolutePath, FileName, Filesize").From("Comic").Join("ComicFile USING (Hash)")

	if query.ComicInfo.ID > 0 {
		sql = sql.Where("Comic.ID = ?", query.ComicInfo.ID)
	}
	if query.ComicInfo.Hash != "" {
		sql = sql.Where("Comic.Hash = ?", query.ComicInfo.Hash)
	}
	if query.ComicInfo.Title != "" {
		sql = sql.Where("Title = ?", query.ComicInfo.Title)
	}
	if query.ComicInfo.Series != "" {
		sql = sql.Where("Series = ?", query.ComicInfo.Title)
	}
	if query.ComicInfo.IssueNumber > 0 {
		sql = sql.Where("IssueNumber = ?", query.ComicInfo.IssueNumber)
	}
	//TODO: Continue with these query operators


	return nil
}

func GetDBHandler() *Dbhandler {
	var handler Dbhandler
	db, err := sql.Open("sqlite3", "./library.mdb")
	if err != nil {
		log.Fatal("Unable to open database: ", err)
	}
	handler.Db = db
	handler.Transaction = nil
	return &handler
}

func (h *Dbhandler) ExecuteSql(sql string, params ...interface{}) error {
	stmt, err := h.Db.Prepare(sql)
	if err != nil {
		log.Println("Unable to begin statement in transaction: ", err)
		return err
	}
	_, err = stmt.Exec(params...)
	return err
}

//Creates a new dbhandler object for running a transaction
func (h *Dbhandler) BeginTransaction() {

	if h.Db == nil {
		db, err := sql.Open("sqlite3", "./library.mdb")
		if err != nil {
			log.Fatal("Unable to open database: ", err)
		}
		h.Db = db
	}

	if h.Transaction != nil {
		log.Printf("Transaction already in progress. Not initializing a new one...")
		return
	}

	tx, err := h.Db.Begin()
	if err != nil {
		log.Printf("Unable to begin transaction: %s", err)
	}
	h.Transaction = tx
}

func (h *Dbhandler) FinishTransaction() error {
	if h.Transaction == nil {
		log.Print("Transaction finished before started. Wut?")
		return nil
	}
	err := h.Transaction.Commit()
	h.Transaction = nil
	return err
}
