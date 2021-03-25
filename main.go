package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {

	resp, err := http.Get("https://www.olx.co.id/item/dijual-avanza-veloz-2013-iid-821749571")
	if err != nil {
		log.Fatal(err)
	}

	// Print the HTTP Status Code and Status Name
	fmt.Println("HTTP Response Status:", resp.StatusCode, http.StatusText(resp.StatusCode))

	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		fmt.Println("HTTP Status is in the 2xx range")
	} else {
		fmt.Println("Argh! Broken")
	}

}
