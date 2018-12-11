package main

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"log/syslog"
	"encoding/json"
	"bytes"
	"time"
)

type OneTouchStatus string

var (
	OneTouchStatusApproved OneTouchStatus = "approved"
	OneTouchStatusPending OneTouchStatus = "pending"
	OneTouchStatusDenied OneTouchStatus = "denied"
	OneTouchStatusExpired OneTouchStatus = "expired"
)

type ApprovalRequest struct {
	Status   OneTouchStatus `json:"status"`
	UUID     string         `json:"uuid"`
	Notified bool           `json:"notified"`
}

type APIResponse struct {
	Success         bool             `json:"success"`
	ApprovalRequest *ApprovalRequest `json:"approval_request"`
	Message         string           `json:"message"`
}

type Authy struct {
	APIKey  string
	BaseURL string
}


func (authy *Authy) DoRequest(request *http.Request) (*http.Response, error) {
	request.Header.Set("X-Authy-API-Key", authy.APIKey)
    request.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(request)
    return resp, err
}


func (authy *Authy) SendApprovalRequest(userID string, message string) (OneTouchStatus, error) {
	sb := bytes.NewBufferString(authy.BaseURL)
	if (authy.BaseURL[len(authy.BaseURL) - 1] != '/') {
		sb.WriteString("/")
	}
	sb.WriteString("onetouch/json/users/")
	sb.WriteString(userID)
	sb.WriteString("/approval_requests")
    req, err := http.NewRequest("GET", sb.String(), nil)
    if err != nil {
        panic(err)
    }

    resp, err := authy.DoRequest(req)
    if err != nil {
        panic(err)
    }

    defer resp.Body.Close()
    body, _ := ioutil.ReadAll(resp.Body)

    var jsonResponse APIResponse
    err = json.Unmarshal(body, &jsonResponse)
	if err != nil {
        panic(err)
    }

    fmt.Println("response Body:", jsonResponse)

    sb.Reset()
    sb.WriteString(authy.BaseURL)
	if (authy.BaseURL[len(authy.BaseURL) - 1] != '/') {
		sb.WriteString("/")
	}
	sb.WriteString("onetouch/json/approval_requests/")
	sb.WriteString(jsonResponse.ApprovalRequest.UUID)

	for i := 0; i < 20; i++ {
		req, err = http.NewRequest("GET", sb.String(), nil)

		resp, err = authy.DoRequest(req)
	    if err != nil {
	        panic(err)
	    }

	    defer resp.Body.Close()
	    body, _ = ioutil.ReadAll(resp.Body)
	    err = json.Unmarshal(body, &jsonResponse)
		if err != nil {
	        panic(err)
	    }

    	fmt.Println("response Body:", jsonResponse.ApprovalRequest)

    	if jsonResponse.ApprovalRequest.Status != OneTouchStatusPending {
    		break
    	}
    	time.Sleep(500 * time.Millisecond)
    }

    return jsonResponse.ApprovalRequest.Status, nil
}

func pamLog(format string, args ...interface{}) {
	l, err := syslog.New(syslog.LOG_AUTH|syslog.LOG_WARNING, "pam-authy")
	if err != nil {
		return
	}
	l.Warning(fmt.Sprintf(format, args...))
}

func main() {
	config := Config{}
	//config.LoadFromFile("/etc/pam_authy.conf")
	config.LoadFromFile("data.conf")
	fmt.Println(config.Users["root2"])
	authy := Authy{
		APIKey:  config.Authy.Token,
		BaseURL: config.Authy.URL,
	}
	authy.SendApprovalRequest("12312332", "say yes plsplspls")
}
