package tgz

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

const (
	purgeTestData = true
	dataDir       = "test_data"
	dataTar       = "x.tar.gz"
)

func TestUnpackAndPackArchiveDot(t *testing.T) {
	err := unpackAndPackArchive("./", "test_unpack_dot")
	if err != nil {
		t.Fatalf("Failed to: %v", err)
	}
}

func TestUnpackAndPackArchivePrefix(t *testing.T) {
	err := unpackAndPackArchive("asd/qwe/", "test_unpack_prefix")
	if err != nil {
		t.Fatalf("Failed to: %v", err)
	}
}

func TestUnpackAndPackArchive(t *testing.T) {
	err := unpackAndPackArchive("", "test_unpack")
	if err != nil {
		t.Fatalf("Failed to: %v", err)
	}
}

func TestPackAndUnpackDirDot(t *testing.T) {
	err := packAndUnpackDir("x", "./", "test_pack_dot")
	if err != nil {
		t.Fatalf("Failed to: %v", err)
	}
}

func TestPackAndUnpackDirPrefix(t *testing.T) {
	err := packAndUnpackDir("x/1", "asd/qwe/", "test_pack_prefix")
	if err != nil {
		t.Fatalf("Failed to: %v", err)
	}
}

func TestPackAndUnpackDir(t *testing.T) {
	err := packAndUnpackDir("x", "", "test_pack")
	if err != nil {
		t.Fatalf("Failed to: %v", err)
	}
}

func unpackAndPackArchive(prefix, testPrefix string) error {
	srcTar := filepath.Join(dataDir, dataTar)
	dstDir := filepath.Join(dataDir, testPrefix)
	tar := dstDir + ".tar.gz"

	err := Unpack(srcTar, dstDir)
	if err != nil {
		return fmt.Errorf("unpack archive: %v", err)
	}
	if purgeTestData {
		defer os.RemoveAll(dstDir)
	}

	err = PackWithPrefix(dstDir, tar, prefix, -1)
	if err != nil {
		return fmt.Errorf("pack archive: %v", err)
	}
	if purgeTestData {
		defer os.Remove(tar)
	}

	return nil
}

func packAndUnpackDir(srcDir, prefix, testPrefix string) error {
	srcDir = filepath.Join(dataDir, srcDir)
	dstDir := filepath.Join(dataDir, testPrefix)
	tar := dstDir + ".tar.gz"

	err := PackWithPrefix(srcDir, tar, prefix, -1)
	if err != nil {
		return fmt.Errorf("pack archive: %v", err)
	}
	if purgeTestData {
		defer os.Remove(tar)
	}

	err = Unpack(tar, dstDir)
	if err != nil {
		return fmt.Errorf("unpack archive: %v", err)
	}
	if purgeTestData {
		defer os.RemoveAll(dstDir)
	}

	return nil
}
