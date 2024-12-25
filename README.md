# TGZ Package

The `tgz` package provides a simple way to create and extract `tar.gz`
archives in Go.

```go
import github.com/woozymasta/tgz

func main() {
  tgz.Pack("source/directory", "archive.tar.gz")
  tgz.Unpack("archive.tar.gz", "output/directory")
}
```

* **Support for Relative Paths**: Archives can include files and directories
  using relative paths, ensuring compatibility with various directory
  structures and simplifying file organization.
* **Customizable Directory Prefix**: Add a prefix to all files and
  directories in the archive, allowing more flexible control over extracted
  folder structures and better organization within archives.
* **Adjustable Compression Level**: Set a custom gzip compression level from
  0 (no compression) to 9 (maximum compression) to balance between archive
  size and compression speed as needed.
* **Cross-Platform Compatibility**: Seamlessly supports both Windows and
  Unix-based systems, automatically handling path separators and ensuring
  consistent performance across environments.
* **Efficient Directory Walking**: Recursively walks through all files and
  directories in the source directory, efficiently gathering content for
  archiving while excluding redundant or hidden paths.
* **Selective File Type Handling**: Supports archiving regular files and
  directories while skipping unsupported file types, ensuring compatibility
  with various filesystem configurations.
* **Simplified Extraction Process**: Extracts all files and directories with
  path integrity, handling any missing directories and preserving original
  file permissions where applicable.

## Installation

Add the package to your Go module:

```bash
go get github.com/woozymasta/tgz
```

Functions:

* `Pack(sourceDir string, targetArchive string) error`  
  Creates a `tar.gz` archive from the contents of the specified `sourceDir`
  and saves it to targetArchive. Uses default compression.
* `PackWithLevel(sourceDir string, targetArchive string, level int) error`  
  Creates a `tar.gz` archive from `sourceDir` and saves it to `targetArchive`,
  with a specified gzip compression `level` (0-9).
* `PackWithPrefix(sourceDir, targetArchive, prefix string, level int) error`  
  Creates a `tar.gz` archive from `sourceDir` and saves it to `targetArchive`,
  with `prefix` added to the archive file paths.
  Allows gzip compression `level` (0-9).
* `Unpack(sourceArchive, targetDir string) error`  
  Extracts a `tar.gz` archive `sourceArchive` into the specified `targetDir`.

## Examples

### Creating archive

```go
err := tgz.Pack("source/directory", "archive.tar.gz")
if err != nil {
  log.Fatal(err)
}
```

### Creating archive with custom compression

```go
// 9 is best compression
// 1 is best speed
// -1 default compression
err := tgz.PackWithLevel("source/directory", "archive.tar.gz", 9)
if err != nil {
  log.Fatal(err)
}
```

[More about compression](https://pkg.go.dev/compress/flate#pkg-constants)

### Creating archive with a path prefix

```go
err := tgz.PackWithPrefix("source/directory", "archive.tar.gz", "some/prefix", -1)
if err != nil {
  log.Fatal(err)
}
```

### Extracting archive

```go
err := tgz.Unpack("archive.tar.gz", "output/directory")
if err != nil {
  log.Fatal(err)
}
```

## Testing

```bash
go test ./...
```

## Other archive packages

* [ZIPp Package](https://github.com/WoozyMasta/zipp) -
  simple way to create and extract .zip archives

