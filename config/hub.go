package config

type Hub struct {
	Enabled bool   `yaml:"enabled" json:"enabled"`
	Addr    string `yaml:"addr" json:"addr"`
	Token   string `yaml:"token" json:"token"`
}
