package parser

import (
	"net"
	"strings"
)

func UniqueSlice(slice []string) []string {
	uniqueMap := make(map[string]bool)
	uniqueSlice := make([]string, 0)

	for _, item := range slice {
		if _, found := uniqueMap[item]; !found {
			uniqueMap[item] = true
			uniqueSlice = append(uniqueSlice, item)
		}
	}

	return uniqueSlice
}

func SeparateDomainsAndSuffixes(domains []string) ([]string, []string) {
	regularDomains := make([]string, 0)
	wildcardDomains := make([]string, 0)

	for _, domain := range domains {
		if len(domain) == 0 {
			continue
		}

		if strings.HasPrefix(domain, "*") {
			wildcardDomains = append(wildcardDomains, strings.Replace(domain, "*", "", 1))
		} else {
			regularDomains = append(regularDomains, domain)
		}
	}

	return regularDomains, wildcardDomains
}

func GetIPNetwork(ip net.IP) net.IPNet {
	var mask net.IPMask
	var network net.IPNet

	// Если IPv4
	if ip.To4() != nil {
		mask = net.CIDRMask(32, 32) // уменьшаем маску на 1 бит для хоста
		network = net.IPNet{IP: ip, Mask: mask}
		return network
	}

	// Если IPv6
	mask = net.CIDRMask(128, 128) // уменьшаем маску на 1 бит для хоста
	network = net.IPNet{IP: ip, Mask: mask}
	return network
}
