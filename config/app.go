package config

type BaseApp struct {
	Name    string   `yaml:"name" json:"name"`
	Command []string `yaml:"command" json:"command"`
	Width   uint     `yaml:"width" json:"width"`
	Height  uint     `yaml:"height" json:"height"`
}

type DbApp struct {
	Id int `yaml:"id" json:"id"`
	BaseApp
}
