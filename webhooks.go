package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type webhookRegistration struct {
	ID    int
	URL   string `json:"url"`
	Event string `json:"event"`
	Time  string
}

var webhooks []webhookRegistration

func webhookHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		// Expects incoming body in terms of WebhookRegistration struct
		webhook := webhookRegistration{}
		jsonFile, err := os.Open("testjson.json")
		if err != nil {
			fmt.Println("Problem opening json file")
		} else {
			fmt.Println("Successfully opened testjson.json")
		}
		defer jsonFile.Close()
		/*byteValue, _ := ioutil.ReadAll(jsonFile)
		json.Unmarshal(byteValue, &webhook)
		json.Unmarshal(byteValue, &webhooks)*/
		ok := json.NewDecoder(r.Body).Decode(&webhook)
		webhook.ID = len(webhooks) + 1
		webhook.Time = time.Now().String()
		if ok != nil {
			http.Error(w, "Something went wrong: "+err.Error(), http.StatusBadRequest)
		}
		webhooks = append(webhooks, webhook)
		file, _ := os.OpenFile("testjson.json", os.O_CREATE, os.ModePerm)
		defer file.Close()
		encoder := json.NewEncoder(file)
		encoder.Encode(webhooks)
		// Note: Approach does not guarantee persistence or permanence of resource id (for CRUD)
		fmt.Fprintln(w, len(webhooks)-1)
		fmt.Println("Webhook " + webhook.URL + " has been registered.")
	case http.MethodGet:
		// For now just return all webhooks, don't respond to specific resource requests
		path := r.URL.Path
		parts := strings.Split(path, "/")
		key := parts[4]
		if key == "" {
			err := json.NewEncoder(w).Encode(webhooks)
			if err != nil {
				http.Error(w, "l50 Something went wrong: "+err.Error(), http.StatusInternalServerError)
			}
		} else {
			intkey, err := strconv.Atoi(key)
			if err != nil {
				http.Error(w, "l45 Something went wrong: "+err.Error(), http.StatusInternalServerError)
			}
			for i := 0; i < len(webhooks); i++ {
				fmt.Println(i)
				fmt.Println(webhooks[i].ID)
				fmt.Println(len(webhooks))
				if intkey == webhooks[i].ID {
					err := json.NewEncoder(w).Encode(webhooks[i])
					if err != nil {
						http.Error(w, "Something went wrong: "+err.Error(), http.StatusInternalServerError)
					}
				}
			}
		}
	case http.MethodDelete:
		path := r.URL.Path
		parts := strings.Split(path, "/")
		key := parts[4]
		if key == "" {
			err := json.NewEncoder(w).Encode(webhooks)
			if err != nil {
				http.Error(w, "l72 Something went wrong: "+err.Error(), http.StatusInternalServerError)
			}
		} else {
			intkey, err := strconv.Atoi(key)
			if err != nil {
				http.Error(w, "l77 Something went wrong: "+err.Error(), http.StatusInternalServerError)
			}
			for i := 0; i < len(webhooks); i++ {
				fmt.Println(i)
				fmt.Println(webhooks[i].ID)
				fmt.Println(len(webhooks))
				if intkey == webhooks[i].ID {
					fmt.Println("Removing element ")
					fmt.Println(webhooks[i].ID)
					webhooks = append(webhooks[:i], webhooks[i+1:]...)
					file, _ := os.OpenFile("testjson.json", os.O_CREATE, os.ModePerm)
					defer file.Close()
					encoder := json.NewEncoder(file)
					encoder.Encode(webhooks)
					if err != nil {
						http.Error(w, "Something went wrong: "+err.Error(), http.StatusInternalServerError)
					}
				}
			}

		}

	default:
		http.Error(w, "Invalid method "+r.Method, http.StatusBadRequest)
	}
}

func serviceHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		fmt.Println("Received POST request...")
		for _, v := range webhooks {
			go callURL(v.URL, "Response on registered event in webhook demo: "+v.Event)
		}
	default:
		http.Error(w, "Invalid method "+r.Method, http.StatusBadRequest)
	}
}

/*
	Calls given URL with given content and awaits response (status and body).
*/
func callURL(url string, content string) {
	fmt.Println("Attempting invocation of url " + url + " ...")
	res, err := http.Post(url, "string", bytes.NewReader([]byte(content)))
	if err != nil {
		fmt.Println("Error in HTTP request: " + err.Error())
	}
	response, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Something is wrong with invocation response: " + err.Error())
	}

	fmt.Println("Webhook invoked. Received status code " + strconv.Itoa(res.StatusCode) +
		" and body: " + string(response))
}
