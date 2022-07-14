package workers

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"authentication.com/helpers"
	"authentication.com/models"
)

func ValidateUserLogin(w http.ResponseWriter, credential models.Credentials) bool {

	userDBPath := "./databases/users.json"
	valid, userDatabase := helpers.ReadDatabase(userDBPath)

	if !valid {
		apiResponse := helpers.PrepareAPIResponse(http.StatusInternalServerError, "Internal Server Error! Please check the logs", nil)
		helpers.SendAPIResponse(apiResponse, w, http.StatusInternalServerError)
		fmt.Println("Error when Reading Data from JSON File, [users.json]")
		return false
	}

	var users []models.User
	var err error

	err = json.Unmarshal(userDatabase, &users)

	if err != nil {
		apiResponse := helpers.PrepareAPIResponse(http.StatusInternalServerError, "Internal Server Error! Please check the logs", nil)
		helpers.SendAPIResponse(apiResponse, w, http.StatusInternalServerError)
		fmt.Println("Error when Assign Data from JSON File to variables, [users.json]")
		return false
	}

	var userExist bool
	userExist = false
	for _, user := range users {
		if user.Email == credential.Email {
			var sha = sha1.New()
			sha.Write([]byte(credential.Password))
			encrypPassword := sha.Sum(nil)
			encryptedString := fmt.Sprintf("%x", encrypPassword)
			if user.Password == string(encryptedString) {
				userExist = true
				break
			}
		}
	}

	if !userExist {
		apiResponse := helpers.PrepareAPIResponse(http.StatusUnauthorized, "Email and Password Incorrect!", nil)
		helpers.SendAPIResponse(apiResponse, w, http.StatusInternalServerError)
		fmt.Println("User Not Exist")
		return false
	}
	return true
}

func ValidateHeader(r *http.Request) (bool, string) {

	reqToken := r.Header.Get("Authorization")
	if reqToken == "" {
		return false, ""
	}

	splitToken := strings.Split(reqToken, "Token ")
	reqToken = splitToken[1]

	if len(splitToken) > 0 {
		return true, reqToken
	}
	return false, ""
}
