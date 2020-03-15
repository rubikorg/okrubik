package config

type ProjectConfig struct {
	Host  string `toml:"asset_host"`
	Port  string `toml:"port"`
	Bench string `toml:"bench"`
}
