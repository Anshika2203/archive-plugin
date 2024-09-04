// Copyright 2024 the Drone Authors. All rights reserved.
// Use of this source code is governed by the Blue Oak Model License
// that can be found in the LICENSE file.

package plugin

import (
	"fmt"
	"os"
	"path/filepath"
)

func validatePath(path string) error {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("path does not exist: %s", path)
		}
		return fmt.Errorf("error checking path: %s", err)
	}
	return nil
}

func ensureDirectory(path string) error {
	return os.MkdirAll(path, os.ModePerm)
}

func isDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return fileInfo.IsDir(), nil
}

func getAbsolutePath(path string) (string, error) {
	return filepath.Abs(path)
}
