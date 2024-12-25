package database

import (
	"database/sql"
	"log"

	"forum/controllers"

	_ "github.com/mattn/go-sqlite3"
)

func IntDB() *controllers.Date {
	db, err := sql.Open("sqlite3", "forum.db")
	if err != nil {
		log.Fatal(err)
	}
	return &controllers.Date{
		DB: db,
	}
}

func CreateTable(db *controllers.Date) error {
	qeury := `CREATE TABLE IF NOT EXISTS user(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_name TEXT UNIQUE NOT NULL, 
		email TEXT UNIQUE NOT NULL, 
		passwd TEXT NOT NULL, 
		create_date DEFAULT CURRENT_TIMESTAMP
	) ; 

	CREATE TABLE IF NOT EXISTS post(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL, 
		title TEXT NOT NULL, 
		contant TEXT NOU NULL, 
		create_date DEFAULT CURRENT_TIMESTAMP ,
		img TEXT,
		categories TEXT  NOU NULL, 
		FOREIGN KEY  (user_id) REFERENCES user(id) 
	) ;

	CREATE TABLE IF NOT EXISTS categories(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name_categorie TEXT UNIQUE NOT NULL
	) ;


	CREATE TABLE IF NOT EXISTS categories_post(
		categorie_id INTEGER NOT NULL, 
		post_id INTEGER NOT NULL, 
		FOREIGN KEY (categorie_id) REFERENCES categories(id)
		FOREIGN KEY (post_id) REFERENCES post(id)
		UNIQUE(categorie_id , post_id)
	) ;

	CREATE TABLE IF NOT EXISTS comment (
		id INTEGER PRIMARY KEY AUTOINCREMENT , 
		post_id INTEGER NOT NULL, 
		user_id INTEGER NOT NULL,
		contant TEXT NOT NULL , 
		create_date DEFAULT CURRENT_TIMESTAMP ,
		FOREIGN KEY (post_id) REFERENCES post(id)
		FOREIGN KEY (user_id) REFERENCES user(id)
	);

	CREATE TABLE IF NOT EXISTS reaction (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    post_id INTEGER,
    comment_id INTEGER,
    user_id INTEGER NOT NULL,
    type TEXT CHECK (type IN ('likes', 'dislikes')) NOT NULL,
    FOREIGN KEY (post_id) REFERENCES post(id),
    FOREIGN KEY (user_id) REFERENCES user(id),
    FOREIGN KEY (comment_id) REFERENCES comment(id),
    CHECK (
        (post_id IS NOT NULL AND comment_id IS NULL) OR 
        (post_id IS NULL AND comment_id IS NOT NULL)
    ),
    UNIQUE (user_id, post_id),
    UNIQUE (user_id, comment_id)
);


	CREATE TABLE IF NOT EXISTS session (
		id INTEGER PRIMARY KEY, 
		uid TEXT UNIQUE NOT NULL, 
		user_id INTEGER UNIQUE NOT NULL, 
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
