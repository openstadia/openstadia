package config

type Local struct {
	Enabled bool   `yaml:"enabled" json:"enabled"`
	Host    string `yaml:"host" json:"host"`
	Port    string `yaml:"port" json:"port"`
}
