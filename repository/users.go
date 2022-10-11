package repository

import (
	"a21hc3NpZ25tZW50/db"
	"a21hc3NpZ25tZW50/model"
	"encoding/json"
	"errors"
)

type UserRepository struct {
	db db.DB
}

func NewUserRepository(db db.DB) UserRepository {
	return UserRepository{db}
}

func (u *UserRepository) ReadUser() ([]model.Credentials, error) {
	records, err := u.db.Load("users")
	if err != nil {
		return nil, err
	}

	var listUser []model.Credentials
	err = json.Unmarshal([]byte(records), &listUser)
	if err != nil {
		return nil, err
	}

	return listUser, nil
}

func (u *UserRepository) AddUser(creds model.Credentials) error {
	credentialData, err := u.ReadUser()
	if err != nil {
		return err
	}

	for _, credItem := range credentialData {
		if credItem.Username == creds.Username {
			err = errors.New("username already taken")
			return err
		}
	}

	credentialData = append(credentialData, creds)

	jsonData, err := json.Marshal(credentialData)
	if err != nil {
		return err
	}
	u.db.Save("users", jsonData)

	return nil
}

func (u *UserRepository) ResetUser() error {
	err := u.db.Reset("users", []byte("[]"))
	if err != nil {
		return err
	}

	return nil
}

func (u *UserRepository) LoginValid(list []model.Credentials, req model.Credentials) bool {
	for _, element := range list {
		if element.Username == req.Username && element.Password == req.Password {
			return true
		}
	}
	return false
}
