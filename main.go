package main 

import (
    "log"
    "net/http"
    "square-8-challenge/router"
)

func main() {
    

    r := router.Router()

    // listener 
    log.Fatal(http.ListenAndServe(":8080", r))
}
