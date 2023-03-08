package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ipaqsa/netcom/cryptoUtils"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"node/pkg/service"
	"sort"
	"strconv"
	"strings"
	"time"
)

func profileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		if !service.Node.Status {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		tmpl, err := template.ParseFiles("web/profile.gohtml")
		if err != nil {
			errorLogger.Println(err.Error())
			return
		}
		err = tmpl.Execute(w, service.Node.Person)
		if err != nil {
			errorLogger.Println(err.Error())
			return
		}
	}
}
func removeHandler(w http.ResponseWriter, r *http.Request) {
	infoLogger.Print("remove request")
	if r.Method == "POST" {
		if !service.Node.Status {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		err := service.Node.RemoveAccount()
		if err != nil {
			sendAnswer(w, "Error", err.Error())
			errorLogger.Println(err.Error())
			return
		}
		service.Node.Switch()
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
	http.Redirect(w, r, "/home", http.StatusSeeOther)
}
func registerHandler(w http.ResponseWriter, r *http.Request) {
	infoLogger.Print("register request")
	if r.Method == "GET" {
		if service.Node.Status {
			service.Node.Switch()
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		} else {
			http.ServeFile(w, r, "web/register.html")
		}
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		sendAnswer(w, "Error", err.Error())
		errorLogger.Println(err.Error())
		return
	}
	r.Body.Close()
	regs, err := parseRegister(body)
	if err != nil {
		sendAnswer(w, "Error", "parse fail")
		http.ServeFile(w, r, "web/register.html")
		return
	}
	room, _ := strconv.Atoi(regs.Room)
	err = service.Node.Register(regs.Login, regs.Password, regs.FirstName, regs.SecondName, room)
	if err != nil {
		sendAnswer(w, "Error", err.Error())
		return
	}
	sendAnswer(w, "OK", "")
}
func loginHandler(w http.ResponseWriter, r *http.Request) {
	infoLogger.Print("login request")
	if r.Method == "GET" {
		http.ServeFile(w, r, "web/login.html")
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		errorLogger.Println(err.Error())
	}
	r.Body.Close()
	regs, err := parseLogin(body)
	if err != nil {
		sendAnswer(w, "Error", "parse fail")
		http.ServeFile(w, r, "web/register.html")
		return
	}
	err = service.Node.Authentication(regs.Login, regs.Password)
	if err != nil {
		sendAnswer(w, "Error", err.Error())
		return
	}
	sendAnswer(w, "OK", "")
}
func homeHandler(w http.ResponseWriter, r *http.Request) {
	infoLogger.Print("home request")
	if r.Method == "GET" {
		if !service.Node.Status {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		var receiver = ""
		q := r.URL.RawQuery
		splits := strings.Split(q, "=")
		if len(splits) == 2 {
			receiver = splits[1]
		}
		tmpl, err := template.ParseFiles("web/main.gohtml")
		if err != nil {
			errorLogger.Println(err.Error())
			return
		}
		data := getData(receiver)
		err = tmpl.Execute(w, data)
		if err != nil {
			errorLogger.Println(err.Error())
			return
		}
	}
}
func sendHandler(w http.ResponseWriter, r *http.Request) {
	infoLogger.Println("send request")
	if r.Method == http.MethodPost {
		if !service.Node.Status {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		decoder := json.NewDecoder(r.Body)
		var t sendFromHTML
		err := decoder.Decode(&t)
		if err != nil {
			errorLogger.Println(err.Error())
			sendAnswer(w, "Error", err.Error())
		}
		splits := strings.Split(r.URL.RawQuery, "=")
		if len(splits) != 2 {
			errorLogger.Println(err.Error())
			sendAnswer(w, "Error", "enter receiver")
		}
		err = service.Node.Send(t.Data, splits[1])
		if err != nil {
			sendAnswer(w, "Error", err.Error())
			errorLogger.Println(err.Error())
			return
		}
		sendAnswer(w, "ok", "")
	}
}
func updateHandler(w http.ResponseWriter, r *http.Request) {
	infoLogger.Print("update request")
	if r.Method == "GET" {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}
	if !service.Node.Status {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		sendAnswer(w, "Error", err.Error())
		return
	}
	r.Body.Close()
	regs, err := parseRegister(body)
	if err != nil {
		sendAnswer(w, "Error", "parse fail")
		http.ServeFile(w, r, "web/profile.html")
		return
	}
	room, _ := strconv.Atoi(regs.Room)
	service.Node.Person.Firstname = regs.FirstName
	service.Node.Person.Lastname = regs.SecondName
	service.Node.Person.Room = room
	passwordChange := "no"
	if regs.Password != service.Node.Person.Hash {
		passwordChange = "yes"
		service.Node.Person.Hash = cryptoUtils.Base64Encode(cryptoUtils.HashSum([]byte(regs.Password)))
	}
	err = service.Node.Update(passwordChange)
	if err != nil {
		sendAnswer(w, "Error", err.Error())
		return
	}
	sendAnswer(w, "OK", "")
}
func fileHandler(w http.ResponseWriter, r *http.Request) {
	infoLogger.Print("file request")
	if r.Method == http.MethodPost {
		if !service.Node.Status {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		err := r.ParseMultipartForm(32 << 20)
		if err != nil {
			errorLogger.Printf("parse from client error: %s", err.Error())
			sendAnswer(w, "Error", err.Error())
			return
		}
		file, handler, err := r.FormFile("file")
		defer file.Close()
		if err != nil {
			errorLogger.Printf("get file error: %s", err.Error())
			sendAnswer(w, "Error", err.Error())
			return
		}
		splits := strings.Split(r.URL.RawQuery, "=")
		if len(splits) != 2 {
			errorLogger.Println(err.Error())
			sendAnswer(w, "Error", "enter receiver")
			return
		}
		bytes, err := io.ReadAll(file)
		err = service.Node.SendFile(bytes, handler.Filename, splits[1])
		if err != nil {
			errorLogger.Printf("get file error: %s", err.Error())
			sendAnswer(w, "Error", err.Error())
			return
		}
		sendAnswer(w, "ok", "")
	}
}
func downloadHandler(w http.ResponseWriter, r *http.Request) {
	infoLogger.Print("download request")
	if r.Method == http.MethodGet {
		if !service.Node.Status {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		q := r.URL.RawQuery
		parts := strings.Split(q, "&")
		if len(parts) != 2 {
			var err = errors.New("wrong query")
			sendAnswer(w, "Error", err.Error())
			errorLogger.Println(err.Error())
			return
		}
		receiverPart := strings.Split(parts[0], "=")
		if len(receiverPart) != 2 {
			var err = errors.New("wrong query")
			sendAnswer(w, "Error", err.Error())
			errorLogger.Println(err.Error())
			return
		}
		receiver := receiverPart[1]

		filenamePart := strings.Split(parts[1], "=")
		if len(filenamePart) != 2 {
			var err = errors.New("wrong query")
			sendAnswer(w, "Error", err.Error())
			errorLogger.Println(err.Error())
			return
		}
		filename := filenamePart[1]

		file, err := service.Node.GetFile(filename, receiver)
		if err != nil {
			sendAnswer(w, "Error", err.Error())
			errorLogger.Println(err.Error())
			return
		}
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Last-Modified", time.Now().UTC().Format(http.TimeFormat))
		_, err = w.Write(cryptoUtils.Base64Decode(file.Data))
		if err != nil {
			errorLogger.Println(err.Error())
			return
		}
		//http.ServeContent(w, r, filename, time.Now(), bytes.NewReader(cryptoUtils.Base64Decode(file.Data)))
	}
}

func updateContact() *ContactsToHTML {
	contacts, err := service.Node.GetContacts()
	if err != nil {
		errorLogger.Println(err.Error())
		return nil
	}
	return toContactsHTML(contacts)
}

func getData(receiver string) *DataToHTML {
	var messages service.Messages
	if receiver != "" {
		messagest, err := service.Node.Messages(receiver)
		if err != nil {
			errorLogger.Println(err.Error())
		} else {
			messages = *messagest
			sort.Slice(messages.Data, func(i, j int) bool {
				return messages.Data[i].Date < messages.Data[j].Date
			})
		}
	} else {
		messages = service.Messages{}
	}
	return &DataToHTML{
		Receiver: receiver,
		Contacts: updateContact(),
		Messages: messages,
	}
}
