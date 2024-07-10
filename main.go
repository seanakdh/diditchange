package diditchange

import (
	"fmt"
	"os"
	"time"
)

const (
	SizeChanged TypeChanged = iota
	ModTimeChanged
)

type TypeChanged int
type FileChangedInfo struct {
	fileName    string
	typeChanged TypeChanged
}

// Returns file name of changed file
func (fci *FileChangedInfo) FileName() string {
	return fci.fileName
}

// Returns type of change occured change
func (fci *FileChangedInfo) TypeChanged() TypeChanged {
	return fci.typeChanged
}

// Watches specified 'file' with specified 'interval', if changed sends a FileChangedInfo
// on 'respch'. Returns an error, if it can't read the file sats
func WatchFileAsync(file string, respch chan<- FileChangedInfo, interval int) error {

	for {
		fci, err := WatchFileSync(file, interval)
		if err != nil {
			return err
		}
		respch <- fci
	}
}

// Watches specified 'file' with specified 'interval', if changed returns a FileChangedInfo.
// Returns empty FileChangedInfo{} and an error, if it can't read the file sats of 'file'
func WatchFileSync(file string, interval int) (FileChangedInfo, error) {
	initialStat, err := os.Stat(file)
	if err != nil {
		return FileChangedInfo{}, err
	}
	for {
		stat, err := os.Stat(file)
		if err != nil {
			return FileChangedInfo{}, err
		}
		if stat.Size() != initialStat.Size() {
			return FileChangedInfo{file, SizeChanged}, nil
		} else if stat.ModTime() != initialStat.ModTime() {
			return FileChangedInfo{file, ModTimeChanged}, nil
		} else {
			time.Sleep(time.Duration(interval) * time.Millisecond)
		}
	}
}

// GetDirFiles returns all files, relative to 'path', in given directory 'path'.
// If 'depth' is greater than 1 or smaller than 0, GetDirFiles gets called recursivley,
// if another dir is encouterd. Returns an nil value and an error if
// at any point a directory can't be read.
func GetDirFiles(path string, depth int) (files []string, err error) {
	dirContent, err := os.ReadDir(path)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	for _, file := range dirContent {
		if !file.IsDir() {
			files = append(files, path+file.Name())
		} else if depth > 1 || depth < 0 {
			if path[len(path)-1] != '/' {
				path += "/"
			}
			subDirFiles, err := GetDirFiles(path+file.Name()+"/", depth-1)
			if err != nil {
				return nil, err
			}
			files = append(files, subDirFiles...)
		}
	}
	return files, nil
}
