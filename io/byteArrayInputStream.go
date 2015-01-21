package io

import (
    "errors"
)

type ByteArrayInputStream struct {
    buffer []byte
    pos    int
    limit  int
}

func NewByteArrayInputStream(buffer []byte) (bais *ByteArrayInputStream) {
    bais = &ByteArrayInputStream{}
    bais.buffer = buffer
    bais.pos = 0
    bais.limit = len(bais.buffer)
    return
}

func (bais *ByteArrayInputStream) Read() int {
    if bais.pos + 1 >= len(bais.buffer) {
        return -1
    }
    i := bais.pos
    bais.pos++
    return int(bais.buffer[i])
}

func (bais *ByteArrayInputStream) Read2Array(buffer []byte) (size int, err error) {
    size, err = bais.Read2ArrayRange(buffer, 0, len(buffer))
    return
}

func (bais *ByteArrayInputStream) Read2ArrayRange(buffer []byte, offset int, length int) (size int, err error) {
    if buffer == nil {
        err = errors.New("数组为空")
        return
    } else if offset < 0 || len(buffer) == 0 || length > len(buffer) - offset {
        err = errors.New("复制数组会导致数组越界")
        return
    } else if length == 0 {
        size = 0
        return
    }
    b := bais.Read()
    if int(b) == -1 {
        return 0, nil
    }
    buffer[offset] = byte(b)
    for i:=1; i<length; i++ {
        b = bais.Read()
        if b == -1 {
            break
        }
        buffer[offset + i] = byte(b)
    }
    return
}

