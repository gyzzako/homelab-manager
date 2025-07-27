package hosts

import (
	"bufio"
	"fmt"
	"homelab-manager/internal"

	"homelab-manager/internal/hosts/providers"
	"os"
	"runtime"
	"strings"
)

const marker = "#by-" + internal.APP_NAME

func UpdateHosts(provider providers.HostProvider) ([]providers.HostEntry, error) {
	hostsPath := getHostsFilePath()
	file, err := os.Open(hostsPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		if strings.Contains(trimmed, marker) {
			continue
		}

		lines = append(lines, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	hostEntries, err := provider.GetHostEntries()
	if err != nil {
		return nil, err
	}

	lines = append(lines, getNewLines(hostEntries)...)

	output := strings.Join(lines, "\n") + "\n"
	return hostEntries, os.WriteFile(hostsPath, []byte(output), 0644)
}

func getHostsFilePath() string {
	if runtime.GOOS == "windows" {
		return "C:\\Windows\\System32\\drivers\\etc\\hosts"
	}
	return "/etc/hosts"
}

func getNewLines(hostEntries []providers.HostEntry) []string {
	var lines []string

	for _, entry := range hostEntries {

		// Add main domain
		if len(entry.Subdomains) == 0 {
			line := fmt.Sprintf("%s\t%s\t%s", entry.Ip, entry.Domain, marker)
			lines = append(lines, line)
		}

		// Add subdomains
		for _, sub := range entry.Subdomains {
			fqdn := fmt.Sprintf("%s.%s", sub, entry.Domain)
			line := fmt.Sprintf("%s\t%s\t%s", entry.Ip, fqdn, marker)
			lines = append(lines, line)
		}
	}

	return lines
}
