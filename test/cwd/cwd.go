package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

const dotCharacter = '.'

func isHidden(path string) (bool, error) {
	// dotfiles also count as hidden (if you want)
	if path[0] == dotCharacter {
		return true, nil
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return false, err
	}

	// Appending `\\?\` to the absolute path helps with
	// preventing 'Path Not Specified Error' when accessing
	// long paths and filenames
	// https://docs.microsoft.com/en-us/windows/win32/fileio/maximum-file-path-limitation?tabs=cmd
	pointer, err := syscall.UTF16PtrFromString(`\\?\` + absPath)
	if err != nil {
		return false, err
	}

	attributes, err := syscall.GetFileAttributes(pointer)
	if err != nil {
		return false, err
	}

	return attributes&syscall.FILE_ATTRIBUTE_HIDDEN != 0, nil
}

func main() {
	abs := filepath.Join("C:", "idea", "sdfsdf")

	fmt.Println(abs)

	dir, err := os.UserHomeDir()
	if err != nil {
		return
	}

	fmt.Println(strings.Split(dir, string(os.PathSeparator)))

	//paths := []string{"C:"}
	//path := filepath.Join(paths...)
	path := "C:\\"

	dirs, err := os.ReadDir(path)
	if err != nil {
		panic(err)
	}

	for _, dir := range dirs {
		info, err := dir.Info()
		if err != nil {
			continue
		}

		hidden, err := isHidden(filepath.Join(path, info.Name()))
		if err != nil {
			continue
		}
		fmt.Println(info.IsDir(), info.Name(), hidden)
	}

}
