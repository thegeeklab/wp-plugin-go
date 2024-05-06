package util

import "os/user"

// GetUserHomeDir returns the home directory path for the current user.
// If the current user cannot be determined, it returns the default "/root" path.
func GetUserHomeDir() string {
	home := "/root"

	if currentUser, err := user.Current(); err == nil {
		home = currentUser.HomeDir
	}

	return home
}
