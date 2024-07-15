package helpers

import "mime/multipart"

var ProcessFileWrapper = ProcessFile

func SetProcessFileWrapper(mock func(file multipart.File) (string, error)) {
    ProcessFileWrapper = mock
}