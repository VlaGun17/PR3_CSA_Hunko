package main

import (
	"log"
	"PR3_CSA_Hunko/server"
)

func main() {
	if err := server.ListenAndServe(":8443"); err != nil {
		log.Fatal(err)
	}
}
