// Code generated by "esc -o ./views/fs.go -ignore fs.go -pkg views -prefix views ./views"; DO NOT EDIT.

package views

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

	"/doc/show.html": {
		name:    "show.html",
		local:   "views/doc/show.html",
		size:    60,
		modtime: 1565038982,
		compressed: `
H4sIAAAAAAAC/6quTklNy8xLVVDKTczMU6qt5VJQUFCwScksswOzQKC6Wg8mrg+WqK5OzUuprQUEAAD/
/+o8KgA8AAAA
`,
	},

	"/flag/create.html": {
		name:    "create.html",
		local:   "views/flag/create.html",
		size:    2159,
		modtime: 1568556972,
		compressed: `
H4sIAAAAAAAC/8xVT2vbThC951MMc/n9epDk5CwJQkugEEpIybmMtWNpYbW77J/8wfi7l5XkRBV2apdS
4otnNbPvvXmz0m63gjdSM2BPUuNudwEAUAr5CI0i7ytsjA6sA9ZDZsh2V4tk1jEJkD5rWAd2WG+3kH+j
nmG3K4vuarZ3hmyj46ydAS8LVLY2z5lyLQylcfrLLhd7hn0b43qgJkijKyw2ilqEnkNnRIXW+IC/0A7V
r1FGSraaBYKgQNlahidygnX2RKHppG4rPMQ58kpWwnM4nB5KpLYxgKaeK/wxakIIL5Yr7KQQrBEeSUWe
hNYXx6GW9qUBOKOy1plo8biGYbOiNSvYGFfhhilEx1jfjAGkaZXFUPEblLEbKd5AptaS53kK8V2A6Tf2
H/g5nFRuFTXcGSX4TT0kxoE8z/N3mi8LIR//hauWXXoDqGWs741SJga4e312rr0ztLnDbgTOfXAUuH3J
Z3WnG+9IJ9xe6gpXCD09V3i5Wp2EMB3W1Uew3DMLrL8zCzAbCB2DIy1MD9ZzFGa/knrMjd6dO4mB5NAM
EnriPsN4Hfs1uz848+msT7z7rpKuPM/hf2PTd4/Up48wEsGN9NLoaysf7m+x/jKt4fruKzzc357r/gLv
3XdhX3trWtlcW3nGYKJT508ldNIPn6H/PDTRB9PDXgKopGHf9F+a0TqGYPSk2Md1LxcX21QwizPrZE/u
BevPjikw3Chqy2JMHuEqi+PXWlmkK3Nxa4+6DyyncLtlLXa7nwEAAP//kRpWMG8IAAA=
`,
	},

	"/flag/index.html": {
		name:    "index.html",
		local:   "views/flag/index.html",
		size:    817,
		modtime: 1564181141,
		compressed: `
H4sIAAAAAAAC/4SS3YrjMAyF7/sUwrCXiaHXrpeFZWFhGIaZJ1BjNRG4dnHUQsf43YekGTpN/+6CvhOd
I8k5O9pwIFBb5KBKWQAAGMcHaDz2/Uo1MQgFUXYkI+2WM1h1hA64rxoKQknZ9+h93Av889jCC/didLe0
i3MLwbWn7y67faLqVDl/Vl1M/BmDoP/hffp5sJvX0mVhEtpX3JLR0t2mb5SGyNg+0PxphGPorwVGz00H
zRhtlm0d3fFSmTMkDC1BDdPOn05zAs7mDPUwFZRitLjHwukQ9YckFGqP9XliKCVn2CUOAuqXetruJhjh
3t+Ho8CzNQhdos1K6Y3H9je71RDv/18oRVlyLEajNdrzAx99z+h28Ov75AwU3HzjRs8OZPT4AqczGu34
YBc5U3ClfAUAAP//oPGN5jEDAAA=
`,
	},

	"/flag/show.html": {
		name:    "show.html",
		local:   "views/flag/show.html",
		size:    5558,
		modtime: 1568555968,
		compressed: `
H4sIAAAAAAAC/9RYXW+jOhB9z68Y+aW7KxHaPidcVUorRaquqlT7fOXgCVhrbGSbKhXiv1/ZQEtoAiTt
/VheAvH4jD1zxmegLBnuuEQgGeWSVNUMAGDB+AvEghqzJLGSFqUlkR/xo+ltbzBIkTLgJohRWtQkKkuY
PwiazP+kGUJVLcL0tgPQgc8LjUHSQe8biGCr9oHQCXjTovkJbnpz/Lyd0hnQ2HIllyTcCZoQyNCmii1J
rowlB2699dtdQAVPJLIjuDU2R8EM2uPD3oTLvLAgaYZL8lftl4B9zXFJUs4YSgIvVBTo/LuInkbqR8gF
WisRJFoVOTm9BD9Z0C0K2Cm9JDlqlxOaIIk2SghVWHh6+28RetMRuHpXnB2gDU5prjoQLgtzXfueG6up
xeR1fiZUHURNZYIEMi6X5JpARvdLcnN9PQmhCfwbM5tozJ/bFb2HBapqIMSLkPGXfyN3Bh0ZnxEZqB3Y
FEFTyVQGucGCqfaJy3qs3s+5KfVOLkymW8Dz1Pl1BmWRbVF/KmGbxutwltxVluGPhcmp/FD1QYbG0AQD
LgWXSKJF6OyiH2FzAv6niWcYc8OVvMv5z80jiVbNM9w9reHn5vHcHPfwPlu6LdyjSnh8l/Mz0l9oMcm6
LIHvThXq6sD/0xqqarS4j8whZYmSDeS7c+WCxpgqwVAviU25AReaKwNxYazKoA0JCIff5gm+qdxJERXf
P3OedFXFZ8TdntCVvvCS6AxczsZQ16sRzG1hrZINiim2Ge9pbmPQuQ9yzTOqX0n0TF8QnB940rhDjTJG
swhrsxMRWoSndXkRulr/ndsEhgIt/i+SqGQsePxrSTTaQkuIldxxnX27utMIr6oAU2j84+r7sXS7I8xt
ZLDUNiiQmpoAA8XSsuEryNAvvea5PzTcDpcl+L4E5k9cKGugqmazQV1IjrV+H6yK4Ca4PdUk1hQ+5OwB
ocPcLeYjkYfy/h68yVKF0slDhtLOHb/GqHWUtn6l8xppYjPhLqdsR91PBaj5HacY/9qq/eRpjTbdN+uF
qvIQyBosZK2yRM52fr+3qCUVfnHDgW11fXZm7AZrfkpOjuVjb9erAdiDbZ0P7zg6iP+A1BYaH7xd7WLY
x6ekx6AdEZqx86V7xsymavw/UvRhIQ1eXPq/P02+miUk8gEdlp4vp0f/v7IElKw9QqbL08VKMdryDOnF
Re9CBx8u3KLuPbkir6vQMgnWq8EXoO73ineQ2ejLTpfPs1HZsLi3g2YHrwz1DtYrsCm1IBGZAatgi9DK
XrgVymtHdPnhcW6M39WTRPdv9/BsqS3M1Bh3QCbHeIrU99T5jLBccooM9cWfU5YNJtxY1OA5cEJkjp0c
3Ya13WvTWPwdAAD//2L9o3O2FQAA
`,
	},

	"/layout.html": {
		name:    "layout.html",
		local:   "views/layout.html",
		size:    2293,
		modtime: 1565156367,
		compressed: `
H4sIAAAAAAAC/7yWUWvrNhTH3/spFN3XKVrZHi4jNoyuZYMWRulbKUOxjp3TypKRjpOFkO8+JDuJkzhQ
unCfcmQd//8/nejI2mw0lGiB8UZVwLfbG8YYm020K2jdAFtQbfLuWQyZUbbKOFi+ewhK5zcpTmNCMpCT
qyoDfia74WC+BlKsWCgfgDLeUim+91KHaatqyLiGUHhsCJ3lrHCWwFLGXzrl39ids+SdYWvXevYAiloP
7BkMqACBrZAWrKMIkwsGS4RV4zwN1FeoaZFpWGIBIg1+YmiRUBkRCmUgu53+zIfrMWg/mAeT8UBrA2EB
QJwtPJQZlyoEoCCb1kMRul9Ro50WIQyZPqMR3y+dJaFWEFwNvcZBZCLEK5bMELC/7tn3t/xrjJVHHYSH
0DgbcAnCGS1wFHvyClZj+SZEfoZRUU8RH1wLZYxBiCOOq9UjAhi1di0FGVCDqMG2fTF+VCFGGT5VgZns
+jLFc6fX+U0/oXHJUGe80+WsMCqEjKetWZ1uKPYEtu3biB3hqx71W0TiSTFGj2g/9pqpYnGdA9y98p+q
nre+As+wcPZIPOWERtl8JtPPwVaqIeFuLR1D7xoLNVeepRW14nYX1Frcil9PUaLEsARJ6jinX+9plogV
Rlvt/rRvPH/ZnXpDyr1Ea841DAYaset3yHk+EtQ831dflkZVEq2Gf/mYeCz9IwZiD0ZVIXLNpMH/61d4
UAQXDe/SdLK8gmODxpEs0eqLhn/HFHb38ngFO+2KIJ/vf//j6X5aX7aMaV92G80/6ioZVqqqwIsWLyEw
Ur6Kn9B/5kZFpMMraf+N6o/SzmRrTrpCalwOuy4NR/pOod3z9d/Qsb775azvNhuCujFxn/Qq0/7mcWp4
FHcXAhZ8cTgk34NscfoeeDwt0vz+ze7g68/DdI/ZbMDq7fa/AAAA//+yy/SV9QgAAA==
`,
	},

	"/login.html": {
		name:    "login.html",
		local:   "views/login.html",
		size:    360,
		modtime: 1560121981,
		compressed: `
H4sIAAAAAAAC/2SQUY7yMAyE33sKy+80F2h6if+/QGgMWJvYUeKyQlXvvmrLSrC8JfONRx4vS6QLCwHm
wILr2gEADJHvMKXQmsdJxUgMx53s9KI1/+IyVzptAkKYjFU8uqRXFoRMdtPosWh7HT8imFJsZO/yjljK
bCAhk0fTLxIEexTyWEJr31ojQklhopumSNXj/83T9z2O3WfYeTZTeQa0+ZzZ8G31p+HlfSqVc6gPHP/x
VYBlcAf4U8F9dhjcdonjP7jI97FbFpK4rj8BAAD///JUCoxoAQAA
`,
	},

	"/pilot/edit.html": {
		name:    "edit.html",
		local:   "views/pilot/edit.html",
		size:    1731,
		modtime: 1564185877,
		compressed: `
H4sIAAAAAAAC/7RUzY6bMBC+8xQjq6eVWKt7rIBLk1R7Wa3U3isHD2DV2NQekq0s3r3CsBFsm0RNW58S
efh+xjNfCPwOPm2/AO+UtsRRKoI7PgxJCBIrZRBYK5Rhw5AAAGRSHaDUwvucldYQGmJFvIm3zcOby7RB
IUH5tERD6FjxPNLAR2vIWQ3PwqD+ACHAfbzYvhA6I/TjBoYh481DkZzAQwAnTI1wv0NBvcOdFrWHWdjy
ZJV1LYiSlDU5m51VWtTcI6XOam17YtAiNVbmrLOe2C8g05nddL3DNKKefqVCq9qgZCAFiXSv6CicRJMe
BZWNMnXO3rOF/LVAhVp6pOIMLUCmTNcT0I8Oc9YoKdEwMKLFnH2dhDM4CN3jrP8mpNiZe3yhx80JLQR4
95vH+BuCsfVrhuUTvsIn5wkWQxf7X07jk9bO9t0FZfFjLfaoobIuZ2jGx2/j0L5R8SRajDMXy69AetRY
Eii5glw3dUF1ESwC2m6c1dcGHRtFqJWncbpCAFUBfp835DMJQliXDMMkCGUIaOQwXGcczzZKRHldHp/0
/bGRvRblt8tGViU3GtlZVyJscNz4g6D/6ag3Uyqe9bMouNHNkyXwHZaqUv/MScYnKUVyuWzfE1kzr7Lv
960ittq8qYAVXhww49O/C8nApTqcy0B+PgQzPiZssYx+NHKM+mQBO3f1ZwAAAP//le3oM8MGAAA=
`,
	},

	"/pilot/find.html": {
		name:    "find.html",
		local:   "views/pilot/find.html",
		size:    460,
		modtime: 1563994919,
		compressed: `
H4sIAAAAAAAC/2yQz27CMAzG7zyF5SNSG+3ecgJNSDtw2Au4jQFLqROlLn9U9d0nCpsY283JL/l++TKO
bgnvm09wSUI0txf1sHTTtBhHz3tRBuxIFKdpAQBQeTlBG6jva2yjGqvhaiYz3cfcAbUmUWt8ikTo2I7R
15hib/idkIbMxXznZyooyEHZI3gyKhqxM2XPWpzJ2qPooca3J+PdKhx8z/Z7e0aiaTCwa+IajS+GoNRx
jfPTSr7Ydo2QArV8jMFzrnF3I7AbmiAtbC7GWSnAdl2WJf4jaAazqA9DPzSdvNR7HHiai5Slo3zF1Uck
D7Owcnf20sz9rVa52zfd15XzclotxpHVT9NXAAAA//+9jaN8zAEAAA==
`,
	},

	"/doc": {
		name:  "doc",
		local: `views/doc`,
		isDir: true,
	},

	"/flag": {
		name:  "flag",
		local: `views/flag`,
		isDir: true,
	},

	"/pilot": {
		name:  "pilot",
		local: `views/pilot`,
		isDir: true,
	},

	"/views": {
		name:  "views",
		local: `./views`,
		isDir: true,
	},
}

var _escDirs = map[string][]os.FileInfo{

	"views/doc": {
		_escData["/doc/show.html"],
	},

	"views/flag": {
		_escData["/flag/create.html"],
		_escData["/flag/index.html"],
		_escData["/flag/show.html"],
	},

	"views/pilot": {
		_escData["/pilot/edit.html"],
		_escData["/pilot/find.html"],
	},

	"./views": {
		_escData["/doc"],
		_escData["/flag"],
		_escData["/layout.html"],
		_escData["/login.html"],
		_escData["/pilot"],
	},
}
