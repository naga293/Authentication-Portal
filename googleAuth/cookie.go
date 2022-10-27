package googleAuth

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"
)

func GenerateStateOauthCookie(w http.ResponseWriter) string {
	var expiration = time.Now().Add(1 * time.Hour) //time.Now().Add(2 * time.Minute)
	b := make([]byte, 16)
	rand.Read(b)

	state := base64.URLEncoding.EncodeToString(b)
	fmt.Println(state)
	cookie := http.Cookie{
		Name:     "oauthstate",
		Value:    state,
		Expires:  expiration,
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)

	return state
}
