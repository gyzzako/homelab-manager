package sql

import (
	"database/sql"
	"fmt"
	"homelab-manager/internal/hosts/providers"
	"strings"
)

type SQLProvider struct {
	DataSource string
	Type       string
	Query      string
}

func (p *SQLProvider) GetHostEntries() ([]providers.HostEntry, error) {
	var db *sql.DB
	var err error

	switch p.Type {
	case "sqlite":
		db, err = sql.Open("sqlite3", p.DataSource)
	default:
		return nil, fmt.Errorf("unsupported DB type: %s", p.Type)
	}
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query(p.Query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []providers.HostEntry
	for rows.Next() {
		var ip string
		var fullDomain string
		if err := rows.Scan(&ip, &fullDomain); err != nil {
			return nil, err
		}

		domain, subdomains := parseDomain(fullDomain)

		entry := providers.HostEntry{
			Ip:         ip,
			Domain:     domain,
			Subdomains: subdomains,
		}
		entries = append(entries, entry)
	}

	return entries, nil
}

func parseDomain(fullDomain string) (domain string, subdomains []string) {
	parts := strings.Split(fullDomain, ".")
	if len(parts) < 2 {
		return fullDomain, []string{}
	}

	domain = strings.Join(parts[len(parts)-2:], ".")
	subdomains = parts[:len(parts)-2]

	return domain, subdomains
}
