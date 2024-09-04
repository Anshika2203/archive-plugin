// Copyright 2024 the Drone Authors. All rights reserved.
// Use of this source code is governed by the Blue Oak Model License
// that can be found in the LICENSE file.

package plugin

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type ZipPlugin struct {
	Source   string `envconfig:"PLUGIN_SOURCE"`
	Target   string `envconfig:"PLUGIN_TARGET"`
	LogLevel string `envconfig:"PLUGIN_LOG_LEVEL"`
}

func (p *ZipPlugin) Exec(ctx context.Context) error {
	sourceInfo, err := os.Stat(p.Source)
	if err != nil {
		return fmt.Errorf("error accessing source: %w", err)
	}

	if sourceInfo.IsDir() || !strings.HasSuffix(strings.ToLower(p.Source), ".zip") {
		// Zipping
		return p.Zip()
	} else {
		// Unzipping
		return p.Unzip()
	}
}

func (p *ZipPlugin) Zip() error {
	zipfile, err := os.Create(p.Target)
	if err != nil {
		return err
	}
	defer zipfile.Close()

	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	info, err := os.Stat(p.Source)
	if err != nil {
		return nil
	}

	var baseDir string
	if info.IsDir() {
		baseDir = filepath.Base(p.Source)
	}

	return filepath.Walk(p.Source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		if baseDir != "" {
			header.Name = filepath.Join(baseDir, strings.TrimPrefix(path, p.Source))
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

func (p *ZipPlugin) Unzip() error {
	reader, err := zip.OpenReader(p.Source)
	if err != nil {
		return err
	}
	defer reader.Close()

	for _, file := range reader.File {
		path := filepath.Join(p.Target, file.Name)
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
