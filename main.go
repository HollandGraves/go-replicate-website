package main

import (
	"fmt"
	"net/http"
	"os"
)

// MAIN FUNCTION

func main() {
	resp, err := http.Get("http://google.com")
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	bs := make([]byte, 99999)
	resp.Body.Read(bs)
	fmt.Println(string(bs))

	err = os.Mkdir("http://google.com", 0666)
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
}

// PERSONAL DEFINED FUNCTIONS
