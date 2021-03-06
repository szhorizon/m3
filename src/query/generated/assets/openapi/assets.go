// Code generated by "esc -modtime 12345 -prefix openapi/ -pkg openapi -ignore .go -o openapi/assets.go ."; DO NOT EDIT.

// Copyright (c) 2019 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package openapi

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
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
	return nil, nil
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

	"/asset-gen.sh": {
		local:   "asset-gen.sh",
		size:    238,
		modtime: 12345,
		compressed: `
H4sIAAAAAAAC/0zKz0rEMBDH8Xue4rfTnBbSsP45LR5EfAHrTUTWdpIO0hlJIgjiu0sr6M5l4PP9dbv4
Khrr7NztMNw/vgwPdzf+4JIVCERB/s8p7o+YzAGAJOzwhDDBC56PaDPrFtYbTZvoB2+QxG2fx9lAmZXL
qYlmpGILvNBvrSPCYlOThUGHi8ura0J4L5zkE+S/pOv28Tuu9pb/gRAkqxVGnw3BzqanWrnVPhuBenKT
KbufAAAA//9BiTev7gAAAA==
`,
	},

	"/index.html": {
		local:   "openapi/index.html",
		size:    636,
		modtime: 12345,
		compressed: `
H4sIAAAAAAAC/0ySQW/bMAyF7/kVjC+9RJaHDtiQyd6wpceuQ9DLblUk2lYrS55IpzC2/ffBUdLlRr4n
fXwgqNa7h2+PP3/cQc+Db1ZqLcTq+8Pj3Rb2U4CnQb8gaCJk0WEQvyZM8xO4FuY4QTbDDKbXoUMCjsC9
I2idx/VKiGalMhZA9ajtUgAoduyxub/dfYU97qJRMivZHZD1QkyEXBcTt+JjIa+9oAesi6PD1zEmLsDE
wBi4Ll6d5b62eHQGxanZgAuOnfaCjPZYvyvOIO/CC/QJ27romUfaStnGwFR2MXYe9eioNHGQhuhzqwfn
5/p+8TElzdvbqtq8r6rNh6r6s4+HyPFaKiChrwvi2SP1iHwZelJyDXCIdobf5wZg0KlzYQvVpzdp1Na6
0F1pfzNHvoGUvKxVLbzznIQ2GqARjZiSr2/iiEGPThJrdkYuRjkP/qZR8vT0Es8kNzJQMv+XYmwon8mi
d8dUBmQZxiF/+uI1I7E8TMF6pCyWxDpY7WPA8pmKZsl6ouawOaOS+Sj+BQAA//8by2IcfAIAAA==
`,
	},

	"/spec.yml": {
		local:   "openapi/spec.yml",
		size:    23216,
		modtime: 12345,
		compressed: `
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
76r3dwAAAP//AQ7F67BaAAA=
`,
	},

	"/": {
		isDir: true,
		local: "openapi",
	},
}
