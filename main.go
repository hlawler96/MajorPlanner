package main

import (
    "log"
    "net/http"


)


func main() {
    router := NewRouter()
//     c := cors.New(cors.Options{
//     AllowedMethods: []string{"GET","POST", "OPTIONS"},
//     AllowedOrigins: []string{"*"},
//     AllowCredentials: true,
//     AllowedHeaders: []string{"Content-Type","Bearer","Bearer ","content-type","Origin","Accept"},
//     OptionsPassthrough: true,
// })
    // handler := c.Handler(router)

    log.Fatal(http.ListenAndServe(":8080", router))
}
