package main

import (
	"fmt"
	"net/http"
	"os"
	"github.com/gorilla/mux"
	"log"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"

	services "api/services"
)

const (
	strPostGre = "host=127.0.0.1 port=5432 user=admin password=4dm!n dbname=glng sslmode=disable"
	strMysql = "username:password@tcp(127.0.0.1:3306)/kampus-merdeka?charset=utf8"
)

func initHandlers() {
	rtr := mux.NewRouter()
	fmt.Println("Start Processs...")	

	connMysq,errConn := ConnectDb("mysql")

	if errConn != nil {
		fmt.Println(errConn)
		return
	}

	// defer connMysq.Close()

	var method string

	rtr.HandleFunc("/hitungLuas", func(w http.ResponseWriter, r *http.Request) {
		rsp := services.GetLuasPersegi(w, r)
		services.SendResponses(w, rsp)
	}).Methods("POST")

	rtr.HandleFunc("/getStudents", func(w http.ResponseWriter, r *http.Request) {
		rsp := services.GetStudents()
		services.SendResponses(w, rsp)
	}).Methods("GET")

	rtr.HandleFunc("/preparedStudent", func(w http.ResponseWriter, r *http.Request) {
		rsp := services.PreparedStudents()
		services.SendResponses(w, rsp)
	}).Methods("GET")

	rtr.HandleFunc("/getSingleStud", func(w http.ResponseWriter, r *http.Request) {
		rsp := services.GetSingleStud(w,r,connMysq)
		services.SendResponses(w, rsp)
	}).Methods("POST")

	rtr.HandleFunc("/modifyData", func(w http.ResponseWriter, r *http.Request) {
		rsp := services.ModifyData(w,r,connMysq)
		services.SendResponses(w, rsp)
	}).Methods("POST")
	
	http.Handle("/", rtr)
	fmt.Println("start service...",method)
	
}

func main() {
	initHandlers()

	var err error

	err = http.ListenAndServe(":8181",nil)

	if err != nil {
		// log.Errorf("Unable to start the server %v", err)
		log.Println("Failed start service")
		os.Exit(1)
	}

}

func ConnectDb(paramDb string)(*sql.DB,error){
	var err error
	var db *sql.DB	

	if paramDb == "postgres"{
		db,err = sql.Open("postgres",strPostGre)
		
		if err != nil {
			log.Println(err)
			return nil,err
		}

		err = db.Ping()

		if err != nil {
			log.Println(err)
			return nil,err
		}
	}else if paramDb == "mysql"{
		db,err = sql.Open("mysql",strMysql)

		if err != nil {
			return nil,err
		}
	}
	// defer db.Close()	
	return db,nil
}

