package main

import (
	"log"
	"net/http"
	"os"
    "fmt"
	//"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
)

func determineListenAddress() (string, error) {
	port := os.Getenv("PORT")
	if port == "" {
	  return "", fmt.Errorf("$PORT not set")
	}
	return ":" + port, nil
  }


  func hello(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/javascript")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method == "POST"{
		r.ParseForm()

		id:=r.Form["id"][0]
        fmt.Fprintln(w, "Here it is: "+id)
	}else{
		fmt.Fprintln(w, "not post")
	}
	
  }



  func main() {

	addr, err := determineListenAddress()
	if err != nil {
	  log.Fatal(err)
	}

	http.HandleFunc("/", hello)

	log.Printf("Listening on %s...\n", addr)

	if err := http.ListenAndServe(addr, nil); err != nil {
	  panic(err)
	}

  }