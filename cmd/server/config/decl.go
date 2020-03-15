package config

type ProjectConfig struct {
	Port  string `toml:"port"`
	Bench string `toml:"bench"`
}
