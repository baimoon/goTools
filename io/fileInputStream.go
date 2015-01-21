package io

import (
	"errors"
	"fmt"
	"io"
	"os"
)

type FileInputStream struct {
	name   string
	buffer []byte
	pos    int
	mark   int
	limit  int
	file   *os.File
}

/**
 * 创建一个指定缓存大小的InputStream对象
 */
func NewFileInputStreamWithSize(fileName string, bufferSize int) (fis *FileInputStream, err error) {
	fis = &FileInputStream{name: fileName, pos: 0, limit: 0, mark: 0}
	fis.file, err = os.Open(fileName)
	if err != nil {
		fmt.Println(err)
		return
	}
	fis.pos = 0
	fis.limit = 0
	fis.buffer = make([]byte, bufferSize)
	return
}

/**
 * 创建一个FileInputStream对象，默认的缓存大小为1024
 */
func NewFileInputStream(fileName string) (fis *FileInputStream, err error) {
	fis, err = NewFileInputStreamWithSize(fileName, 1024)
	return
}

/**
 * 按字节读取数据
 */
func (fis *FileInputStream) Read() int {
	if fis.pos == fis.limit {
		err := fis._read()
		if err != nil {
			panic(err)
		}
	}
	if fis.pos < fis.limit {
		b := fis.buffer[fis.pos]
		fis.pos++
		return int(b)
	}
	return -1
}

/**
 * 读取数据到指定的数组中
 **/
func (fis *FileInputStream) Read2Array(buffer []byte) (size int, err error) {
	size, err = fis.Read2ArrayRange(buffer, 0, len(buffer))
	return
}

/**
 * 读取数据到指定数组的指定范围中
 **/
func (fis *FileInputStream) Read2ArrayRange(buffer []byte, offset int, length int) (size int, err error) {
	if buffer == nil {
		err = errors.New("要复制的目标数组为nil")
		return
	} else if offset < 0 || length < 0 || offset > len(buffer) || length > len(buffer)-offset {
		errMsg := fmt.Sprintf("要复制的目标数组会越界, offset:%d, len(buffer):%d, length:%d", offset, len(buffer), length)
		err = errors.New(errMsg)
		return
	} else if length == 0 {
		return 0, nil
	}
	b := fis.Read()
	if b == -1 {
		return
	}
	buffer[offset] = byte(b)
	i := 1
	for ; i < length; i++ {
		b = fis.Read()
		if b == -1 {
			break
		}
		buffer[offset+i] = byte(b)
	}
	size = i
	return
}

/**
 * 数据读取的内部方法，会从文件中读取数据到内部的buffer中
 **/
func (fis *FileInputStream) _read() error {
	if fis.pos != fis.limit {
		return errors.New("数据产生覆盖")
	}
	n, err := fis.file.Read(fis.buffer)
	if err != nil && err != io.EOF {
		return err
	}
	fis.pos = 0
	fis.limit = n
	return nil
}

func (fis *FileInputStream) Reset() {
	fis.pos = fis.mark
}

/**
 * 标记操作
 * 标记后，在使用Reset方法后，读取数据时
 **/
func (fis *FileInputStream) Mark(pos int) (err error) {
	err = UNSUPPORTED_OPERATION
	return
}

/**
 * 判断是否支持mark操作
 **/
func (fis FileInputStream) MarkSupported() (supported bool) {
	return false
}

func (fis *FileInputStream) Close() {
	fis.file.Close()
}
