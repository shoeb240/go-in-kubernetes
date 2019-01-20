package main

import (
    "fmt"
    "github.com/gorilla/mux"
    "log"
    "net/http"
)

func main(){
    fmt.Println("Waiting for request to serve...")

    myRouter := mux.NewRouter()
    myRouter.HandleFunc("/users", getUsers).Methods("GET")

    log.Fatal(http.ListenAndServe(":8081", myRouter))
}

func getUsers(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "I am responding to your API call")
}
