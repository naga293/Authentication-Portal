package googleAuth

import (
	"authApp/auth"
	"authApp/db"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type Mail struct {
	Id            string
	Email         string
	VerifiedEmail string
	Name          string
	GivenName     string
	FamilyName    string
	Picture       string
	Locale        string
}

func GoogleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Create oauthState cookie
	oauthState := GenerateStateOauthCookie(w)
	/*
		AuthCodeURL receive state that is a token to protect the user
		from CSRF attacks. You must always provide a non-empty string
		and validate that it matches the the state query parameter
		on your redirect callback.
	*/
	u := AppConfig.GoogleLoginConfig.AuthCodeURL(oauthState)
	http.Redirect(w, r, u, http.StatusTemporaryRedirect)
}

func GoogleCallback(w http.ResponseWriter, r *http.Request) {
	// check is method is correct
	if r.Method != "GET" {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// get oauth state from cookie for this user
	oauthState, _ := r.Cookie("oauthstate")
	state := r.FormValue("state")
	code := r.FormValue("code")
	//w.Header().Add("content-type", "application/json")

	// ERROR : Invalid OAuth State
	if state != oauthState.Value {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		fmt.Fprintf(w, "invalid oauth google state")
		return
	}

	// Exchange Auth Code for Tokens
	token, err := AppConfig.GoogleLoginConfig.Exchange(
		context.Background(), code)

	// ERROR : Auth Code Exchange Failed
	if err != nil {
		http.Redirect(w, r, "/profile", http.StatusSeeOther)
		fmt.Println(w, "falied code exchange: %s", err.Error())
	}

	// Fetch User Data from google server
	response, err := http.Get(OauthGoogleUrlAPI + token.AccessToken)

	// ERROR : Unable to get user data from google
	if err != nil {
		fmt.Fprintf(w, "failed getting user info: %s", err.Error())
		return
	}

	// Parse user data JSON Object
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Fprintf(w, "failed read response: %s", err.Error())
		return
	}

	// send back response to browser
	var mail Mail
	json.Unmarshal(contents, &mail)
	res, err := db.Db.Query("SELECT Id,Password,Email from USERS WHERE EMAIL=?", mail.Email)
	if err != nil {
		fmt.Fprintln(w, err)
	}
	if res.Next() {
		var id int
		var pd string
		var amail string
		res.Scan(&id, &pd, &amail)
		token1, err := auth.GenerateJWT(mail.Email, strconv.Itoa(id))
		if err != nil {
			fmt.Fprintln(w, err.Error())
		}
		expiration := time.Now().Add(24 * time.Hour)
		cookie := http.Cookie{Name: "token", Value: token1, Expires: expiration}
		http.SetCookie(w, &cookie)
		http.Redirect(w, r, "/profile", http.StatusSeeOther)
		return

	} else {
		_, err = db.Db.Query("INSERT INTO Users VALUES(Null,?,?,'gmailAccount','city','yes')", mail.Name, mail.Email)
		if err != nil {
			fmt.Fprintln(w, err.Error())
		}
		res, err := db.Db.Query("SELECT ID from USERS WHERE EMAIL=?", mail.Email)

		if err != nil {
			fmt.Fprintln(w, err.Error())
		}
		var id int
		if res.Next() {
			res.Scan(&id)
		} else {
			fmt.Fprintln(w, "Please try to signup again")
		}

		token1, err := auth.GenerateJWT(mail.Email, strconv.Itoa(id))
		if err != nil {
			fmt.Fprintln(w, err.Error())
		}
		expiration := time.Now().Add(24 * time.Hour)
		cookie := http.Cookie{Name: "token", Value: token1, Expires: expiration}
		http.SetCookie(w, &cookie)
		t, err := template.ParseFiles("templates/profile_form.gtpl")
		if err != nil {
			fmt.Fprintln(w, err)
		}
		err = t.Execute(w, mail)

		if err != nil {
			fmt.Fprintln(w, err)
		}
	}

}
