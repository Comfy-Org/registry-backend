package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
)

var pubID = flag.String("publisher-id", "", "Pubisher ID that is represented by the JWT token")
var token = flag.String("token", os.Getenv("JWT_TOKEN"), "JWT token to use for authentication")
var baseURL = flag.String("base-url", os.Getenv("REGISTRY_BASE_URL"), "Base url of registry service")

func main() {
	flag.Parse()

	if pubID == nil || *pubID == "" {
		flag.PrintDefaults()
		log.Fatalf("Flag `--publisher-id` must be set to non-empty string.\n")
	}
	if token == nil || *token == "" {
		flag.PrintDefaults()
		log.Fatalf("Flag `--token` or environment variable `JWT_TOKEN` must be set to non-empty string.\n")
	}
	if baseURL == nil || *baseURL == "" {
		flag.PrintDefaults()
		log.Fatalf("Flag `--base-url` or environment variable `BASE_URL` must be set to non-empty string.\n")
	}

	u, err := url.Parse(*baseURL)
	if err != nil {
		log.Fatalf("Invalid base url :%v .\n", err)
	}
	u = u.JoinPath("publishers", *pubID, "ban")

	req, _ := http.NewRequest(http.MethodPost, u.String(), nil)
	req.Header.Set("Authorization", "Bearer "+*token)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Failed to send HTTP request :%v\n", err)
	}
	b, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("Failed to read HTTP response :%v\n", err)
	}
	if res.StatusCode > 299 {
		log.Fatalf("Received non-success response: \nStatus: %d\nBody: \n\n%v\n", res.StatusCode, string(b))
	}
	fmt.Printf("Publisher '%s' has been banned\n", *pubID)

}
