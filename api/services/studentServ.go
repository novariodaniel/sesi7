package services

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"encoding/json"
	// "fmt"
	"io/ioutil"
	"log"
	"net/http"	
	
	dt "api/datastruct"	
	// logger "api/logging"
)

const (
	strPostGre = "host=127.0.0.1 port=5432 user=admin password=4dm!n dbname=glng sslmode=disable"
	strMysql = "username:password@tcp(127.0.0.1:3306)/kampus-merdeka?charset=utf8"
)

func GetStudents()dt.ResponseData {
	var resp dt.ResponseData

	connMysq,errConn := ConnectDb("mysql")

	if errConn != nil {
		log.Println(errConn)
		resp.Status = 502
		resp.Data = "Connection to database failed"
		return resp
	}

	defer connMysq.Close()

	var listStudent []dt.Student

	age := 27

	activeRecord,errQuery := connMysq.Query("select id,name,grade from tb_student where age = ?",age)

	defer activeRecord.Close()

	if errQuery != nil {
		log.Println(errQuery)
		resp.Status = 500
		resp.Data = "Query Error"
		return resp
	}

	for activeRecord.Next(){
		var eachStudent dt.Student

		var errScan = activeRecord.Scan(&eachStudent.ID,&eachStudent.Name,&eachStudent.Grade)

		if errScan != nil {
			log.Println(errScan)
			resp.Status = 500
			resp.Data = "Error scan data"
			return resp
		}

		listStudent = append(listStudent,eachStudent)		
	}

	resp.Status = 200
	resp.Data = listStudent
	return resp
}

func PreparedStudents()dt.ResponseData{
	var resp dt.ResponseData

	connMysq,errConn := ConnectDb("mysql")

	defer connMysq.Close()

	if errConn != nil {
		log.Println(errConn)
		resp.Status = 502
		resp.Data = "Connection to database failed"
	}	

	var singleStudent dt.Student	

	activeRecord, errQuery := connMysq.Prepare("select id,name,grade from tb_student where id = ?")

	defer activeRecord.Close()

	if errQuery != nil{
		log.Println(errQuery)
		resp.Status = 500
		resp.Data = "Query Error"
	}
	
	activeRecord.QueryRow("W003").Scan(&singleStudent.ID,&singleStudent.Name,&singleStudent.Grade)

	resp.Status = 200
	resp.Data = singleStudent
	return resp
}

func GetSingleStud(w http.ResponseWriter, r *http.Request, connMysq *sql.DB)dt.ResponseData{
	c, errRead := ioutil.ReadAll(r.Body)

	defer r.Body.Close()

	var response dt.ResponseData

	if errRead != nil {
		log.Println(errRead)
		response.Status = 500
		response.Data = errRead
		return response
	}

	var msg dt.RequestSingle

	errUnmarshal := json.Unmarshal(c,&msg)

	if errUnmarshal != nil {
		log.Println(errUnmarshal)
		response.Status = 502
		response.Data = "Unmarshall failed"
		return response
	}

	// connMysq,errConn := ConnectDb("mysql")

	// defer connMysq.Close()

	// if errConn != nil {
	// 	log.Println(errConn)
	// 	response.Status = 502
	// 	response.Data = "Connection to database failed"
	// }	

	var result dt.Student1

	id := msg.Param
	
	connMysq.QueryRow("select id,name,age,grade from tb_student where id = ?",id).Scan(&result.ID,&result.Name,&result.Age,&result.Grade)
	
	// log.Println(result.ID)

	if !result.ID.Valid {
		// log.Println(err)		
		response.Status = 500
		response.Data = "no result"
		return response
	}

	var resp dt.Student
	resp.ID = result.ID.String
	resp.Name = result.Name.String
	resp.Age = result.Age.String
	resp.Grade = result.Grade.String

	response.Status = 200
	response.Data = resp
	return response

}

func ModifyData(w http.ResponseWriter, r *http.Request,connMysq *sql.DB)dt.ResponseData{
	
	//1

	c, errRead := ioutil.ReadAll(r.Body)

	defer r.Body.Close()

	var response dt.ResponseData

	if errRead != nil {
		log.Println(errRead)
		response.Status = 500
		response.Data = errRead
		return response
	}
	//2
	var msg dt.RequestModify

	errUnmarshal := json.Unmarshal(c,&msg)

	if errUnmarshal != nil {
		log.Println(errUnmarshal)
		response.Status = 502
		response.Data = "Unmarshall failed"
		return response
	}

	//3
	// connMysq,errConn := ConnectDb("mysql")	

	// if errConn != nil {
	// 	log.Println(errConn)
	// 	response.Status = 502
	// 	response.Data = "Connection to database failed"
	// 	return response
	// }	

	var errExec error

	if msg.Method == "insert"{
		_,errExec = connMysq.Exec("insert into tb_student value(?,?,?,?)",msg.Data.ID,msg.Data.Name,msg.Data.Age,msg.Data.Grade)
	}else if msg.Method == "update"{
		log.Println("sini")
		_,errExec = connMysq.Exec("update tb_student set name= ?, age = ?, grade = ? where id = ?",msg.Data.Name,msg.Data.Age,msg.Data.Grade,msg.Data.ID)
	}else if msg.Method == "delete"{
		_,errExec = connMysq.Exec("delete from tb_student where id = ?",msg.Data.ID)
	}else{
		log.Println("Error, method info not available")
		response.Status = 500
		// response.Data = "error when running modify data"
		response.Error = errExec
		return response
	}

	if errExec != nil {
		log.Println(errExec)
		response.Status = 502
		// response.Data= "error when running modify statement"
		response.Error = errExec		
	}

	response.Status = 200
	response.Data = "sukses "+msg.Method
	
	return response	

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
	return db,nil
}