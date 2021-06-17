package cmd

import (
	"io/ioutil"
	"log"
	"net/http"
)

func Expenses(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	log.Printf("body: %s, err: %v", b, err)

	// TODO: implement DB per customer
	// HandleAdminSQL(w, r)
}
