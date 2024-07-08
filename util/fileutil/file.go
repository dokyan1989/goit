package fileutil

import (
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func ScanRecursive(dir_path string, includes []string, ignores []string) []string {
	files := []string{}

	err := filepath.Walk(dir_path, func(path string, f os.FileInfo, ferr error) error {
		skip := true

		for _, i := range includes {
			if strings.Contains(path, i) {
				skip = false
				break
			}
		}

		if skip {
			return nil
		}

		// Loop : Ignore Files & Folders
		for _, i := range ignores {
			if strings.Contains(path, i) {
				// skip = true
				return nil
			}
		}

		fpath, err := os.Stat(path)
		if err != nil {
			log.Fatal(err)
		}

		// File & Folder Mode
		fMode := fpath.Mode()

		if fMode.IsRegular() {
			// Append to Files Array
			files = append(files, path)
		}

		return nil
	})
	if err != nil {
		log.Println(err)
	}
	return files
}

// SameFileError is error about samefile.
type SameFileError struct {
	Src string
	Dst string
}

func (e SameFileError) Error() string {
	return fmt.Sprintf("%s and %s are the same file", e.Src, e.Dst)
}

// SpecialFileError is an error about special file.
type SpecialFileError struct {
	File     string
	FileInfo os.FileInfo
}

func (e SpecialFileError) Error() string {
	return fmt.Sprintf("`%s` is a named pipe", e.File)
}

// NotADirectoryError is not a directory.
type NotADirectoryError struct {
	Src string
}

func (e NotADirectoryError) Error() string {
	return fmt.Sprintf("`%s` is not a directory", e.Src)
}

// AlreadyExistsError folder/file.
type AlreadyExistsError struct {
	Dst string
}

func (e AlreadyExistsError) Error() string {
	return fmt.Sprintf("`%s` already exists", e.Dst)
}

func samefile(src string, dst string) bool {
	//nolint: errcheck
	srcInfo, _ := os.Stat(src)
	//nolint: errcheck
	dstInfo, _ := os.Stat(dst)
	return os.SameFile(srcInfo, dstInfo)
}

func specialfile(fi os.FileInfo) bool {
	return (fi.Mode() & os.ModeNamedPipe) == os.ModeNamedPipe
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// CopyFile data from src to dst
func CopyFile(src, dst string) error {
	if samefile(src, dst) {
		return &SameFileError{src, dst}
	}

	// Make sure src exists and neither are special files
	srcStat, err := os.Lstat(src)
	if err != nil {
		return err
	}
	if specialfile(srcStat) {
		return &SpecialFileError{src, srcStat}
	}

	dstStat, err := os.Stat(dst)
	if err != nil && !os.IsNotExist(err) {
		return err
	} else if err == nil {
		if specialfile(dstStat) {
			return &SpecialFileError{dst, dstStat}
		}
	}

	// Do the actual copy
	fsrc, err := os.Open(src)
	if err != nil {
		return err
	}
	defer fsrc.Close()

	fdst, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer fdst.Close()

	size, err := io.Copy(fdst, fsrc)
	if err != nil {
		return err
	}

	if size != srcStat.Size() {
		return fmt.Errorf("%s: %d/%d copied", src, size, srcStat.Size())
	}

	return nil
}

// CopyMode copy mode from src to dst.
func CopyMode(src, dst string) error {
	_, err := os.Lstat(src)
	if err != nil {
		return err
	}

	_, err = os.Lstat(dst)
	if err != nil {
		return err
	}

	// Atleast one is not a symlink, get the actual file stats
	srcStat, err := os.Stat(src)
	if err != nil {
		return err
	}
	err = os.Chmod(dst, srcStat.Mode())
	return err
}

// Copy data and mode bits ("cp src dst").
// Return the file's destination.
func Copy(src, dst string) (string, error) {
	dstInfo, err := os.Stat(dst)

	if err == nil && dstInfo.Mode().IsDir() {
		dst = filepath.Join(dst, filepath.Base(src))
	}

	if err != nil && !os.IsNotExist(err) {
		return dst, err
	}

	err = CopyFile(src, dst)
	if err != nil {
		return dst, err
	}

	err = CopyMode(src, dst)
	if err != nil {
		return dst, err
	}

	return dst, nil
}

// CopyTreeOptions options.
type CopyTreeOptions struct {
	Ignore func(string, []fs.DirEntry) []string
}

// CopyTree is recursively copy a directory tree.
func CopyTree(src, dst string, options *CopyTreeOptions) error {
	if options == nil {
		options = &CopyTreeOptions{Ignore: nil}
	}

	// Check src must be exist and is a directory
	srcStat, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !srcStat.IsDir() {
		return &NotADirectoryError{src}
	}
	// Check dest must be exist and is a directory
	dstStat, err := os.Stat(dst)
	if err != nil {
		return err
	}
	if !dstStat.IsDir() {
		return &NotADirectoryError{dst}
	}
	// Read all src
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}
	ignoredNames := []string{}
	if options.Ignore != nil {
		ignoredNames = options.Ignore(src, entries)
	}

	for _, entry := range entries {
		if stringInSlice(entry.Name(), ignoredNames) {
			continue
		}
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		entryFileInfo, err := os.Lstat(srcPath)
		if err != nil {
			return err
		}

		if entryFileInfo.IsDir() {
			err = CopyTree(srcPath, dstPath, options)
			if err != nil {
				return err
			}
		} else {
			_, err = Copy(srcPath, dstPath)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
