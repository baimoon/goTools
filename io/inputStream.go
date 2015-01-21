package io

import ()

type InputStream interface {
	Read() int
	Read2Array(buffer []byte) (size int, err error)
	Read2ArrayRange(buffer []byte, offset int, length int) (size int, err error)
}
