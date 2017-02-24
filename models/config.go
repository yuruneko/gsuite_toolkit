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
	Type string
	Ip []string
}

func (config *TomlConfig) GetAllIps() []string {
	var allIp []string
	for _, network := range config.Networks {
		for _, n := range network {
			allIp = append(allIp, n.Ip...)
		}
	}
	return allIp
}