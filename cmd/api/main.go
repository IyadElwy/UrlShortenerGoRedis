package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/teris-io/shortid"
	"urlShortener.IyadElwy/internal/db"
)

type application struct {
	idGenerator *shortid.Shortid
	redis       *db.Redis
	infoLog     *log.Logger
	errorLog    *log.Logger
	envs        map[string]string
}

func main() {
	port := flag.String("port", ":4039", "Port for webserver")
	flag.Parse()

	infoLog := log.New(os.Stdout, "info: ", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "error: ", log.Ldate|log.Ltime)

	sid, err := shortid.New(1, shortid.DefaultABC, 17)
	if err != nil {
		errorLog.Fatal(err)
	}
	redis, err := db.NewRedisClient("localhost:6379", "", 0)
	if err != nil {
		errorLog.Fatal(err)
	}

	envMap, err := godotenv.Read(".env")
	if err != nil {
		errorLog.Fatal(err)
	}

	app := &application{
		idGenerator: sid,
		redis:       redis,
		infoLog:     infoLog,
		errorLog:    errorLog,
		envs:        envMap,
	}
	infoLog.Printf("Server started on port %s", *port)
	err = http.ListenAndServe(*port, app.routes())
	log.Fatal(err)
}
