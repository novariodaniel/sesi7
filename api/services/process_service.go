package services

import (
	// "bytes"
	// "crypto/sha256"
	// "crypto/tls"
	// "database/sql"
	// b64 "encoding/base64"
	// "encoding/hex"
	"encoding/json"
	// "errors"
	"fmt"
	// "io"
	"io/ioutil"
	// "log"
	"net/http"
	"os"	
	"bufio"			
	// "runtime"
	"strconv"
	// "sync"
	"strings"
	// "time"
	// "math"

	// "github.com/pkg/sftp"

	// "golang.org/x/crypto/ssh"

	// conf "api/config"	
	dt "api/datastruct"
	// lib "api/lib"	
	log "api/logging"
)

var ListStoredData []*dt.StoredData

func StartProcess(w http.ResponseWriter, r *http.Request)dt.ResponseProcess{

	var resp dt.ResponseProcess

	fmt.Println("tes")

	// //get invoice from FBL5N
	// listInvoice := getFixAmount()

	// //cek each invoice is overdue?
	// overdueInvoice := cekOverdue(listInvoice)

	// fmt.Println(len(overdueInvoice))

	// // insertInvoice := bulkInsertInvoice(overdueInvoice)

	// // fmt.Println(insertInvoice)

	// // os.Exit(71)

	// //get config denda based on invoice/cluster/cmd/channel
	// configDenda := getConfDenda(overdueInvoice)

	// fmt.Println(len(configDenda))

	// //cek if invoice is exist or not in custom_invoice table?
	// listCustom := getCustomInvoice(configDenda)

	// //get detail transaction of denda
	// listTransDenda := getTransDenda()

	// dendaProcessing := calculateDenda(listTransDenda,listCustom)

	// fmt.Println(dendaProcessing)

	// //collect data must be processing
	// // dendaProcessing := mustProcessing(listMax,listCustom)

	// resp.Status = dt.SuccessResponse
	// resp.Desc = dt.DescSuccess
	// // resp.Desc = dendaProcessing
	return resp

}

func GetEmpty()dt.ResponseProcess{
	var resp dt.ResponseProcess
	fmt.Println("get empty")

	// resp.Status = 200
	resp.Desc = "Hello Go developers"
	
	return resp
}

func GetLanguage()dt.DetailLanguage{
	fmt.Println("get language")
	var structInput dt.DetailLanguage
	var influenStru dt.Influen
	
	influenStru.InfluencedBy = append(influenStru.InfluencedBy,"B")
	influenStru.InfluencedBy = append(influenStru.InfluencedBy,"ALGOL 68")
	influenStru.InfluencedBy = append(influenStru.InfluencedBy,"Assembly")
	influenStru.InfluencedBy = append(influenStru.InfluencedBy,"FORTRAN")

	influenStru.Influences = append(influenStru.Influences,"C++")
	influenStru.Influences = append(influenStru.Influences,"Objective-C")
	influenStru.Influences = append(influenStru.Influences,"C#")
	influenStru.Influences = append(influenStru.Influences,"Java")
	influenStru.Influences = append(influenStru.Influences,"Javascript")
	influenStru.Influences = append(influenStru.Influences,"PHP")
	influenStru.Influences = append(influenStru.Influences,"Go")

	structInput.Language = "C"
	structInput.Appeared = 1972
	structInput.Created = append(structInput.Created,"Dennis Ritchie")
	structInput.Functional = true
	structInput.Objectorient = false
	structInput.Relation.InfluencedBy = influenStru.InfluencedBy
	structInput.Relation.Influences = influenStru.Influences

	return structInput
}

func Palindrome()dt.ResponsePalindrome{
	var resp dt.ResponsePalindrome

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Insert argument here : ")
	scanner.Scan() 
	arg := scanner.Text()
	result := isPalindrome(arg)

	resp.Status = 404
	resp.Desc = "Not palindrome"

	if result {
		resp.Status = 200
		resp.Desc = "Palindrome"
	}

	return resp
	
}

func isPalindrome(arg string)bool{
	argPalin := strings.ToLower(arg)
	arrStr := strings.Split(argPalin,"")
	countFail := true
	for i,rangeArg := range arrStr{
		if i >= len(arrStr)-1-i {
			break
		}
		pairStr := arrStr[len(arrStr)-1-i]
		if rangeArg != pairStr {
			countFail = false
			break
		}
	}
	if !countFail {
		return false
	}
	return true
}

func StoreData(w http.ResponseWriter, r *http.Request)dt.ResponseStore{
	var resp dt.ResponseStore

	w.Header().Set("Access-Control-Allow-Origin", "*")
	
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	c, errRead := ioutil.ReadAll(r.Body)
	
	defer r.Body.Close()

	if errRead!=nil{
		log.Errorf("Error", errRead)
		resp.Status = 404
		resp.Desc = "error when read parameter request"
		return resp
	}

	var msg []dt.DetailLanguage

	errUnmarshal := json.Unmarshal(c, &msg)

	if errUnmarshal != nil{
		log.Errorf("Error",errUnmarshal)
	}

	idx := 1

	lengthStored := len(ListStoredData)

	idx += lengthStored
	

	for _,rangeMsg := range msg{
		detailStored := new(dt.StoredData)
		detailStored.ID = idx
		detailStored.ListDetailLang.Language = rangeMsg.Language
		detailStored.ListDetailLang.Appeared = rangeMsg.Appeared
		detailStored.ListDetailLang.Created  = rangeMsg.Created
		detailStored.ListDetailLang.Functional  = rangeMsg.Functional
		detailStored.ListDetailLang.Objectorient = rangeMsg.Objectorient
		detailStored.ListDetailLang.Relation = rangeMsg.Relation
		ListStoredData = append(ListStoredData,detailStored)
		idx++
	}

	resp.Status = 200
	resp.Desc = "Sukses insert data!"
	return resp
}

func GetLang(w http.ResponseWriter, r *http.Request, paramID string)dt.DetailLanguage{
	var resp dt.DetailLanguage
	
	intID,errConv := strconv.Atoi(paramID)
	
	if errConv != nil {
		log.Errorf("Error",errConv)
		return resp
	}

	for _,rangeList := range ListStoredData{
		if intID == rangeList.ID{
			resp.Language = rangeList.ListDetailLang.Language
			resp.Appeared = rangeList.ListDetailLang.Appeared
			resp.Created = rangeList.ListDetailLang.Created
			resp.Functional = rangeList.ListDetailLang.Functional
			resp.Objectorient = rangeList.ListDetailLang.Objectorient
			resp.Relation = rangeList.ListDetailLang.Relation
			break
		}
	}

	return resp
}

func GetLanguages(w http.ResponseWriter, r *http.Request)[]dt.DetailLanguage{
	var resp []dt.DetailLanguage
	var singleLang dt.DetailLanguage

	for _,rangeList := range ListStoredData{
		singleLang.Language = rangeList.ListDetailLang.Language
		singleLang.Appeared = rangeList.ListDetailLang.Appeared
		singleLang.Created = rangeList.ListDetailLang.Created
		singleLang.Functional = rangeList.ListDetailLang.Functional
		singleLang.Objectorient = rangeList.ListDetailLang.Objectorient
		singleLang.Relation = rangeList.ListDetailLang.Relation
		resp = append(resp,singleLang)
	}

	return resp
}

func UpdateLang(w http.ResponseWriter, r *http.Request, paramID string)dt.ResponseStore{
	var resp dt.ResponseStore
	intID,errConv := strconv.Atoi(paramID)

	resp.Status = 404
	
	if errConv != nil {
		log.Errorf("Error",errConv)
		resp.Desc = "Bad request"
		return resp
	}

	fmt.Println(intID)

	w.Header().Set("Access-Control-Allow-Origin", "*")
	
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	c, errRead := ioutil.ReadAll(r.Body)
	
	defer r.Body.Close()

	if errRead!=nil{
		log.Errorf("Error", errRead)
		resp.Desc = "error when read parameter request"
		return resp
	}

	var msg dt.DetailLanguage

	errUnmarshal := json.Unmarshal(c, &msg)

	if errUnmarshal != nil{
		log.Errorf("Error",errUnmarshal)
	}

	isUpdate := false

	for _,rangeList := range ListStoredData{
		if intID == rangeList.ID {
			isUpdate = true
			rangeList.ListDetailLang.Language = msg.Language
			rangeList.ListDetailLang.Appeared = msg.Appeared
			rangeList.ListDetailLang.Created = msg.Created
			rangeList.ListDetailLang.Functional = msg.Functional
			rangeList.ListDetailLang.Objectorient = msg.Objectorient
			rangeList.ListDetailLang.Relation = msg.Relation
		}
	}

	resp.Desc = "Failed when update data, id not exist in data list..."

	if isUpdate{
		resp.Status = 200
		resp.Desc = "Sukses update data!"
	}
	return resp
}

func DeleteLang(w http.ResponseWriter, r *http.Request, paramID string)dt.ResponseStore{
	var resp dt.ResponseStore
	intID,errConv := strconv.Atoi(paramID)

	resp.Status = 404
	
	if errConv != nil {
		log.Errorf("Error",errConv)
		resp.Desc = "Bad request"
		return resp
	}

	var newList []*dt.StoredData

	isDelete := false

	for _,rangeList := range ListStoredData {
		if rangeList.ID == intID {
			isDelete = true
			continue
		}
		newList = append(newList,rangeList)
	}

	if !isDelete {
		resp.Desc = "Id not exist in data list!"
		return resp
	}

	var lastList []*dt.StoredData

	for i,rangeNew := range newList{
		singleData := new(dt.StoredData)
		singleData.ID = i+1
		singleData.ListDetailLang = rangeNew.ListDetailLang
		lastList = append(lastList,singleData)
	}

	ListStoredData = lastList

	for i,rangeList := range ListStoredData{
		fmt.Println(i,rangeList.ListDetailLang.Language)
	}
		resp.Status = 200
		resp.Desc = "Sukses delete data!"
	return resp
}

func SendResponses(rw http.ResponseWriter, dt interface{}) {
	js, err := json.Marshal(dt)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	rw.Header().Set("Content-Type", "application/json")
	rw.Write(js)
}