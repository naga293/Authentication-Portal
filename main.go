package main

import (
	"authApp/auth"
	"authApp/db"
	"authApp/googleAuth"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/mail"
	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	gotenv "github.com/subosito/gotenv"
)

// User struct declaration
type User struct {
	Id       int
	Name     string
	Email    string
	City     string
	Password string
	//Password string
}

type Message struct {
	Email            string
	password         string
	confirm_password string
	Errors           map[string]string
}

type Token struct {
	Cookie  string
	Message string
}

type Login struct {
	Username string
	Password string
}

type Signup struct {
	Username        string
	Password        string
	ConfirmPassword string
}

const JSON = "application/json"

func (msg *Message) Validate() bool {
	msg.Errors = make(map[string]string)
	_, err := mail.ParseAddress(msg.Email)
	//match := rxEmail.Match([]byte(msg.Email))
	if err != nil {
		msg.Errors["Email"] = "Please enter a valid email address"
	} else {
		id, err := db.Db.Query("SELECT ID from USERS WHERE EMAIL=?", msg.Email)
		if err != nil {
			msg.Errors["db"] = err.Error()
		}
		if id.Next() {
			msg.Errors["signup"] = "User Already Exists"
		}
	}

	if msg.password != msg.confirm_password {
		msg.Errors["confirm_password"] = "Passwords Do not match"
	}

	return len(msg.Errors) == 0
}

func login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) //get request method
	contentType := r.Header.Get("Content-type")
	if r.Method == "GET" {
		switch contentType {
		case JSON:
			var user Login
			user.Username = "example@email.com"
			user.Password = "examplepassword"
			w.Header().Set("Content-Type", JSON)
			json.NewEncoder(w).Encode(user)

		default:
			t, _ := template.ParseFiles("templates/login.gtpl")
			t.Execute(w, nil)
		}
	} else {
		var user Login
		switch contentType {
		case JSON:

			if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
				fmt.Println(err)
			}

		default:
			r.ParseForm()
			user.Username = r.PostFormValue("username")
			user.Password = r.PostFormValue("password")
		}
		res, err := db.Db.Query("SELECT Id,Password,gmail from USERS WHERE EMAIL=?", user.Username)
		if err != nil {
			fmt.Fprintln(w, err)
		}
		if res.Next() {
			var pwd string
			var id int
			var gmail string
			if gmail == "yes" {
				msg := "You have registered with your gmailaccount please login with it"
				t, _ := template.ParseFiles("templates/login.gtpl")
				t.Execute(w, msg)
				return
			}
			res.Scan(&id, &pwd)
			if pwd == user.Password {
				token, err := auth.GenerateJWT(user.Username, strconv.Itoa(id))
				if err != nil {
					fmt.Fprintln(w, err.Error())
				}
				expiration := time.Now().Add(24 * time.Hour)
				cookie := http.Cookie{Name: "token", Value: token, Expires: expiration}
				http.SetCookie(w, &cookie)
				switch contentType {
				case JSON:
					var currentSession Token
					currentSession.Cookie = token
					currentSession.Message = "logged in successfully"
					w.Header().Set("Content-Type", JSON)
					json.NewEncoder(w).Encode(currentSession)

				default:
					http.Redirect(w, r, "/profile", http.StatusSeeOther)
				}

			} else {
				msg := "Please Input Valid Credentials"
				switch contentType {

				case JSON:

					var currentSession Token
					currentSession.Cookie = "invalid"
					currentSession.Message = msg
					w.Header().Set("Content-Type", JSON)
					json.NewEncoder(w).Encode(currentSession)

				default:
					t, _ := template.ParseFiles("templates/login.gtpl")
					t.Execute(w, msg)
				}

			}
		} else {
			msg := "User account does not exist"
			switch contentType {

			case JSON:

				var currentSession Token
				currentSession.Cookie = "invalid"
				currentSession.Message = msg
				w.Header().Set("Content-Type", JSON)
				json.NewEncoder(w).Encode(currentSession)

			default:
				t, _ := template.ParseFiles("templates/login.gtpl")
				t.Execute(w, msg)
			}
		}
	}
}

func signup(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) //get request method
	contentType := r.Header.Get("Content-type")
	if r.Method == "GET" {
		switch contentType {
		case JSON:
			var user Signup
			user.Username = "example@email.com"
			user.Password = "examplepassword"
			user.ConfirmPassword = "examplepassword"
			w.Header().Set("Content-Type", JSON)
			json.NewEncoder(w).Encode(user)

		default:
			t, _ := template.ParseFiles("templates/signup.gtpl")
			t.Execute(w, nil)
		}

	} else {
		var user Signup
		switch contentType {
		case JSON:

			if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
				fmt.Println(err)
			}

		default:
			r.ParseForm()
			user.Username = r.PostFormValue("username")
			user.Password = r.PostFormValue("password")
			user.ConfirmPassword = r.PostFormValue("confirm_password")
		}
		msg := &Message{
			Email:            user.Username,
			password:         user.Password,
			confirm_password: user.ConfirmPassword,
		}

		if !msg.Validate() {
			switch contentType {

			case JSON:

				var currentSession Token
				currentSession.Cookie = "invalid"
				currentSession.Message = msg.Errors["Email"] + msg.Errors["db"] + msg.Errors["signup"] + msg.Errors["confirm_password"]
				w.Header().Set("Content-Type", JSON)
				json.NewEncoder(w).Encode(currentSession)

			default:
				t, err := template.ParseFiles("templates/signup.gtpl")
				if err != nil {
					fmt.Fprintln(w, err)
				}
				err = t.Execute(w, msg)
				if err != nil {
					fmt.Fprintln(w, err)
				}
				return
			}

		} else {
			_, err := db.Db.Query("INSERT INTO Users VALUES(Null,'name',?,?,'city','no')", msg.Email, msg.password)
			if err != nil {
				fmt.Fprintln(w, err.Error())
			}
			res, err := db.Db.Query("SELECT ID from USERS WHERE EMAIL=?", msg.Email)

			if err != nil {
				fmt.Fprintln(w, err.Error())
			}
			var id int
			if res.Next() {
				res.Scan(&id)
			} else {
				fmt.Fprintln(w, "Please try to signup again")
			}
			token, err := auth.GenerateJWT(msg.Email, strconv.Itoa(id))
			if err != nil {
				fmt.Fprintln(w, err.Error())
			}
			expiration := time.Now().Add(24 * time.Hour)
			cookie := http.Cookie{Name: "token", Value: token, Expires: expiration}
			http.SetCookie(w, &cookie)
			switch contentType {
			case JSON:
				var currentSession Token
				currentSession.Cookie = token
				currentSession.Message = "signed up successfully"
				w.Header().Set("Content-Type", JSON)
				json.NewEncoder(w).Encode(currentSession)

			default:
				http.Redirect(w, r, "/profile_form", http.StatusSeeOther)
			}

		}
	}
}

func profileForm(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method)
	tokenCookie, err := r.Cookie("token")
	contentType := r.Header.Get("Content-type")
	if err != nil || tokenCookie.Value == "" {
		switch contentType {
		case JSON:
			var currentSession Token
			currentSession.Message = "Session Expired Please login again" + err.Error()
			w.Header().Set("Content-Type", JSON)
			json.NewEncoder(w).Encode(currentSession)

		default:
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		}
		return
	}
	if r.Method == "GET" {
		//tokenCookie, err := r.Cookie("token")
		// if err != nil {
		// 	fmt.Fprintln(w, err.Error())
		// }
		token := tokenCookie.Value
		id, err := auth.ValidateToken(token)
		if err != nil {
			fmt.Fprintln(w, err.Error())
		}
		res, err := db.Db.Query("SELECT Name,City,Email from USERS WHERE Id=?", id)
		if err != nil {
			fmt.Fprintln(w, err.Error())
		}
		if res.Next() {
			var current User
			res.Scan(&current.Name, &current.City, &current.Email)
			switch contentType {
			case JSON:
				w.Header().Set("Content-Type", JSON)
				json.NewEncoder(w).Encode(current)

			default:
				t, _ := template.ParseFiles("templates/profileForm.gtpl")
				t.Execute(w, current)
			}

		}

	} else {
		//tokenCookie, err := r.Cookie("token")
		var user User
		switch contentType {
		case JSON:

			if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
				fmt.Println(err)
			}

		default:
			r.ParseForm()
			user.City = r.PostFormValue("city")
			user.Name = r.PostFormValue("name")
		}
		if err != nil {
			fmt.Fprintln(w, err.Error())
		}
		token := tokenCookie.Value
		id, err := auth.ValidateToken(token)
		if err != nil {
			fmt.Fprintln(w, err.Error())
		}
		//id, err = strconv.Atoi(id)
		_, err = db.Db.Query("UPDATE Users SET Name = ?, City= ? WHERE ID = ?", user.Name, user.City, id)
		if err != nil {
			fmt.Fprintln(w, err.Error())
		}
		switch contentType {
		case JSON:
			var currentSession Token
			currentSession.Cookie = token
			currentSession.Message = "data updated successfully"
			w.Header().Set("Content-Type", JSON)
			json.NewEncoder(w).Encode(currentSession)

		default:
			http.Redirect(w, r, "/profile", http.StatusSeeOther)
		}
	}
}

func profile(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) //get request method
	contentType := r.Header.Get("Content-type")
	tokenCookie, err := r.Cookie("token")
	if err != nil || tokenCookie.Value == "" {
		switch contentType {
		case JSON:
			var currentSession Token
			currentSession.Message = "Session Expired Please login again" + err.Error()
			w.Header().Set("Content-Type", JSON)
			json.NewEncoder(w).Encode(currentSession)

		default:
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		}
		return
	}
	if r.Method == "GET" {
		tokenCookie, err := r.Cookie("token")
		if err != nil {
			fmt.Fprintln(w, err.Error())
		}
		token := tokenCookie.Value
		id, err := auth.ValidateToken(token)
		if err != nil {
			fmt.Fprintln(w, err.Error())
		}
		res, err := db.Db.Query("SELECT ID,name,email,city FROM Users WHERE id = ?", id)
		if err != nil {
			fmt.Fprintln(w, err)
		}
		defer res.Close()

		if res.Next() {

			var current User

			err := res.Scan(&current.Id, &current.Name, &current.Email, &current.City)

			if err != nil {
				log.Fatal(err)
			}
			switch contentType {
			case JSON:
				w.Header().Set("Content-Type", JSON)
				json.NewEncoder(w).Encode(current)

			default:
				t, err := template.ParseFiles("templates/profile.gtpl")
				if err != nil {
					fmt.Fprintln(w, err)
				}
				err = t.Execute(w, current)

				if err != nil {
					fmt.Fprintln(w, err)
				}
			}

		} else {
			msg := "profile not found"
			switch contentType {
			case JSON:
				var current Token
				current.Message = msg
				w.Header().Set("Content-Type", JSON)
				json.NewEncoder(w).Encode(current)

			default:
				t, _ := template.ParseFiles("templates/login.gtpl")
				t.Execute(w, msg)
			}
		}
	} else {
		tokenCookie, err := r.Cookie("token")
		if err != nil {
			fmt.Fprintln(w, err.Error())
		}
		token := tokenCookie.Value
		id, err := auth.ValidateToken(token)
		if err != nil {
			fmt.Fprintln(w, err.Error())
		}
		res, _ := db.Db.Query("SELECT gmail FROM Users WHERE id = ?", id)
		if res.Next() {
			cookie := http.Cookie{
				Name:   "token",
				Value:  "",
				Path:   "/",
				MaxAge: -1,
			}
			cookie1 := http.Cookie{
				Name:   "oauthstate",
				Value:  "",
				Path:   "/",
				MaxAge: -1,
			}
			http.SetCookie(w, &cookie)
			http.SetCookie(w, &cookie1)
			var gmail string
			res.Scan(&gmail)
			if gmail == "yes" {
				http.Redirect(w, r, "https://mail.google.com/mail/u/0/?logout&hl=en", http.StatusSeeOther)
			} else {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
			}
		}

		//http.SetCookie(, cookie)
	}
}

func main() {
	port := os.Getenv("PORT")
	fmt.Println(port)
	googleAuth.LoadConfig()
	gotenv.Load()
	db.Db = db.OpenDB()
	defer db.Db.Close()
	http.HandleFunc("/", login)
	http.HandleFunc("/gmail_login", googleAuth.GoogleLogin)
	http.HandleFunc("/signup", signup)
	http.HandleFunc("/profile_form", profileForm)
	http.HandleFunc("/profile", profile)
	http.HandleFunc("/google_callback", googleAuth.GoogleCallback)
	err := http.ListenAndServe(":"+port, nil) // setting listening port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
