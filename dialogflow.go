package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var projectID = "PROJECT_ID"
var baseURL = "https://dialogflow.googleapis.com/v2/projects/" + projectID + "/agent"

func NewDialogflowClient() (*http.Client, error) {
	var scope = "https://www.googleapis.com/auth/dialogflow"

	data, err := ioutil.ReadFile("keys.json")
	if err != nil {
		log.Fatal(err)
	}
	conf, err := google.JWTConfigFromJSON(data, scope)
	if err != nil {
		log.Fatal(err)
	}
	//ts := conf.TokenSource(oauth2.NoContext)
	DialogflowClient := conf.Client(oauth2.NoContext)

	return DialogflowClient, nil
}

func getAgent(client *http.Client) {
	type Agent struct {
		Parent                  string  `json:"parent"`
		DisplayName             string  `json:"displayName"`
		DefaultLanguageCode     string  `json:"defaultLanguageCode"`
		TimeZone                string  `json:"timeZone"`
		EnableLogging           bool    `json:"enableLogging"`
		MatchMode               string  `json:"matchMode"`
		ClassificationThreshold float64 `json:"classificationThreshold"`
		APIVersion              string  `json:"apiVersion"`
		Tier                    string  `json:"tier"`
	}

	resp, err := client.Get(baseURL)
	if err != nil {
		log.Fatal(err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	agent := &Agent{}
	json.Unmarshal(body, agent)
	fmt.Println(agent)

	defer resp.Body.Close()
}

func listIntents(client *http.Client) {
	type Intent struct {
		Intents []struct {
			Name         string `json:"name"`
			DisplayName  string `json:"displayName"`
			Priority     int    `json:"priority"`
			WebhookState string `json:"webhookState,omitempty"`
			Parameters   []struct {
				Name                  string `json:"name"`
				DisplayName           string `json:"displayName"`
				Value                 string `json:"value"`
				EntityTypeDisplayName string `json:"entityTypeDisplayName"`
				Mandatory             bool   `json:"mandatory"`
			} `json:"parameters,omitempty"`
			Messages []struct {
				Text struct {
				} `json:"text"`
			} `json:"messages,omitempty"`
			DefaultResponsePlatforms []string `json:"defaultResponsePlatforms,omitempty"`
			Events                   []string `json:"events,omitempty"`
			Action                   string   `json:"action,omitempty"`
			IsFallback               bool     `json:"isFallback,omitempty"`
		} `json:"intents"`
	}

	url := baseURL + "/intents"

	resp, err := client.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	intents := &Intent{}
	json.Unmarshal(body, intents)
	fmt.Println(intents)

	defer resp.Body.Close()
}

func listEntities(client *http.Client) {
	type Entity struct {
		EntityTypes []struct {
			Name        string `json:"name"`
			DisplayName string `json:"displayName"`
			Kind        string `json:"kind"`
			Entities    []struct {
				Value    string   `json:"value"`
				Synonyms []string `json:"synonyms"`
			} `json:"entities,omitempty"`
		} `json:"entityTypes"`
	}

	url := baseURL + "/entityTypes"
	resp, err := client.Get(url)
	if err != nil {
		log.Println("Error on response.\n[ERRO] -", err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	entities := &Entity{}
	json.Unmarshal(body, entities)
	fmt.Println(entities)

	defer resp.Body.Close()
}

func createIntent(client *http.Client, name string) {
	url := baseURL + "/intents"
	var jsonbody = []byte(`{
		"displayName": "` + name + `",
		"webhookState": "WEBHOOK_STATE_ENABLED",
		"priority": 500000,
		"events": [
			"asdfghjkssdf"
		],
		"defaultResponsePlatforms": [
			"ACTIONS_ON_GOOGLE"
		],
		"trainingPhrases": [
			{
				"type": "TYPE_UNSPECIFIED",
				"parts": [
					{
						"text": "yeet1"
					}
				]
			}
		],
		"parameters": [
			{
				"displayName": "yeetyeet",
				"value": "yeet yeet yeet",
				"entityTypeDisplayName": "yeet",
				"mandatory": true,
				"prompts": [
					"What do you want?"
				],
				"isList": false
			}
		]
	}`)

	resp, err := client.Post(url, "application/json", bytes.NewReader(jsonbody))
	if err != nil {
		log.Fatal(err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))

	defer resp.Body.Close()
}

func createEntity(client *http.Client, name string) {
	url := baseURL + "/entityTypes"
	var jsonbody = []byte(`{
		"displayName": "` + name + `",
		"kind": "KIND_MAP"
	}`)

	resp, err := client.Post(url, "application/json", bytes.NewReader(jsonbody))
	if err != nil {
		log.Fatal(err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))

	defer resp.Body.Close()
}

func deleteIntent(client *http.Client, intentID string) {
	url := baseURL + "/intents/" + intentID

	req, err := http.NewRequest("DELETE", url, nil)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		return
	}
	statusOK := resp.StatusCode >= 200 && resp.StatusCode < 300
	if !statusOK {
		fmt.Println("Non-OK HTTP status:", resp.StatusCode)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println(body)

	defer resp.Body.Close()
}

func deleteEntity(client *http.Client, entityID string) {
	url := baseURL + "/entityTypes/" + entityID

	req, err := http.NewRequest("DELETE", url, nil)
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERRO] -", err)
		return
	}
	statusOK := resp.StatusCode >= 200 && resp.StatusCode < 300
	if !statusOK {
		fmt.Println("Non-OK HTTP status:", resp.StatusCode)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
	fmt.Println("congratulations! success!")

	defer resp.Body.Close()
}

func main() {
	client, err := NewDialogflowClient()
	if err != nil {
		log.Fatal(err)
	}

	//getAgent(client)

	listIntents(client)
	listEntities(client)

	createIntent(client, "apitest12345678")
	createEntity(client, "yeet123123245")

	deleteIntent(client, IntentID)
	deleteEntity(client, EntityID)
}
