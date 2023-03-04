package db

import (
	"database/sql"
	"errors"
	"github.com/ipaqsa/netcom/cryptoUtils"
	"github.com/ipaqsa/netcom/logger"
	_ "github.com/mattn/go-sqlite3"
)

var infoLogger = logger.NewLogger("INFO")
var errorLogger = logger.NewLogger("ERROR")

func DataBaseInit(path string) *DB {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		errorLogger.Println(err.Error())
		return nil
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS storage (id INTEGER PRIMARY KEY AUTOINCREMENT,
																login VARCHAR(25) UNIQUE, 
    															hash VARCHAR(300), 
    															firstname VARCHAR(20), 
    															lastname VARCHAR(25), 
    															room INT);`)
	if err != nil {
		errorLogger.Println(err.Error())
		return nil
	}
	infoLogger.Println("database init is successful")
	return &DB{
		ptr: db,
	}
}

func (db *DB) GetPerson(login string, hash string) (*Person, error) {
	db.mtx.Lock()
	defer db.mtx.Unlock()
	var person Person
	var id int
	row := db.ptr.QueryRow(`SELECT * FROM storage WHERE login=$1 LIMIT 1`, login)
	err := row.Scan(&id, &person.Login, &person.Hash, &person.Firstname, &person.Lastname, &person.Room)
	if err != nil {
		errorLogger.Println(err.Error())
		return nil, errors.New("person`s not found")
	}
	hash = cryptoUtils.Base64Encode(cryptoUtils.HashSum(cryptoUtils.Base64Decode(hash)))
	if person.Hash != hash {
		return nil, errors.New("wrong password")
	}
	return &person, nil
}

func (db *DB) RegisterPerson(person *Person) error {
	p, err := db.GetPerson(person.Login, person.Hash)
	if err == nil || p != nil {
		return errors.New("login exists")
	}
	db.mtx.Lock()
	defer db.mtx.Unlock()
	_, err = db.ptr.Exec(`INSERT INTO storage (login, hash, firstname, lastname, room) VALUES ($1, $2, $3, $4, $5)`,
		person.Login, cryptoUtils.Base64Encode(cryptoUtils.HashSum(cryptoUtils.Base64Decode(person.Hash))), person.Firstname, person.Lastname, person.Room)
	if err != nil {
		return err
	}
	return nil
}
func (db *DB) UpdatePerson(person *Person, passwordChange string) error {
	p, err := db.GetPerson(person.Login, person.Hash)
	if err == nil || p != nil {
		if err.Error() != "wrong password" && p != nil {
			return errors.New("login doesnt exists")
		}
	}
	db.mtx.Lock()
	defer db.mtx.Unlock()
	if passwordChange == "yes" {
		_, err = db.ptr.Exec(`UPDATE storage SET hash=$1, firstname=$2, lastname=$3, room=$4 WHERE login = $5`,
			cryptoUtils.Base64Encode(cryptoUtils.HashSum(cryptoUtils.Base64Decode(person.Hash))), person.Firstname, person.Lastname, person.Room, person.Login)
	} else {
		_, err = db.ptr.Exec(`UPDATE storage SET firstname=$1, lastname=$2, room=$3 WHERE login = $4`, person.Firstname, person.Lastname, person.Room, person.Login)
	}
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) GetAllMember(room int) []User {
	db.mtx.Lock()
	defer db.mtx.Unlock()
	data := make([]User, 0)
	rows, err := db.ptr.Query(`SELECT login, firstname, lastname FROM storage WHERE room = $1`, room)
	if err != nil {
		errorLogger.Printf("GetAllMember: %s", err.Error())
		return nil
	}
	defer rows.Close()
	for rows.Next() {
		var login string
		var name string
		var secondname string
		if err = rows.Scan(&login, &name, &secondname); err != nil {
			errorLogger.Printf("GetAllMember: %s", err.Error())
			return data
		}
		user := User{
			Login:     login,
			FirstName: name,
			LastName:  secondname,
		}
		data = append(data, user)
	}
	return data
}

func (db *DB) RemovePerson(personId string) error {
	db.mtx.Lock()
	defer db.mtx.Unlock()
	_, err := db.ptr.Exec(`DELETE FROM storage WHERE login=$1`, personId)
	if err != nil {
		errorLogger.Printf("remove: %s", err.Error())
		return nil
	}
	return nil
}
