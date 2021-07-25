package handlers

import (
	"database/sql"
	"encoding/json"
	db "github.com/Shubhaankar-sharma/todoapp/db/sqlc"
	"github.com/Shubhaankar-sharma/todoapp/utils"
	"html"
	"log"
	"net/http"
	"strconv"
)

type loginParams struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Login(queries *db.Queries, secret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]string{"error": "something unexpected occurred"}
		decoder := json.NewDecoder(r.Body)
		params := &loginParams{}
		err := decoder.Decode(params)
		if err != nil {
			log.Println(err)
			utils.JSON(w, 500, resp)
			return
		}
		if params.Username == "" || params.Password == "" {
			resp["error"] = "Please Fill All The Fields"
			utils.JSON(w, 400, resp)
			return
		}
		params.Username = html.EscapeString(params.Username)
		user, err := queries.GetUser(r.Context(), params.Username)
		if err != nil {
			if err == sql.ErrNoRows{
				resp["error"] = "no user exists with the username: "+params.Username
				utils.JSON(w, 400, resp)
				return
			}
			log.Println(err)
			utils.JSON(w, 500, resp)
			return
		}
		err = CheckPasswordHash(params.Password, user.Password)
		if err != nil {
			resp["error"] = err.Error()
			utils.JSON(w, 400, resp)
			return
		}
		token, err := utils.EncodeAuthToken(strconv.Itoa(int(user.ID)), secret)
		if err != nil {
			log.Println(err)
			utils.JSON(w, 500, resp)
			return
		}
		resp = map[string]string{"token": token}
		utils.JSON(w, 200, resp)
		return
	}
}
