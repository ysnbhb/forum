package controllers

import (
	"fmt"
	"net/http"

	"forum/utils"
)

func (db *Date) GetReaction(r *http.Request, id int, which int, typ string) utils.Reaction {
	reaction := utils.Reaction{}
	query := `
		SELECT count(*) FROM reaction WHERE %s = ? AND type = likes
	`
	query = fmt.Sprintf(query, typ)
	db.DB.QueryRow(query, which).Scan(&reaction.NumLike)
	query = `
		SELECT count(*) FROM reaction WHERE %s = ? AND type = dislikes
	`
	query = fmt.Sprintf(query, typ)
	db.DB.QueryRow(query, which).Scan(&reaction.NumDisLike)

	return reaction
}
