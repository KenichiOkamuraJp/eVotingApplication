package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"bytes"
	"./blockchainutils"
	"fmt"
	"strconv"
)

const remoteUrl = "https://689ed6712dc4436f901ffbddfa7a9d22-vp0.us.blockchain.ibm.com:5002"
const localUrl = "http://localhost:7050"
const sourcePath = "http://github.com/KenichiOkamuraJp/learn-chaincode/finished"
const localAppName = "mycc"
const dbUserID = "user_type1_1"
const dbUserPass = "b9fa1606ac"

const debug = true
const env = "local"

var url string

var chaincode string = "ccc9ccb8b10e8dc9faa54b7e0e58011d0c764f710900fe19b7dcd9047ec167bc266e184b4eb1567979c4a9f4823b5d50d6bec9a966775c649ad0ea9dee720ccc"

////////////////////////////////////////////////////////////////////////////////////////////////////
// eventHandler
////////////////////////////////////////////////////////////////////////////////////////////////////

func indexHandler (c *gin.Context)  {
	c.HTML(http.StatusOK, "admin01_login.tmpl", nil)	
}

func loginHandler (c *gin.Context)  {
	c.Request.ParseForm()
	id := c.Request.Form["id"]
	pass := c.Request.Form["pass"]	
	
	if id[0] == "jinji" && pass[0] == "jinji" {
		c.HTML(http.StatusOK, "admin02_deploy.tmpl", nil)		
	} else if id[0] == "admin" && pass[0] == "admin" {
		//Chaincode "Query"-"getAll"
		c.HTML(http.StatusOK, "admin03_master.tmpl", nil)		
	} else {
		c.HTML(http.StatusOK, "admin01_login.tmpl", gin.H{"errorMessage": "Invalid id or password",})	
	}
}

func logoutHandler (c *gin.Context)  {
	c.HTML(http.StatusOK, "admin04_logout.tmpl", nil)	
}

func deployHandler (c *gin.Context)  {
	
	var log []string
	
	//Login to blockchain
	if env == "remote" {
		log = append(log,"Trying to log in...")
		
		flgRegOk, regResp := postRegistrar()
		if flgRegOk {
			log = append(log,"Success.")
			log = append(log,"Message :" + regResp.Message)
		} else {
			log = append(log,"Failure.")
			log = append(log,"Message :" + regResp.Message)
			return			
		}
		log = append(log,"................")	
	}
	
	//Chaincode "deploy"
	log = append(log,"Trying to deploy...")
	flgDepOk, depResp := actChaincode("deploy",[]string{"init","arg"})
	if flgDepOk {
		log = append(log,"Success.")
		log = append(log,"Message :" + depResp.Message)
		chaincode = depResp.Message		
	} else {
		log = append(log,"Failure.")
		log = append(log,"Code :" + strconv.FormatFloat(depResp.ErrorCode, 'f', 4, 64))
		log = append(log,"Message :" + depResp.ErrorMessage)
		log = append(log,"Data :" + depResp.ErrorData)
	}
	
	c.HTML(http.StatusOK, "admin02_deploy.tmpl", gin.H{"response":log,})	
}

func maintenanceHandler (c *gin.Context)  {
	
	var flgOk	bool
	var resp	blockchainutils.ChaincodeResponse
	var message []string
	var screen	int
	
	c.Request.ParseForm()
	command := c.Request.Form["command"]
	
	switch command[0] {
		case "insert" :
		//1.Chaincode "invoke"-"add"-"ins_id"-"ins_pass"
		flgOk, resp = actChaincode("invoke",[]string{"add",c.Request.Form["ins_id"][0],c.Request.Form["ins_pass"][0]})
		screen = 3
		case "delete" :
		//2.Chaincode "invoke"-"delete"-"del_id"-"del_pass"
		flgOk, resp = actChaincode("invoke",[]string{"delete",c.Request.Form["del_id"][0],c.Request.Form["del_pass"][0]})
		screen = 3
		case "reload" :
		screen = 3
		case "initialize" :
		flgOk, resp = actChaincode("invoke",[]string{"initialize"})
		screen = 3	
		case "logout" :
		screen = 4
		default :
		screen = 5
	}
	
	if command[0] == "insert" || command[0] == "delete" || command[0] == "initialize" {
		if flgOk {
			message = append(message, "** Successfully operated **") 
		} else { 
			message = append(message, "** Error Occured in operation **")
			message = append(message, "Code : " + strconv.FormatFloat(resp.ErrorCode, 'f', 4, 64))
			message = append(message, "Message : " + resp.ErrorMessage)
			message = append(message, "Data : " + resp.ErrorData)
		}		
	}
	
	switch screen {
		case 3 :
		flgOk, resp = actChaincode("query",[]string{"getAll"})
		if flgOk {
			message = append(message, "** Successfully queried **") 
			// ここにgetAllの結果(resp.Message)を解析するルーチンを入れる
		} else { 
			message = append(message, "** Error Occured in querying **")
			message = append(message, "Code : " + strconv.FormatFloat(resp.ErrorCode, 'f', 4, 64))
			message = append(message, "Message : " + resp.ErrorMessage)
			message = append(message, "Data : " + resp.ErrorData)
		}			
		c.HTML(http.StatusOK, "admin03_master.tmpl", gin.H{"result": message,})
		case 4 :
		c.HTML(http.StatusOK, "admin04_logout.tmpl", nil)
		default :
		c.HTML(http.StatusOK, "admin05_error.tmpl", nil)
	}
	
}

////////////////////////////////////////////////////////////////////////////////////////////////////
// accessBlockchain
////////////////////////////////////////////////////////////////////////////////////////////////////

func postRegistrar () (bool, blockchainutils.RegistrarResponse){

	// Request

	regReq := blockchainutils.RegistrarRequest{EnrollId:dbUserID,EnrollSecret:dbUserPass}
	
	req ,err := http.NewRequest (
		"POST",
		remoteUrl + "/registrar",
		bytes.NewBuffer(regReq.EncodeJSON()),
	)
	if err != nil {
		return false, blockchainutils.RegistrarResponse{}
	}
	
	// Content-Type setting
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, blockchainutils.RegistrarResponse{}
	}
	defer resp.Body.Close()

	// Response
	regResp := new (blockchainutils.RegistrarResponse)
	regResp.DecodeJSON(*resp)
	
	// forDebug
	if !regResp.Ok {
		fmt.Println("Failure to login")
		fmt.Println(regResp.Message)		
	} else {
		if debug { fmt.Println(regResp.Message) }
	} 		
	
	return regResp.Ok, *regResp

} 

func actChaincode (action string, command []string) (bool, blockchainutils.ChaincodeResponse) {
	
	// Request
	actReq := new(blockchainutils.ChaincodeRequest)
	
	// Set parameter for remote access
	aPath := ""
	aChaincode := ""
	if action == "deploy" {
		aPath = sourcePath
	} else {
		aChaincode = chaincode
	}
	
	// Set parameter
	switch env {
		case "local" :
		actReq.SetValues("2.0", action, 1,"", localAppName ,"",command ,""       ,1)
		case "remote" :
		actReq.SetValues("2.0", action, 1,aPath, aChaincode    ,"",command , dbUserID,1)
		default :
		fmt.Println("Environment setting error!")
		return false, blockchainutils.ChaincodeResponse{}
	}

	req, err := http.NewRequest(
		"POST",
		url + "/chaincode",
		bytes.NewBuffer(actReq.EncodeJSON()),
	)
	if err != nil {
		return false, blockchainutils.ChaincodeResponse{}
	}

	// Content-Type setting
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, blockchainutils.ChaincodeResponse{}
	} 
	defer resp.Body.Close()

	// Response
	actResp := new (blockchainutils.ChaincodeResponse)
	actResp.DecodeJSON(*resp)
	
	// forDebug
	if !actResp.Ok {
		fmt.Println("Failure to " + action)
		fmt.Println(strconv.FormatFloat(actResp.ErrorCode, 'f', 4, 64))	
		fmt.Println(actResp.ErrorMessage)	
		fmt.Println(actResp.ErrorData)		
	} else {
		if debug { fmt.Println(actResp.Message) }
	} 		
	
	return actResp.Ok, *actResp
}

////////////////////////////////////////////////////////////////////////////////////////////////////
// Main
////////////////////////////////////////////////////////////////////////////////////////////////////

func main() {

	router := gin.Default()
	router.LoadHTMLGlob("html/*")
	
	switch env {
		case "local" :
		url = localUrl
		case "remote" :
		url = remoteUrl
		default :
		fmt.Println("Environment setting error!")
		return
	}
	
	router.GET("/index", indexHandler)
	router.POST("/login", loginHandler)
	router.POST("/deploy", deployHandler)
	router.POST("/logout", logoutHandler)
	router.POST("/maintenance", maintenanceHandler)

	port := os.Getenv("VCAP_APP_PORT")
	if port == "" {
		port = "8080"
	}
	router.Run(":" + port)
}
