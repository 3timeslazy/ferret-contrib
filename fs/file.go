package fs

import (
	"bufio"
	"encoding/json"
	"hash/fnv"
	"io/ioutil"
	"os"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

type File struct {
	File *os.File
}

var TypeFile = core.NewType("file")

func (f *File) Type() core.Type {
	return TypeFile
}

func (f *File) String() string {
	return f.File.Name()
}

func (f *File) Compare(other core.Value) int64 {
	if other, ok := other.(*File); !ok {
		return types.Compare(f.Type(), other.Type())
	}

	otherf := other.Unwrap().(*os.File)

	otherinfo, othererr := otherf.Stat()
	srcinfo, srcerr := f.File.Stat()

	if othererr != nil && srcerr != nil {
		return 0
	}
	if othererr != nil {
		return 1
	}
	if srcerr != nil {
		return -1
	}

	sizediff := srcinfo.Size() - otherinfo.Size()

	if sizediff > 0 {
		return 1
	}

	if sizediff < 0 {
		return -1
	}

	return 0
}

func (f *File) Unwrap() interface{} {
	return f.File
}

func (f *File) Hash() uint64 {
	content, err := ioutil.ReadFile(f.File.Name())
	if err != nil {
		return 0
	}

	h := fnv.New64a()
	buf := bufio.NewWriter(h)

	buf.Write([]byte(f.Type().String()))
	buf.Write([]byte(":"))
	buf.Write(content)
	if buf.Flush() != nil {
		return 0
	}

	return h.Sum64()
}

func (f *File) Copy() core.Value {
	cpy, err := os.Open(f.File.Name())
	if err != nil {
		return values.None
	}

	return &File{cpy}
}

func (f *File) MarshalJSON() ([]byte, error) {
	data, err := ioutil.ReadFile(f.File.Name())
	if err != nil {
		return nil, err
	}

	return json.Marshal(data)
}
