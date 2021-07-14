package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/paketo-buildpacks/dep-server/pkg/dependency/licenses"
)

type DepMetadata struct {
	Version         string   `json:"version"`
	URI             string   `json:"uri"`
	SHA256          string   `json:"sha256"`
	Source          string   `json:"source"`
	SourceSHA256    string   `json:"source_sha256"`
	DeprecationDate string   `json:"deprecation_date"`
	CPE             string   `json:"cpe,omitempty"`
	Licenses        []string `json:"licenses"`
}

type DispatchDepMetadata struct {
	Version         string `json:"version"`
	URI             string `json:"uri"`
	SHA256          string `json:"sha256"`
	Source          string `json:"source_uri"`
	SourceSHA256    string `json:"source_sha256"`
	DeprecationDate string `json:"deprecation_date"`
	CPE             string `json:"cpe,omitempty"`
	Licenses        string `json:"licenses"`
}

func main() {
	// Takes in the name of 1 dep => dispatches to the test-upload workflow with all metadata
	dependencyName := "bundler"

	// reach out to the api.deps..../<dep-name> get all the metadata for all versions
	resp, err := http.Get(fmt.Sprintf("https://api.deps.paketo.io/v1/dependency?name=%s&per_page=100", dependencyName))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// translate JSON
	var deps []DepMetadata
	err = json.NewDecoder(resp.Body).Decode(&deps)
	if err != nil {
		log.Fatal(err)
	}

	// for each dep version ... get all metadata except licenses
	licenseRetriever := licenses.NewLicenseRetriever()
	for _, dep := range deps {
		if len(dep.Licenses) == 0 || dep.CPE == "" {
			fmt.Println(dep.Version)
			// pass the dep name and source URL and whatever else to pkg/dependency/licenses to get licenses
			licenses, err := licenseRetriever.LookupLicenses(dependencyName, dep.Source)
			if err != nil {
				log.Fatal(err)
			}

			dep.Licenses = licenses
			if dep.CPE == "" {
				dep.CPE = fmt.Sprintf("cpe:2.3:a:bundler:bundler:%s:*:*:*:*:ruby:*:*", dep.Version)
			}

			// dispatchDep is an exact copy of the dep, but the license are a string instead of slice.
			dispatchDep := DispatchDepMetadata{}
			dispatchDep.Version = dep.Version
			dispatchDep.URI = dep.URI
			dispatchDep.SHA256 = dep.SHA256
			dispatchDep.Source = dep.Source
			dispatchDep.SourceSHA256 = dep.SourceSHA256
			dispatchDep.DeprecationDate = dep.DeprecationDate
			dispatchDep.CPE = dep.CPE
			dispatchDep.Licenses = strings.Join(dep.Licenses, ",")

			payload, err := json.Marshal(dispatchDep)
			if err != nil {
				log.Fatal(err)
			}

			var dispatch struct {
				EventType     string          `json:"event_type"`
				ClientPayload json.RawMessage `json:"client_payload"`
			}

			dispatch.EventType = fmt.Sprintf("%s-test", dependencyName)
			dispatch.ClientPayload = json.RawMessage(payload)

			payloadData, err := json.Marshal(&dispatch)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println(string(payloadData))
			req, err := http.NewRequest("POST", "https://api.github.com/repos/paketo-buildpacks/dep-server/dispatches", bytes.NewBuffer(payloadData))
			if err != nil {
				log.Fatal(err)
			}

			req.Header.Set("Authorization", fmt.Sprintf("token %s", os.Getenv("GITHUB_TOKEN")))

			resp2, err := http.DefaultClient.Do(req)
			if err != nil {
				log.Fatal(err)
			}
			defer resp2.Body.Close()

			if resp2.StatusCode != http.StatusOK && resp2.StatusCode != 204 {
				fmt.Println(resp2.StatusCode)
				log.Fatal(err)
			}

			fmt.Printf("Success version %s!\n", dep.Version)

		} else {
			fmt.Printf("Skipped %s %s because license and CPE are already present\n", dependencyName, dep.Version)
		}
	}

	fmt.Println("Success!")

}
