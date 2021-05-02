package config

type ProjectConfig struct {
	Port  int    `toml:"port"`
	Bench string `toml:"bench"`
}
