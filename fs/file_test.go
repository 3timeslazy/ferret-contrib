package fs_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/MontFerret/ferret-contrib/fs"

	. "github.com/smartystreets/goconvey/convey"
)

func TestFile(t *testing.T) {
	Convey(".Type", t, func() {
		f := fs.File{}
		typ := f.Type()

		So(typ.Equals(fs.TypeFile), ShouldBeTrue)
	})

	Convey(".String", t, func() {
		file, err := ioutil.TempFile("", "*.String")
		So(err, ShouldBeNil)

		f := fs.File{File: file}

		So(f.String(), ShouldEqual, file.Name())
	})

	Convey(".Compare", t, func() {

		kb := make([]byte, 1024)
		mb := make([]byte, len(kb)*1024)

		Convey("Should return -1", func() {
			f1, err := ioutil.TempFile("", "*.Compare")
			So(err, ShouldBeNil)
			f2, err := ioutil.TempFile("", "*.Compare")
			So(err, ShouldBeNil)

			_, err = f1.Write(kb)
			So(err, ShouldBeNil)
			_, err = f2.Write(mb)
			So(err, ShouldBeNil)

			file1 := fs.File{File: f1}
			file2 := fs.File{File: f2}

			So(file1.Compare(&file2), ShouldEqual, -1)

			f1.Close()
			f2.Close()
		})

		Convey("Should return 1", func() {
			f1, err := ioutil.TempFile("", "*.Compare")
			So(err, ShouldBeNil)
			f2, err := ioutil.TempFile("", "*.Compare")
			So(err, ShouldBeNil)

			_, err = f1.Write(mb)
			So(err, ShouldBeNil)
			_, err = f2.Write(kb)
			So(err, ShouldBeNil)

			file1 := fs.File{File: f1}
			file2 := fs.File{File: f2}

			So(file1.Compare(&file2), ShouldEqual, 1)

			f1.Close()
			f2.Close()
		})

		Convey("Should return 0", func() {
			f1, err := ioutil.TempFile("", "*.Compare")
			So(err, ShouldBeNil)
			f2, err := ioutil.TempFile("", "*.Compare")
			So(err, ShouldBeNil)

			_, err = f1.Write(kb)
			So(err, ShouldBeNil)
			_, err = f2.Write(kb)
			So(err, ShouldBeNil)

			file1 := fs.File{File: f1}
			file2 := fs.File{File: f2}

			So(file1.Compare(&file2), ShouldEqual, 0)

			f1.Close()
			f2.Close()
		})

		Convey("Should return -1 when other file is broken", func() {
			f1, err := ioutil.TempFile("", "*.Compare")
			So(err, ShouldBeNil)
			f2, err := ioutil.TempFile("", "*.Compare")
			So(err, ShouldBeNil)

			err = f2.Close()
			So(err, ShouldBeNil)
			err = os.Remove(f2.Name())
			So(err, ShouldBeNil)

			file1 := fs.File{File: f1}
			file2 := fs.File{File: f2}

			So(file1.Compare(&file2), ShouldEqual, -1)

			f1.Close()
		})

		Convey("Should return 1 when src file is broken", func() {
			f1, err := ioutil.TempFile("", "*.Compare")
			So(err, ShouldBeNil)
			f2, err := ioutil.TempFile("", "*.Compare")
			So(err, ShouldBeNil)

			err = f1.Close()
			So(err, ShouldBeNil)
			err = os.Remove(f1.Name())
			So(err, ShouldBeNil)

			file1 := fs.File{File: f1}
			file2 := fs.File{File: f2}

			So(file1.Compare(&file2), ShouldEqual, 1)

			f2.Close()
		})

		Convey("Should return ??? when both files are broken", func() {})
	})

	Convey("Unwrap", t, func() {
		srcf := new(os.File)
		file := fs.File{File: srcf}

		f := file.Unwrap()
		typ := fmt.Sprintf("%T", f)

		So(typ, ShouldEqual, "*os.File")
		So(f, ShouldResemble, srcf)
	})

	Convey(".Hash", t, func() {})

	Convey(".Copy", t, func() {})

	Convey(".MarshalJSON", t, func() {})
}
