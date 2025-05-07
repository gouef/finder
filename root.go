package finder

import "os"

var root string

func GetProjectRoot() string {
	if root == "" {
		path, err := os.Getwd()

		if err != nil {
			return ""
		}

		root = path
	}

	return root
}
