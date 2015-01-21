package io

import (
//"fmt"
//"errors"
//"os"
//"io"
//"uio"
)

type FileReader struct {
	buffer []byte
	fis    *FileInputStream
	pos    int
	limit  int
}

func NewFileReaderWithCapacity(fileName string, bufferSize int) (fr *FileReader, err error) {
	fr = &FileReader{}
	fr.fis, err = NewFileInputStreamWithSize(fileName, bufferSize)
	if err != nil {
		return nil, err
	}
	fr.buffer = make([]byte, 1)
	fr.pos = 0
	fr.limit = 0
	return
}

func NewFileReader(fileName string) (fr *FileReader, err error) {
	return NewFileReaderWithCapacity(fileName, 1024)
}

func (fr *FileReader) ReadLine() (str string, err error) {
	buf := make([]byte, 0, 1024)
	for {
		c := fr.fis.Read()
		if c == -1 {
			err = NO_MORE_DATA
			return
		}
		if c == int('\r') {
			otherC := fr.fis.Read()
			if otherC == int('\n') {
				break
			}
			buf = append(buf, byte(c))
			buf = append(buf, byte(otherC))
		}
		if c == int('\n') {
			break
		}
		buf = append(buf, byte(c))
	}
	str = string(buf)
	return
}

func (fr *FileReader) Close() {
	fr.fis.Close()
}
