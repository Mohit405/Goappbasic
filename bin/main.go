package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	// "github.com/julienschmidt/httprouter"
	"github.com/mohit405/Redis"
	"github.com/mohit405/config"
	"github.com/mohit405/mysql"
)

type Conn struct {
	port int
	env  string
}

type application struct {
	conn     Conn
	logger   *log.Logger
	ErrorLog *log.Logger
	sqlconn  *mysql.Storage
	rediconn *Redis.Storage
}

func main() {
	var con Conn

	flag.IntVar(&con.port, "port", 8080, "Api server port")
	flag.StringVar(&con.env, "env", "development", "Environment(development|staging|production)")
	flag.Parse()

	logger1 := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	redis := config.Redis{
		Addr: "localhost:6379",
		Db:   0,
	}

	sqlConn := config.Mysql{
		Dialect: "mysql",
		DSN:     "root:mypassword@tcp(127.0.0.1)/practise",
	}

	mysqlConn, err := mysql.NewStorage(sqlConn)
	if err != nil {
		return
	}

	redisConn, err := Redis.MakeConnection(redis)
	if err != nil {
		return
	}

	app := &application{
		conn:     con,
		logger:   logger1,
		sqlconn:  mysqlConn,
		rediconn: redisConn,
	}
	srv := &http.Server{
		Addr:        fmt.Sprintf(":%d", con.port),
		Handler:     app.routes(),
		IdleTimeout: time.Minute,
	}

	err = srv.ListenAndServe()
	logger1.Fatal(err)
}
