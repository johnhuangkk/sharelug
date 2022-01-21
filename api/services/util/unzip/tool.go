package unzip

import (
	"bytes"
	"errors"
	"github.com/yeka/zip"
	"io/ioutil"
)

func UnzipDataWithPassword(in []byte, password string) ([]byte,error) {
	buf := bytes.NewReader(in)
	reader,err := zip.NewReader(buf,buf.Size())
	if err != nil {
		return nil, err
	}

	if len(reader.File) == 0 {
		return nil, errors.New("File not found")
	}
	f := reader.File[0]

	if f.IsEncrypted() {
		f.SetPassword(password)
	}


	r, err := f.Open()
	if err != nil {
		return nil,err
	}

	unzipData, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return unzipData, nil
}
