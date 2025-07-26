package hosts

type Config struct {
	Ip         string   `yaml:"ip"`
	Domain     string   `yaml:"domain"`
	Subdomains []string `yaml:"subdomains"`
}
