package slack

import (
	"encoding/json"
	"net/http"
	"fmt"
	"io/ioutil"
	"bytes"
	"gopkg.in/yaml.v2"
)

const baseUrl string = "https://slack.com/api/"

// Presence.Set() method 
// Presence.doRequest method
type Presence struct {
	Value string `json:"presence"`
}

type SlackWorkspace struct {
	Name  string `yaml:"name"`
	Token string `yaml:"token"`
}

type SlackConfig struct {
	Default   string           `yaml:"default"`
	Workspace []SlackWorkspace `yaml:"workspaces"`
}


func (p *Presence) Set(w SlackWorkspace) (string, error) { 
	method := "POST"
	endpoint := "users.setPresence"

	url := baseUrl+endpoint
	presenceData := Presence{p.Value}
	payload, err := json.Marshal(presenceData)
	
	if err != nil {
		return "", err
	}
	
	req, err := http.NewRequest(method, url, bytes.NewBuffer(payload))
	if err != nil {
		return "", err
	}

	return doRequest(req, w), nil

}

func doRequest(req *http.Request, workspace SlackWorkspace) string {
	client := &http.Client{}

	token := "Bearer "
	token += workspace.Token

	req.Header.Add("Content-type", "application/json; charset=utf-8")
	req.Header.Add("Authorization", token)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	return string(body)
}

func GetWorkspace(workspaceName string, config SlackConfig) SlackWorkspace {
	var workspace SlackWorkspace

	for index, _ := range config.Workspace {
		if config.Workspace[index].Name == workspaceName {
			workspace = config.Workspace[index]
			break
		}
	}
	return workspace
}

func ParseConfig(configFile string) SlackConfig {

	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		fmt.Println(err)
	}
	var config SlackConfig
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		fmt.Printf("Error parsing YAML file: %s\n", err)
	}

	return config
}