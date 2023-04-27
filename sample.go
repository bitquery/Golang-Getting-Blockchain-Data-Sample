package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/go-retryablehttp"
)

type Response struct {
	Data interface{} `json:"data"`
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	url := "https://streaming.bitquery.io/graphql"
	method := "POST"
	apiKey := "YOUR KEY"

	payload := strings.NewReader(`{"query":"{\n  EVM(dataset: archive, network: eth) {\n    Transfers(\n      where: {Transfer: {Currency: {SmartContract: {is: \"0xBC4CA0EdA7647A8aB7C2061c2E118A18a936f13D\"}}}}\n      limit: {count: 10, offset: 0}\n      orderBy: {descending: Block_Number}\n    ) {\n      Transfer {\n        Currency {\n          SmartContract\n          Name\n          Decimals\n          Fungible\n          HasURI\n          Symbol\n        }\n        Id\n        URI\n        Data\n        owner: Receiver\n      }\n    }\n  }\n}\n","variables":"{}"}`)

	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 5
	retryClient.RetryWaitMin = 5 * time.Second
	retryClient.RetryWaitMax = 30 * time.Second

	req, err := retryablehttp.NewRequest(method, url, payload)
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-API-KEY", apiKey)

	res, err := retryClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Print the JSON response in the browser
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response.Data)
}
