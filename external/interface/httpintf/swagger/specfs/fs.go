// Code generated by "esc -o ./specfs/fs.go -pkg specfs api.json"; DO NOT EDIT.

package specfs

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

// FS returns a http.Filesystem for the embedded assets. If useLocal is true,
// the filesystem's contents are instead used.
func FS(useLocal bool) http.FileSystem {
	if useLocal {
		return _escLocal
	}
	return _escStatic
}

// Dir returns a http.Filesystem for the embedded assets on a given prefix dir.
// If useLocal is true, the filesystem's contents are instead used.
func Dir(useLocal bool, name string) http.FileSystem {
	if useLocal {
		return _escDirectory{fs: _escLocal, name: name}
	}
	return _escDirectory{fs: _escStatic, name: name}
}

// FSByte returns the named file from the embedded assets. If useLocal is
// true, the filesystem's contents are instead used.
func FSByte(useLocal bool, name string) ([]byte, error) {
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

// FSMustByte is the same as FSByte, but panics if name is not present.
func FSMustByte(useLocal bool, name string) []byte {
	b, err := FSByte(useLocal, name)
	if err != nil {
		panic(err)
	}
	return b
}

// FSString is the string version of FSByte.
func FSString(useLocal bool, name string) (string, error) {
	b, err := FSByte(useLocal, name)
	return string(b), err
}

// FSMustString is the string version of FSMustByte.
func FSMustString(useLocal bool, name string) string {
	return string(FSMustByte(useLocal, name))
}

var _escData = map[string]*_escFile{

	"/api.json": {
		name:    "api.json",
		local:   "api.json",
		size:    21995,
		modtime: 1587738245,
		compressed: `
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
`,
	},
}

var _escDirs = map[string][]os.FileInfo{}
