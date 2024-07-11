package watchit

import (
	"io/fs"
	"os"
	"time"
)

const (
	TypeSizeChanged TypeChanged = iota
	TypeModTimeChanged
)

type TypeChanged int

type WatchItInfo struct {
	fileName    string
	typeChanged TypeChanged
	err         error
}

// Returns file name of changed file
func (fci *WatchItInfo) FileName() string {
	return fci.fileName
}

// Returns type of occured change
func (fci *WatchItInfo) TypeChanged() TypeChanged {
	return fci.typeChanged
}

// returns contained error, will be nil if no error occured
func (fci *WatchItInfo) GetError() error {
	return fci.err
}

// Watches specified 'file' with specified 'interval' by calling WatchFileSync in a go routine,
// which contains a loop if changed sends a FileChangedInfo on 'respch', if os.Stat fails wraps
// the error in WatchItInfo and stops goroutine
func WatchFileAsync(file string, respch chan<- WatchItInfo, interval int) {
	go func() {
		for {
			fci := WatchFileSync(file, interval)
			respch <- fci
			if fci.err != nil {
				break
			}
		}
	}()
}

// Ranges over `files` and calls WatchFileAsync for every file
func WatchMultipleFilesAsync(files []string, respch chan<- WatchItInfo, interval int) {
	for _, file := range files {
		WatchFileAsync(file, respch, interval)
	}
}

// Watches specified 'file' with specified 'interval', if changed returns a WatchItInfo.
// If os.Stat fails, returns an error wrapped in WatchItInfo return value
func WatchFileSync(file string, interval int) WatchItInfo {
	initialStat, err := os.Stat(file)
	if err != nil {
		return WatchItInfo{file, -1, err}
	}

	for {
		stat, err := os.Stat(file)
		if err != nil {
			return WatchItInfo{file, -1, err}
		}
		if stat.Size() != initialStat.Size() {
			return WatchItInfo{file, TypeSizeChanged, nil}
		} else if stat.ModTime() != initialStat.ModTime() {
			return WatchItInfo{file, TypeModTimeChanged, nil}
		} else {
			time.Sleep(time.Duration(interval) * time.Millisecond)
		}
	}
}

// GetDirFiles returns all files, relative to 'path', in given directory 'path'.
// If 'depth' is greater than 1 or smaller than 0, GetDirFiles gets called recursivley,
// if another dir is encouterd. 'depth' smaller than zero is effectivly infinite recursivity
// Returns an nil value and an error if at any point a directory can't be read.
// Ignores all FileModes except Dir
func GetDirFiles(path string, depth int) (files []string, err error) {

	if file, err := os.Stat(path); err == nil && file.IsDir() {
		if path[len(path)-1] != '/' {
			path += "/"
		}
	} else if err != nil {
		return nil, err
	}

	dirContent, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	for _, file := range dirContent {
		// Check if File has any other mode than ModeDir, if so skip file processing
		if file.Type()&(fs.ModeType^fs.ModeDir) > 0 {
			continue
		}
		if !file.IsDir() {
			files = append(files, path+file.Name())
		} else if depth > 1 || depth < 0 {
			// fmt.Println(file)
			subDirFiles, err := GetDirFiles(path+file.Name()+"/", depth-1)
			if err != nil {
				return nil, err
			}
			files = append(files, subDirFiles...)
		}
	}
	return files, nil
}
