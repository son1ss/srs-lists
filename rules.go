package main

type Rule struct {
	Domain       []string `json:"domain,omitempty"`
	DomainSuffix []string `json:"domain_suffix,omitempty"`
	IPCIDR       []string `json:"ip_cidr,omitempty"`
}

// RuleSet структура для представления всего JSON файла
type RuleSet struct {
	Version int    `json:"version"`
	Rules   []Rule `json:"rules"`
}
