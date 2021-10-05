package services

import(
    // "database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"	
	// "fmt"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	// "errors"
	dt "api/datastruct"
)


func TestModifyData(t *testing.T) {
	bodyData := new(dt.RequestModify)

	bodyData.Method = "insert"
	bodyData.Data.ID = "W009"
	bodyData.Data.Name = "Daniel"
	bodyData.Data.Age = "2"
	bodyData.Data.Grade = "3"

	marshalBody,_ := json.Marshal(bodyData)
	
	mockRequest,_ := http.NewRequest("POST","/modifyData",bytes.NewBuffer(marshalBody))

	//insert
	t.Run("Insert success",func(t *testing.T){
		db, mock, err := sqlmock.New()	
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		// mock.ExpectBegin()
		// mock.ExpectExec("UPDATE products").WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec("insert into tb_student").WithArgs(bodyData.Data.ID,bodyData.Data.Name,bodyData.Data.Age,bodyData.Data.Grade).WillReturnResult(sqlmock.NewResult(1, 1))
		// jika mau balikin error, pakai WillReturnError
		// mock.ExpectCommit()

		// now we execute our method		 
		expectedResult := dt.ResponseData{Status:200,Data : "sukses insert"}

		actualResp := ModifyData(httptest.NewRecorder(),mockRequest,db)

		// we make sure that all expectations were met
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
		
		assert.Equal(t,expectedResult,actualResp) 
	})
}

func TestGetSingleStud(t *testing.T){
	// bodyData := new(dt.RequestModify)

	// mockRequest,_ := http.NewRequest("POST","/getSingleStud",bytes.NewBuffer(marshalBody))
	// fmt.Println("cek0")
	// fmt.Println(mockRequest)

	t.Run("Get single student",func(t *testing.T){

		param := map[string]interface{}{
			"param" : "E001",
		}
		
		marshalBody,_ := json.Marshal(param)

		mockRequest,_ := http.NewRequest("POST","/getSingleStud",bytes.NewBuffer(marshalBody))

		db,mock,err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		rows := sqlmock.NewRows([]string{"id", "name","age","grade"}).
		AddRow("E001", "daniel","4","3")

		// mock.ExpectQuery("select id,name,age,grade from tb_student where id = ?").WithArgs(param["param"]).WillReturnRows(rows)

		mock.ExpectQuery("select id,name,age,grade from tb_student where id = ?").WithArgs(param["param"]).WillReturnRows(rows)

		expectedResult := dt.ResponseData{Status:200,Data : dt.Student{
			ID : "E001",
			Name : "daniel",
			Age : "4",
			Grade : "3",
		}}

		// errStr := errors.New("0 - z")
		// expectedResult := dt.ResponseData{Status:500,Data : "no result",Error:nil}
		
		actualResp := GetSingleStud(httptest.NewRecorder(),mockRequest,db)

		// fmt.Println(mockRequest,"11")
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}		
		assert.Equal(t,expectedResult,actualResp)		
	})

	t.Run("Get single student no rows",func(t *testing.T){
		param := map[string]interface{}{
			"param" : "z",
		}
		
		marshalBody,_ := json.Marshal(param)

		mockRequest,_ := http.NewRequest("POST","/getSingleStud",bytes.NewBuffer(marshalBody))
		// fmt.Println("cek")
		db,mock,err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		// rows := sqlmock.NewRows([]string{"id", "name","age","grade"}).
		// AddRow("E001", "daniel","4","3")

		// mock.ExpectQuery("select id,name,age,grade from tb_student where id = ?").WithArgs(param["param"]).WillReturnRows(rows)

		mock.ExpectQuery("select id,name,age,grade from tb_student where id = ?").WithArgs(param["param"]).WillReturnError(nil)

		// expectedResult := dt.ResponseData{Status:200,Data : dt.Student{
		// 	ID : "E001",
		// 	Name : "daniel",
		// 	Age : "4",
		// 	Grade : "3",
		// }}

		// errStr := errors.New("0 - z")
		expectedResult := dt.ResponseData{Status:500,Data : "no result"}

		// fmt.Println(mockRequest)
		// fmt.Println(param)
		actualResp := GetSingleStud(httptest.NewRecorder(),mockRequest,db)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
		
		assert.Equal(t,expectedResult,actualResp)

	})

}

// func TestModifyData(t *testing.T) {
// 	bodyData := new(dt.RequestModify)

// 	bodyData.Method = "update"
// 	bodyData.Data.ID = "W009"
// 	bodyData.Data.Name = "Daniel"
// 	bodyData.Data.Age = "2"
// 	bodyData.Data.Grade = "3"

// 	marshalBody,_ := json.Marshal(bodyData)
	
// 	mockRequest,_ := http.NewRequest("POST","/modifyData",bytes.NewBuffer(marshalBody))

// 	//update
// 	t.Run("Update success",func(t *testing.T){
// 		db, mock, err := sqlmock.New()	
// 		if err != nil {
// 			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
// 		}
// 		defer db.Close()

// 		mock.ExpectExec("insert into tb_student").WithArgs("W009","Daniel","2","3").WillReturnResult(sqlmock.NewResult(1, 1))
// 		// jika mau balikin error, pakai WillReturnError
// 		// mock.ExpectCommit()

// 		// now we execute our method
// 		expectedResult := dt.ResponseData{Status:200,Data : "sukses insert"}

// 		actualResp := ModifyData(httptest.NewRecorder(),mockRequest,db)

// 		// we make sure that all expectations were met
// 		if err := mock.ExpectationsWereMet(); err != nil {
// 			t.Errorf("there were unfulfilled expectations: %s", err)
// 		}
		
// 		assert.Equal(t,expectedResult,actualResp) 
// 	})
// }

