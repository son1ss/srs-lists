package parser

import (
	"fmt"
	"regexp"
	"strings"

	"golang.org/x/text/encoding/charmap"
)

var (
	rgxIPv4   = regexp.MustCompile(`^((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(/(3[0-2]|2[0-9]|1[0-9]|[0-9]))?$`)
	rgxIPv6   = regexp.MustCompile(`^s*((([0-9A-Fa-f]{1,4}:){7}([0-9A-Fa-f]{1,4}|:))|(([0-9A-Fa-f]{1,4}:){6}(:[0-9A-Fa-f]{1,4}|((25[0-5]|2[0-4]d|1dd|[1-9]?d)(.(25[0-5]|2[0-4]d|1dd|[1-9]?d)){3})|:))|(([0-9A-Fa-f]{1,4}:){5}(((:[0-9A-Fa-f]{1,4}){1,2})|:((25[0-5]|2[0-4]d|1dd|[1-9]?d)(.(25[0-5]|2[0-4]d|1dd|[1-9]?d)){3})|:))|(([0-9A-Fa-f]{1,4}:){4}(((:[0-9A-Fa-f]{1,4}){1,3})|((:[0-9A-Fa-f]{1,4})?:((25[0-5]|2[0-4]d|1dd|[1-9]?d)(.(25[0-5]|2[0-4]d|1dd|[1-9]?d)){3}))|:))|(([0-9A-Fa-f]{1,4}:){3}(((:[0-9A-Fa-f]{1,4}){1,4})|((:[0-9A-Fa-f]{1,4}){0,2}:((25[0-5]|2[0-4]d|1dd|[1-9]?d)(.(25[0-5]|2[0-4]d|1dd|[1-9]?d)){3}))|:))|(([0-9A-Fa-f]{1,4}:){2}(((:[0-9A-Fa-f]{1,4}){1,5})|((:[0-9A-Fa-f]{1,4}){0,3}:((25[0-5]|2[0-4]d|1dd|[1-9]?d)(.(25[0-5]|2[0-4]d|1dd|[1-9]?d)){3}))|:))|(([0-9A-Fa-f]{1,4}:){1}(((:[0-9A-Fa-f]{1,4}){1,6})|((:[0-9A-Fa-f]{1,4}){0,4}:((25[0-5]|2[0-4]d|1dd|[1-9]?d)(.(25[0-5]|2[0-4]d|1dd|[1-9]?d)){3}))|:))|(:(((:[0-9A-Fa-f]{1,4}){1,7})|((:[0-9A-Fa-f]{1,4}){0,5}:((25[0-5]|2[0-4]d|1dd|[1-9]?d)(.(25[0-5]|2[0-4]d|1dd|[1-9]?d)){3}))|:)))(%.+)?s*(\/([0-9]|[1-9][0-9]|1[0-1][0-9]|12[0-8]))?$`)
	rgxDomain = regexp.MustCompile(`^(([a-zA-Z0-9А-яёЁ\*]|[a-zA-Z0-9А-яёЁ][a-zA-Z0-9А-яёЁ\-]*[a-zA-Z0-9А-яёЁ])\.)*([A-Za-z0-9А-яёЁ]|[A-Za-z0-9А-яёЁ][A-Za-z0-9А-яёЁ\-]*[A-Za-z0-9А-яёЁ])$`)
)

func ParseCsvDumpAntizapret(input string) ([]string, []string) {
	var ipAddresses []string
	var domains []string

	// Декодируем входную строку из Windows-1251 в UTF-8
	decoder := charmap.Windows1251.NewDecoder()
	decodedInput, _ := decoder.String(input)

	lines := strings.Split(decodedInput, "\n")
	for _, line := range lines {
		// Разделяем строку на столбцы по символу ";"
		columns := strings.Split(line, ";")

		// Пропускаем первую строку (в ней один столбец)
		if len(columns) == 1 {
			continue
		}

		// Извлекаем IP-адреса из первого столбца
		ips := strings.Split(columns[0], "|")
		for _, ip := range ips {
			if ip == "" {
				continue
			}
			if !strings.Contains(ip, "/") {
				ip += "/32"
			}
			ipAddresses = append(ipAddresses, ip)
		}

		// Если есть второй столбца, извлекаем домены из нее
		if len(columns) > 1 {
			domainMatches := strings.Split(columns[1], "|")
			domains = append(domains, domainMatches...)
		}
	}

	// Убираем дубликаты
	ipAddresses = UniqueSlice(ipAddresses)
	domains = UniqueSlice(domains)

	return ipAddresses, domains
}

func ParseDefaultList(input string) ([]string, []string) {
	var ipAddresses []string
	var domains []string

	lines := strings.Split(input, "\n")

	for _, line := range lines {

		// Извлекаем домен
		if rgxDomain.MatchString(line) {
			domains = append(domains, line)
			continue
		}

		// Извлекаем IPv4-адрес
		if rgxIPv4.MatchString(line) {
			ipAddresses = append(ipAddresses, line)
			continue
		}

		// Извлекаем IPv6-адрес
		if rgxIPv6.MatchString(line) {
			ipAddresses = append(ipAddresses, line)
			continue
		}

		// Если строка не была комментарием или пустой строкой, выводим предупреждение, что не удалось распарсить строку
		// (так как такие строки отсекаются регулярками, нет смысла делать проверку перед регулярками)
		if strings.TrimSpace(line) != "" && !strings.HasPrefix(line, "#") {
			fmt.Printf("Failed to parse '%s' as an IPv4, IPv6, or domain address \n", line)
		}
	}

	// Убираем дубликаты
	ipAddresses = UniqueSlice(ipAddresses)
	domains = UniqueSlice(domains)

	return ipAddresses, domains
}
