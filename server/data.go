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

package main

import "crypto/md5"
import "golang.org/x/crypto/bcrypt"
import "os"
import "time"
import "net/http"
import "encoding/hex"
import "database/sql"

import "gopkg.in/gomail.v2"

import "github.com/teris-io/shortid"

import "github.com/gorilla/sessions"

const PasswordCost = 15

// Sessions
// =====================================================================================================================

var SessionStore sessions.Store

func init() {
	rawkey := []byte(os.Getenv("RSN2_SESSIONS_KEY"))
	key := make([]byte, hex.DecodedLen(len(rawkey)))
	_, err := hex.Decode(key, rawkey)
	if err != nil {
		panic("Could not load session key.\n" + err.Error())
	}

	SessionStore = sessions.NewCookieStore(key)
}

// GetSession returns the user id for the current user and 200, or an empty string and a HTTP error code.
func GetSession(l *SessionLogger, w http.ResponseWriter, r *http.Request) (string, int) {
	session, _ := SessionStore.Get(r, "rsn2-session")
	if auth, ok := session.Values["auth"].(bool); !ok || !auth {
		l.W.Printf("Error loading auth state from session.\n")
		return "", http.StatusForbidden
	}
	user, ok := session.Values["user"].(string)
	if !ok || user == "" {
		l.W.Printf("Error loading user from session.\n")
		return "", http.StatusBadRequest
	}
	return user, http.StatusOK
}

// Background Updates
// =====================================================================================================================

func GetAllFeeds(l *SessionLogger) [][2]string {
	rows, err := Queries["GetAllFeeds"].Preped.Query()
	if err != nil {
		l.E.Printf("Feed list failed for background update, error: %v\n", err)
		return nil
	}
	defer rows.Close()

	feeds := [][2]string{}
	for rows.Next() {
		f, id := "", ""
		err := rows.Scan(&f, &id)
		if err != nil {
			l.E.Printf("Feed list failed for background update, error: %v\n", err)
			return nil
		}
		feeds = append(feeds, [2]string{f, id})
	}
	return feeds
}

func ArticleExists(l *SessionLogger, url string) (exists, ok bool) {
	article := ""
	err := Queries["ArticleExistsByURL"].Preped.QueryRow(url).Scan(&article)
	if err != nil {
		l.E.Printf("DB existence check failed for new article %v, error: %v\n", url, err)
		return false, false
	}
	return article != "", true
}

var articleIDService <-chan string

func init() {
	go func() {
		c := make(chan string)
		articleIDService = c

		idsource := shortid.MustNew(5, shortid.DefaultABC, uint64(time.Now().UnixNano()))

		for {
			c <- idsource.MustGenerate()
		}
	}()
}

func ArticleAdd(l *SessionLogger, feed, title, url string, published time.Time) {
	article := <-articleIDService
	_, err := Queries["ArticleAdd"].Preped.Exec(article, feed, title, url, published.Unix())
	if err != nil {
		l.E.Printf("Cannot insert article %v into db, error: %v\n", url, err)
	}
}

func FeedListSubs(l *SessionLogger, feed string) []string {
	rows, err := Queries["FeedListSubs"].Preped.Query(feed)
	if err != nil {
		l.E.Printf("Feed subscribed user list failed for background update, error: %v\n", err)
		return nil
	}
	defer rows.Close()

	users := []string{}
	for rows.Next() {
		user := ""
		err := rows.Scan(&user)
		if err != nil {
			l.E.Printf("Feed subscribed user list failed for background update, error: %v\n", err)
			return nil
		}
		users = append(users, user)
	}
	return users
}

// /api/user/confirm-email
// =====================================================================================================================

// EmailConfirm takes the given confirmation string and id and makes the user as confirmed in the DB if
// the confirmation code is correct.
func EmailConfirm(l *SessionLogger, confirmation string) int {
	// This isn't really a high security thing, so the confirmation code is just an MD5ed version of the email
	// with a user ID appended to the end. If a user really wants to fake an email, there really isn't anything
	// stopping them from using a throwaway account.
	if len(confirmation) <= 32 {
		l.W.Println("Invalid email confirmation code, too short.")
		return http.StatusBadRequest // Too short.
	}

	rawhash := confirmation[:32]
	id := confirmation[32:]

	hash := make([]byte, hex.DecodedLen(len(rawhash)))
	_, err := hex.Decode(hash, []byte(rawhash))
	if err != nil {
		l.E.Printf("Error decoding hash for %v from DB, error: %v\n", id, err)
		return http.StatusBadRequest
	}

	if len(hash) != 16 {
		l.E.Println("You done goofed (impossible condition happened).")
		return http.StatusInternalServerError
	}

	// Now to get the user's email from the DB.
	email := ""
	err = Queries["GetEmail"].Preped.QueryRow(id).Scan(&email)
	if err == sql.ErrNoRows {
		l.E.Printf("User %v does not exist.\n", id)
		return http.StatusBadRequest
	}
	if err != nil {
		l.E.Printf("Error fetching email for %v from DB, error: %v\n", id, err)
		return http.StatusInternalServerError
	}

	truth := md5.Sum([]byte(email))

	// Make sure the hash matches the truth.
	for i := 0; i < 16; i++ {
		if hash[i] != truth[i] {
			l.W.Println("Invalid email confirmation token provided.")
			return http.StatusBadRequest
		}
	}

	// And finally mark the email confirmed.
	_, err = Queries["ConfirmEmail"].Preped.Exec(id)
	if err != nil {
		l.E.Printf("Error confirming email for %v (%v), error: %v\n", email, id, err)
		return http.StatusInternalServerError
	}

	return http.StatusOK
}

// EmailForceConfirm forcably confirms an email for the given user.
func EmailForceConfirm(l *SessionLogger, id string) int {
	_, err := Queries["ConfirmEmail"].Preped.Exec(id)
	if err != nil {
		l.E.Printf("Error force-confirming email for %v, error: %v\n", id, err)
		return http.StatusInternalServerError
	}
	return http.StatusOK
}

// EmailDelete takes the given confirmation string and id and deletes the user from the DB if
// the confirmation code is correct.
func EmailDelete(l *SessionLogger, confirmation string) int {
	// This isn't really a high security thing, so the confirmation code is just an MD5ed version of the email
	// with a user ID appended to the end. If a user really wants to fake an email, there really isn't anything
	// stopping them from using a throwaway account.
	if len(confirmation) <= 32 {
		l.W.Println("Invalid email confirmation code, too short.")
		return http.StatusBadRequest // Too short.
	}

	rawhash := confirmation[:32]
	id := confirmation[32:]

	hash := make([]byte, hex.DecodedLen(len(rawhash)))
	_, err := hex.Decode(hash, []byte(rawhash))
	if err != nil {
		l.E.Printf("Error decoding hash for %v from DB, error: %v\n", id, err)
		return http.StatusBadRequest
	}

	if len(hash) != 16 {
		l.E.Println("You done goofed (impossible condition happened).")
		return http.StatusInternalServerError
	}

	// Now to get the user's email from the DB.
	email := ""
	err = Queries["GetEmail"].Preped.QueryRow(id).Scan(&email)
	if err == sql.ErrNoRows {
		l.E.Printf("User %v does not exist.\n", id)
		return http.StatusBadRequest
	}
	if err != nil {
		l.E.Printf("Error fetching email for %v from DB, error: %v\n", id, err)
		return http.StatusInternalServerError
	}

	truth := md5.Sum([]byte(email))

	// Make sure the hash matches the truth.
	for i := 0; i < 16; i++ {
		if hash[i] != truth[i] {
			l.W.Println("Invalid email confirmation token provided.")
			return http.StatusBadRequest
		}
	}

	// And finally remove the user from the DB.
	_, err = Queries["DeleteEmail"].Preped.Exec(id)
	if err != nil {
		l.E.Printf("Error deleting email for %v (%v), error: %v\n", email, id, err)
		return http.StatusInternalServerError
	}

	return http.StatusOK
}

// /api/user/login (one row)
// =====================================================================================================================
// also used for user creation
type UserLoginData struct {
	Email    string
	Password string
}

// UserLogin returns false for valid if the username or password is wrong, and false for canlogin if the
// email is not confirmed.
func UserLogin(l *SessionLogger, email, password string) (code int, id string, canlogin bool) {
	dbpass := ""
	err := Queries["UserLogin"].Preped.QueryRow(email).Scan(&id, &dbpass, &canlogin)
	if err != nil {
		l.W.Printf("Cannot find user %v in db, error: %v\n", email, err)
		return http.StatusBadRequest, "", false
	}

	err = bcrypt.CompareHashAndPassword([]byte(dbpass), []byte(password))
	if err != nil {
		l.W.Printf("Password check failed for user %v (%v), error: %v\n", email, id, err)
		return http.StatusBadRequest, "", false
	}

	return http.StatusOK, id, canlogin
}

// /api/user/new
// =====================================================================================================================

var userIDService <-chan string
var SMTPPassword string
var Domain string

func init() {
	go func() {
		c := make(chan string)
		userIDService = c

		idsource := shortid.MustNew(6, shortid.DefaultABC, uint64(time.Now().UnixNano()))

		for {
			c <- idsource.MustGenerate()
		}
	}()

	SMTPPassword = os.Getenv("RSN2_SMTP_PASSWORD")
	Domain = os.Getenv("RSN2_DOMAIN")
}

func UserNew(l *SessionLogger, email, password string) int {
	id := <-userIDService

	// Make sure the user doesn't exist.
	ok := 0
	err := Queries["UserEmailExists"].Preped.QueryRow(email).Scan(&ok)
	if err != nil {
		l.E.Printf("DB existence check failed for new user %v, error: %v\n", email, err)
		return http.StatusInternalServerError
	}
	if ok == 1 {
		l.W.Printf("User %v already exists.\n", email)
		return http.StatusBadRequest
	}

	// Hash password
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), PasswordCost)
	if err != nil {
		l.W.Printf("Cannot insert user %v (%v) into db, error: %v\n", email, id, err)
		return http.StatusInternalServerError
	}

	_, err = Queries["UserNew"].Preped.Exec(id, email, string(hashed))
	if err != nil {
		l.E.Printf("Cannot insert user %v (%v) into db, error: %v\n", email, id, err)
		return http.StatusInternalServerError
	}

	// Generate confirmation token.
	src := md5.Sum([]byte(email))
	token := make([]byte, hex.EncodedLen(len(src)))
	hex.Encode(token, src[:])
	url := Domain + "/confirm-email?token=" + string(token) + id
	durl := Domain + "/delete-email?token=" + string(token) + id

	// Send confirmation email
	m := gomail.NewMessage()
	m.SetHeader("From", "noreply@httpcolonslashslashwww.com")
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Verify your Email")
	m.SetBody("text/html", `
<h2>Welcome to RSN2!</h2>
<p>Before you can start using your new account you need to verify your email by clicking the following link:</p>
<a href="`+url+`">`+url+`</a>
<p>If you did not make this account, you can <a href="`+durl+`">delete it</a> instead.
	`)

	go func() {
		d := gomail.NewDialer("smtp-relay.sendinblue.com", 587, "milo@httpscolonslashslashwww.com", SMTPPassword)
		if err := d.DialAndSend(m); err != nil {
			l.E.Printf("Could not send confirmation email for %v (%v): %v\n", email, id, err)
			return
		}
		l.I.Printf("Confirmation email for %v (%v) sent!\n", email, id)
	}()

	return http.StatusOK
}

type UserNewPassData struct {
	OldPassword string
	Password    string
}

func UserNewPass(l *SessionLogger, user, oldpassword, newpassword string) int {
	// Hash password
	hashed, err := bcrypt.GenerateFromPassword([]byte(newpassword), PasswordCost)
	if err != nil {
		l.W.Printf("Cannot update user %v with new password, error: %v\n", user, err)
		return http.StatusInternalServerError
	}

	dbpass := ""
	err = Queries["UserGetPass"].Preped.QueryRow(user).Scan(&dbpass)
	if err != nil {
		l.W.Printf("Cannot find user %v in db, error: %v\n", user, err)
		return http.StatusBadRequest
	}

	err = bcrypt.CompareHashAndPassword([]byte(dbpass), []byte(oldpassword))
	if err != nil {
		l.W.Printf("Password check failed for user %v, error: %v\n", user, err)
		return http.StatusBadRequest
	}

	_, err = Queries["UserNewPass"].Preped.Exec(user, string(hashed))
	if err != nil {
		l.E.Printf("Cannot update user %v with new password, error: %v\n", user, err)
		return http.StatusInternalServerError
	}
	return http.StatusOK
}

func UserNewName(l *SessionLogger, user, password, email string) int {
	// Make sure the user doesn't exist.
	ok := 0
	err := Queries["UserEmailExists"].Preped.QueryRow(email).Scan(&ok)
	if err != nil {
		l.E.Printf("DB existence check failed for user email change %v, error: %v\n", email, err)
		return http.StatusInternalServerError
	}
	if ok == 1 {
		l.W.Printf("User with email %v already exists.\n", email)
		return http.StatusBadRequest
	}

	dbpass := ""
	err = Queries["UserGetPass"].Preped.QueryRow(user).Scan(&dbpass)
	if err != nil {
		l.W.Printf("Cannot find user %v in db, error: %v\n", user, err)
		return http.StatusBadRequest
	}

	err = bcrypt.CompareHashAndPassword([]byte(dbpass), []byte(password))
	if err != nil {
		l.W.Printf("Password check failed for user %v, error: %v\n", user, err)
		return http.StatusBadRequest
	}

	_, err = Queries["UserNewName"].Preped.Exec(user, email)
	if err != nil {
		l.E.Printf("Cannot update user %v (%v) email in db, error: %v\n", email, user, err)
		return http.StatusInternalServerError
	}

	// Generate confirmation token.
	src := md5.Sum([]byte(email))
	token := make([]byte, hex.EncodedLen(len(src)))
	hex.Encode(token, src[:])
	url := Domain + "/confirm-email?token=" + string(token) + user

	// Send confirmation email
	m := gomail.NewMessage()
	m.SetHeader("From", "noreply@httpcolonslashslashwww.com")
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Verify your Email")
	m.SetBody("text/html", `
<h2>Thank you for using RSN2!</h2>
<p>Before you can start using your account again you will need to verify your email by clicking the following link:</p>
<a href="`+url+`">`+url+`</a>
	`)

	go func() {
		d := gomail.NewDialer("smtp-relay.sendinblue.com", 587, "milo@httpscolonslashslashwww.com", SMTPPassword)
		if err := d.DialAndSend(m); err != nil {
			l.E.Printf("Could not send confirmation email for %v (%v): %v\n", email, user, err)
			return
		}
		l.I.Printf("Confirmation email for %v (%v) sent!\n", email, user)
	}()

	return http.StatusOK
}

// /api/feed/list
// =====================================================================================================================

type Feed struct {
	ID     string
	Name   string
	URL    string
	Paused bool
}

func FeedList(l *SessionLogger, id string) []*Feed {
	rows, err := Queries["FeedList"].Preped.Query(id)
	if err != nil {
		l.E.Printf("Feed list failed for user %v, error: %v\n", id, err)
		return nil
	}
	defer rows.Close()

	feeds := []*Feed{}
	for rows.Next() {
		f := &Feed{}
		err := rows.Scan(&f.ID, &f.Name, &f.URL, &f.Paused)
		if err != nil {
			l.E.Printf("Feed list failed for user %v, error: %v\n", id, err)
			return nil
		}
		feeds = append(feeds, f)
	}
	return feeds
}

// /api/feed/details (one row)
// =====================================================================================================================

func FeedDetails(l *SessionLogger, user, feed string) *Feed {
	f := &Feed{}
	err := Queries["FeedDetails"].Preped.QueryRow(user, feed).Scan(&f.ID, &f.Name, &f.URL, &f.Paused)
	if err != nil {
		l.W.Printf("Error reading feed %v for user %v, error: %v\n", feed, user, err)
		return nil
	}
	return f
}

// /api/feed/articles
// =====================================================================================================================

type Article struct {
	ID        string
	Title     string
	URL       string
	Published time.Time
	Read      bool
}

func FeedArticles(l *SessionLogger, user, feed string) []*Article {
	rows, err := Queries["FeedArticles"].Preped.Query(user, feed)
	if err != nil {
		l.E.Printf("Feed article list failed for feed %v, user %v. Error: %v\n", feed, user, err)
		return nil
	}
	defer rows.Close()

	articles := []*Article{}
	for rows.Next() {
		a := &Article{}
		var stamp int64
		err := rows.Scan(&a.ID, &a.Title, &a.URL, &stamp, &a.Read)
		if err != nil {
			l.E.Printf("Feed article list failed for feed %v, user %v. Error: %v\n", feed, user, err)
			return nil
		}
		a.Published = time.Unix(stamp, 0)
		articles = append(articles, a)
	}
	return articles
}

// /api/feed/subscribe
// =====================================================================================================================

var feedIDService <-chan string

func init() {
	go func() {
		c := make(chan string)
		feedIDService = c

		idsource := shortid.MustNew(7, shortid.DefaultABC, uint64(time.Now().UnixNano()))

		for {
			c <- idsource.MustGenerate()
		}
	}()
}

type FeedSubscribeData struct {
	URL  string
	Name string
}

func FeedSubscribe(l *SessionLogger, id, url, name string) int {
	// First things first: Check to see if a feed with this url already esists.
	feed := ""
	err := Queries["FeedExistsByURL"].Preped.QueryRow(url).Scan(&feed)
	if err != nil {
		l.E.Printf("DB existence check failed for new feed %v, error: %v\n", url, err)
		return http.StatusInternalServerError
	}
	if feed == "" {
		// Create new feed.
		feed = <-feedIDService
		_, err = Queries["FeedAdd"].Preped.Exec(feed, url)
		if err != nil {
			l.E.Printf("Cannot insert feed %v into db, error: %v\n", url, err)
			return http.StatusInternalServerError
		}
	}

	ok := 0
	err = Queries["FeedAlreadySubscibed"].Preped.QueryRow(id, feed).Scan(&ok)
	if err != nil {
		l.E.Printf("DB existence check failed for subscribed feed %v by user %v, error: %v\n", feed, id, err)
		return http.StatusInternalServerError
	}
	if ok == 1 {
		l.W.Printf("Feed %v already subscribed by user %v.\n", feed, id)
		// This isn't a straight up error, but it isn't OK either.
		return http.StatusAccepted
	}

	_, err = Queries["FeedSubscibe"].Preped.Exec(id, feed, name)
	if err != nil {
		l.E.Printf("Failed subscribing feed %v as user %v, error: %v\n", feed, id, err)
		return http.StatusInternalServerError
	}
	return http.StatusOK
}

// /api/feed/unsubscribe
// =====================================================================================================================

func FeedUnsub(l *SessionLogger, user, feed string) int {
	_, err := Queries["FeedUnsub1"].Preped.Exec(user, feed)
	if err != nil {
		l.E.Printf("Failed unsubscribing feed %v as user %v, error: %v\n", feed, user, err)
		return http.StatusInternalServerError
	}

	// Now check if the feed has no subscribers.
	hassub := 0
	err = Queries["FeedHasSubs"].Preped.QueryRow(feed).Scan(&hassub)
	if err != nil {
		l.E.Printf("Feed subscriber check failed for feed %v, error: %v\n", feed, err)
		return http.StatusInternalServerError
	}
	if hassub == 1 {
		// If the feed still has other subscribers delete our paused flags and slink off into the night.
		_, err = Queries["FeedUnsub2"].Preped.Exec(user, feed)
		if err != nil {
			l.E.Printf("Failed unsubscribing feed %v as user %v, error: %v\n", feed, user, err)
			return http.StatusInternalServerError
		}
		return http.StatusOK
	}

	// No subscribers left, delete feed for real.
	_, err = Queries["FeedDelete"].Preped.Exec(feed)
	if err != nil {
		l.E.Printf("Failed deleting feed %v, error: %v\n", feed, err)
		return http.StatusInternalServerError
	}
	return http.StatusOK
}

// /api/feed/pause
// =====================================================================================================================

func FeedPause(l *SessionLogger, user, feed string) int {
	_, err := Queries["FeedPause"].Preped.Exec(user, feed)
	if err != nil {
		l.E.Printf("Failed pausing feed %v, error: %v\n", feed, err)
		return http.StatusInternalServerError
	}
	return http.StatusOK
}

// //api/feed/unpause
// =====================================================================================================================

func FeedUnpause(l *SessionLogger, user, feed string) int {
	_, err := Queries["FeedUnpause"].Preped.Exec(user, feed)
	if err != nil {
		l.E.Printf("Failed unpausing feed %v, error: %v\n", feed, err)
		return http.StatusInternalServerError
	}
	return http.StatusOK
}

// /api/article/read
// =====================================================================================================================

func ArticleMarkRead(l *SessionLogger, user, article string) int {
	_, err := Queries["ArticleRead"].Preped.Exec(user, article)
	if err != nil {
		l.E.Printf("Failed marking article (%v) read, error: %v\n", article, err)
		return http.StatusInternalServerError
	}
	return http.StatusOK
}

// /api/article/unread
// =====================================================================================================================

func ArticleMarkUnread(l *SessionLogger, user, article string) int {
	_, err := Queries["ArticleUnread"].Preped.Exec(user, article)
	if err != nil {
		l.E.Printf("Failed marking article (%v) unread, error: %v\n", article, err)
		return http.StatusInternalServerError
	}
	return http.StatusOK
}

// /api/article/feed
// =====================================================================================================================

type UnreadArticle struct {
	ID        string
	Title     string
	URL       string
	FeedName  string // Feed *name*, not ID.
	Published time.Time
}

func GetUnread(l *SessionLogger, user string) []*UnreadArticle {
	rows, err := Queries["GetUnread"].Preped.Query(user)
	if err != nil {
		l.E.Printf("Unread article list failed for user %v. Error: %v\n", user, err)
		return nil
	}
	defer rows.Close()

	articles := []*UnreadArticle{}
	for rows.Next() {
		a := &UnreadArticle{}
		var stamp int64
		err := rows.Scan(&a.ID, &a.Title, &a.URL, &a.FeedName, &stamp)
		if err != nil {
			l.E.Printf("Unread article list failed for user %v. Error: %v\n", user, err)
			return nil
		}
		a.Published = time.Unix(stamp, 0)
		articles = append(articles, a)
	}
	return articles
}
