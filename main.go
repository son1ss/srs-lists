package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/son1ss/srs-lists/parser"
)

type Source struct {
	URL         string `json:"url"`
	ContentType string `json:"contentType"`
	Category    string `json:"category"`
}

func main() {
	var ips []string
	var domains []string

	data, err := os.ReadFile("./source.json")
	if err != nil {
		fmt.Println("Error reading source.json:", err)
		return
	}

	var sources []Source
	err = json.Unmarshal(data, &sources)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	for _, source := range sources {
		data, err := downloadURL(source.URL)
		if err != nil {
			fmt.Printf("Failed to download %s: %v\n", source.URL, err)
			continue
		}

		switch strings.ToLower(source.ContentType) {
		case "csvdumpantizapret":
			ips, domains = parser.ParseCsvDumpAntizapret(string(data))
		case "defaultlist":
			ips, domains = parser.ParseDefaultList(string(data))
		default:
			fmt.Printf("Unsupported content type: %s\n", source.ContentType)
			continue
		}

		// Process parsed data as needed
		fmt.Printf("Category: %s, IPs Count: %d, Domains Count: %d\n", source.Category, len(ips), len(domains))
	}

	ruleSetDomains, ruleSetDomainSuffixes := parser.SeparateDomainsAndSuffixes(domains)

	ruleSet := RuleSet{
		Version: 1,
		Rules: []Rule{
			{
				Domain:       parser.UniqueSlice(ruleSetDomains),
				DomainSuffix: parser.UniqueSlice(ruleSetDomainSuffixes),
				IPCIDR:       parser.UniqueSlice(ips),
			},
		},
	}

	srsError := GenerateSRSList(ruleSet)
	if srsError != nil {
		fmt.Println("Ошибка при генерации SRS списка:", srsError)
	}
}
