package db

import (
	"io/ioutil"
)

type JsonDB struct{}

func NewJsonDB() *JsonDB {
	return &JsonDB{}
}

func (db *JsonDB) Load(dbName DBName) ([]byte, error) {
	jsonData, err := ioutil.ReadFile("data/" + dbName + ".json")
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}

func (db *JsonDB) Save(dbName DBName, data Data) error {
	err := ioutil.WriteFile("data/"+dbName+".json", data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (db *JsonDB) Reset(fileDb string, defVal []byte) error {
	err := ioutil.WriteFile("data/"+fileDb+".json", defVal, 0644)
	if err != nil {
		return err
	}
	return nil
}
