package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

type Voter struct {
	Id			string
	Password	string
	Frequency	int64
	Destination	int64
}

const dumKey int32 = 129


// Init resets all the things
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	return nil, nil
}

// Invoke isur entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	var err error
	var message []byte

	// Handle different functions
	switch function {
		case "setup" :
			//Create table
			err = createVotersList(stub)
			if err != nil {
				return nil, fmt.Errorf("Error in creating table. %s", err)
			}
		
			//Initialize table data
			err = prepareVotersList(stub)
			if err != nil {
				return nil, fmt.Errorf("Error in initializing data. %s", err)
			}	
					
		case "add" :
			err = insertRow(stub ,args[0],args[1]) ; 
			if err != nil { 
				return nil, fmt.Errorf("Error in inserting record.. %s", err)
			}
			
		case "delete" :
			err = deleteRow(stub ,args[0]) ; 
			if err != nil { 
				return nil, fmt.Errorf("Error in deleting record.. %s", err)
			}
		
		case "update" :
			err = updateRow(stub ,args[0],args[1],0,0) ; 
			if err != nil { 
				return nil, fmt.Errorf("Error in updating record.. %s", err)
			}
		
		case "initialize" :
			err = stub.DeleteTable("votersList")
			if err != nil {
				return nil, fmt.Errorf("Error in dropping table.. %s", err)
			}

			//Create table
			err = createVotersList(stub)
			if err != nil {
				return nil, fmt.Errorf("Error in creating table. %s", err)
			}
			
			//Initialize table data
			err = prepareVotersList(stub)
			if err != nil {
				return nil, fmt.Errorf("Error in initializing data. %s", err)
			}	
		default :
			return nil, errors.New("Received unknown function invocation: " + function)
	}
	
	return message, nil
}

// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	var message []byte
	
	// Handle different functions
	switch function {
		case "getAll" :

			voters ,err := getAllRows(stub)
			if err != nil { 
				return nil, fmt.Errorf("Error in getting all records.. %s", err)
			}
	
			jsonVoters, jsonErr := json.Marshal(voters)
			if jsonErr != nil {
				return nil, fmt.Errorf("Error in marshaling JSON: %s", err)
			}
			message = jsonVoters
			
		case "authentication" :
		
			exist, voter , err := getRow(stub, args[0])
			if err != nil { 
				return nil, fmt.Errorf("Error in getting all records.. %s", err)
			}
			
			if !exist {
				message = []byte("Not Exist")
			} else {
				if (voter.Password == args[1]) {
					message = []byte("Authenticated")
				} else {
					message = []byte("Not Authenticated")
				}				
			}
			
		case "existCheck" :
		
			exist, _ , err := getRow(stub, args[0])
			if err != nil { 
				return nil, fmt.Errorf("Error in getting all records.. %s", err)
			}
		
			if !exist {
				message = []byte("Not Exist")
			} else {
				message = []byte("Exist")
			} 		
			
		default :
			return nil, errors.New("Received unknown function invocation: " + function)
	}
	
		return message, nil
}

// ============================================================================================================================
// Main
// ============================================================================================================================
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

func createVotersList(stub shim.ChaincodeStubInterface) error {
	// Create voters list
	var columnDefsVotersList []*shim.ColumnDefinition
	
	columnDefDummy := shim.ColumnDefinition{Name: "dummy",
		Type: shim.ColumnDefinition_INT32, Key: true}	
	columnDefId := shim.ColumnDefinition{Name: "id",
		Type: shim.ColumnDefinition_STRING, Key: true}
	columnDefPassword := shim.ColumnDefinition{Name: "password",
		Type: shim.ColumnDefinition_STRING, Key: false}
	columnDefFrequency := shim.ColumnDefinition{Name: "frequency",
		Type: shim.ColumnDefinition_INT64, Key: false}
	columnDefDestination := shim.ColumnDefinition{Name: "destination",
		Type: shim.ColumnDefinition_INT64, Key: false}
		
	columnDefsVotersList = append(columnDefsVotersList, &columnDefDummy)		
	columnDefsVotersList = append(columnDefsVotersList, &columnDefId)
	columnDefsVotersList = append(columnDefsVotersList, &columnDefPassword)
	columnDefsVotersList = append(columnDefsVotersList, &columnDefFrequency)
	columnDefsVotersList = append(columnDefsVotersList, &columnDefDestination)
	
	return stub.CreateTable("votersList", columnDefsVotersList)
}

func prepareVotersList (stub shim.ChaincodeStubInterface) error {
	
	var err error
	
	err = insertRow(stub ,"Makoto_Uchida" 		,"a76z10") ; if err != nil { return err }
	err = insertRow(stub ,"Kenichi_Okamura" 	,"b85y29") ; if err != nil { return err }
	err = insertRow(stub ,"Takuro_Asanuma" 		,"c94x38") ; if err != nil { return err }
	err = insertRow(stub ,"Daichi_Kanazawa" 	,"d03w47") ; if err != nil { return err }
	err = insertRow(stub ,"Toshiki_Nobutsune" 	,"e12v56") ; if err != nil { return err }
	err = insertRow(stub ,"Shota_Saito" 		,"f21u65") ; if err != nil { return err }
	err = insertRow(stub ,"Toshiyuki_Kuramochi","g30t74") ; if err != nil { return err }
	err = insertRow(stub ,"Teruaki_Matsunami" 	,"h49s83") ; if err != nil { return err }
	err = insertRow(stub ,"Takahito_Koyama" 	,"i58r92") ; if err != nil { return err }
	err = insertRow(stub ,"Makoto_Oku" 			,"j67q01") ; if err != nil { return err }

	return nil
}

func insertRow (stub shim.ChaincodeStubInterface, userId string, password string) error {
	// Insert a voter
	var columns []*shim.Column
	
	col0 := shim.Column{Value: &shim.Column_Int32{Int32: dumKey}}
	col1 := shim.Column{Value: &shim.Column_String_{String_: userId}}
	col2 := shim.Column{Value: &shim.Column_String_{String_: password}}
	col3 := shim.Column{Value: &shim.Column_Int64{Int64: 0}}
	col4 := shim.Column{Value: &shim.Column_Int64{Int64: 0}}
	columns = append(columns, &col0)	
	columns = append(columns, &col1)
	columns = append(columns, &col2)
	columns = append(columns, &col3)
	columns = append(columns, &col4)

	row := shim.Row{Columns: columns}
	
	ok, err := stub.InsertRow("votersList", row)
	if err != nil {
		return fmt.Errorf("insertRow operation failed. %s", err)
	}
	if !ok {
		return fmt.Errorf("insertRow operation failed. Row with given id already exists")
	}
	
	return nil
}

func deleteRow (stub shim.ChaincodeStubInterface, userId string) error {
	// Insert a voter
	var columns []shim.Column
	
	col0 := shim.Column{Value: &shim.Column_Int32{Int32: dumKey}}	
	col1 := shim.Column{Value: &shim.Column_String_{String_: userId}}
	columns = append(columns, col0)
	columns = append(columns, col1)
	
	err := stub.DeleteRow("votersList", columns)
	if err != nil {
		return fmt.Errorf("deleteRow operation failed. %s", err)
	}
		
	return nil
}

func updateRow (stub shim.ChaincodeStubInterface, userId string, password string, freq int64, dest int64) error {
	// Insert a voter
	var columns []*shim.Column
	
	col0 := shim.Column{Value: &shim.Column_Int32{Int32: dumKey}}
	col1 := shim.Column{Value: &shim.Column_String_{String_: userId}}
	col2 := shim.Column{Value: &shim.Column_String_{String_: password}}
	col3 := shim.Column{Value: &shim.Column_Int64{Int64: freq}}
	col4 := shim.Column{Value: &shim.Column_Int64{Int64: dest}}
	columns = append(columns, &col0)
	columns = append(columns, &col1)
	columns = append(columns, &col2)
	columns = append(columns, &col3)
	columns = append(columns, &col4)

	row := shim.Row{Columns: columns}
	
	ok, err := stub.ReplaceRow("votersList", row)
	if err != nil {
		return fmt.Errorf("replaceRow operation failed. %s", err)
	}
	if !ok {
		return errors.New("replaceRow operation failed. Row with given id already exists")
	}
	
	return nil
}

func getRow (stub shim.ChaincodeStubInterface, userId string) (bool,Voter,error) {

	var columns []shim.Column
	col0 := shim.Column{Value: &shim.Column_Int32{Int32: dumKey}}
	colID := shim.Column{Value: &shim.Column_String_{String_: userId}}
	columns = append(columns, col0)
	columns = append(columns, colID)

	row , err := stub.GetRow("votersList", columns)
	if err != nil {
		return false, Voter{}, fmt.Errorf("getRow operation failed. %s", err)
	}
	
	if len(row.Columns) == 0 {
		return false, Voter{}, nil
	}

	repVoter := Voter {
		Id			: row.Columns[1].GetString_(),
		Password	: row.Columns[2].GetString_(),
		Frequency	: row.Columns[3].GetInt64(),
		Destination	: row.Columns[4].GetInt64(),
	}

	return true, repVoter, nil	
}

func getAllRows (stub shim.ChaincodeStubInterface) ([]Voter,error) {
	
	var repVoters []Voter
	var columns []shim.Column

	col0 := shim.Column{Value: &shim.Column_Int32{Int32: dumKey}}
	columns = append(columns, col0)

	rowChannel, err := stub.GetRows("votersList", columns)
	if err != nil {
		return nil, fmt.Errorf("getAllRow operation failed. %s", err)
	}
	
	for {		
		select {
		case row, ok := <-rowChannel:
			if !ok {
				rowChannel = nil
			} else {
				repVoters = append(repVoters, Voter{
					Id			: row.Columns[1].GetString_(),
					Password	: row.Columns[2].GetString_(),
					Frequency	: row.Columns[3].GetInt64(),
					Destination	: row.Columns[4].GetInt64(), 	
				})				
			}
		}
		if rowChannel == nil {
			break
		}
	}

	return repVoters, nil
	
}

