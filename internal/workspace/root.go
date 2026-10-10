package workspace

import "os"

func Root() (string, error) {
	current, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return Locate(current)
}
