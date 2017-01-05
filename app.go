package main

import (	
	"os"
	"fmt"
	"bytes"
	"strconv"
	"net/http"

	"encoding/json"
	
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/contrib/sessions"
)

type Voter struct {
	Id			string
	Password	string
	Frequency	int64
	Destination	int64
}

var candidates = []string{
	"未投票","伊藤博文","黒田清隆","山縣有朋","松方正義","大隈重信","桂太郎","西園寺公望","山本権兵衛","寺内正毅","原敬","高橋是清","加藤友三郎",
	"清浦奎吾","加藤高明","若槻礼次郎","田中義一","浜口雄幸","犬養毅","斎藤実","岡田啓介","広田弘毅","林銑十郎","近衛文麿","平沼騏一郎","阿部信行",
	"米内光政","東条英機","小磯国昭","鈴木貫太郎","東久邇宮稔彦王","幣原喜重郎","吉田茂","片山哲","芦田均","鳩山一郎","石橋湛山","岸信介","池田勇人",
	"佐藤栄作","田中角栄","三木武夫","福田赳夫","大平正芳","鈴木善幸","中曽根康弘","竹下登","宇野宗佑","海部俊樹","宮澤喜一","細川護煕","羽田孜",
	"村山富市","橋本龍太郎","小渕恵三","森喜朗","小泉純一郎","安倍晋三","福田康夫","麻生太郎","鳩山由紀夫","菅直人","野田佳彦","白票",
}

const remoteUrl = "https://689ed6712dc4436f901ffbddfa7a9d22-vp0.us.blockchain.ibm.com:5002"
const localUrl = "http://localhost:7050"
const sourcePath = "http://github.com/KenichiOkamuraJp/eVotingApplication/server"
const localAppName = "mycc"
const dbUserID = "user_type1_1"
const dbUserPass = "b9fa1606ac"

const debug = false
const env = "remote"

var url string

var chaincode string = "99319c2c5a19b814a56455d5683ccac981c50363f1e52e0eb53ab1c8bb451f4821bf5b086bda1ddae1315211c11196cc0bce350a314838dcfc54312bd16249e8"

////////////////////////////////////////////////////////////////////////////////////////////////////
// eventHandler
////////////////////////////////////////////////////////////////////////////////////////////////////
// Administration
func adminIndexHandler (c *gin.Context)  {
	c.HTML(http.StatusOK, "admin01_login.tmpl", nil)	
}

func adminLoginHandler (c *gin.Context)  {
	c.Request.ParseForm()
	id := c.Request.Form["id"]
	pass := c.Request.Form["pass"]
	
	var message []string
	var voters []Voter
	
	if id[0] == "jinji" && pass[0] == "jinji" {
		c.HTML(http.StatusOK, "admin02_deploy.tmpl", nil)		
	} else if (id[0] == "debug" && pass[0] == "debug") || (id[0] == "admin" && pass[0] == "admin") {
		// Get_All
		flgOk, resp := actChaincode("query",[]string{"getAll"})
		if flgOk {
			// If no voters are registerd
			if resp.Message == "null" {
				message = append(message, "No voters are registered.")	
			} else {
				// Decode votersList				
				voters = make([]Voter,0,0)
				err := json.Unmarshal([]byte(resp.Message), &voters)
				if err != nil {
					fmt.Println("JSON Unmarshal error:", err)
					return					
				}
			}
			
		} else { 
			message = append(message, "** Error Occured in querying **")
			message = append(message, "Code : " + strconv.FormatFloat(resp.ErrorCode, 'f', 4, 64))
			message = append(message, "Message : " + resp.ErrorMessage)
			message = append(message, "Data : " + resp.ErrorData)
		}	
		if id[0] == "debug" {
			c.HTML(http.StatusOK, "admin03_master.tmpl", gin.H{"result": message, "votersList" : voters})		
		} else {
			c.HTML(http.StatusOK, "admin0301_master_limited.tmpl", gin.H{"result": message, "votersList" : voters})	
		}
			
	} else {
		c.HTML(http.StatusOK, "admin01_login.tmpl", gin.H{"errorMessage": "Invalid id or password",})	
	}
}

func adminLogoutHandler (c *gin.Context)  {
	clearSession(c)
	c.HTML(http.StatusOK, "admin04_logout.tmpl", nil)	
}

func adminDeployHandler (c *gin.Context)  {
	
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

func adminSetupHandler (c *gin.Context)  {
	
	var log []string
		
	//Chaincode "deploy"
	log = append(log,"Trying to setup...")
	flgDepOk, setupResp := actChaincode("invoke",[]string{"setup"})
	if flgDepOk {
		log = append(log,"Success.")
		log = append(log,"Message :" + setupResp.Message)
		chaincode = setupResp.Message		
	} else {
		log = append(log,"Failure.")
		log = append(log,"Code :" + strconv.FormatFloat(setupResp.ErrorCode, 'f', 4, 64))
		log = append(log,"Message :" + setupResp.ErrorMessage)
		log = append(log,"Data :" + setupResp.ErrorData)
	}
	c.HTML(http.StatusOK, "admin02_deploy.tmpl", gin.H{"response":log,})	
}

func adminMaintenanceHandler (c *gin.Context)  {
	
	var flgOk	bool
	var resp	ChaincodeResponse
	var message []string
	var voters []Voter = make([]Voter,0,0)
	var screen	int
	
	c.Request.ParseForm()
	command := c.Request.Form["command"]
	
	switch command[0] {
		case "insert" :
		flgOk, resp = actChaincode("query",[]string{"existCheck",c.Request.Form["ins_id"][0]})
		if flgOk {
			if resp.Message == "Not Exist" {
				//1.Chaincode "invoke"-"add"-"ins_id"-"ins_pass"
				flgOk, resp = actChaincode("invoke",[]string{"add",c.Request.Form["ins_id"][0],c.Request.Form["ins_pass"][0]})				
			} else {
				message = append(message,c.Request.Form["ins_id"][0] + " already exists")
				flgOk = false
			}
		}
		screen = 3
		case "delete" :
		flgOk, resp = actChaincode("query",[]string{"authentication",c.Request.Form["del_id"][0],c.Request.Form["del_pass"][0]})
		if flgOk {
			if resp.Message == "Authenticated" {
				//2.Chaincode "invoke"-"delete"-"del_id"-"del_pass"
				flgOk, resp = actChaincode("invoke",[]string{"delete",c.Request.Form["del_id"][0],c.Request.Form["del_pass"][0]})				
			} else {
				message = append(message, "Invalid id or password")
				flgOk = false
			}
		}		
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

			// If no voters are registerd
			if resp.Message == "null" {
				message = append(message, "But, no voters are registered.")	
			} else {
				// Decode votersList				
				voters = make([]Voter,0,0)
				err := json.Unmarshal([]byte(resp.Message), &voters)
				if err != nil {
					fmt.Println("JSON Unmarshal error:", err)
					return					
				}
			}
			
		} else { 
			message = append(message, "** Error Occured in querying **")
			message = append(message, "Code : " + strconv.FormatFloat(resp.ErrorCode, 'f', 4, 64))
			message = append(message, "Message : " + resp.ErrorMessage)
			message = append(message, "Data : " + resp.ErrorData)
		}			
		c.HTML(http.StatusOK, "admin03_master.tmpl", gin.H{"result": message, "votersList" : voters})
		case 4 :
		c.HTML(http.StatusOK, "admin04_logout.tmpl", nil)
		default :
		c.HTML(http.StatusOK, "admin05_error.tmpl", nil)
	}
	
}

func adminMaintenanceLimitedHandler (c *gin.Context)  {
	
	var flgOk	bool
	var resp	ChaincodeResponse
	var message []string
	var voters []Voter = make([]Voter,0,0)
	var screen	int
	
	c.Request.ParseForm()
	command := c.Request.Form["command"]
	
	switch command[0] {
		case "insert" :
		flgOk, resp = actChaincode("query",[]string{"existCheck",c.Request.Form["ins_id"][0]})
		if flgOk {
			if resp.Message == "Not Exist" {
				//1.Chaincode "invoke"-"add"-"ins_id"-"ins_pass"
				flgOk, resp = actChaincode("invoke",[]string{"add",c.Request.Form["ins_id"][0],c.Request.Form["ins_pass"][0]})				
			} else {
				message = append(message,c.Request.Form["ins_id"][0] + " already exists")
				flgOk = false
			}
		}
		screen = 3
		case "delete" :
		flgOk, resp = actChaincode("query",[]string{"authentication",c.Request.Form["del_id"][0],c.Request.Form["del_pass"][0]})
		if flgOk {
			if resp.Message == "Authenticated" {
				//2.Chaincode "invoke"-"delete"-"del_id"-"del_pass"
				flgOk, resp = actChaincode("invoke",[]string{"delete",c.Request.Form["del_id"][0],c.Request.Form["del_pass"][0]})				
			} else {
				message = append(message, "Invalid id or password")
				flgOk = false
			}
		}		
		screen = 3
		case "reload" :
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

			// If no voters are registerd
			if resp.Message == "null" {
				message = append(message, "But, no voters are registered.")	
			} else {
				// Decode votersList				
				voters = make([]Voter,0,0)
				err := json.Unmarshal([]byte(resp.Message), &voters)
				if err != nil {
					fmt.Println("JSON Unmarshal error:", err)
					return					
				}
			}
			
		} else { 
			message = append(message, "** Error Occured in querying **")
			message = append(message, "Code : " + strconv.FormatFloat(resp.ErrorCode, 'f', 4, 64))
			message = append(message, "Message : " + resp.ErrorMessage)
			message = append(message, "Data : " + resp.ErrorData)
		}			
		c.HTML(http.StatusOK, "admin0301_master_limited.tmpl", gin.H{"result": message, "votersList" : voters})
		case 4 :
		c.HTML(http.StatusOK, "admin04_logout.tmpl", nil)
		default :
		c.HTML(http.StatusOK, "admin05_error.tmpl", nil)
	}
	
}

// Voters
func indexHandler (c *gin.Context)  {
	c.HTML(http.StatusOK, "vote01_login.tmpl", nil)	
}

func loginHandler (c *gin.Context)  {
	c.Request.ParseForm()
	id := c.Request.Form["id"]
	pass := c.Request.Form["pass"]
	
	flgOk, resp := actChaincode("query",[]string{"authentication",id[0],pass[0]})
	if flgOk {
		if resp.Message == "Authenticated" {
			createSession(c,id[0])
			info := getSessionInfo(c)
			if info.IsSessionAlive {
				c.HTML(http.StatusOK, "vote02_vote.tmpl", gin.H{"user" : info})			
			} else {
				clearSession(c)
				c.HTML(http.StatusOK, "vote01_login.tmpl", gin.H{"errorMessage": "Session invalid",})	
			}
		} else {
			c.HTML(http.StatusOK, "vote01_login.tmpl", gin.H{"errorMessage": "Invalid id or password",})	

		}
	}
}

func voteHandler (c *gin.Context) {
	c.Request.ParseForm()	
	destinationStr := c.Request.Form["dest"][0]
	destination, _ := strconv.ParseInt(destinationStr,10,64)
	command := c.Request.Form["command"][0]
	
	// command check
	if command == "logout" {
		clearSession(c)
		c.HTML(http.StatusOK, "admin04_logout.tmpl", nil)
	} else if command == "vote" {
	
		// If not selected return error
		if destination == 255 {
			info := getSessionInfo(c)
			c.HTML(http.StatusOK, "vote02_vote.tmpl", gin.H{"user" : info})
			return	
		}
	
		// session check
		info := getSessionInfo(c)
		userid := info.UserID
		if !info.IsSessionAlive {
			clearSession(c)
			c.HTML(http.StatusOK, "vote05_error.tmpl", gin.H{"errorMessage": "Session invalid",})	
			return	
		}
	
		// get old value
		voter := Voter{}
		flgOk, resp := actChaincode("query",[]string{"getOne", userid})
		if flgOk {
			// If no voters are registerd				
			err := json.Unmarshal([]byte(resp.Message), &voter)
			if err != nil {
				fmt.Println("JSON Unmarshal error:", err)
				return
			}
		} else {
			clearSession(c)
			c.HTML(http.StatusOK, "vote05_error.tmpl", gin.H{"errorMessage": "Error in querying data",})	
			return	
		}
		newFrequency := voter.Frequency + 1
		setSession(c, userid, voter.Password, newFrequency, destination)
		
		c.HTML(http.StatusOK, "vote03_confirm.tmpl", gin.H{"user": info, "destination" : candidates[destination]})	
		
	}
}

func confirmHandler (c * gin.Context) {
	c.Request.ParseForm()	
	command := c.Request.Form["command"][0]
	
	//session check
	info := getSessionInfo(c)
		if !info.IsSessionAlive {
		clearSession(c)
		c.HTML(http.StatusOK, "vote05_error.tmpl", gin.H{"errorMessage": "Session invalid",})	
		return	
	}
	
	if command == "revote" {
		c.HTML(http.StatusOK, "vote02_vote.tmpl", gin.H{"user" : info})
	} else {
		// set new value
		id := info.UserID
		pass := info.Password
		freq := strconv.FormatInt(info.Frequency,10)
		dest := strconv.FormatInt(info.Destination,10)
		flgOk, _ := actChaincode("invoke",[]string{"update", id, pass, freq, dest})
		if flgOk {
			clearSession(c)
			c.HTML(http.StatusOK, "vote04_logout.tmpl", gin.H{"Message": "正常に投票が完了しました",})	
			return				
		} else {
			clearSession(c)
			c.HTML(http.StatusOK, "vote05_error.tmpl", gin.H{"errorMessage": "Error in querying data",})	
			return	
		}
	}
	
}

type CountStruct struct {
	C0 int; C1 int; C2 int; C3 int; C4 int; C5 int; C6 int; C7 int; C8 int; C9 int
	C10 int; C11 int; C12 int; C13 int; C14 int; C15 int; C16 int; C17 int; C18 int; C19 int
	C20 int; C21 int; C22 int; C23 int; C24 int; C25 int; C26 int; C27 int; C28 int; C29 int	
	C30 int; C31 int; C32 int; C33 int; C34 int; C35 int; C36 int; C37 int; C38 int; C39 int
	C40 int; C41 int; C42 int; C43 int; C44 int; C45 int; C46 int; C47 int; C48 int; C49 int
	C50 int; C51 int; C52 int; C53 int; C54 int; C55 int; C56 int; C57 int; C58 int; C59 int
	C60 int; C61 int; C62 int; C63 int
}

func viewHandler (c * gin.Context) {
	
	var count [64]int
	flgOk, resp := actChaincode("query",[]string{"count"})
	if flgOk {
		// If no voters are registerd				
		err := json.Unmarshal([]byte(resp.Message), &count)
		if err != nil {
			fmt.Println("JSON Unmarshal error:", err)
			return
		}
	} else {
		c.HTML(http.StatusOK, "vote05_error.tmpl", gin.H{"errorMessage": "Error in querying data",})	
		return	
	}	
	cstr := CountStruct{
		count[0],count[1],count[2],count[3],count[4],count[5],count[6],count[7],count[8],count[9],
		count[10],count[11],count[12],count[13],count[14],count[15],count[16],count[17],count[18],count[19],
		count[20],count[21],count[22],count[23],count[24],count[25],count[26],count[27],count[28],count[29],
		count[30],count[31],count[32],count[33],count[34],count[35],count[36],count[37],count[38],count[39],
		count[40],count[41],count[42],count[43],count[44],count[45],count[46],count[47],count[48],count[49],
		count[50],count[51],count[52],count[53],count[54],count[55],count[56],count[57],count[58],count[59],
		count[60],count[61],count[62],count[63],
	}
	
	var sum_vote, sum int
	//var vote_rate float64
	for i := 0 ; i < 62 ; i ++ {
		sum_vote = sum_vote + count[i]
	}
	sum = sum_vote + count[63]
	vote_rate := fmt.Sprintf("%2.1f",float64(sum_vote) / float64(sum) * 100)
	
	fmt.Println(vote_rate)
	
	c.HTML(http.StatusOK, "view01_result.tmpl",gin.H{
		"sum" : sum, "sumVote" : sum_vote, "rate" : vote_rate , "count" : cstr,
	})
	
}

////////////////////////////////////////////////////////////////////////////////////////////////////
// session management
////////////////////////////////////////////////////////////////////////////////////////////////////
type SessionInfo struct {
	UserID         	string
	UserName       	string
	Password		string
	Frequency		int64
	Destination		int64
	IsSessionAlive 	bool
}

func createSession(c *gin.Context, userID string) {
	setSession(c, userID, "", 0, 0)
}

func setSession(c *gin.Context, userID string,password string, frequency int64,destination int64) {
	session := sessions.Default(c)
	session.Set("userID", userID)
	session.Set("uName", userID)
	session.Set("pass", password)
	session.Set("freq", frequency)
	session.Set("dest", destination)
	session.Set("alive", true)
	session.Save()	
}

func clearSession(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
}

func getSessionInfo(c *gin.Context) SessionInfo {
	var info SessionInfo
	session := sessions.Default(c)
	userid := session.Get("userID")
	uname := session.Get("uName")
	password := session.Get("pass")
	freq := session.Get("freq")
	dest := session.Get("dest")
	alive := session.Get("alive")

	info = SessionInfo{
		UserID:         userid.(string),
		UserName:       uname.(string),
		Password:		password.(string),
		Frequency:		freq.(int64),
		Destination:	dest.(int64),
		IsSessionAlive: alive.(bool),
	}

	return info
}

////////////////////////////////////////////////////////////////////////////////////////////////////
// accessBlockchain
////////////////////////////////////////////////////////////////////////////////////////////////////

func postRegistrar () (bool, RegistrarResponse){

	// Request

	regReq := RegistrarRequest{EnrollId:dbUserID,EnrollSecret:dbUserPass}
	
	req ,err := http.NewRequest (
		"POST",
		remoteUrl + "/registrar",
		bytes.NewBuffer(regReq.EncodeJSON()),
	)
	if err != nil {
		return false, RegistrarResponse{}
	}
	
	// Content-Type setting
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, RegistrarResponse{}
	}
	defer resp.Body.Close()

	// Response
	regResp := new (RegistrarResponse)
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

func actChaincode (action string, command []string) (bool, ChaincodeResponse) {
	
	// Request
	actReq := new(ChaincodeRequest)
	
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
		return false, ChaincodeResponse{}
	}

	req, err := http.NewRequest(
		"POST",
		url + "/chaincode",
		bytes.NewBuffer(actReq.EncodeJSON()),
	)
	if err != nil {
		return false, ChaincodeResponse{}
	}

	// Content-Type setting
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, ChaincodeResponse{}
	} 
	defer resp.Body.Close()

	// Response
	actResp := new (ChaincodeResponse)
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
	
	store := sessions.NewCookieStore([]byte("secret"))
	router.Use(sessions.Sessions("mysession", store))
	
	switch env {
		case "local" :
		url = localUrl
		case "remote" :
		url = remoteUrl
		default :
		fmt.Println("Environment setting error!")
		return
	}
	
	router.GET("/admin_index", 			adminIndexHandler)
	router.POST("/admin_login", 		adminLoginHandler)
	router.POST("/admin_deploy", 		adminDeployHandler)
	router.POST("/admin_setup", 		adminSetupHandler)
	router.POST("/admin_logout", 		adminLogoutHandler)
	router.POST("/admin_maintenance",	adminMaintenanceHandler)
	router.POST("/admin_maintenance_limited",	adminMaintenanceLimitedHandler)
	
	router.GET("/index", 	indexHandler)
	router.POST("/login", 	loginHandler)
	router.POST("/vote", 	voteHandler)
	router.POST("/confirm",	confirmHandler)
	router.POST("/view", 	viewHandler)

	port := os.Getenv("VCAP_APP_PORT")
	if port == "" {
		port = "8080"
	}
	router.Run(":" + port)
}
