package cmd

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

func Zip(dir string) (string, error) {
	output, err := outputName(dir)

	if err != nil {
		return "", err
	}

	newZipFile, err := os.Create(output)
	defer newZipFile.Close()

	if err != nil {
		return "", err
	}

	zipWriter := zip.NewWriter(newZipFile)
	defer zipWriter.Close()

	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		defer file.Close()

		if err != nil {
			return err
		}

		info, err = file.Stat()
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)

		if err != nil {
			return err
		}

		header.Name = path
		header.Method = zip.Deflate

		writer, err := zipWriter.CreateHeader(header)

		if err != nil {
			return err
		}

		if _, err := io.Copy(writer, file); err != nil {
			return err
		}

		return nil
	})

	return output, err

}

func outputName(path string) (string, error) {
	file, err := os.Open(path)

	defer file.Close()

	if err != nil {

		return "", err
	}

	info, err := file.Stat()

	if err != nil {
		return "", err
	}

	return info.Name(), nil
}
