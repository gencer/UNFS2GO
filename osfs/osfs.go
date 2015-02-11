package osfs

import (
	"../minfs"
	"errors"
	"os"
	pathpkg "path"
	"path/filepath"
	"time"
)

type osFS struct {
	realpath string //Real path being exported
}

func New(realpath string) (minfs.MinFS, error) {
	_, err := os.Stat(realpath)
	if err == nil {
		return &osFS{realpath}, nil
	}
	return nil, err
}

func (f *osFS) ReadFile(name string, b []byte, off int64) (int, error) {
	realname := f.translate(name)
	fh, err := os.Open(realname)
	if err != nil {
		return -1, err
	}
	defer fh.Close()
	return fh.ReadAt(b, off)
}

func (f *osFS) WriteFile(name string, b []byte, off int64) (int, error) {
	realname := f.translate(name)
	fh, err := os.OpenFile(realname, os.O_RDWR, 0644)
	if err != nil {
		return -1, err
	}
	defer fh.Close()
	return fh.WriteAt(b, off)
}

func (f *osFS) CreateFile(name string) error {
	realname := f.translate(name)
	fil, err := os.Create(realname)
	if err != nil {
		return err
	}
	fil.Close()
	return nil
}

func (f *osFS) CreateDirectory(name string) error {
	realname := f.translate(name)
	return os.Mkdir(realname, 0777)
}

func (f *osFS) Move(oldpath string, newpath string) error {
	orname := f.translate(oldpath)
	nrname := f.translate(newpath)
	return os.Rename(orname, nrname)
}

func (f *osFS) Remove(name string, recursive bool) error {
	realname := f.translate(name)
	if recursive {
		return os.RemoveAll(realname)
	}
	return os.Remove(realname)
}

func (f *osFS) ReadDirectory(name string) ([]os.FileInfo, error) {
	realname := f.translate(name)
	fh, err := os.Open(realname)
	if err != nil {
		return []os.FileInfo{}, err
	}
	defer fh.Close()
	return fh.Readdir(0)
}

func (f *osFS) GetAttribute(path string, attribute string) (interface{}, error) {
	realname := f.translate(path)
		fi, err := os.Stat(realname)
		if err != nil {
		return nil, errors.New("GetAttribute Error Stat'n " + path + "(translated as " + realname + "):" + err.Error())
	}
	switch attribute {
	case "modtime":
		return fi.ModTime(), nil
	case "size":
		return fi.Size(), nil
	}
	return nil, errors.New("GetAttribute Error: Unsupported attribute " + attribute)
}

func (f *osFS) SetAttribute(path string, attribute string, newvalue interface{}) error {
	realname := f.translate(path)
	switch attribute {
	case "modtime":
		return os.Chtimes(realname, time.Now(), newvalue.(time.Time))
	case "size":
		return os.Truncate(realname, newvalue.(int64))
	}
	return errors.New("SetAttribute Error: Unsupported attribute " + attribute)
}

func (f *osFS) Stat(name string) (os.FileInfo, error) {
	realname := f.translate(name)
	return os.Stat(realname)
}

func (f *osFS) String() string { return "os(" + f.realpath + ")" }

func (f *osFS) translate(path string) string {
	path = pathpkg.Clean("/" + path)
	return filepath.Join(f.realpath, path)
}
