package githubapi

import (
	"encoding/json"
	"log"
	"net/http"
)

type GitHubObject struct {
	Sha  string `json:"sha"`
	Type string `json:"type"`
	Url  string `json:"url"`
}

type GitHubTag struct {
	Ref    string `json:"ref"`
	NodeId string `json:"node_id"`
	Url    string `json:"url"`
	Object GitHubObject
}

func GetRepoTags(org string, repo string) []GitHubTag {
	var apiEndpoint = "https://api.github.com/repos/" + org + "/" + repo + "/git/refs/tags"
	var tags = &[]GitHubTag{}

	req, err := http.NewRequest("GET", apiEndpoint, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("X-GitHub-Api-Version", "2022-11-28")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	if res.StatusCode != http.StatusOK {
		log.Fatalf("%s: %s", apiEndpoint, res.Status)
	}
	err = json.NewDecoder(res.Body).Decode(tags)
	if err != nil {
		log.Fatal(err)
	}
	return *tags
}
