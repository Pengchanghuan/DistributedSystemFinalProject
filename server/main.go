package main

import (
	//"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"sort"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
)

var db map[string]*User
var HTMLADDRESS = "../static/view/"

var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))

type User struct {
	Username string
	Password string
	Followee map[string]bool
	Follower map[string]bool
	Tweets   []Tweet
}

func (user *User) addTweet(content string, timestamp time.Time) {
	message := Tweet{
		Username:        user.Username,
		Timestamp:       timestamp,
		Content:         content,
		StringTimeStamp: timestamp.Format("Mon Jan _2 15:15:15 2020"),
	}
	user.Tweets = append(user.Tweets, message)
}

func newUser(username, password string) *User {
	user1 := &User{
		Username: username,
		Password: password,
		Followee: make(map[string]bool),
		Follower: make(map[string]bool),
	}
	return user1
}

type Tweet struct {
	Username        string
	Timestamp       time.Time
	Content         string
	StringTimeStamp string
}

type HtmlMessages struct {
	Username   string
	Followee   map[string]bool
	Follower   map[string]bool
	Tweets     []Tweet
	UserTweets []Tweet
	Show       bool
}

func initial() {
	db = make(map[string]*User) //username -> user

	//-------below are all test content-------
	user1 := newUser("user1", "user1")
	user1.addTweet("wo fa de", time.Now())
	db["user1"] = user1
}

func TestHandlers() *mux.Router {
	fmt.Println("Start testing web server")
	initial()

	var router = mux.NewRouter()

	router.HandleFunc("/", loginHandler)
	router.HandleFunc("/login", loginHandler)
	router.HandleFunc("/register", registerHandler)
	router.HandleFunc("/home", homeHandler)
	router.HandleFunc("/follow", followHandler)
	router.HandleFunc("/cancel", cancelHandler)
	//http.ListenAndServe(":8080", nil)

	return router
}

func registerHandler(response http.ResponseWriter, request *http.Request) {
	fmt.Println("register page")
	if request.Method == "GET" {
		t, _ := template.ParseFiles(HTMLADDRESS + "register.html")
		t.Execute(response, nil)
	} else if request.Method == "POST" {
		request.ParseForm()

		username := request.FormValue("username")
		password := request.FormValue("password")
		test := request.FormValue("test")

		if _, exist := db[username]; exist {
			response.WriteHeader(409)
			fmt.Fprintf(response, "<script>alert('Duplicate username')</script>")

			t, _ := template.ParseFiles(HTMLADDRESS + "register.html")
			t.Execute(response, nil)
		} else {
			if test == "true" {
				response.WriteHeader(200)
			}
			user := newUser(username, password)
			db[username] = user
			setUsernameSession(username, response)
			http.Redirect(response, request, "/home", http.StatusSeeOther)
		}
	}
}

func loginHandler(response http.ResponseWriter, request *http.Request) {
	fmt.Println("login page")
	if request.Method == "GET" {
		t, _ := template.ParseFiles(HTMLADDRESS + "login.html")
		t.Execute(response, nil)
	} else if request.Method == "POST" {
		request.ParseForm()

		username := request.FormValue("username")
		password := request.FormValue("password")
		test := request.FormValue("test")
		_, exist := db[username]

		if !exist || db[username].Password != password {
			response.WriteHeader(409)
			fmt.Fprintf(response, "<script>alert('Username or password is incorrect')</script>")
			t, _ := template.ParseFiles(HTMLADDRESS + "login.html")
			t.Execute(response, nil)
		} else {
			if test == "true" {
				response.WriteHeader(200)
			}
			setUsernameSession(username, response)
			http.Redirect(response, request, "/home", http.StatusSeeOther)
		}
	}
}

func setUsernameSession(username string, response http.ResponseWriter) {
	fmt.Printf("set Username:%v\n", username)
	value := map[string]string{
		"username": username,
	}
	if encoded, err := cookieHandler.Encode("UsernameSession", value); err == nil {
		cookie := &http.Cookie{
			Name:  "UsernameSession",
			Value: encoded,
			Path:  "/",
		}
		http.SetCookie(response, cookie)
	}
}

func getUserName(request *http.Request) (username string) {
	if cookie, err := request.Cookie("UsernameSession"); err == nil {
		cookieValue := make(map[string]string)
		if err = cookieHandler.Decode("UsernameSession", cookie.Value, &cookieValue); err == nil {
			username = cookieValue["username"]
		}
	}
	return username
}

func getUser(request *http.Request) *User {
	username := getUserName(request)
	if username != "" {
		return db[username]
	}

	return nil
}

func clearUsernameSession(response http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   "UsernameSession",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(response, cookie)
}

func setSearchedUsernameSession(username string, response http.ResponseWriter) {
	value := map[string]string{
		"username": username,
	}
	if encoded, err := cookieHandler.Encode("SearchedUsernameSession", value); err == nil {
		cookie := &http.Cookie{
			Name:  "SearchedUsernameSession",
			Value: encoded,
			Path:  "/",
		}
		http.SetCookie(response, cookie)
	}
}

func getSearchedUsername(request *http.Request) (username string) {
	if cookie, err := request.Cookie("SearchedUsernameSession"); err == nil {
		cookieValue := make(map[string]string)
		if err = cookieHandler.Decode("SearchedUsernameSession", cookie.Value, &cookieValue); err == nil {
			username = cookieValue["username"]
		}
	}
	return username
}

func clearSearchedUsernameSession(response http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   "SearchedUsernameSession",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(response, cookie)
}

func homeHandler(response http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	username := request.FormValue("username")
	test := request.FormValue("test")

	if test != "true" {
		if getUserName(request) == "" {
			response.WriteHeader(401)
			http.Redirect(response, request, "/login", http.StatusSeeOther)
			return
		}
	}

	clearSearchedUsernameSession(response)
	if request.Method == "GET" {
		fmt.Println("home get")
		t, _ := template.ParseFiles(HTMLADDRESS + "home.html")

		userInfo := getUser(request)

		tweets := []Tweet{}
		for followee := range userInfo.Followee {
			user := db[followee]
			for _, message := range user.Tweets {
				tweets = append(tweets, message)
			}
		}
		sort.Slice(tweets, func(i, j int) bool {
			return tweets[i].Timestamp.After(tweets[j].Timestamp)
		})
		if len(tweets) > 10 {
			tweets = tweets[0:10]
		}

		userTweets := []Tweet{}
		for _, message := range userInfo.Tweets {
			userTweets = append(userTweets, message)
		}
		sort.Slice(userTweets, func(i, j int) bool {
			return userTweets[i].Timestamp.After(userTweets[j].Timestamp)
		})
		if len(userTweets) > 10 {
			userTweets = userTweets[0:10]
		}

		htmlMessages := HtmlMessages{
			Username:   userInfo.Username,
			Followee:   userInfo.Followee,
			Follower:   userInfo.Follower,
			Tweets:     tweets,
			UserTweets: userTweets,
		}
		t.Execute(response, htmlMessages)
	} else if request.Method == "POST" {
		fmt.Println("home post")

		var userInfo *User
		if test == "true" {
			userInfo = db[username]
			response.WriteHeader(200)
		} else {
			userInfo = getUser(request)
		}

		userInfo.addTweet(request.FormValue("postcontent"), time.Now())

		http.Redirect(response, request, "/home", http.StatusSeeOther)
	}
}

func followHandler(response http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	test := request.FormValue("test")

	if test != "true" {
		if getUserName(request) == "" {
			response.WriteHeader(401)
			http.Redirect(response, request, "/login", http.StatusSeeOther)
			return
		}
	}

	if request.Method == "GET" {
		fmt.Println("follow get")
		t, _ := template.ParseFiles(HTMLADDRESS + "follow.html")
		t.Execute(response, nil)
	} else if request.Method == "POST" {
		fmt.Println("follow post")

		followeeName := request.FormValue("username")
		fmt.Println("followeeName:", followeeName)
		// add request
		if followeeName == "" {
			followeeName := getSearchedUsername(request)
			clearSearchedUsernameSession(response)
			username := getUserName(request)

			user := db[username]
			followee := db[followeeName]

			user.Followee[followeeName] = true
			followee.Follower[username] = true

			http.Redirect(response, request, "/home", http.StatusSeeOther)
			return
		}

		followee, exist := db[followeeName]

		if exist {
			htmlMessages := HtmlMessages{}
			tweets := []Tweet{}

			for _, message := range followee.Tweets {
				tweets = append(tweets, message)
			}

			var userInfo *User
			if test == "true" {
				userInfo = db[followeeName]
			} else {
				userInfo = getUser(request)
			}

			show := true
			if _, ok := userInfo.Followee[followeeName]; ok {
				show = false
			} else if followeeName == getUserName(request) {
				show = false
			} else if test == "true" {
				show = false
			}

			htmlMessages.Username = followeeName
			htmlMessages.Tweets = tweets
			htmlMessages.Show = show

			clearSearchedUsernameSession(response)
			setSearchedUsernameSession(followeeName, response)

			t, _ := template.ParseFiles(HTMLADDRESS + "follow.html")
			t.Execute(response, htmlMessages)
		} else {
			response.WriteHeader(404)
			fmt.Fprintf(response, "<script>alert('No users found')</script>")
			t, _ := template.ParseFiles(HTMLADDRESS + "follow.html")
			t.Execute(response, nil)
		}
	}
}

func cancelHandler(response http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	test := request.FormValue("test")

	if test != "true" {
		if getUserName(request) == "" {
			response.WriteHeader(401)
			http.Redirect(response, request, "/login", http.StatusSeeOther)
			return
		}
	}

	if request.Method == "GET" {
		t, _ := template.ParseFiles(HTMLADDRESS + "cancel.html")
		t.Execute(response, nil)
	} else if request.Method == "POST" {
		var username string
		if test == "true" {
			username = request.FormValue("username")
		} else {
			username = getUserName(request)
		}

		delete(db, username)
		clearUsernameSession(response)
		http.Redirect(response, request, "/login", http.StatusSeeOther)
	}
}

func main() {
	fmt.Println("Start web server")
	initial()

	var router = mux.NewRouter()
	http.Handle("/", router)
	fs := http.FileServer(http.Dir(HTMLADDRESS))
	router.PathPrefix("/js/").Handler(fs)
	router.PathPrefix("/css/").Handler(fs)
	router.PathPrefix("/images/").Handler(fs)

	router.HandleFunc("/", loginHandler)
	router.HandleFunc("/login", loginHandler)
	router.HandleFunc("/register", registerHandler)
	router.HandleFunc("/home", homeHandler)
	router.HandleFunc("/follow", followHandler)
	router.HandleFunc("/cancel", cancelHandler)

	http.ListenAndServe(":8080", nil)
}
