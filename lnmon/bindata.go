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

var _lnmonTmpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x9c\x55\xcd\x6e\xe4\x36\x0c\x3e\xc7\x4f\xc1\xe6\xd2\x4b\xd7\xc2\xee\xde\x0a\x55\x40\xb2\x40\xd1\x00\x69\xba\x68\x66\x5b\xf4\xa8\x19\x71\x46\x6a\x64\xd2\x90\xe8\x4d\x02\x23\xef\x5e\x48\x63\x67\x3c\x99\xa4\x58\xec\xc5\xb0\x48\x7e\xe4\xc7\x3f\x49\x7b\xe9\xa2\x69\xf4\x9a\xdd\xa3\x69\xb4\x7f\x6f\x7e\xc3\x18\xf9\x07\xad\xfc\x7b\xd3\xe8\xde\xac\x7c\xc8\x10\x32\x58\x02\x7c\xe8\x31\x85\x0e\x49\xe0\x3e\x88\x87\xcb\x20\x1b\x0e\x04\x96\x1c\x5c\x87\x9d\x17\x0a\xb4\x83\x1b\x94\x7b\x4e\x77\x90\x06\xaa\x67\x2b\xa0\x37\xec\xd0\x44\x6a\xfd\xdd\xbf\xd4\x76\xa8\x55\x15\xc0\xfa\x11\xb4\x05\x9f\x70\xfb\xcb\xb9\x17\xe9\xf3\xcf\x4a\x4d\x26\xe7\xa6\xfc\x68\x65\x4d\xab\x55\x6f\x9a\x3d\x15\x84\x7e\x48\x3d\x67\x2c\x8c\x84\x21\xa2\x4d\x04\x1d\x27\x04\xbb\xe6\x41\x40\x3c\x42\x20\xc2\x04\x85\x43\xa0\x5d\x06\xde\x1e\x11\x2d\x16\xd7\x37\x25\x74\x67\x8b\x01\x04\x01\xb4\x39\x60\x2a\x0e\xb9\xc7\x64\x05\x1b\xf1\x98\x11\x9c\xc5\x8e\x29\xff\x54\x81\x0b\xa3\x75\xc6\xf4\x15\xe1\xde\xdb\x1a\xf1\xf1\xc7\x84\xe0\x38\xd0\xee\xc0\xf5\x1f\x1e\xa0\xb3\x8f\x60\x63\x66\x58\x17\x52\x82\x09\xb3\xa0\x83\x40\x95\xc4\x49\xe6\x1d\x3f\x17\x48\x45\x52\x1b\xa6\xcc\x11\xb3\x0a\xe4\xf0\xa1\x2d\x7d\x6a\xf1\xc1\x76\x7d\xc4\x73\xf3\x39\x71\x87\xe2\x71\xc8\xe0\x6c\xf6\x6b\xb6\xc9\x95\x62\x41\xde\x24\xdb\x97\xb4\x9c\x15\xdb\x6c\x13\x77\x20\xbe\x16\x8b\xe3\xcc\xce\x7f\x30\x53\x47\xe6\xa6\xb9\xb9\x23\x81\xb6\xac\x95\xff\x60\x9a\x71\x84\xb0\x85\xf6\x2a\xff\x39\xf5\xf1\xdd\xd3\x53\x49\xec\x4d\x64\x06\x9d\x25\x31\xed\xcc\xd4\x79\xad\xa6\x73\x0d\xab\x7b\xf3\xc7\x90\x20\xb8\x6a\x59\x31\xe3\x08\xed\x15\x6d\xb9\xbd\x61\x87\x57\xee\xe9\x69\xf2\xb5\xe7\x39\x8e\xfb\x29\x6b\x2f\x62\xb0\x79\x0e\x5f\x7c\xd8\x2a\x38\x72\x03\xa7\x60\x8c\x19\x61\x0f\xfa\x42\x77\xc4\xf7\x34\x01\xb7\x9c\x80\x87\x04\x45\x42\xec\x70\x46\xbc\x03\x24\x07\x8b\x30\xce\x25\xcc\x35\xd0\x29\xe5\x8b\x49\xf9\x22\xee\x04\xfd\x8a\x29\x07\xa6\xd7\xa1\x7f\x4d\xca\xd7\xa1\xeb\xc8\x9b\x3b\x8f\xa5\xbe\xaf\xc3\x2f\x17\x06\xa7\x2e\xfe\x46\x28\xb9\x96\xa9\x1f\x47\x88\x48\x50\xab\x5b\x88\xd6\x64\x73\x1d\x3f\x16\x1b\xff\x0f\xf2\xc9\x5b\x22\x8c\x15\xb5\x99\xff\xbf\x05\x58\x63\xb5\x9f\x11\xd3\x22\xa2\x94\x3d\xb1\x09\x6b\xd5\xfb\xa2\x7b\x9e\xc4\x8f\xa6\xda\x6a\xe5\x3f\x96\xb3\xd8\x75\x44\xd3\x9c\x69\x49\xa6\x39\x3b\xd3\xe2\xcd\x95\xd3\x4a\xfc\x7c\xaa\xb3\xb0\x14\x5c\xdb\x2c\x90\x11\x69\x29\x9c\xe8\x43\x16\x2b\xb8\x54\xfc\x3a\x90\xab\x17\xc7\x90\x5f\x15\x8b\xc7\x6e\x52\x68\x55\x39\x8c\x23\x24\x4b\x3b\x3c\x4e\xad\x0c\xe3\x81\xa4\x33\x87\x16\xed\x47\xb9\xbd\xf5\x9c\x16\xed\xd1\x4a\xdc\x6c\x5b\xac\xf6\x23\x5d\xd4\xc7\xf2\x92\xcd\x2a\x74\x98\xc5\x76\x7d\x7b\x1b\x68\x83\x2f\xac\x0e\x91\xe6\x1e\xb5\xb7\x25\xcb\x37\x63\x95\x25\x5e\xb6\xf3\x08\xfa\x7b\x88\x31\xdc\x5a\xe1\xec\xc3\x8a\xbf\xe4\xf6\x22\x5f\xae\x3e\x4d\x66\xfb\x6d\xf8\x7e\x77\x2b\x8f\xdd\xdb\x0e\x9f\xeb\x7b\x58\x3b\x35\xb5\xbf\xd1\x43\x79\x92\x54\xf9\xbe\xd8\xe4\x6f\xb8\x7d\x88\x05\x5e\xde\x40\x27\x2b\xde\x68\x35\xbd\x78\x6a\xff\x00\xfe\x17\x00\x00\xff\xff\x9c\xac\x27\x1e\x08\x07\x00\x00")

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

	info := bindataFileInfo{name: "lnmon.tmpl", size: 1800, mode: os.FileMode(436), modTime: time.Unix(1518021072, 0)}
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

