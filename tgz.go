// Package tgz provides functions to create and extract .tar.gz archives.
// Supports custom compression levels, file path prefixes, and relative paths.
// Compatible with cross-platform systems including Windows path structures.

package tgz

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Pack creates a .tar.gz archive from the source directory (sourceDir)
// and saves it to the specified targetArchive path.
func Pack(sourceDir, targetArchive string) error {
	return PackWithPrefix(sourceDir, targetArchive, "", -1)
}

// PackWitLevel creates a .tar.gz archive from the source directory (sourceDir)
// and saves it to the specified targetArchive path, with a specified compression
// level (0-9). Refer to https://pkg.go.dev/compress/flate#pkg-constants for level options.
func PackWitLevel(sourceDir, targetArchive string, level int) error {
	return PackWithPrefix(sourceDir, targetArchive, "", level)
}

// PackWithPrefix creates a .tar.gz archive from the source directory (sourceDir),
// adds a custom prefix to the paths within the archive (such as './' or filepath.Abs(sourceDir)),
// and saves it to targetArchive with the specified gzip compression level (0-9).
func PackWithPrefix(sourceDir, targetArchive, prefix string, level int) error {
	// Check if the source directory exists
	info, err := os.Stat(sourceDir)
	if err != nil {
		return fmt.Errorf("source directory does not exist: %v", err)
	}

	if !info.IsDir() {
		return fmt.Errorf("source is not a directory")
	}

	// Create a file to write the archive
	tarFile, err := os.Create(targetArchive)
	if err != nil {
		return fmt.Errorf("could not create archive file: %v", err)
	}
	defer tarFile.Close()

	// Initialize gzip and tar writers
	gzipWriter, err := gzip.NewWriterLevel(tarFile, level)
	if err != nil {
		return fmt.Errorf("could not create gzip writer: %v", err)
	}
	defer gzipWriter.Close()

	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()

	// Archive each file from the source directory
	err = filepath.Walk(sourceDir, func(filePath string, fileInfo os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Determine relative path for each file in the archive
		relPath, err := filepath.Rel(sourceDir, filePath)
		if err != nil {
			return err
		}

		// Remove any leading "./" from relPath, if present
		relPath = strings.TrimPrefix(relPath, "./")

		// Add specified prefix to the archive paths
		if prefix != "" {
			if strings.HasPrefix(prefix, "./") {
				// If prefix starts with "./", just prepend relPath
				relPath = prefix + relPath
			} else {
				// Otherwise, join using filepath.Join
				relPath = filepath.Join(prefix, relPath)
			}
		}

		// Explicitly exclude unnecessary prefixes
		if relPath == "." || relPath == "./." {
			return nil
		}

		// Convert all paths to use forward slashes
		relPath = filepath.ToSlash(relPath)

		// Obtain file header for the archive entry
		header, err := tar.FileInfoHeader(fileInfo, fileInfo.Name())
		if err != nil {
			return err
		}
		header.Name = relPath

		// Write file header to the archive
		if err := tarWriter.WriteHeader(header); err != nil {
			return err
		}

		// Skip writing content for non-regular files
		if !fileInfo.Mode().IsRegular() {
			return nil
		}

		// Open the file to read its content
		file, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer file.Close()

		// Copy file content to the archive
		if _, err := io.Copy(tarWriter, file); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

// Unpack extracts a .tar.gz archive (sourceArchive) into the target directory (targetDir).
func Unpack(sourceArchive, targetDir string) error {
	// Open archive file for reading
	file, err := os.Open(sourceArchive)
	if err != nil {
		return fmt.Errorf("could not open archive file: %v", err)
	}
	defer file.Close()

	// Initialize gzip and tar readers
	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		return fmt.Errorf("could not create gzip reader: %v", err)
	}
	defer gzipReader.Close()

	tarReader := tar.NewReader(gzipReader)

	// Iterate over each file in the archive
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break // End of archive
		}
		if err != nil {
			return err
		}

		// Determine the file path for extraction
		filename := header.Name

		// Remove any leading "./" from path if present
		filename = strings.TrimPrefix(filename, "./")
		targetPath := filepath.Join(targetDir, filename)

		// Handle extraction based on file type
		switch header.Typeflag {
		case tar.TypeDir:
			// Create directory
			if err := os.MkdirAll(targetPath, os.FileMode(header.Mode)); err != nil {
				return err
			}

		case tar.TypeReg:
			// Create all necessary directories
			if err := os.MkdirAll(filepath.Dir(targetPath), os.FileMode(header.Mode)); err != nil {
				return err
			}
			// Create file
			outFile, err := os.Create(targetPath)
			if err != nil {
				return err
			}
			// Copy file content
			if _, err := io.Copy(outFile, tarReader); err != nil {
				outFile.Close()
				return err
			}
			outFile.Close()

		default:
			// Skip other file types
		}
	}

	return nil
}
