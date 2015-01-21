package io

import "errors"

var (
	NO_MORE_DATA          = errors.New("没有更多的数据错误")
	NULL_POINTER          = errors.New("空指针错误")
	UNSUPPORTED_OPERATION = errors.New("不支持的操作错误")
)
