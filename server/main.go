/*
Copyright 2020-2021 by Milo Christiansen

This software is provided 'as-is', without any express or implied warranty. In
no event will the authors be held liable for any damages arising from the use of
this software.

Permission is granted to anyone to use this software for any purpose, including
commercial applications, and to alter it and redistribute it freely, subject to
the following restrictions:

1. The origin of this software must not be misrepresented; you must not claim
that you wrote the original software. If you use this software in a product, an
acknowledgment in the product documentation would be appreciated but is not
required.

2. Altered source versions must be plainly marked as such, and must not be
misrepresented as being the original software.

3. This notice may not be removed or altered from any source distribution.
*/

// RSN2: Multi-user RSS feed tracker.
package main

import "fmt"
import "mime"
import "net/url"
import "net/http"
import "encoding/json"

import "github.com/milochristiansen/axis2"
import "github.com/milochristiansen/axis2/sources"

const MaxBodyBytes = int64(65536)

func main() {
	// /api/user/confirm-email
	http.HandleFunc("/api/user/confirm-email", func(w http.ResponseWriter, r *http.Request) {
		l := newSessionLogger("/api/user/confirm-email")

		v, ok := r.URL.Query()["token"]
		if !ok || len(v) == 0 {
			l.W.Println("No email confirmation token provided.")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if status := EmailConfirm(l, v[0]); status != http.StatusOK {
			w.WriteHeader(status)
			return
		}

		w.WriteHeader(http.StatusOK)
	})

	// /api/user/delete-email
	http.HandleFunc("/api/user/delete-email", func(w http.ResponseWriter, r *http.Request) {
		l := newSessionLogger("/api/user/delete-email")

		v, ok := r.URL.Query()["token"]
		if !ok || len(v) == 0 {
			l.W.Println("No email deletion token provided.")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if status := EmailDelete(l, v[0]); status != http.StatusOK {
			w.WriteHeader(status)
			return
		}

		w.WriteHeader(http.StatusOK)
	})

	// /api/user/logged-in
	http.HandleFunc("/api/user/logged-in", func(w http.ResponseWriter, r *http.Request) {
		l := newSessionLogger("/api/user/login")

		user, status := GetSession(l, w, r)
		if user == "" {
			w.WriteHeader(status)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	// /api/user/login
	http.HandleFunc("/api/user/login", func(w http.ResponseWriter, r *http.Request) {
		l := newSessionLogger("/api/user/login")

		r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)

		data := &UserLoginData{}
		err := json.NewDecoder(r.Body).Decode(data)
		if err != nil {
			l.W.Printf("Error parsing login body. Error: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		status, user, canlogin := UserLogin(l, data.Email, data.Password)
		if status != http.StatusOK {
			w.WriteHeader(status)
			return
		}

		if !canlogin {
			l.W.Printf("Login attempt for unconfirmed user %v\n", data.Email)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		session, _ := SessionStore.Get(r, "rsn2-session")
		session.Values["user"] = user
		session.Values["auth"] = true
		err = session.Save(r, w)
		if err != nil {
			l.W.Printf("Error saving session. Error: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	})

	// /api/user/logout
	http.HandleFunc("/api/user/logout", func(w http.ResponseWriter, r *http.Request) {
		l := newSessionLogger("/api/user/login")

		session, _ := SessionStore.Get(r, "rsn2-session")
		session.Values["user"] = ""
		session.Values["auth"] = false
		err := session.Save(r, w)
		if err != nil {
			l.W.Printf("Error saving session. Error: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	})

	// /api/user/new
	http.HandleFunc("/api/user/new", func(w http.ResponseWriter, r *http.Request) {
		l := newSessionLogger("/api/user/new")

		r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)

		data := &UserLoginData{}
		err := json.NewDecoder(r.Body).Decode(data)
		if err != nil {
			l.W.Printf("Error parsing user creation body. Error: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(UserNew(l, data.Email, data.Password))
	})

	// /api/feed/list
	http.HandleFunc("/api/feed/list", func(w http.ResponseWriter, r *http.Request) {
		l := newSessionLogger("/api/feed/list")

		user, status := GetSession(l, w, r)
		if user == "" {
			w.WriteHeader(status)
			return
		}

		feeds := FeedList(l, user)
		if feeds == nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err := json.NewEncoder(w).Encode(feeds)
		if err != nil {
			l.E.Printf("Error encoding payload. Error: %v\n", err)
			return
		}
	})

	// /api/feed/details
	http.HandleFunc("/api/feed/details", func(w http.ResponseWriter, r *http.Request) {
		l := newSessionLogger("/api/feed/details")

		user, status := GetSession(l, w, r)
		if user == "" {
			w.WriteHeader(status)
			return
		}

		feed := r.FormValue("id")
		if feed == "" {
			l.W.Printf("Missing feed ID.\n")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		details := FeedDetails(l, user, feed)
		if details == nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err := json.NewEncoder(w).Encode(details)
		if err != nil {
			l.E.Printf("Error encoding payload. Error: %v\n", err)
			return
		}
	})

	// /api/feed/articles
	http.HandleFunc("/api/feed/articles", func(w http.ResponseWriter, r *http.Request) {
		l := newSessionLogger("/api/feed/articles")

		user, status := GetSession(l, w, r)
		if user == "" {
			w.WriteHeader(status)
			return
		}

		feed := r.FormValue("id")
		if feed == "" {
			l.W.Printf("Missing feed ID.\n")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		articles := FeedArticles(l, user, feed)
		if articles == nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err := json.NewEncoder(w).Encode(articles)
		if err != nil {
			l.W.Printf("Error encoding payload. Error: %v\n", err)
			return
		}
	})

	// /api/feed/subscribe
	http.HandleFunc("/api/feed/subscribe", func(w http.ResponseWriter, r *http.Request) {
		l := newSessionLogger("/api/feed/subscribe")

		user, status := GetSession(l, w, r)
		if user == "" {
			w.WriteHeader(status)
			return
		}

		r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)

		data := &FeedSubscribeData{}
		err := json.NewDecoder(r.Body).Decode(data)
		if err != nil {
			l.W.Printf("Error parsing feed subscribe body. Error: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if data.Name == "" {
			l.W.Printf("No feed name given.\n")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		_, err = url.ParseRequestURI(data.URL)
		if err != nil {
			l.W.Printf("Malformed URL. Error: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(FeedSubscribe(l, user, data.URL, data.Name))
	})

	// /api/feed/unsubscribe
	http.HandleFunc("/api/feed/unsubscribe", func(w http.ResponseWriter, r *http.Request) {
		l := newSessionLogger("/api/feed/unsubscribe")

		user, status := GetSession(l, w, r)
		if user == "" {
			w.WriteHeader(status)
			return
		}

		feed := r.FormValue("id")
		if feed == "" {
			l.W.Printf("Missing feed ID.\n")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(FeedUnsub(l, user, feed))
	})

	// /api/feed/pause
	http.HandleFunc("/api/feed/pause", func(w http.ResponseWriter, r *http.Request) {
		l := newSessionLogger("/api/feed/pause")

		user, status := GetSession(l, w, r)
		if user == "" {
			w.WriteHeader(status)
			return
		}

		feed := r.FormValue("id")
		if feed == "" {
			l.W.Printf("Missing feed ID.\n")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		s := FeedPause(l, user, feed)
		if s == http.StatusOK {
			Feeds.BroadcastTo(l, user)
		}
		w.WriteHeader(s)
	})

	// /api/feed/unpause
	http.HandleFunc("/api/feed/unpause", func(w http.ResponseWriter, r *http.Request) {
		l := newSessionLogger("/api/feed/unpause")

		user, status := GetSession(l, w, r)
		if user == "" {
			w.WriteHeader(status)
			return
		}

		feed := r.FormValue("id")
		if feed == "" {
			l.W.Printf("Missing feed ID.\n")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		s := FeedUnpause(l, user, feed)
		if s == http.StatusOK {
			Feeds.BroadcastTo(l, user)
		}
		w.WriteHeader(s)
	})

	// /api/article/read
	http.HandleFunc("/api/article/read", func(w http.ResponseWriter, r *http.Request) {
		l := newSessionLogger("/api/article/read")

		user, status := GetSession(l, w, r)
		if user == "" {
			w.WriteHeader(status)
			return
		}

		article := r.FormValue("id")
		if article == "" {
			l.W.Printf("Missing article ID.\n")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		s := ArticleMarkRead(l, user, article)
		if s == http.StatusOK {
			Feeds.BroadcastTo(l, user)
		}
		w.WriteHeader(s)
	})

	// /api/article/unread
	http.HandleFunc("/api/article/unread", func(w http.ResponseWriter, r *http.Request) {
		l := newSessionLogger("/api/article/unread")

		user, status := GetSession(l, w, r)
		if user == "" {
			w.WriteHeader(status)
			return
		}

		article := r.FormValue("id")
		if article == "" {
			l.W.Printf("Missing article ID.\n")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		s := ArticleMarkUnread(l, user, article)
		if s == http.StatusOK {
			Feeds.BroadcastTo(l, user)
		}
		w.WriteHeader(s)
	})

	// /api/article/feed
	http.HandleFunc("/api/article/feed", func(w http.ResponseWriter, r *http.Request) {
		l := newSessionLogger("/api/article/feed")

		user, status := GetSession(l, w, r)
		if user == "" {
			w.WriteHeader(status)
			return
		}

		Feeds.Upgrade(l, w, r, user)
	})

	ml.I.Println("Initializing AXIS VFS.")
	fs := new(axis2.FileSystem)
	fs.Mount("", sources.NewOSDir("/app/html"), false)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		l := newSessionLogger("/")

		l.I.Println("Page Request:" + r.URL.Path)

		typ := mime.TypeByExtension(GetExt(r.URL.Path))
		content, err := fs.ReadAll(r.URL.Path[1:])
		if err != nil {
			l.E.Println("  Error:", err)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if typ != "" {
			w.Header().Set("Content-Type", typ)
		}
		fmt.Fprintf(w, "%s", content)
	})

	go Background()

	// err := http.ListenAndServe(":3366", nil)
	// if err != nil {
	// panic(err)
	// }

	err := http.ListenAndServeTLS(":443", "/app/cert/server.crt", "/app/cert/server.key", nil)
	if err != nil {
		panic(err)
	}
}

func GetExt(name string) string {
	// Find the last part of the extension
	i := len(name) - 1
	for i >= 0 {
		if name[i] == '.' {
			return name[i:]
		}
		i--
	}
	return ""
}
