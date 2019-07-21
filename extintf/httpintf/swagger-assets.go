// Code generated by "esc -private -o ./swagger-assets.go -pkg httpintf ./swagger.json"; DO NOT EDIT.

package httpintf

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"sync"
	"time"
)

type _escLocalFS struct{}

var _escLocal _escLocalFS

type _escStaticFS struct{}

var _escStatic _escStaticFS

type _escDirectory struct {
	fs   http.FileSystem
	name string
}

type _escFile struct {
	compressed string
	size       int64
	modtime    int64
	local      string
	isDir      bool

	once sync.Once
	data []byte
	name string
}

func (_escLocalFS) Open(name string) (http.File, error) {
	f, present := _escData[path.Clean(name)]
	if !present {
		return nil, os.ErrNotExist
	}
	return os.Open(f.local)
}

func (_escStaticFS) prepare(name string) (*_escFile, error) {
	f, present := _escData[path.Clean(name)]
	if !present {
		return nil, os.ErrNotExist
	}
	var err error
	f.once.Do(func() {
		f.name = path.Base(name)
		if f.size == 0 {
			return
		}
		var gr *gzip.Reader
		b64 := base64.NewDecoder(base64.StdEncoding, bytes.NewBufferString(f.compressed))
		gr, err = gzip.NewReader(b64)
		if err != nil {
			return
		}
		f.data, err = ioutil.ReadAll(gr)
	})
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (fs _escStaticFS) Open(name string) (http.File, error) {
	f, err := fs.prepare(name)
	if err != nil {
		return nil, err
	}
	return f.File()
}

func (dir _escDirectory) Open(name string) (http.File, error) {
	return dir.fs.Open(dir.name + name)
}

func (f *_escFile) File() (http.File, error) {
	type httpFile struct {
		*bytes.Reader
		*_escFile
	}
	return &httpFile{
		Reader:   bytes.NewReader(f.data),
		_escFile: f,
	}, nil
}

func (f *_escFile) Close() error {
	return nil
}

func (f *_escFile) Readdir(count int) ([]os.FileInfo, error) {
	if !f.isDir {
		return nil, fmt.Errorf(" escFile.Readdir: '%s' is not directory", f.name)
	}

	fis, ok := _escDirs[f.local]
	if !ok {
		return nil, fmt.Errorf(" escFile.Readdir: '%s' is directory, but we have no info about content of this dir, local=%s", f.name, f.local)
	}
	limit := count
	if count <= 0 || limit > len(fis) {
		limit = len(fis)
	}

	if len(fis) == 0 && count > 0 {
		return nil, io.EOF
	}

	return fis[0:limit], nil
}

func (f *_escFile) Stat() (os.FileInfo, error) {
	return f, nil
}

func (f *_escFile) Name() string {
	return f.name
}

func (f *_escFile) Size() int64 {
	return f.size
}

func (f *_escFile) Mode() os.FileMode {
	return 0
}

func (f *_escFile) ModTime() time.Time {
	return time.Unix(f.modtime, 0)
}

func (f *_escFile) IsDir() bool {
	return f.isDir
}

func (f *_escFile) Sys() interface{} {
	return f
}

// _escFS returns a http.Filesystem for the embedded assets. If useLocal is true,
// the filesystem's contents are instead used.
func _escFS(useLocal bool) http.FileSystem {
	if useLocal {
		return _escLocal
	}
	return _escStatic
}

// _escDir returns a http.Filesystem for the embedded assets on a given prefix dir.
// If useLocal is true, the filesystem's contents are instead used.
func _escDir(useLocal bool, name string) http.FileSystem {
	if useLocal {
		return _escDirectory{fs: _escLocal, name: name}
	}
	return _escDirectory{fs: _escStatic, name: name}
}

// _escFSByte returns the named file from the embedded assets. If useLocal is
// true, the filesystem's contents are instead used.
func _escFSByte(useLocal bool, name string) ([]byte, error) {
	if useLocal {
		f, err := _escLocal.Open(name)
		if err != nil {
			return nil, err
		}
		b, err := ioutil.ReadAll(f)
		_ = f.Close()
		return b, err
	}
	f, err := _escStatic.prepare(name)
	if err != nil {
		return nil, err
	}
	return f.data, nil
}

// _escFSMustByte is the same as _escFSByte, but panics if name is not present.
func _escFSMustByte(useLocal bool, name string) []byte {
	b, err := _escFSByte(useLocal, name)
	if err != nil {
		panic(err)
	}
	return b
}

// _escFSString is the string version of _escFSByte.
func _escFSString(useLocal bool, name string) (string, error) {
	b, err := _escFSByte(useLocal, name)
	return string(b), err
}

// _escFSMustString is the string version of _escFSMustByte.
func _escFSMustString(useLocal bool, name string) string {
	return string(_escFSMustByte(useLocal, name))
}

var _escData = map[string]*_escFile{

	"/swagger.json": {
		name:    "swagger.json",
		local:   "./swagger.json",
		size:    9154,
		modtime: 1563750743,
		compressed: `
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
`,
	},
}

var _escDirs = map[string][]os.FileInfo{}
