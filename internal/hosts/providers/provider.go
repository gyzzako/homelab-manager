package providers

type HostEntry struct {
	Ip         string
	Domain     string
	Subdomains []string
}

type HostProvider interface {
	GetHostEntries() ([]HostEntry, error)
}

type Provider string

const (
	ProviderConfig Provider = "CONFIG"
	ProviderSql    Provider = "SQL"
	ProviderUrl    Provider = "URL"
)
