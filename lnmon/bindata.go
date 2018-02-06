// Code generated by go-bindata.
// sources:
// lnmon.tmpl
// DO NOT EDIT!

package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _lnmonTmpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x9c\x95\x4d\x4f\x23\x39\x10\x86\xcf\xf4\xaf\xa8\xe5\x4e\x9b\x8f\x13\xc8\xeb\x15\x20\xad\x16\x69\x97\x45\x4b\xd8\xd1\x1c\x9d\xb8\x12\x7b\xe2\x76\x45\x76\xf5\x00\x6a\xe5\xbf\x8f\xec\x74\xbe\x48\x02\xcc\x5c\xa2\xd8\x2e\x3f\xf5\xbe\x55\x6d\x5b\x5a\x6e\xbc\xaa\xe4\x90\xcc\xab\xaa\xa4\x3d\x53\x7f\xa1\xf7\xf4\x9b\x14\xf6\x4c\x55\x72\xa6\x06\xd6\x25\x70\x09\x74\x00\x7c\x99\x61\x74\x0d\x06\x86\xd8\x86\xe0\xc2\x04\x34\x83\x0f\xb5\x9d\x7e\x0b\x75\x83\xf0\xec\xd8\xc2\xdf\x6e\x62\xb9\x2c\xde\x23\x3f\x53\x9c\x82\x0e\x06\x6e\x1c\x8f\xc8\x05\x18\xbe\x82\xd4\x60\x23\x8e\x7f\x3f\xb6\xcc\xb3\x74\x25\x44\xbf\xfd\x58\xe5\x3f\x52\x68\x55\x4b\x31\x53\x55\xce\xfe\x95\x5a\x68\x32\x10\x86\x08\x2e\x30\x46\x4c\x8c\x06\x5c\x00\xb6\xb8\x8d\xba\x12\xa2\xa1\x95\x98\xab\xcb\xd3\xcb\x53\x31\xa2\x90\xc8\x63\x12\x2e\x18\x7c\xa9\xb3\xd9\x63\xf5\x10\xa9\x41\xb6\xd8\x26\x30\x3a\xd9\x21\xe9\x68\x72\x5a\x48\xa3\xa8\x67\x59\xb9\xd1\xac\xab\x71\xa4\x06\x38\xdb\x67\x22\xbf\xd4\x64\xcf\x95\x1c\x91\x41\x35\x5c\x38\x32\x52\x94\x21\xb8\x30\x26\x29\xec\xb9\xaa\xba\x0e\xdc\x18\xea\xde\xb2\xa9\xef\xd2\x7f\x7d\xb9\x4e\xe6\xf3\xec\xea\x00\x20\x81\x4c\x1c\x29\x4c\x54\x5f\x5e\x29\xfa\xf1\x22\x79\xd7\x01\xfa\x84\xf0\x49\x48\xa0\x55\x9f\x76\x40\x27\x80\xc1\x64\xd0\x86\x21\xbf\x6c\xdc\x61\x4b\xab\xde\x1e\x34\xb5\x07\x72\xd8\x56\x11\x23\x67\xea\xdf\x36\x82\x33\x25\xb2\xec\xe9\xba\xed\x54\x61\x4c\xf5\x3d\x19\xbc\x33\xf3\x79\x8f\x5d\x55\xa4\x7c\x73\x9b\xd1\xd7\xde\xe9\xb4\x14\x95\xc9\xba\x4c\x6c\xc1\x61\x97\xb3\x51\xd9\xa7\x30\x0d\xf4\x1c\xfa\x8d\x63\x8a\x40\x6d\x84\x3c\x13\xc8\xe0\x4e\x09\x97\x69\x8c\x89\x98\x4a\xa2\x77\x8d\x5c\xf7\x71\x6f\x24\xf4\x94\xef\x18\x93\xa3\xf0\x21\xe5\xff\x3e\x6e\x3f\x65\xe8\x69\x34\xb5\x58\x0e\xce\x47\xa4\x9b\x8d\xd8\x5d\xda\x17\x84\x5c\x0c\xa0\x31\x74\x1d\x78\x0c\x5b\x80\xdc\x94\xec\xa4\x14\x26\x95\x53\x49\xac\xfd\x27\x77\xdf\x5a\x1d\x02\xfa\x02\x18\x2d\xff\xff\x24\xa3\x28\xa8\x1f\x10\xe3\x86\x0e\xb6\x9a\x41\x47\x2c\x7d\x9b\xe5\xb5\xd5\xe9\xbd\x50\x25\x56\x0a\x7b\x91\xc7\xac\x87\x1e\x55\x75\x24\x39\xaa\xea\xe8\x48\xb2\x55\xf7\xba\x41\x29\xd8\x2e\xc7\xb7\x14\x02\x8e\x18\xcd\x1f\x5b\xb3\x0b\xbd\x90\x58\xf3\x56\xf8\x9f\x6d\x30\xf9\xc6\x80\x36\xed\x9d\x66\x8b\x4d\xbf\x20\x45\xc9\xda\x75\x10\x75\x98\xe0\x41\x5f\xf9\x5b\x5e\x2b\x34\xa5\x8b\x59\x65\x69\x17\x9b\xcd\xe9\x95\xd8\x37\x6b\xeb\xf6\x2f\x8b\x5e\x3f\x66\xe5\xeb\x8e\xbf\x21\xe5\xd3\xbe\xd9\x9f\xad\xad\xff\x38\xef\xdd\xa3\x66\x4a\xd6\x0d\xe8\x29\xd5\xd7\xe9\x66\x70\xdb\x87\x2d\x4e\xc5\xaf\xe3\x06\x16\x9b\xc3\xc0\x55\xcd\xd6\xc7\x4f\xf4\x4d\xac\x64\x9b\x5f\x32\x91\x7f\xf7\xde\x95\xef\xde\x4d\x9f\xbc\x2d\x45\xff\x50\x8a\xc5\xbb\xf9\x23\x00\x00\xff\xff\x49\x46\x7e\xf8\x3f\x07\x00\x00")

func lnmonTmplBytes() ([]byte, error) {
	return bindataRead(
		_lnmonTmpl,
		"lnmon.tmpl",
	)
}

func lnmonTmpl() (*asset, error) {
	bytes, err := lnmonTmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "lnmon.tmpl", size: 1855, mode: os.FileMode(436), modTime: time.Unix(1517919783, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"lnmon.tmpl": lnmonTmpl,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}
var _bintree = &bintree{nil, map[string]*bintree{
	"lnmon.tmpl": &bintree{lnmonTmpl, map[string]*bintree{}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}

