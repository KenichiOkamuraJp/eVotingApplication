package blockchainutils

import (
	"encoding/json"
	"net/http"
	"io/ioutil"
	"fmt"
)

/////////////////////////////////////////////////////////////////////////////
//Request structs and decoders
/////////////////////////////////////////////////////////////////////////////

type RegistrarRequest struct {
	EnrollId		string			`json:"enrollId"`
	EnrollSecret	string			`json:"enrollSecret"`
}

func (rr * RegistrarRequest) EncodeJSON() []byte {

	repBytes, err := json.Marshal(rr)
	if err != nil {
		fmt.Println("JSON Marshal error:", err)
		return nil
	}
	return repBytes
}

type ChaincodeId struct {
	Path		string				`json:"path,omitempty"`
	Name		string				`json:"name,omitempty"`
}

type ChaincodeCtorMsg struct {
	Function	string				`json:"function,omitempty"`
	Args		[]string			`json:"args"`
}

type ChaincodeParams struct {
	ParamType	float64				`json:"type"`
	Id			ChaincodeId			`json:"chaincodeID"`
	CtorMsg		ChaincodeCtorMsg	`json:"ctorMsg`
	SecureContext string			`json:"secureContext,omitempty"`
}

type ChaincodeRequest struct {
	Jsonrpc		string				`json:"jsonrpc"`
	Method		string  			`json:"method"`
	Params		ChaincodeParams 	`json:"params"`
	Id			float64				`json:"id"`
}

func (cr *ChaincodeRequest) SetValues(aRpc string, aMethod string, aType float64,
 aPath string, aName string, aFunction string, anArgs []string, aSecure string,anId float64) {
	cr.Jsonrpc = aRpc
	cr.Method = aMethod
	aCCID := ChaincodeId{Path : aPath, Name : aName}
	aCtorMsg := ChaincodeCtorMsg{Function : aFunction, Args : anArgs}
	cr.Params = ChaincodeParams{ParamType : aType, Id : aCCID, CtorMsg : aCtorMsg, SecureContext : aSecure}
	cr.Id = anId
}

func (cr *ChaincodeRequest) EncodeJSON() []byte {
	
	repBytes, err := json.Marshal(cr)
	if err != nil {
		fmt.Println("JSON Marshal error:", err)
		return nil
	}
	
	return repBytes	
}

/////////////////////////////////////////////////////////////////////////////
//Response structs and decoders
/////////////////////////////////////////////////////////////////////////////

type RegistrarResponse struct {
	Ok				bool
	Message			string
}

func (rr *RegistrarResponse) DecodeJSON(resp http.Response) bool {
	
	byteArray, _ := ioutil.ReadAll(resp.Body)
	str := string(byteArray)
	
	var data map[string]interface{}
	err := json.Unmarshal([]byte(str), &data)
	if err != nil {
		panic(err)
	}
	
	_, flgSucess := data["OK"]
	if flgSucess {
		rr.Ok = true
		rr.Message = data["OK"].(string)
	} else {
		rr.Ok = false
		rr.Message = data["Error"].(string)
	}
	
	return rr.Ok
		
}

type ChaincodeResponse struct {
	Jsonrpc			string
	Ok				bool
	ErrorCode		float64
	ErrorMessage	string
	ErrorData		string
	Status			string
	Message			string
	Id				float64
}

func (cr *ChaincodeResponse) DecodeJSON(resp http.Response) bool {
	
	byteArray, _ := ioutil.ReadAll(resp.Body)
	str := string(byteArray)
	
	var data map[string]interface{}
	err := json.Unmarshal([]byte(str), &data)
	if err != nil {
		panic(err)
	}
	
	cr.Jsonrpc = data["jsonrpc"].(string)
	_, flgSucess := data["result"]
	if flgSucess {
		cr.Ok = true
		var subData map[string]interface{}
		subData = data["result"].(map[string]interface{})
		cr.Status = subData["status"].(string)
		cr.Message = subData["message"].(string)		
	} else {
		cr.Ok = false
		var subData map[string]interface{}
		subData = data["error"].(map[string]interface{})
		cr.ErrorCode = subData["code"].(float64)
		cr.ErrorMessage = subData["message"].(string)
		cr.ErrorData = subData["data"].(string)	
	}
	cr.Id = data["id"].(float64)
	
	return cr.Ok
}