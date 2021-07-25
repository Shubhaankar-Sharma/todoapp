package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	db "github.com/Shubhaankar-sharma/todoapp/db/sqlc"
	"github.com/Shubhaankar-sharma/todoapp/utils"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"html"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func CreateUserHandler(queries *db.Queries, secret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]string{"error": "something unexpected occurred"}
		decoder := json.NewDecoder(r.Body)
		params := &db.CreateUserParams{}
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
		err = BeforeSave(params)
		if err != nil {
			log.Println(err)
			utils.JSON(w, 500, resp)
			return
		}
		userId, err := queries.CreateUser(r.Context(), *params)
		if err != nil {
			if pqErr, ok := err.(*pq.Error); ok {
				switch pqErr.Code.Name() {
				case "foreign_key_violation", "unique_violation" :
					resp["error"] = fmt.Sprintf("user with the username: %s already exists",
						params.Username)
					utils.JSON(w, 400, resp)
					return
				}
			}
			log.Println(err)
			utils.JSON(w, 500, resp)
			return
		}
		token, err := utils.EncodeAuthToken(strconv.Itoa(int(userId)), secret)
		if err != nil {
			log.Println(err)
			utils.JSON(w, 500, resp)
			return
		}
		resp = map[string]string{"message": "user created: " + params.Username,
								 "token": token}
		utils.JSON(w, 201, resp)
		return
	}
}

// HashPassword hashes password from user input
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14) // 14 is the cost for hashing the password.
	return string(bytes), err
}

// CheckPasswordHash checks password hash and password from user input if they match
func CheckPasswordHash(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return errors.New("password is incorrect")
	}
	return nil
}

func BeforeSave(params *db.CreateUserParams) error {
	password := strings.TrimSpace(params.Password)
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return err
	}
	params.Password = hashedPassword
	params.Username = html.EscapeString(params.Username)
	return nil
}
