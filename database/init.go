package database

import (
	"database/sql"
	"log"

	"forum/handul"

	_ "github.com/mattn/go-sqlite3"
)

func IntDB() *handul.Date {
	db, err := sql.Open("sqlite3", "forum.db")
	if err != nil {
		log.Fatal(err)
	}
	return &handul.Date{
		DB: db,
	}
}

func CreateTable(db *handul.Date) error {
	qeury := `CREATE TABLE IF NOT EXISTS user(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_name TEXT UNIQUE , 
		email TEXT UNIQUE , 
		passwd TEXT , 
		create_date DEFAULT CURRENT_TIMESTAMP
	) ; 

	CREATE TABLE IF NOT EXISTS post(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER , 
		title TEXT, 
		contant TEXT , 
		create_date DEFAULT CURRENT_TIMESTAMP ,
		FOREIGN KEY  (user_id) REFERENCES user(id) 
	) ;

	CREATE TABLE IF NOT EXISTS categories(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name_categorie TEXT UNIQUE NOT NULL
	) ;


	CREATE TABLE IF NOT EXISTS categories_post(
		categorie_id INTEGER , 
		post_id INTEGER , 
		FOREIGN KEY (categorie_id) REFERENCES categories(id)
		FOREIGN KEY (post_id) REFERENCES post(id)
	) ;

	CREATE TABLE IF NOT EXISTS comment (
		id INTEGER PRIMARY KEY AUTOINCREMENT , 
		post_id INTEGER , 
		user_id INTEGER ,
		contant TEXT  , 
		FOREIGN KEY (post_id) REFERENCES post(id)
		FOREIGN KEY (user_id) REFERENCES user(id)
	);

	CREATE TABLE IF NOT EXISTS reaction(
		id INTEGER PRIMARY KEY AUTOINCREMENT , 
		post_id INTEGER , 
		comment_id INTEGER ,
		user_id INTEGER ,
		type TEXT  , 
		FOREIGN KEY (post_id) REFERENCES post(id),
		FOREIGN KEY (user_id) REFERENCES user(id)
		FOREIGN KEY (comment_id) REFERENCES comment(id),
		CHECK (
		(post_id IS NOT NULL AND comment_id IS NULL)
		OR	(post_id IS  NULL AND comment_id IS NOT  NULL)
		)
	);

	CREATE TABLE IF NOT EXISTS session (
		id INTEGER PRIMARY KEY, 
		uid TEXT UNIQUE, 
		user_id INTEGER UNIQUE, 
		create_date DATE DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES user (id)
	);
	INSERT INTO categories (name_categorie)
	VALUES 
    ('Game'),
    ('Tecnolge'),
    ('Ecomerc'),
    ('Natur'),
    ('Viset')
ON CONFLICT (name_categorie) DO NOTHING;
	`
	_, err := db.DB.Exec(qeury)
	return err
}
