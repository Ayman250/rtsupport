package main

import(
"github.com/gorilla/websocket"
r "github.com/dancannon/gorethink"
"net/http"
"fmt"
"time"
"github.com/dgrijalva/jwt-go"

)

type Claims struct {
    Username string `json:"username"`
    // recommended having
    jwt.StandardClaims
}

var upgrader = websocket.Upgrader {
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {return true},
}

type Handler func(*Client, interface{})

type Router struct {
	rules map[string]Handler
	//Database Session
	session *r.Session
}

func NewRouter(session *r.Session) *Router{
	router := &Router{
		rules: make(map[string]Handler),
		session: session,
	}
	return router
}

func (r *Router) Handle(msgName string, handler Handler){
	r.rules[msgName] = handler
}

func (r *Router) FindHandler(msgName string) (Handler, bool) {
	handler, found := r.rules[msgName]
	return handler, found
}

// create a JWT and put in the clients cookie
func setToken(res http.ResponseWriter, req *http.Request) {

    // Expires the token and cookie in 1 hour
    expireToken := time.Now().Add(time.Hour * 1).Unix()
    expireCookie := time.Now().Add(time.Hour * 1)

    // We'll manually assign the claims but in production you'd insert values from a database 
    claims := Claims {
        "myusername",
        jwt.StandardClaims {
            ExpiresAt: expireToken,
            Issuer:    "localhost:4000",
        },
    }

    // Create the token using your claims
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

    // Signs the token with a secret.    
    signedToken, _ := token.SignedString([]byte("secret"))

    // Place the token in the client's cookie 
    cookie := http.Cookie{Name: "Auth", Value: signedToken, Expires: expireCookie, HttpOnly: true}
    http.SetCookie(res, &cookie)
    
    // Redirect the user to his profile
    http.Redirect(res, req, "/profile", 307)
}

// middleware to protect private pages
func validate(protectedPage http.HandlerFunc) http.HandlerFunc {
    return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request){
        
        // If no Auth cookie is set then return a 404 not found
        cookie, err := req.Cookie("Auth")
        if err != nil {
            http.NotFound(res, req)
            return
        }
       
        // Return a Token using the cookie
        token, err := jwt.ParseWithClaims(cookie.Value, &Claims{}, func(token *jwt.Token) (interface{}, error){
            // Make sure token's signature wasn't changed
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, fmt.Errorf("Unexpected siging method")    
            }    
            return []byte("secret"), nil
        })
        if err != nil {
            http.NotFound(res, req)
            return
        }
       
        // Grab the tokens claims and pass it into the original request
        if claims, ok := token.Claims.(*Claims); ok && token.Valid {
            ctx := context.WithValue(req.Context(), MyKey, *claims)
            page(res, req.WithContext(ctx))
        } else {
            http.NotFound(res, req)
            return
        }
    })    
}

// only viewable if the client has a valid token
func protectedProfile(res http.ResponseWriter, req *http.Request){

}

// deletes the cookie
func logout(res http.ResponseWriter, req *http.Request) {

}

func (rout *Router) ServeHTTP(w http.ResponseWriter, r *http.Request){
	socket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		// w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		return
	}
	client := NewClient(socket, rout.FindHandler, rout.session)
	defer client.Close(rout.session)
	go client.Write()
	client.Read()
	
}