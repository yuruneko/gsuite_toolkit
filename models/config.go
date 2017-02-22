package models

type TomlConfig struct {
	Owner DomainOwner
	Network map[string]Network
}

type DomainOwner struct {
	DomainName string
	Organization string
}

type Network struct {
	Name string
	Ips []string
}
