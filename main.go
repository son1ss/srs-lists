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
	var antizapretIps []string
	var antizapretDomains []string

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
			newIps, newDomains := parser.ParseCsvDumpAntizapret(string(data))
			antizapretIps = append(antizapretIps, newIps...)
			antizapretDomains = append(antizapretDomains, newDomains...)
		case "defaultlist":
			newIps, newDomains := parser.ParseDefaultList(string(data))
			ips = append(ips, newIps...)
			domains = append(domains, newDomains...)
		case "list":
			newIps := parser.ParseList(string(data))

			ips = append(ips, newIps...)
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
				Domain:       ruleSetDomains,
				DomainSuffix: ruleSetDomainSuffixes,
				IPCIDR:       ips,
			},
		},
	}

	antizapretRuleSetDomains, antizapretRuleSetDomainSuffixes := parser.SeparateDomainsAndSuffixes(antizapretDomains)

	antizapretRuleSet := RuleSet{
		Version: 1,
		Rules: []Rule{
			{
				Domain:       antizapretRuleSetDomains,
				DomainSuffix: antizapretRuleSetDomainSuffixes,
				IPCIDR:       antizapretIps,
			},
		},
	}

	srsError := GenerateSRSList(ruleSet, "af")
	antizapretSrsError := GenerateSRSList(antizapretRuleSet, "az")

	if srsError != nil {
		fmt.Println("Ошибка при генерации SRS списка:", srsError)
	}

	if antizapretSrsError != nil {
		fmt.Println("Ошибка при генерации SRS списка для Antizapret:", antizapretSrsError)
	}
}
