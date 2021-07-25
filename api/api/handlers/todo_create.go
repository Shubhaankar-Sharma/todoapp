package handlers

import (
	"encoding/json"
	db "github.com/Shubhaankar-sharma/todoapp/db/sqlc"
	"github.com/Shubhaankar-sharma/todoapp/utils"
	"html"
	"log"
	"net/http"
	"time"
)

func CreateToDoHandler(queries *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const shortForm = "2006-Jan-02"
		resp :=  map[string]string{"error": "something unexpected occurred"}
		decoder := json.NewDecoder(r.Body)
		params := &todoCreateParams{}
		err := decoder.Decode(params)
		if err != nil {
			log.Println(err)
			utils.JSON(w, 500, resp)
			return
		}
		if params.Body == "" || params.EndDate == "" {
			resp["error"] = "body and deadline are required"
			utils.JSON(w, 400, resp)
			return
		}
		queryParams := db.CreateToDoParams{}
		queryParams.Body = html.EscapeString(params.Body)
		queryParams.UserID = int32(r.Context().Value("uid").(int))
		queryParams.EndDate, err = time.Parse(shortForm, params.EndDate)
		queryParams.Done = 0
		if err != nil{
			resp["error"] = "date is wrong"
			utils.JSON(w, 400, resp)
			return
		}
		err = queries.CreateToDo(r.Context(), queryParams)
		if err != nil {
			log.Println(err)
			utils.JSON(w, 500, resp)
			return
		}
		resp = map[string]string{"message": "created"}
		utils.JSON(w, 201, resp)
		return
	}
}

type todoCreateParams struct{
	Body string `json:"body"`
	EndDate string `json:"endDate"`
}