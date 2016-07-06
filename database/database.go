package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	"git.chiefnoah.tech/chiefnoah/gocomics/models"
	_ "github.com/mattn/go-sqlite3"
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
	'CategoryID'	INTEGER,
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
	//I'm guessing the "OR REPLACE" part will cause issues later one when people don't want changing the directory
	//structure of their comics to lose their categories. It might not be a problem if categories are purely based on
	//the actual directory structure, but if the are arbitrarily set (which I plan on letting you do) it will have to
	//rely on the CategoryID field of the Comic table, which won't get updated when a row is replaced because it generates
	//a new ID. A possible solution is to use SQL triggers, but I'm unfamiliar with them, I'll have to look into it.
	INSERT_CATEGORY := `INSERT OR REPLACE INTO Category(Name, Parent, IsRoot, Full) VALUES(?, (SELECT ID FROM Category WHERE Name = ?), ?, (SELECT Full FROM Category WHERE Name = ?) || ? || ?)`

	stmt, err := h.Transaction.Prepare(INSERT_CATEGORY)
	if err != nil {
		log.Println("Error preparing statement: ", err)
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(category.Name, category.Parent, category.IsRoot, category.Parent, fmt.Sprintf("%c", os.PathSeparator), category.Name)
	if err != nil {
		log.Println("Error adding category: ", err)
	}
	return err
}

func (h *Dbhandler) AddComic(comic *models.ComicInfo, file *models.ComicFile) error {

	//The Hash field is shared because I felt like it. It doesn't really need to be shared but it helps to have a
	//uniquely identifying field besides the ID generated by the database
	INSERT_COMIC_FILE_INFO := `INSERT INTO ComicFile(RelativePath, AbsolutePath, FileName, Hash, Filesize) VALUES(?, ?, ?, ?, ?)`
	INSERT_COMIC_INFO := `INSERT INTO Comic(Title, Series, IssueNumber, PageCount, ComicFileID, Hash, Volume, DateAdded, PublishDate, Synopsis, Rating, Status) VALUES (?, ?, ?, ?, (SELECT ID FROM ComicFile WHERE ComicFile.Hash = ?), ?, ?, ?, ?, ?, ?, ?)`

	stmt, err := h.Transaction.Prepare(INSERT_COMIC_FILE_INFO)
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
			st, err := h.Transaction.Prepare(sql)
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

	stmt2, err := h.Transaction.Prepare(INSERT_COMIC_INFO)
	defer stmt2.Close()
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = stmt2.Exec(comic.Title, comic.Series, comic.IssueNumber, comic.PageCount, file.Hash, file.Hash, comic.Volume,
		comic.DateAdded, comic.PublishDate, comic.Synopsis, comic.Rating, comic.Status)
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

func (h *Dbhandler) ExecuteSql(sql string, params ...interface{}) error {
	stmt, err := h.Transaction.Prepare(sql)
	if err != nil {
		log.Println("Unable to begin statement in transaction: ", err)
		return err
	}
	_, err = stmt.Exec(params...)
	return err
}

//Creates a new dbhandler object for running a transaction
func BeginTransaction() *Dbhandler {
	var handler Dbhandler

	db, err := sql.Open("sqlite3", "./library.mdb")
	if err != nil {
		log.Fatal("Unable to open database: ", err)
	}
	handler.Db = db

	tx, err := db.Begin()
	handler.Transaction = tx
	return &handler
}

func (h *Dbhandler) FinishTransaction() error {
	err := h.Transaction.Commit()
	defer h.Db.Close()
	return err
}
