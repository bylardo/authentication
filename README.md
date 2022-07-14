# authentication
Service that handle the authentication 

# Description
There are 3 endpoints, such as:
/login [POST]
/authenticated [GET]
/refresh-token [GET]

# /login endpoints
Payload:
    {
        "email": [your email], 
        "password" : [your password] 
    }

You can modify, add your data in databases/users.json

# /authenticated endpoints
Check if the token valid or not. Add new params in Header:
Authentication : Token [your_token]

# /refresh-token endpoints
Refresh new token. Add the new param in Header:
Authentication : Token [your_token]

# How to Run
1. go mod authentication.com
2. go get github.com/golang-jwt/jwt/v4
3. go get github.com/gorilla/mux
4. go build
5. go run .