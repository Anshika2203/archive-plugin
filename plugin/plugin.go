// Copyright 2024 the Drone Authors. All rights reserved.
// Use of this source code is governed by the Blue Oak Model License
// that can be found in the LICENSE file.

package plugin

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type ZipPlugin struct{}

func (p *ZipPlugin) Zip(source, target string) error {
	zipfile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer zipfile.Close()

	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	info, err := os.Stat(source)
	if err != nil {
		return nil
	}

	var baseDir string
	if info.IsDir() {
		baseDir = filepath.Base(source)
	}

	return filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		if baseDir != "" {
			header.Name = filepath.Join(baseDir, strings.TrimPrefix(path, source))
		}

		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(writer, file)
		return err
	})
}

func (p *ZipPlugin) Unzip(source, target string) error {
	reader, err := zip.OpenReader(source)
	if err != nil {
		return err
	}
	defer reader.Close()

	for _, file := range reader.File {
		path := filepath.Join(target, file.Name)
		if file.FileInfo().IsDir() {
			os.MkdirAll(path, file.Mode())
			continue
		}

		fileReader, err := file.Open()
		if err != nil {
			return err
		}
		defer fileReader.Close()

		targetFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}
		defer targetFile.Close()

		_, err = io.Copy(targetFile, fileReader)
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	source := os.Getenv("SOURCE")
	target := os.Getenv("TARGET")

	if source == "" || target == "" {
		fmt.Println("Both SOURCE and TARGET environment variables must be set")
		os.Exit(1)
	}

	plugin := &ZipPlugin{}

	// Determine if we're zipping or unzipping based on the source file
	sourceInfo, err := os.Stat(source)
	if err != nil {
		fmt.Printf("Error accessing source: %v\n", err)
		os.Exit(1)
	}

	if sourceInfo.IsDir() || !strings.HasSuffix(strings.ToLower(source), ".zip") {
		// Zipping
		err = plugin.Zip(source, target)
		if err != nil {
			fmt.Printf("Error zipping: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Zip operation completed successfully")
	} else {
		// Unzipping
		err = plugin.Unzip(source, target)
		if err != nil {
			fmt.Printf("Error unzipping: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Unzip operation completed successfully")
	}
}
