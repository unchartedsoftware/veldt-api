package conf

import (
	"fmt"
	"strings"
)

// StaticDir represents a directory to static serve from.
type StaticDir struct {
	Path  string
	Route string
}

// StaticDirs represents a list of public directories to static serve from.
type StaticDirs []StaticDir

// String will return the string representation of the slice.
func (p *StaticDirs) String() string {
	var strs []string
	for _, s := range *p {
		strs = append(strs, fmt.Sprintf("%s -> %s", s.Path, s.Route))
	}
	return strings.Join(strs, "\n")
}

// CheckExisting returns an error if a key/value does not result in a 1:1 mapping.
func (p *StaticDirs) CheckExisting(route string, path string) error {
	for _, sd := range *p {
		if sd.Route == route {
			return fmt.Errorf("Route '%s' has already been registered, duplicate routes are not supported.",
				sd.Route)
		}
		if sd.Path == path {
			return fmt.Errorf("Path '%s' has already been registered under '%s', duplicate pathes are not supported.",
				sd.Path,
				sd.Route)
		}
	}
	return nil
}

// Set will add a directory to the slice.
func (p *StaticDirs) Set(arg string) error {
	args := strings.Split(arg, "=")
	var route string
	var path string
	if len(args) == 1 {
		route = ""
		path = args[0]
	} else {
		route = args[0]
		path = args[1]
	}
	err := p.CheckExisting(route, path)
	if err != nil {
		return err
	}
	*p = append(*p, StaticDir{
		Route: route,
		Path:  path,
	})
	return nil
}
