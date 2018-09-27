package main

// Yeah, this isn't a proper Go test. I just needed something to iterate on
// to get the formatting of the email.

// To run this, fromthis dir just do:
//   go test

import (
	"encoding/json"
	"fmt"
	"github"
	"io/ioutil"
	"path/filepath"
	"testing"
)

func TestBuildEmailMessageBody(t *testing.T) {
	g, err := ioutil.ReadFile(filepath.Join("example.json"))
	if err != nil {
		t.Fatalf("failed reading example.json: %s", err)
	}

	var githubPushEvent github.PushEvent
	err = json.Unmarshal([]byte(g), &githubPushEvent)
	if nil != err {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(githubPushEvent)

	fmt.Println(BuildEmailSubject(githubPushEvent))
	emailBody := BuildEmailMessageBody(githubPushEvent)
	fmt.Println(emailBody)
}
