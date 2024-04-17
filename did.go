package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing an Asset
type SmartContract struct {
	contractapi.Contract
}

//智能合约api接口

type Message struct {
	ID                  string `json:"ID"`
	Name                string `json:"姓名"`
	Gender              string `json:"性别"`
	Age                 int    `json:"年龄"`
	Illness_history     string `json:"疾病史"`
	Allergy_information string `json:"过敏信息"`
	Symptom             string `json:"症状"`
	Health_examination  string `json:"体检信息"`
	Diagnosis           string `json:"诊断"`
	MAR                 string `json:"用药记录"`
}

// 用户信息定义,ID作为查询时的键值。
// InitLedger adds a base set of assets to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	messages := []Message{
		{ID: "用户1", Name: "张三", Gender: "男", Age: 50, Illness_history: "高血压", Allergy_information: "青霉素", Symptom: "头晕,呕吐", Health_examination: "略", Diagnosis: "略", MAR: "略"},
		{ID: "用户2", Name: "李四", Gender: "男", Age: 60, Illness_history: "略", Allergy_information: "略", Symptom: "略", Health_examination: "略", Diagnosis: "略", MAR: "略"},
		{ID: "用户3", Name: "王五", Gender: "男", Age: 45, Illness_history: "略", Allergy_information: "略", Symptom: "略", Health_examination: "略", Diagnosis: "略", MAR: "略"},
		{ID: "用户4", Name: "赵六", Gender: "男", Age: 65, Illness_history: "略", Allergy_information: "略", Symptom: "略", Health_examination: "略", Diagnosis: "略", MAR: "略"},
	}

	for _, message := range messages {
		messageJSON, err := json.Marshal(message)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(message.ID, messageJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}

	return nil
}

// CreateAsset issues a new asset to the world state with given details.
func (s *SmartContract) Create(ctx contractapi.TransactionContextInterface, id string, name string, gender string, age int, illness_history string, allergy_information string, syptom string, health_examination string, diagnosis string, mar string) error {
	exists, err := s.Exists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the asset %s already exists", id)
	}

	message := Message{
		ID:                  id,
		Name:                name,
		Gender:              gender,
		Age:                 age,
		Illness_history:     illness_history,
		Allergy_information: allergy_information,
		Symptom:             syptom,
		Health_examination:  health_examination,
		Diagnosis:           diagnosis,
		MAR:                 mar,
	}
	messageJSON, err := json.Marshal(message)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, messageJSON)
}

// ReadAsset returns the asset stored in the world state with given id.
func (s *SmartContract) Read(ctx contractapi.TransactionContextInterface, id string) (*Message, error) {
	messageJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if messageJSON == nil {
		return nil, fmt.Errorf("the asset %s does not exist", id)
	}

	var message Message
	err = json.Unmarshal(messageJSON, &message)
	if err != nil {
		return nil, err
	}

	return &message, nil
}

// UpdateAsset updates an existing asset in the world state with provided parameters.
func (s *SmartContract) Update(ctx contractapi.TransactionContextInterface, id string, name string, gender string, age int, illness_history string, allergy_information string, syptom string, health_examination string, diagnosis string, mar string) error {
	exists, err := s.Exists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the asset %s does not exist", id)
	}

	// overwriting original asset with new asset
	message := Message{
		ID:                  id,
		Name:                name,
		Gender:              gender,
		Age:                 age,
		Illness_history:     illness_history,
		Allergy_information: allergy_information,
		Symptom:             syptom,
		Health_examination:  health_examination,
		Diagnosis:           diagnosis,
		MAR:                 mar,
	}
	messageJSON, err := json.Marshal(message)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, messageJSON)
}

// DeleteAsset deletes an given asset from the world state.
func (s *SmartContract) Delete(ctx contractapi.TransactionContextInterface, id string) error {
	exists, err := s.Exists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the asset %s does not exist", id)
	}

	return ctx.GetStub().DelState(id)
}

// AssetExists returns true when asset with given ID exists in world state
func (s *SmartContract) Exists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	messageJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return messageJSON != nil, nil
}

// GetAllAssets returns all assets found in world state
func (s *SmartContract) GetAllMessages(ctx contractapi.TransactionContextInterface) ([]*Message, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all assets in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var messages []*Message
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var message Message
		err = json.Unmarshal(queryResponse.Value, &message)
		if err != nil {
			return nil, err
		}
		messages = append(messages, &message)
	}

	return messages, nil
}

func main() {
	MyChaincode, err := contractapi.NewChaincode(&SmartContract{})
	if err != nil {
		log.Panicf("Error creating my chaincode: %v", err)
	}

	if err := MyChaincode.Start(); err != nil {
		log.Panicf("Error starting my chaincode: %v", err)
	}
}
