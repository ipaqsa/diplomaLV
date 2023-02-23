package service

import "agent/pkg"

func isService(name string) bool {
	for _, service := range pkg.Config.Services {
		if name == service {
			return true
		}
	}
	return false
}
