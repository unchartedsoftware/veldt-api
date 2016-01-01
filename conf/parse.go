package conf

import (
	"flag"
)

// ParseCommandLine parses the commandline arguments and returns a Conf object.
func ParseCommandLine() *Conf {
	// port the server runs on
	port := flag.String("port", "8080", "Port to bind HTTP server")
	// public dir that will be statically served
	public := flag.String("public", "./build/public", "The public directory to static serve from")
	// alias flags
	aliases := make(AliasMap)
	flag.Var(&aliases, "alias", "Alias a string value by another string.")

	// Parse the commandline into the above flags
	flag.Parse()

	// Set, save, and return config
	config := &Conf{
		Port:       *port,
		Public:     *public,
		Aliases:    aliases,
		InvAliases: invertAliases(aliases),
	}
	SaveConf(config)
	return config
}
