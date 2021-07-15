package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/paketo-buildpacks/dep-server/pkg/dependency/licenses"
)

// USAGE
// `go run main.go --name <dependency-name>`
// This will update the license and CPE fields for each dependency if they are not already present.

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
	var (
		dependencyName string
	)

	// Takes in the name of 1 dep => dispatches to the test-upload workflow with all metadata
	flag.StringVar(&dependencyName, "name", "", "Dependency name")
	flag.Parse()
	if dependencyName == "" {
		fmt.Println("`name` s required")
		os.Exit(1)
	}

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
		// don't touch anything if the metadata is complete
		if len(dep.Licenses) == 0 || dep.CPE == "" {
			fmt.Println(dep.Version)

			// pass the dep name and source URL and whatever else to pkg/dependency/licenses to get licenses
			if len(dep.Licenses) == 0 {
				licenses, err := licenseRetriever.LookupLicenses(dependencyName, dep.Source)
				if err != nil {
					log.Fatal(err)
				}

				dep.Licenses = licenses
			}

			if dep.CPE == "" {
				dep.CPE = GetCPE(dependencyName, dep.Version)
			}

			// dispatchDep is an exact copy of the dep, but the licenses are a string instead of slice.
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

func GetCPE(depName, version string) string {
	switch depName {
	case "bundler":
		return fmt.Sprintf("cpe:2.3:a:bundler:bundler:%s:*:*:*:*:ruby:*:*", version)
	case "composer":
		return ""
	case "curl":
		return fmt.Sprintf("cpe:2.3:a:haxx:curl:%s:*:*:*:*:*:*:*", version)
	case "dotnet-aspnetcore":
		dotnetAspnetCoreVersion := strings.Join(strings.Split(version, ".")[0:2], ".")
		return fmt.Sprintf("cpe:2.3:a:microsoft:asp.net_core:%s:*:*:*:*:*:*:*", dotnetAspnetCoreVersion)

	case "dotnet-runtime":
		dotnetCoreRuntimeVersion, err := semver.NewVersion(version)
		if err != nil {
			log.Fatal(err)
		}
		productName := ".net"
		if dotnetCoreRuntimeVersion.LessThan(semver.MustParse("5.0.0-0")) { // use 5.0.0-0 to ensure 5.0.0 previews/RCs use the new `.net` product name
			productName = ".net_core"
		}
		return fmt.Sprintf("cpe:2.3:a:microsoft:%s:%s:*:*:*:*:*:*:*", productName, version)

	case "dotnet-sdk":
		dotnetCoreSDKVersion, err := semver.NewVersion(version)
		if err != nil {
			log.Fatal(err)
		}
		productName := ".net"
		if dotnetCoreSDKVersion.LessThan(semver.MustParse("5.0.0-0")) { // use 5.0.0-0 to ensure 5.0.0 previews/RCs use the new `.net` product name
			productName = ".net_core"
		}
		return fmt.Sprintf("cpe:2.3:a:microsoft:%s:%s:*:*:*:*:*:*:*", productName, version)
	case "go":
		return fmt.Sprintf("cpe:2.3:a:golang:go:%s:*:*:*:*:*:*:*", strings.TrimPrefix(version, "go"))
	case "httpd":
		return fmt.Sprintf("cpe:2.3:a:apache:http_server:%s:*:*:*:*:*:*:*", version)
	case "icu":
		return fmt.Sprintf(`cpe:2.3:a:icu-project:international_components_for_unicode:%s:*:*:*:*:c\/c\+\+:*:*`, version)
	case "nginx":
		return fmt.Sprintf("cpe:2.3:a:nginx:nginx:%s:*:*:*:*:*:*:*", version)
	case "node":
		return fmt.Sprintf("cpe:2.3:a:nodejs:node.js:%s:*:*:*:*:*:*:*", strings.TrimPrefix(version, "v"))
	case "php":
		return fmt.Sprintf("cpe:2.3:a:php:php:%s:*:*:*:*:*:*:*", version)
	case "pip":
		return fmt.Sprintf("cpe:2.3:a:pypa:pip:%s:*:*:*:*:python:*:*", version)
	case "pipenv":
		return ""
	case "python":
		fmt.Sprintf("cpe:2.3:a:python:python:%s:*:*:*:*:*:*:*", version)
	case "ruby":
		return fmt.Sprintf("cpe:2.3:a:ruby-lang:ruby:%s:*:*:*:*:*:*:*", version)
	case "rust":
		return fmt.Sprintf("cpe:2.3:a:rust-lang:rust:%s:*:*:*:*:*:*:*", version)
	case "tini":
		return fmt.Sprintf("cpe:2.3:a:tini_project:tini:%s:*:*:*:*:*:*:*", strings.TrimPrefix(version, "v"))
	case "yarn":
		return fmt.Sprintf("cpe:2.3:a:yarnpkg:yarn:%s:*:*:*:*:*:*:*", version)
	}
	return ""
}
