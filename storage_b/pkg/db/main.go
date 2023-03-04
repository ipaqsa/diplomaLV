package db

import (
	"database/sql"
	"github.com/ipaqsa/netcom/logger"
	_ "github.com/mattn/go-sqlite3"
)

var infoLogger = logger.NewLogger("INFO")
var errorLogger = logger.NewLogger("ERROR")

func DataBaseInit(path string) *DB {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil
	}
	return &DB{
		ptr: db,
	}
}

func (db *DB) AddMessage(sender, receiver, data, date, mtype string, t int) error {
	err := db.CreateDialogReceiverTable(receiver)
	if err != nil {
		return err
	}
	err = db.CreateDialogSenderTable(sender)
	if err != nil {
		return err
	}
	if t == 0 {
		err = db.addMessageToReceiver(sender, receiver, data, date, mtype)
		if err != nil {
			errorLogger.Printf("%s", err.Error())
			return err
		}
	} else {
		err = db.addMessageToSender(sender, receiver, data, date, mtype)
		if err != nil {
			errorLogger.Printf("%s", err.Error())
			return err
		}
	}
	return nil
}
func (db *DB) CreateDialogSenderTable(sender string) error {
	db.mtx.Lock()
	defer db.mtx.Unlock()
	query := "CREATE TABLE IF NOT EXISTS `" + sender + "_send` (id INTEGER PRIMARY KEY AUTOINCREMENT, date VARCHAR(25), receiver VARCHAR(30), data TEXT, mtype VARCHAR(10), check_s INT);"
	_, err := db.ptr.Exec(query)
	if err != nil {
		errorLogger.Printf("%s", err.Error())
		return err
	}
	return nil
}
func (db *DB) CreateDialogReceiverTable(receiver string) error {
	db.mtx.Lock()
	defer db.mtx.Unlock()
	query := "CREATE TABLE IF NOT EXISTS `" + receiver + "_rec` (id INTEGER PRIMARY KEY AUTOINCREMENT, date VARCHAR(25), sender VARCHAR(30), data TEXT, mtype VARCHAR(10), check_s INT);"
	_, err := db.ptr.Exec(query)
	if err != nil {
		errorLogger.Printf("%s", err.Error())
		return err
	}
	return nil
}
func (db *DB) addMessageToSender(sender, receiver, data, date, mtype string) error {
	db.mtx.Lock()
	defer db.mtx.Unlock()

	query := "INSERT INTO" + " `" + sender + "_send` (date, receiver, data, mtype, check_s) VALUES ($1, $2, $3, $4, $5)"

	_, err := db.ptr.Exec(query, date, receiver, data, mtype, 0)
	if err != nil {
		return err
	}
	return nil
}
func (db *DB) addMessageToReceiver(sender, receiver, data, date, mtype string) error {
	db.mtx.Lock()
	defer db.mtx.Unlock()

	query := "INSERT INTO" + " `" + receiver + "_rec` (date, sender, data, mtype, check_s) VALUES ($1, $2, $3, $4, $5)"
	_, err := db.ptr.Exec(query, date, sender, data, mtype, 0)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) GetMessages(sender, receiver string) *Messages {
	db.mtx.Lock()
	defer db.mtx.Unlock()

	var messages Messages
	queryS := "SELECT date, data, mtype, check_s FROM" + " `" + sender + "_send`  WHERE receiver=$1"
	queryR := "SELECT date, data, mtype, check_s FROM" + " `" + sender + "_rec`  WHERE sender=$1"
	rowsS, err := db.ptr.Query(queryS, receiver)
	if err != nil {
		errorLogger.Printf("%s", err.Error())
	} else {
		for rowsS.Next() {
			var message = Message{}
			if err = rowsS.Scan(&message.Date, &message.Data, &message.Type, &message.Check); err != nil {
				break
			}
			message.Meta = "me"
			messages.Data = append(messages.Data, message)
		}
	}
	rowsR, err := db.ptr.Query(queryR, receiver)
	if err != nil {
		errorLogger.Printf("%s", err.Error())
		return &messages
	}
	for rowsR.Next() {
		var message = Message{}
		if err = rowsR.Scan(&message.Date, &message.Data, &message.Type, &message.Check); err != nil {
			return &messages
		}
		message.Meta = "nome"
		println(message.Type)
		messages.Data = append(messages.Data, message)
	}
	err = db.updateStatus(sender, receiver)
	if err != nil {
		return nil
	}
	return &messages
}

func (db *DB) updateStatus(sender, receiver string) error {
	queryS := "UPDATE" + " `" + sender + "_rec` SET check_s = 1  WHERE sender=$1"
	queryR := "UPDATE" + " `" + receiver + "_send` SET check_s = 1  WHERE receiver=$1"
	_, err := db.ptr.Exec(queryS, receiver)
	if err != nil {
		return err
	}
	_, err = db.ptr.Exec(queryR, sender)
	if err != nil {
		return err
	}
	return nil
}
