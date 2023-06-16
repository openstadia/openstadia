package config

type App struct {
	Name    string   `yaml:"name" json:"name"`
	Command []string `yaml:"command" json:"command"`
	Width   uint     `yaml:"width" json:"width"`
	Height  uint     `yaml:"height" json:"height"`
}
