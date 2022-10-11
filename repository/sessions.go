package repository

import (
	"a21hc3NpZ25tZW50/db"
	"a21hc3NpZ25tZW50/model"
	"encoding/json"
	"fmt"
	"time"
)

type SessionsRepository struct {
	db db.DB
}

func NewSessionsRepository(db db.DB) SessionsRepository {
	return SessionsRepository{db}
}

func (u *SessionsRepository) ReadSessions() ([]model.Session, error) {
	records, err := u.db.Load("sessions")
	if err != nil {
		return nil, err
	}

	var listSessions []model.Session
	err = json.Unmarshal([]byte(records), &listSessions)
	if err != nil {
		return nil, err
	}

	return listSessions, nil
}

func (u *SessionsRepository) DeleteSessions(tokenTarget string) error {
	listSessions, err := u.ReadSessions()
	if err != nil {
		return err
	}

	var filteredListSession = []model.Session{}

	// Select target token and delete from listSessions
	for _, chartItem := range listSessions {
		if chartItem.Token != tokenTarget {
			filteredListSession = append(filteredListSession, chartItem)
		}
	}

	listSessions = filteredListSession

	jsonData, err := json.Marshal(listSessions)
	if err != nil {
		return err
	}

	err = u.db.Save("sessions", jsonData)
	if err != nil {
		return err
	}

	return nil
}

func (u *SessionsRepository) AddSessions(session model.Session) error {
	sessionData, err := u.ReadSessions()
	if err != nil {
		return err
	}

	sessionData = append(sessionData, session)

	jsonData, err := json.Marshal(sessionData)
	if err != nil {
		return err
	}

	u.db.Save("sessions", jsonData)

	return nil // TODO: replace this
}

func (u *SessionsRepository) CheckExpireToken(token string) (model.Session, error) {
	var listSession []model.Session
	var sessionData model.Session

	jsonData, err := u.db.Load("sessions")
	if err != nil {
		return model.Session{}, err
	}

	err = json.Unmarshal(jsonData, &listSession)
	if err != nil {
		panic(err)
	}

	sessionData, err = u.TokenExist(listSession, token)

	if err != nil {
		return model.Session{}, err
	}

	if u.TokenExpired(sessionData) {
		u.DeleteSessions(sessionData.Token)
		return model.Session{}, fmt.Errorf("Token is Expired!")
	}

	return sessionData, nil
}

func (u *SessionsRepository) ResetSessions() error {
	err := u.db.Reset("sessions", []byte("[]"))
	if err != nil {
		return err
	}

	return nil
}

func (u *SessionsRepository) TokenExist(list []model.Session, req string) (model.Session, error) {
	for _, element := range list {
		if element.Token == req {
			return element, nil
		}
	}
	return model.Session{}, fmt.Errorf("Token Not Found!")
}

func (u *SessionsRepository) TokenExpired(s model.Session) bool {
	return s.Expiry.Before(time.Now())
}
