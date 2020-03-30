package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	//"reflect"
)

func TestLoginHandler(t *testing.T) {
	server := httptest.NewServer(TestHandlers())
	fmt.Println("Test login handler")

	// Test for user exist
	resp, err := http.PostForm(server.URL+"/login", url.Values{"username": {"user1"}, "password": {"user1"}, "test": {"true"}})

	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode == 200 {
		fmt.Println("Test for user exist passed")
	} else {
		t.Fatalf("Test for user exist failed")
	}

	// Test for user not exist
	resp1, err1 := http.PostForm(server.URL+"/login", url.Values{"username": {"fakeuser1"}, "password": {"fakepwd1"}, "test": {"true"}})

	if err1 != nil {
		t.Error(err1)
	}

	if resp1.StatusCode == 409 {
		fmt.Println("Test for user not exist passed")
	} else {
		t.Fatalf("Test for user not exist failed")
	}
}

func TestRegisterHandler(t *testing.T) {
	server := httptest.NewServer(TestHandlers())
	fmt.Println("Test register handler")

	// Test for duplicated user
	resp, err := http.PostForm(server.URL+"/register", url.Values{"username": {"user1"}, "password": {"user1"}, "test": {"true"}})

	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode == 409 {
		fmt.Println("Test for duplicated user passed")
	} else {
		t.Fatalf("Test for duplicated user failed")
	}

	// Test for non-duplicated user
	resp1, err1 := http.PostForm(server.URL+"/register", url.Values{"username": {"newuser1"}, "password": {"newpwd1"}, "test": {"true"}})

	if err1 != nil {
		t.Error(err1)
	}

	if resp1.StatusCode == 200 {
		fmt.Println("Test for non-duplicated user passed")
	} else {
		t.Fatalf("Test for non-duplicated user failed")
	}
}

func TestHomeHandler(t *testing.T) {
	server := httptest.NewServer(TestHandlers())
	fmt.Println("Test home handler")

	// test user exist
	resp, err := http.PostForm(server.URL+"/home", url.Values{"username": {"user1"}, "postcontent": {"This is a test."}, "test": {"true"}})

	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode == 200 {
		fmt.Println("Test for user who login passed")
	} else {
		t.Fatalf("Test for user who login failed")
	}

	// test user not exit
	resp1, err1 := http.PostForm(server.URL+"/home", url.Values{"username": {"user1"}, "postcontent": {"This is a test too."}, "test": {"false"}})

	if err1 != nil {
		t.Error(err1)
	}

	if resp1.StatusCode == 401 {
		fmt.Println("Test for user who not login passed")
	} else {
		t.Fatalf("Test for user who not login failed")
	}
}

func TestFollowHandler(t *testing.T) {
	server := httptest.NewServer(TestHandlers())
	fmt.Println("Test follow user")

	// test exist user
	resp, err := http.PostForm(server.URL+"/follow", url.Values{"username": {"user1"}, "test": {"true"}})

	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode == 200 {
		fmt.Println("Test for follow exist user passed")
	} else {
		t.Fatalf("Test for follow exist user failed")
	}

	// test not exist user
	resp1, err1 := http.PostForm(server.URL+"/follow", url.Values{"username": {"fakeuser1"}, "test": {"true"}})

	if err1 != nil {
		t.Error(err)
	}
    
	if resp1.StatusCode == 404 {
		fmt.Println("Test for follow not exist user passed")
	} else {
		t.Fatalf("Test for follow not exist user failed")
	}
}

func TestCancelHandler(t *testing.T) {
	server := httptest.NewServer(TestHandlers())
	fmt.Println("Test delete user account")

	resp, err := http.PostForm(server.URL+"/cancel", url.Values{"username": {"user1"}, "test": {"true"}})

	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode == 200 {
		fmt.Println("Test for delete user account passed")
	} else {
		t.Fatalf("Test for delete user account failed")
	}
}
