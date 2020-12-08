package main

import (
	"OwnProvider/apple"
	"OwnProvider/jwt"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	fmt.Println("OwnProvider inner only starting...")

	p8 := os.Getenv("OWNPROVIDERP8")
	if p8 == "" {
		fmt.Println("Error: ENV - OWNPROVIDERP8 is empty")
		return
	}

	logger, err := os.OpenFile("/var/log/ownprovider.log", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if nil != err {
		fmt.Println("Can not open log file")
	}
	log.SetOutput(logger)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.Println("------------------------------------------")

	// Home Page
	http.HandleFunc("/api/notify", Push)

	// Server
	http.ListenAndServe(":27953", nil)

	// This doesn't print until after we quit
	// fmt.Println("OwnProvider inner only server ready!!!")
}

func Push(w http.ResponseWriter, r *http.Request) {
	env := r.FormValue("env")
	pushType := r.FormValue("voip")
	deviceToken := r.FormValue("token")
	payload := r.FormValue("payload")
	topic := r.FormValue("bundleid")
	expiration := r.FormValue("exp")
	priority := r.FormValue("priority")
	collapseid := r.FormValue("collapseid")
	iss := r.FormValue("teamid")
	key := r.FormValue("keyid")

	if "" != pushType {
		pushType = "alert"
	}

	now := time.Now()
	duration, _ := time.ParseDuration("20m")
	roundedTime := now.Round(duration) // Restrict token renewal to every 20 minutes (Apple's minimum)

	jwtHeader := jwt.Header{
		Alg: "ES256",
		Kid: key,
	}
	jwtPayload := jwt.Payload{
		Iss: iss,
		Iat: roundedTime.Unix(),
	}

	jwToken, err := jwt.Token(jwtHeader, jwtPayload, "")
	if nil != err {
		log.Println("Build JWT token failure before push")
		w.Write([]byte("Error before push."))
		return
	}

	apnsId := uuid.New().String()
	var result map[string]interface{}
	json.Unmarshal([]byte(payload), &result)
	aps, ok := result["aps"].(map[string]interface{})
	if ok {
		if nil != aps["type"] {
			tp := aps["type"].(float64)
			if 2 == tp {
				//kind = "CHAT"
			} else if 3 == tp {
				//kind = "ADMIRER"
			} else if 4 == tp {
				//kind = "GIFTSENT"
			}
		}

		if nil != aps["apnsid"] {
			id := aps["apnsid"].(string)
			if 36 == len(id) {
				apnsId = id
			}
		}
	}

	h := apple.Header{
		Method:         "POST",
		Path:           "/3/device/" + deviceToken,
		Authorization:  "bearer " + jwToken,
		ApnsPushType:   pushType, // alert | background | voip | complication | fileprovider | mdm
		ApnsId:         apnsId,
		ApnsExpiration: expiration,
		ApnsPriority:   priority,
		ApnsTopic:      topic,
		ApnsCollapseId: collapseid,
	}
	httpHeader, err := h.Build()

	//log.Printf("HTTP HEADER : %v", httpHeader)

	server := apple.Gold
	if "sandbox" == env {
		server = apple.Dev
	}
	t := apple.Target{server, httpHeader, []byte(payload), deviceToken}
	resp, err := t.Notify()

	if nil != err {
		log.Printf("Network erro: %v", err)
	}

	respHttpBody, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	if 200 == resp.StatusCode {
		log.Println(deviceToken + "|" + strings.Join(resp.Header["Apns-Id"], " "))
		w.Write([]byte("SUCCESS|" + deviceToken + "|" + strings.Join(resp.Header["Apns-Id"], " ") + "\n"))
	} else {
		log.Println(deviceToken + "|" + resp.Status + ":" + string(respHttpBody[:]))
		w.Write([]byte("FAIL|" + deviceToken + "|" + resp.Status + "|" + string(respHttpBody[:]) + "\n"))
	}
}