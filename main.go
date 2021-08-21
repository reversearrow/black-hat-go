package main

import (
	"context"
	"fmt"
	"github.com/urfave/negroni"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type logger struct{
	Inner http.Handler
}

func (l *logger) ServeHTTP(w http.ResponseWriter, req *http.Request){
	log.Println("Start of the Request")
	l.Inner.ServeHTTP(w, req)
	log.Printf("End of the Request")
}

func hello(w http.ResponseWriter, req *http.Request){
	fmt.Fprintf(w, "hello %s\n", req.URL.Query().Get("name"))
}

type router struct{
}

func(r *router) ServeHTTP(w http.ResponseWriter, req *http.Request){
	switch req.URL.Path{
	case "/a":
		fmt.Fprint(w, "executing a")
	default:
		http.Error(w,"not found", 404)
	}
}

type trivial struct{
}

func (t *trivial) ServeHTTP(w http.ResponseWriter, req *http.Request, next http.HandlerFunc){
	fmt.Println("Executing Trivial Middleware")
	fmt.Println(req.Context().Value("username"))
	next(w, req)
}

type badAuth struct{
	username string
	password string
}

func (a *badAuth) ServeHTTP(w http.ResponseWriter, req *http.Request, next http.HandlerFunc){
	fmt.Println("executing auth middleware")
	username := req.URL.Query().Get("username")
	password := req.URL.Query().Get("password")

	if username != a.username && password != a.password{
		http.Error(w, "unauthorized",401)
		return
	}
	ctx := context.WithValue(req.Context(), "username", username)
	req = req.WithContext(ctx)

	next(w, req)
}

func main(){
	r := mux.NewRouter()
	r.HandleFunc("/hello", hello).Methods("GET")
	n := negroni.Classic()
	n.Use(&badAuth{
		username: "admin",
		password: "password",
	})
	n.Use(&trivial{})
	n.UseHandler(r)
	http.ListenAndServe(":8080",n)
}
