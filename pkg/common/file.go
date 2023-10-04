package common

import (
	"errors"
	"io"
	"os"
)

func NewFile(path string) File {
	desc, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, os.ModeAppend)
	if err != nil {
		panic(err)
	}
	return File{
		Path:     path,
		Desc:     desc,
		Content:  nil,
		IsOpened: true,
	}
}

var ErrClosedFile error = errors.New("File is closed")

func (f *File) Close() {
	if f.IsOpened {
		if err := f.Desc.Close(); err != nil {
			panic(err)
		}
		f.IsOpened = false
	}
}

func (f *File) ReadBytes() ([]byte, error) {
	if !f.IsOpened {
		return nil, ErrClosedFile
	}
	f.Desc.Seek(0, 0)
	bytes, err := io.ReadAll(f.Desc)
	if err != nil {
		panic(err)
	}
	f.Content = bytes
	return bytes, nil
}

func (f File) WriteBytes(bytes []byte) error {
	if !f.IsOpened {
		return ErrClosedFile
	}
	if err := f.Desc.Truncate(0); err != nil {
		panic(err)
	}
	f.Append(bytes)
	return nil
}

func (f File) WriteString(content string) error {
	if !f.IsOpened {
		return ErrClosedFile
	}
	if err := f.Desc.Truncate(0); err != nil {
		panic(err)
	}
	if _, err := f.Desc.WriteString(content); err != nil {
		panic(err)
	}
	return nil
}

func (f File) Append(bytes []byte) error {
	if !f.IsOpened {
		return ErrClosedFile
	}
	if _, err := f.Desc.Write(bytes); err != nil {
		panic(err)
	}
	return nil
}
