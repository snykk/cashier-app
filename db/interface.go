package db

type DBName = string
type Data = []byte

type DB interface {
	Load(DBName) ([]byte, error)
	Save(DBName, Data) error
	Reset(DBName, Data) error
}
