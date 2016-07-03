package database

import (
	"database/sql"
	"log"
	_ "github.com/mattn/go-sqlite3"
	"git.chiefnoah.tech/chiefnoah/gocomics/models"
)

var db *sql.DB

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
	'Other'	TEXT,
	'ComicID'	INTEGER,
	FOREIGN KEY('ComicID') REFERENCES Comic(ID)
);`
	var CREATE_COMIC_FILE string = `CREATE TABLE IF NOT EXISTS 'ComicFile' (
	'ID'	INTEGER PRIMARY KEY AUTOINCREMENT,
	'RelativePath'	TEXT NOT NULL,
	'AbsolutePath'	TEXT NOT NULL,
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
	'Volume'	TEXT,
	'DateAdded'	INTEGER DEFAULT 0,
	'PublishDate'	INTEGER,
	'Synopsis'	TEXT,
	'Rating'	REAL DEFAULT 0.0,
	'Status'	TEXT,
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

func AddComic(comic models.ComicInfo) bool {
	db, err := sql.Open("sqlite3", "./library.mdb")
	if err != nil {
		log.Fatal("Unable to open database: ", err)
	}

	INSERT_COMIC_INFO := "INSERT OR IGNORE INTO comic_info(filesize, date_added, hash) VALUES(?, ?, ?);"
	//INSERT_COMIC := "INSERT OR IGNORE INTO comic(parentId, comicInfoId, fileName, path) VALUES ((SELECT id FROM folder WHERE path = ?), ?, ?, ?)"

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := tx.Prepare(INSERT_COMIC_INFO)
	if err != nil {
		log.Fatal(err)
	}

	defer stmt.Close()

	return false
}