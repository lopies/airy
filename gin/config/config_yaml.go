// Code generated for package config by go-bindata DO NOT EDIT. (@generated)
// sources:
// config.dev.yaml
// config.prod.yaml
// config.test.yaml
package config

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

// Name return file name
func (fi bindataFileInfo) Name() string {
	return fi.name
}

// Size return file size
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}

// Mode return file mode
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}

// Mode return file modify time
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}

// IsDir return file whether a directory
func (fi bindataFileInfo) IsDir() bool {
	return fi.mode&os.ModeDir != 0
}

// Sys return file is sys mode
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _configDevYaml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x84\x53\x4d\x4f\xdb\x40\x14\xbc\xfb\x57\xac\xe0\xe0\x53\x12\x43\x20\x34\xbe\xa1\xb6\xaa\x10\x88\xa6\xa5\xad\x10\x17\xe4\xc4\x1b\x67\x89\xbd\x9b\xec\xda\x90\x70\xe2\x90\xf2\x51\x29\x05\x95\xa4\xa1\x94\x16\x90\x40\xd0\x03\x81\x56\x2d\x20\x08\xf0\x67\x62\x3b\x39\xf1\x17\xaa\xb5\x1d\x22\x10\x6a\x6f\x7e\xcf\xf3\x66\x67\x66\xdf\x32\x48\x67\x21\x95\x05\x00\x0c\xa2\x42\x19\x88\x2a\x4c\x5a\x9a\x28\x00\xa0\xc2\xb4\x62\xe9\x66\x42\xd1\xe0\x04\x9a\x87\x32\xe8\x93\x40\x2f\x68\x5f\x6e\xb4\xea\x7b\xf6\xf2\x62\x7b\xf7\x8f\xf3\x6d\xd7\xa9\x9e\xf0\x51\xa5\xd0\x85\x0d\x4a\x12\x00\x00\xf4\x02\x67\x6b\xc1\xde\x3b\x78\x00\x4d\x23\x1d\xbe\xcd\xe9\x44\x51\x13\x8a\x99\x91\x81\x18\x8e\x88\x1c\xfb\x79\xa9\x79\x79\xda\x3c\xff\xd0\x6c\xec\xb4\xce\x8e\xed\xeb\x92\x30\x33\x67\x72\x59\x26\xc9\x42\xfc\xbc\x90\x43\x14\xca\xe0\x49\x6c\x40\xf2\xe8\x7b\xfd\x7e\xeb\x66\xc9\xd9\xda\x76\x6a\xa7\xed\xda\xef\x0e\x76\x14\x16\x65\x20\x0e\x5b\x66\x86\x50\x34\xaf\x98\x88\x60\xb1\x83\x77\x37\x4b\xed\x4a\xdd\x5e\x2b\x77\xc0\x23\x8c\x59\x90\xca\x40\xd4\x10\x0e\xa5\xa9\x62\xc0\x39\x42\xb3\xdc\x3e\x83\x29\x0a\x4d\x19\x88\xd1\x99\x57\x13\xaf\xa7\x26\xe1\xe4\x3c\x95\xe2\x99\xbc\xf4\x82\x4e\x59\x63\xf1\x67\x9a\x99\xed\x9f\x7c\x17\x1f\x2e\x88\x82\x4e\x34\xd9\x0b\x2c\x69\x69\x32\x30\xa9\x05\x03\xa3\xe3\x8a\xc1\x23\x35\xf2\x10\x43\xaa\x15\x43\x1a\xf1\xbc\xd6\xf6\xed\x9b\x9a\xef\xd8\x5e\x2b\xbb\x07\x3c\x17\x15\xd1\x20\x10\x6a\x61\x13\x19\x30\xa2\x13\x8d\x75\xe1\x41\xee\x47\x1b\x4e\xe5\xda\xfd\x5a\xb7\xaf\xaa\x82\x51\x64\x79\x9d\x1f\x1c\x02\x19\xc2\xb8\xd4\xbe\xfe\xa1\xb0\x14\x96\xc2\x7d\xdc\x00\x00\x39\x42\x3d\x03\x51\x29\xe6\x37\x2c\xe6\x79\xa5\x84\x98\x01\x42\x61\x6c\x8e\x50\xd5\x9b\x8d\x0e\x0c\x06\x38\x35\x89\x7d\xe5\x1a\xc2\xd3\xf7\x52\x01\x20\x47\x61\x1a\x15\x64\xd0\xd3\xe3\x95\x86\x52\x18\x51\x75\xf8\x94\x60\xcc\x82\x15\x69\xd5\xaf\xdd\xab\xba\xfb\xe3\xa2\x5d\xfb\xd5\xba\xf9\xee\x7c\xdc\x77\x7e\xee\x34\xcf\x8f\xfc\x6f\x77\xb3\xe4\x2f\x86\x53\x3d\x69\x2f\xad\x76\x48\x5e\xe6\x20\xbe\x23\xe9\xb2\x38\x2b\xeb\x76\x63\xc1\xa9\x9e\x38\xe5\xba\x7d\xb1\xfe\x4f\x8a\x31\x94\x86\x6f\x10\xd7\x1d\xeb\x12\x34\x2f\x16\xfd\x21\x7b\xf5\xd8\xde\x2b\xbb\x95\xc3\xee\xb4\xb7\x36\xb7\x8d\x65\xbe\xa2\x9f\xb6\x6f\x1b\x2b\x02\x85\x2a\x62\x3c\xd2\xc7\x02\x0d\xe2\x8c\x45\x87\xe2\x5e\xf9\x58\x76\x6a\x72\xdc\x32\x64\x20\x09\x00\xe8\x44\x43\x38\x11\xa4\x75\x3f\xc9\x69\xef\xdf\x34\xbf\x5d\xf7\xcb\xa5\x7d\x55\x75\x2b\x87\xce\xf2\x99\x7d\xb4\x61\xbf\xdf\xf7\x34\xf8\x5b\xea\x6e\x96\xec\x95\xb2\xdb\x58\x10\x14\x23\x9f\xfb\x8f\xae\xc1\xd8\x50\x3f\x2f\x83\x4b\xd6\x2c\xc8\xcc\x07\x3a\x79\x39\x1b\x50\x88\x02\x61\x9e\x55\x88\xd5\x04\x41\x98\xf7\x32\xa6\x99\x63\x72\x24\x42\x18\x0b\xa5\x70\x88\x65\x14\xac\x65\x14\x14\x56\x74\x54\xb4\x70\x8a\x85\x53\xc4\xe0\x1c\x4a\x2a\x05\x19\x1b\x85\xc5\x91\x0e\xeb\x5d\x67\xa2\xf3\x6c\x78\x37\x69\xa5\xb2\xd0\x0c\x5e\x81\x28\xfc\x0d\x00\x00\xff\xff\x8c\x34\x9a\x3d\x6c\x04\x00\x00")

func configDevYamlBytes() ([]byte, error) {
	return bindataRead(
		_configDevYaml,
		"config.dev.yaml",
	)
}

func configDevYaml() (*asset, error) {
	bytes, err := configDevYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "config.dev.yaml", size: 1132, mode: os.FileMode(420), modTime: time.Unix(1660748364, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _configProdYaml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x84\x53\x41\x4f\xe3\x46\x14\xbe\xfb\x57\x8c\x96\x83\x4f\x49\x0c\x59\x42\xe3\xdb\xaa\xad\x2a\xb4\xab\x6d\x5a\xda\x0a\xed\x05\x4d\xe2\x17\x67\x36\xf6\x4c\x32\x63\x2f\x09\x27\x0e\xe9\x02\x95\x52\x50\x49\x1a\x4a\x69\x01\x09\x04\x3d\x10\x68\xd5\x02\x82\x00\x7f\x26\xb6\x93\x13\x7f\xa1\x1a\xdb\x21\x02\xa1\xf6\xe6\xf7\xfc\xbd\x6f\xbe\xef\x9b\x37\x02\xf8\x07\xe0\xba\x82\x90\xcd\x0c\xd0\x91\xca\xc1\x02\x2c\x40\x55\x10\x32\xa0\x88\x5d\xcb\xc9\x61\x13\xe6\xc8\x12\xe8\x68\x52\x43\x13\x68\x78\xbd\x35\xe8\x1e\x78\xab\x1f\x87\xfb\xff\xf8\xbf\xed\xfb\xed\x33\x39\x8c\x6b\x63\xd8\xb4\xa6\x21\x84\xd0\x04\xf2\x77\x96\xbd\x83\xa3\x27\xd0\x22\xb1\xe0\xdb\x8a\xc5\xb0\x91\xc3\x4e\x49\x47\x6a\x32\xa5\x4a\xec\xcf\x2b\xfd\xeb\xf3\xfe\xe5\x0f\xfd\xde\xde\xe0\xe2\xd4\xbb\x6d\x28\xef\x17\x1d\x29\xcc\x61\x65\xa0\x9f\xd7\x2a\x84\x83\x8e\x3e\xc9\xbc\xd4\x42\xfa\x89\xa8\x3f\xb8\x5b\xf1\x77\x76\xfd\xce\xf9\xb0\xf3\xf7\x08\xfb\x1a\xea\x3a\x52\x5f\xb9\x4e\x89\x71\xb2\x84\x1d\xc2\xa8\x3a\xc2\x07\xdb\x8d\x61\xab\xeb\x6d\x34\x47\xe0\x59\x21\x5c\xe0\x3a\x52\x4d\x42\x13\x45\x8e\x6d\x58\x64\xbc\x2c\xed\x0b\x28\x70\x70\x74\xa4\xa6\xdf\x7f\x35\xf7\xf5\xbb\x79\x98\x5f\xe2\x5a\xb6\x54\xd5\xbe\xe0\xef\xdc\x37\xd9\xcf\x4c\xa7\x3c\x35\xff\x5d\xf6\x55\x4d\x55\x2c\x66\xea\x61\x60\x79\xd7\xd4\x91\xc3\x5d\x88\x8d\xbe\xc5\xb6\x0c\xd5\xae\x02\x05\x6e\xd6\x13\x26\x0b\xbd\x76\x0e\xbd\xbb\x4e\xe4\xd8\xdb\x68\x06\x47\x32\x17\x83\xf0\x38\x10\xee\x52\x87\xd8\x90\xb2\x98\x29\xc6\xf0\x38\xf7\x93\x2d\xbf\x75\x1b\xfc\xda\xf5\x6e\xda\x8a\x5d\x17\x55\x4b\x1e\x9c\x40\x25\x26\xa4\xd4\xc9\xa9\x99\xa4\x96\xd4\x92\x93\xd2\x00\x42\x15\xc6\x43\x03\x69\x2d\x13\x35\x5c\x11\x7a\xe5\x8c\x39\x31\x02\x0b\xb1\xc8\xb8\x11\xce\xa6\x5f\x4e\xc7\x38\x23\x4f\x23\xe5\x26\xa1\x0b\x8f\x52\x41\xa8\xc2\xa1\x48\x6a\x3a\x7a\xf1\x22\x2c\x6d\x5c\x9b\x35\x2c\xf8\x94\x51\x2a\xe2\x15\x19\x74\x6f\x83\x9b\x6e\xf0\xc7\xd5\xb0\xf3\xd7\xe0\xee\x77\xff\xc7\x43\xff\xcf\xbd\xfe\xe5\x49\xf4\x1d\x6c\x37\xa2\xc5\xf0\xdb\x67\xc3\x95\xf5\x11\xc9\x97\x15\xa0\x0f\x24\x63\x16\x7f\x6d\xd3\xeb\x2d\xfb\xed\x33\xbf\xd9\xf5\xae\x36\xff\x93\xe2\x0d\x29\xc2\x37\x44\xea\xce\x8c\x09\xfa\x57\x1f\xa3\x21\x6f\xfd\xd4\x3b\x68\x06\xad\xe3\xf1\x74\xb8\x36\xf7\xbd\x55\xb9\xa2\x3f\xed\xde\xf7\xd6\x14\x0e\x06\x11\x32\xd2\xe7\x02\x8d\xe3\xcc\xa4\x67\xb2\x61\xf9\x5c\x76\x46\xfe\xad\x6b\xeb\x48\x53\x10\xb2\x98\x49\x68\x2e\x4e\xeb\x71\x92\x0b\xe1\xbf\x05\x79\xbb\xc1\x2f\xd7\xde\x4d\x3b\x68\x1d\xfb\xab\x17\xde\xc9\x96\xf7\xfd\x61\xa8\x21\xda\xd2\x60\xbb\xe1\xad\x35\x83\xde\xb2\x82\xed\x6a\xe5\x7f\x74\x4d\x67\x66\xa6\x64\x19\x5f\xb2\xe9\x82\x70\x9e\xe8\x94\xe5\x87\x98\x42\x55\x98\x08\xad\x02\x35\x72\x8c\x50\xd9\x2b\x39\x4e\x45\xe8\xa9\x14\x13\x22\x51\xa0\x09\x51\xc2\xd4\x2c\x61\x92\xc4\x16\xa9\xbb\xb4\x20\x92\x05\x66\x4b\x0e\x5c\x28\x80\x10\xaf\xa1\x3e\x3b\x62\x7d\xe8\xcc\x8d\x9e\x8d\xec\xe6\xdd\x42\x19\x9c\xf8\x15\xa8\xca\xbf\x01\x00\x00\xff\xff\x87\xc4\x8b\xa4\x6e\x04\x00\x00")

func configProdYamlBytes() ([]byte, error) {
	return bindataRead(
		_configProdYaml,
		"config.prod.yaml",
	)
}

func configProdYaml() (*asset, error) {
	bytes, err := configProdYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "config.prod.yaml", size: 1134, mode: os.FileMode(420), modTime: time.Unix(1660751535, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _configTestYaml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x84\x53\xcf\x4f\xe3\x46\x14\xbe\xfb\xaf\x18\xc1\xc1\xa7\x24\x86\x40\x68\x7c\x43\x6d\x55\x21\x10\x4d\x4b\x5b\x21\x2e\xc8\x89\x27\xce\x10\x7b\x26\x99\x19\x43\xc2\x89\x43\xca\x8f\x4a\x29\xa8\x24\x0d\xcb\xb2\x0b\x48\x20\xd8\x03\x81\x5d\xed\x02\x82\x00\xff\x4c\x6c\x27\x27\xfe\x85\xd5\xd8\x0e\x59\x56\x68\xf7\x96\xf7\xf2\xbd\xcf\xdf\xf7\xcd\x7b\x0c\xd2\x45\x48\x55\x09\x00\x8b\xe8\x50\x05\x32\x87\x8c\xcb\x12\x00\x3a\xcc\x6a\xb6\xc9\x53\x9a\x01\x67\xd0\x32\x54\xc1\x90\x02\x06\x41\xf7\x76\xa7\xd3\x3c\x72\xd6\x57\xbb\x87\x9f\xdc\x37\x87\x6e\xfd\x42\x4c\x6a\xa5\x3e\x6c\x54\x51\x00\x00\x60\x10\xb8\x7b\x2b\xce\xd1\xc9\x57\xd0\x2c\x32\xe1\x9f\x05\x93\x68\x7a\x4a\xe3\x39\x15\xc8\xd1\x98\x2c\xb0\xff\xaf\xb5\x6f\x2f\xdb\xd7\xff\xb4\x5b\x07\x9d\xab\x73\xe7\xbe\x22\x2d\x2c\x71\xa1\x8a\x93\x3c\xc4\x3f\x97\x0a\x88\x42\x15\xfc\x90\x18\x51\x7c\xfa\xc1\xa0\xdf\x79\x58\x73\xf7\xf6\xdd\xc6\x65\xb7\xf1\xb1\x87\x9d\x84\x65\x15\xc8\xe3\x36\xcf\x11\x8a\x96\x35\x8e\x08\x96\x7b\x78\x6f\xb7\xd2\xad\x35\x9d\xad\x6a\x0f\x3c\xc1\x98\x0d\xa9\x0a\x64\x03\xe1\x48\x96\x6a\x16\x5c\x22\x34\x2f\xec\x33\x98\xa1\x90\xab\x40\x8e\x2f\xfc\x36\xf3\xfb\xdc\x2c\x9c\x5d\xa6\x4a\x32\x57\x54\x7e\xa1\x73\xf6\x54\xf2\x27\x83\xe7\x87\x67\xff\x4a\x8e\x97\x64\xc9\x24\x86\xea\x07\x96\xb6\x0d\x15\x70\x6a\xc3\xd0\xe8\xb4\x66\x89\x44\xad\x22\xc4\x90\x1a\xe5\x88\x41\x7c\xaf\x8d\x63\xe7\xa1\x11\x38\x76\xb6\xaa\xde\x89\xc8\x45\x47\x34\x0c\x84\xda\x98\x23\x0b\xc6\x4c\x62\xb0\x3e\x3c\xcc\xfd\x6c\xc7\xad\xdd\x7b\xaf\x9b\xce\x5d\x5d\xb2\xca\xac\x68\x8a\x0f\x47\x40\x8e\x30\x21\x75\x68\x78\x2c\xaa\x44\x95\xe8\x90\x30\x00\x40\x81\x50\xdf\x40\x5c\x49\x04\x0d\x9b\xf9\x5e\x29\x21\x3c\x44\x68\x8c\x2d\x11\xaa\xfb\xb3\xf1\x91\xd1\x10\xa7\xa7\x71\xa0\xdc\x40\x78\xfe\x59\x2a\x00\x14\x28\xcc\xa2\x92\x0a\x06\x06\xfc\xd2\xd2\x4a\x13\xba\x09\x7f\x24\x18\xb3\x70\x45\x3a\xcd\x7b\xef\xae\xe9\xbd\xbb\xe9\x36\x3e\x74\x1e\xde\xba\xff\x1e\xbb\xef\x0f\xda\xd7\x67\xc1\x6f\x6f\xb7\x12\x2c\x86\x5b\xbf\xe8\xae\x6d\xf6\x48\x7e\x2d\x40\xfc\x44\xd2\x67\x71\x37\xb6\x9d\xd6\x8a\x5b\xbf\x70\xab\x4d\xe7\x66\xfb\x9b\x14\x53\x28\x0b\xff\x40\x42\x77\xa2\x4f\xd0\xbe\x59\x0d\x86\x9c\xcd\x73\xe7\xa8\xea\xd5\x4e\xfb\xd3\xfe\xda\x3c\xb6\xd6\xc5\x8a\xfe\xb7\xff\xd8\xda\x90\x28\xd4\x11\x13\x91\xbe\x14\x68\x18\x67\x22\x3e\x96\xf4\xcb\x97\xb2\xd3\xd3\xd3\xb6\xa5\x02\x45\x02\xc0\x24\x06\xc2\xa9\x30\xad\xe7\x49\xce\xfb\xff\xcd\x8b\xd7\xf5\x5e\xdd\x3a\x77\x75\xaf\x76\xea\xae\x5f\x39\x67\x3b\xce\xdf\xc7\xbe\x86\x60\x4b\xbd\xdd\x8a\xb3\x51\xf5\x5a\x2b\x92\x66\x15\x0b\xdf\xd1\x35\x9a\x18\x1b\x16\x65\xf8\xc8\x86\x1d\xde\xf1\x17\x3a\x45\xb9\x18\x52\xc8\x12\x61\xbe\x55\x88\xf5\x14\x41\x58\xf4\x72\x9c\x17\x98\x1a\x8b\x11\xc6\x22\x19\x1c\x61\x39\x0d\x1b\x39\x0d\x45\x35\x13\x95\x6d\x9c\x61\xd1\x0c\xb1\x04\x87\x96\xc9\x40\xc6\x26\x61\x79\xa2\xc7\xfa\xd4\x99\xe9\x9d\x8d\xe8\xa6\xed\x4c\x1e\xf2\xf0\x0a\x64\xe9\x73\x00\x00\x00\xff\xff\x3d\xf9\x4e\x3b\x6b\x04\x00\x00")

func configTestYamlBytes() ([]byte, error) {
	return bindataRead(
		_configTestYaml,
		"config.test.yaml",
	)
}

func configTestYaml() (*asset, error) {
	bytes, err := configTestYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "config.test.yaml", size: 1131, mode: os.FileMode(420), modTime: time.Unix(1660751545, 0)}
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
	"config.dev.yaml":  configDevYaml,
	"config.prod.yaml": configProdYaml,
	"config.test.yaml": configTestYaml,
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
	"config.dev.yaml":  &bintree{configDevYaml, map[string]*bintree{}},
	"config.prod.yaml": &bintree{configProdYaml, map[string]*bintree{}},
	"config.test.yaml": &bintree{configTestYaml, map[string]*bintree{}},
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
