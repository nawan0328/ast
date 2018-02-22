package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	cron "gopkg.in/robfig/cron.v2"

	"github.com/ant0ine/go-json-rest/rest"
)

func main() {

	var err error

	db, err = sql.Open("mysql", "dbname:passwd@tcp(127.0.0.1:3306)/stock")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	c := cron.New()
	c.AddFunc("0 0 17 * * 1-5", getPrice_history)
	go c.Start()

	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	router, err := rest.MakeRouter(
		rest.Get("/rule", GetStock),
	)
	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(router)
	log.Fatal(http.ListenAndServe(":8080", api.MakeHandler()))

}

func GetStock(w rest.ResponseWriter, r *rest.Request) {
	stock = Stock{} //Initialization
	w.Header().Set("Access-Control-Allow-Origin", "*")
	stock.Date = r.URL.Query().Get("date")
	if stock.Date == "" {
		stock.Date = time.Now().Format("20060102")
	}
	getRule()
	w.WriteJson(&stock)
}
