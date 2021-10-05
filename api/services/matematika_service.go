package services

import(
    "net/http"
    "io/ioutil"
    "encoding/json"    
    // "fmt"
    logger "api/logging"
    dt "api/datastruct"
    "strconv"
    // "os"
)

func GetLuasPersegi(w http.ResponseWriter, r *http.Request)dt.ResponseMath{
    w.Header().Set("Access-Control-Allow-Origin", "*")
	
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	c, errRead := ioutil.ReadAll(r.Body)
	
	defer r.Body.Close()

    var response dt.ResponseMath

	if errRead!=nil{
		logger.Errorf("Error", errRead)
		response.Status = 500
		response.Desc = "error when read parameter request"
		return response
	}

	var msg dt.ParamRequest
		
	errUnmarshal := json.Unmarshal(c, &msg)

	if errUnmarshal != nil{
		logger.Errorf("Error", errUnmarshal)
		response.Status = 502
		response.Desc = "error when unmarshall json param"
		return response 
	}
    
    luas := msg.Sisi * msg.Sisi

    response.Status = 200
    response.Desc = "Luas persegi :" + strconv.Itoa(luas)    

    return response
}