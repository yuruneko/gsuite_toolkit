package models

type TomlConfig struct {
	Owner DomainOwner
	Networks map[string][]Network
}

type DomainOwner struct {
	DomainName string
	Organization string
}

type Network struct {
	Name string
	Ip []string
}
