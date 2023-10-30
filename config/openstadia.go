package config

type Openstadia struct {
	Hub   *Hub      `yaml:"hub" json:"hub"`
	Local *Local    `yaml:"local" json:"local"`
	Apps  []BaseApp `yaml:"apps" json:"apps"`
}

type DbOpenstadia struct {
	Openstadia
	Apps []DbApp `yaml:"apps" json:"apps"`
}
