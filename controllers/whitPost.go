package controllers

import (
	"fmt"

	"forum/utils"
)

func (db *Date) GetReaction(userid, id int, typ string) utils.Reaction {
	reaction := utils.Reaction{}
	query := `
		SELECT count(*) FROM reaction WHERE %s = ? AND type = 'likes'
	`
	query = fmt.Sprintf(query, typ)
	db.DB.QueryRow(query, id).Scan(&reaction.NumLike)
	query = `
		SELECT count(*) FROM reaction WHERE %s = ? AND type = 'dislikes'
	`
	query = fmt.Sprintf(query, typ)
	db.DB.QueryRow(query, id).Scan(&reaction.NumDisLike)
	if userid < 1 {
		return reaction
	}
	query = `
		SELECT type FROM reaction WHERE %s = ? AND user_id = ?
	`
	query = fmt.Sprintf(query, typ)
	db.DB.QueryRow(query, id, userid).Scan(&reaction.Type)
	return reaction
}

func (db *Date) SaveCategories(postid int, categories []int) {
	query := `
		INSERT INTO categories_post (categorie_id , post_id)
		VALUES(? , ?)
	`
	for _, id := range categories {
		db.DB.Exec(query, id, postid)
	}
}
