package datastruct

import(
	"database/sql"
)

type Student struct{
	ID		string `json:"id"`
	Name	string `json:"name"`
	Age 	string `json:"age"`
	Grade 	string `json:"grade"`
}

type Student1 struct{
	ID		sql.NullString `json:"id"`
	Name	sql.NullString `json:"name"`
	Age 	sql.NullString `json:"age"`
	Grade 	sql.NullString `json:"grade"`
}

type RequestModify struct {
	Method	string `json:"method"`
	Data	Student `json:"data"`
}

type RequestSingle struct{
	Param string `json:"param"`
}