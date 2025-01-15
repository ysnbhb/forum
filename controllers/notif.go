package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"forum/utils"
)

func (db *Date) Notif(w http.ResponseWriter, r *http.Request) {
	userId := db.GetIdFromReq(r)
	if userId < 1 {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": http.StatusText(http.StatusUnauthorized)})
		return
	}
	query := ` SELECT   
				user.user_name ,
				notif.post_id , 
				notif.action
				FROM 
					user
				INNER JOIN 
					notif
				ON 
					user.id = notif.user_id
				WHERE 
					notif.to_user_id = ?
				ORDER BY 
					notif.id 
				DESC
				LIMIT ? OFFSET ? ;
				`
	notifs := []utils.Notif{}
	rows, err := db.DB.Query(query, userId, 5, 0)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": http.StatusText(http.StatusInternalServerError)})
		return
	}
	for rows.Next() {
		notif := utils.Notif{}
		rows.Scan(&notif.UserName, &notif.Post_id, &notif.Action)
		notifs = append(notifs, notif)
	}
	json.NewEncoder(w).Encode(notifs)
}
