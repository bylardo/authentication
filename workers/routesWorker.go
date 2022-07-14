package workers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"authentication.com/helpers"
	"authentication.com/models"
	"github.com/golang-jwt/jwt/v4"
)

// Secret Key used
var jwtKey = []byte("9182yuihkjni1p2o3jl123kmuo1i2j3jl123kn")

/**
DoLogin is method to handle the login process. The flow are:
1. validate the request body structure (email and password)
2. Validate the request body with data in database that store in JSON File (users.json)
3. Create the JWT Token with 60 mins expiration time
4. Send the JWT token into response if all the process success
**/
func DoLogin(w http.ResponseWriter, r *http.Request) {
	// Declare credential variable to decoded from Request Body
	var credential models.Credentials

	// Get the request body payloads and decode into credential variables
	err := json.NewDecoder(r.Body).Decode(&credential)

	// If the request body structure is wrong, then send the response as Bad Gateway
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// User Validation
	// If Exist in database, then return true, If Not return false
	valid := ValidateUserLogin(w, credential)

	if !valid {
		return
	}

	// Create Token Expiration Time
	// Expiration Time is 60 minutes
	tokenExpiredTime := time.Now().Add(60 * time.Minute)

	// Initialize JWT Claims
	// Contains Email and Expiration Time
	claims := &models.Claims{
		Email: credential.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: tokenExpiredTime.Unix(),
		},
	}

	//Declare new JWT Token with HS256 algorithm for signing
	tokenNew := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT Token with Secret Key
	tokenStr, err := tokenNew.SignedString(jwtKey)

	if err != nil {
		// Prepare the API Response Data
		apiResponse := helpers.PrepareAPIResponse(http.StatusInternalServerError, "Internal Server Error ! Please check log", nil)
		// Sending the API Response
		helpers.SendAPIResponse(apiResponse, w, http.StatusInternalServerError)
		// Log the error
		fmt.Println("Error when Generate JWT")
		return
	}

	// Send the response to client if Succeed to create JWT
	// Send the JWT into Response Body
	apiResponse := helpers.PrepareAPIResponse(http.StatusOK, "success", tokenStr)
	helpers.SendAPIResponse(apiResponse, w, http.StatusOK)
	return
}

/**
DoAuthenticationRequest is method to handle the login process. The flow are:
1. validate the header. If header has Authentication Param
2. Validate the Token from Authentication Header Value
3. JWT Token Validation using claims
**/
func DoAuthenticationRequest(w http.ResponseWriter, r *http.Request) {

	// Header Sanitation
	// If Header has Authorization Param, then validated. If not, then invalid
	headerValidation, clientToken := ValidateHeader(r)

	if !headerValidation {
		apiResponse := helpers.PrepareAPIResponse(http.StatusUnauthorized, "Unauthorized", nil)
		helpers.SendAPIResponse(apiResponse, w, http.StatusUnauthorized)
		fmt.Println("No Authorization Header Param")
		return
	}

	// Declare jwt claims
	claims := &models.Claims{}

	// Parse the JWT String from Header and assigned to claims struct
	// Parse process result the token and errors
	jwtToken, err := jwt.ParseWithClaims(clientToken, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		// Condition when Signature Invalid
		if err == jwt.ErrSignatureInvalid {
			apiResponse := helpers.PrepareAPIResponse(http.StatusUnauthorized, "Unauthorized", nil)
			helpers.SendAPIResponse(apiResponse, w, http.StatusUnauthorized)
			fmt.Println("Signature Invalid")
			return
		}

		apiResponse := helpers.PrepareAPIResponse(http.StatusBadRequest, "Unauthorized", nil)
		helpers.SendAPIResponse(apiResponse, w, http.StatusBadRequest)
		fmt.Println("Bad Request")
		return
	}

	// Condition when JWT Token Invalid
	if !jwtToken.Valid {
		apiResponse := helpers.PrepareAPIResponse(http.StatusUnauthorized, "Unauthorized", nil)
		helpers.SendAPIResponse(apiResponse, w, http.StatusUnauthorized)
		fmt.Println("Token Invalid")
		return
	}

	// Prepare API Response and Send the Response to Client
	apiResponse := helpers.PrepareAPIResponse(http.StatusOK, "Welcome Authorized User", nil)
	helpers.SendAPIResponse(apiResponse, w, http.StatusOK)
	fmt.Println("Token Valid!")
}

/**
DoHandleRefreshToken is method to handle the refresing token process. The flow are:
1. validate the header. If header has Authentication Param
2. Validate the Token from Authentication Header Value
3. If token valid from Header, then create new token
4. Create new token after 15 seconds before old token expired
**/
func DoHandleRefreshToken(w http.ResponseWriter, r *http.Request) {
	// Header Sanitation
	// If Header has Authorization Param, then validated. If not, then invalid
	headerValidation, clientToken := ValidateHeader(r)

	if !headerValidation {
		apiResponse := helpers.PrepareAPIResponse(http.StatusUnauthorized, "Unauthorized", nil)
		helpers.SendAPIResponse(apiResponse, w, http.StatusUnauthorized)
		fmt.Println("No Authorization Header Param")
		return
	}

	// Declare jwt claims
	claims := &models.Claims{}

	// Parse the JWT String from Header and assigned to claims struct
	// Parse process result the token and errors
	jwtToken, err := jwt.ParseWithClaims(clientToken, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		// Condition when Signature Invalid
		if err == jwt.ErrSignatureInvalid {
			apiResponse := helpers.PrepareAPIResponse(http.StatusUnauthorized, "Unauthorized", nil)
			helpers.SendAPIResponse(apiResponse, w, http.StatusUnauthorized)
			fmt.Println("Signature Invalid")
			return
		}

		apiResponse := helpers.PrepareAPIResponse(http.StatusBadRequest, "Unauthorized", nil)
		helpers.SendAPIResponse(apiResponse, w, http.StatusBadRequest)
		fmt.Println("Bad Request")
		return
	}

	// Condition when JWT Token Invalid
	if !jwtToken.Valid {
		apiResponse := helpers.PrepareAPIResponse(http.StatusUnauthorized, "Unauthorized", nil)
		helpers.SendAPIResponse(apiResponse, w, http.StatusUnauthorized)
		fmt.Println("Token Invalid")
		return
	}

	// We ensure that a new token is not issued until enough time has elapsed
	// In this case, a new token will only be issued if the old token is within  15 seconds of expiry. Otherwise, return a bad request status
	// This is optional, If you want to remove this, then it is fine
	if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) > 15*time.Second {
		// Return the Response if an error happened when create the new JWT Token
		apiResponse := helpers.PrepareAPIResponse(http.StatusInternalServerError, "Internal Server Error", nil)
		helpers.SendAPIResponse(apiResponse, w, http.StatusInternalServerError)
		fmt.Println("Old Token is active more than 15 seconds")
		return
	}

	// Create a new token for the current use, with a renewed expiration time
	expirationTime := time.Now().Add(60 * time.Minute)
	claims.ExpiresAt = expirationTime.Unix()
	JWTtoken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := JWTtoken.SignedString(jwtKey)

	if err != nil {
		// Return the Response if an error happened when create the new JWT Token
		apiResponse := helpers.PrepareAPIResponse(http.StatusInternalServerError, "success", tokenStr)
		helpers.SendAPIResponse(apiResponse, w, http.StatusInternalServerError)
		fmt.Println("Error when create the new token")
		return
	}

	// Send the response to client if Succeed to create JWT
	// Send the JWT into Response Body
	apiResponse := helpers.PrepareAPIResponse(http.StatusOK, "success", tokenStr)
	helpers.SendAPIResponse(apiResponse, w, http.StatusOK)
}
