// Code generated by go-bindata.
// sources:
// bcmon.tmpl
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

var _bcmonTmpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x8c\x52\xbd\x8e\x13\x31\x10\xee\xf7\x29\x3e\xd2\xd0\x70\x6b\xdd\x95\xc8\xb8\xa0\x82\x0e\x21\x1a\x4a\x6f\x3c\x89\xcd\xd9\x33\x2b\xcf\x84\x24\x8a\xee\xdd\xd1\x26\x7b\x07\x02\x21\x5d\xb7\x3b\x1a\x7f\xbf\xe3\xb3\xb5\x1a\x06\x3f\x49\x3a\x87\xc1\xe7\xfb\xf0\x89\x6a\x95\x37\xde\xe5\xfb\x30\xf8\x39\x7c\xcb\x45\x51\x14\x91\x41\xa7\x99\x7a\x69\xc4\x86\x63\xb1\x8c\x8f\xc5\xb6\x52\x18\xfd\xc0\x5c\x78\x8f\x68\xf0\x5b\x49\x14\x2a\x8f\xf9\xf1\x07\x8f\x8d\xbc\xbb\x0e\x30\x9d\xe1\x23\x72\xa7\xdd\x87\x4d\x36\x9b\xf5\xbd\x73\xeb\xca\x26\x2c\x1f\xde\xc5\x30\x7a\x37\x87\xe1\x46\x4a\x98\x0f\x7d\x16\xa5\x85\xdb\x04\x95\x62\x67\x34\xe9\x84\x38\xc9\xc1\x60\x99\x50\x98\xa9\xe3\x28\xfd\xb1\xf0\x5e\x21\xbb\x17\x49\xd3\x19\x2d\x2e\x53\x14\x03\x45\x2d\xd4\x17\x14\x99\xa9\x47\xa3\xc1\x32\x29\x21\x45\x6a\xc2\xfa\x0e\x91\xd3\x9f\x4b\x93\x52\xff\x49\x38\xe6\x78\xa5\x39\xbf\xed\x84\x24\x85\xf7\xbf\x05\x7e\x97\x03\x5a\x3c\x23\x56\x15\x4c\x8b\x12\xa3\x4e\x6a\x94\x50\xf8\xaa\xed\x1f\xbb\x4d\x5e\x52\x71\x95\xdd\x56\x58\xa5\x92\xba\xc2\x89\x4e\xe3\x52\xc3\x48\xa7\xd8\xe6\x4a\x9b\xf0\xa5\x4b\x23\xcb\x74\x50\xa4\xa8\x79\x92\xd8\xd3\x92\x10\x74\xdb\xe3\xbc\xd8\x4a\xd1\xe2\xb0\xeb\xd2\x60\xf9\x9a\x90\xd4\x67\x75\xf9\x21\xdc\x6a\x98\x6e\x61\xa4\xe7\x12\x0a\xef\xc4\xbb\xfc\x10\x86\xcb\x05\x65\x87\xf1\xb3\x7e\x5d\xab\xbb\x7b\x7a\x5a\x6c\xfd\xe7\x9d\xc2\xab\x75\xe1\x7d\x58\xab\xf6\x6e\xfd\xbf\x52\x5e\x2e\xa0\xaa\x84\x57\x62\xb0\x18\xfe\xc6\x19\x57\xa0\x3b\x10\xa7\x05\x68\xf0\x6e\x3d\x49\x77\xbb\xd0\x5f\x01\x00\x00\xff\xff\x70\x3e\xba\x3c\xa9\x02\x00\x00")

func bcmonTmplBytes() ([]byte, error) {
	return bindataRead(
		_bcmonTmpl,
		"bcmon.tmpl",
	)
}

func bcmonTmpl() (*asset, error) {
	bytes, err := bcmonTmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "bcmon.tmpl", size: 681, mode: os.FileMode(436), modTime: time.Unix(1518021084, 0)}
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
	"bcmon.tmpl": bcmonTmpl,
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
	"bcmon.tmpl": &bintree{bcmonTmpl, map[string]*bintree{}},
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

