package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/aghaghiamh/gocast/QAGame/controller"
)

func main() {
	http.HandleFunc("/user/register", controller.UserRegisterHandler)
	http.HandleFunc("/user/login", controller.UserLoginHandler)
	http.HandleFunc("/user/auth", controller.UserAuthHandler)

	fmt.Print("Server is running on port 8080...")
	if err := http.ListenAndServe("localhost:8080", nil); err != nil {
		log.Fatalf("Couldn't Listen to the 8080 port.")
	}
}
