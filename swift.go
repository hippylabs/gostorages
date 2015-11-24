package gostorages

import (
	"github.com/ncw/swift"
	"time"
	"mime"
	"io/ioutil"
)

type SwiftStorage struct {
	*BaseStorage
	Connection *swift.Connection
	ContainerName string
}

type SwiftStorageFile struct {
	*swift.ObjectOpenFile
}

type SwiftConnection swift.Connection

func NewSwiftStorage(swiftConnection *SwiftConnection, ContainerName string, location string, baseURL string) Storage {

	Connection := swift.Connection(*swiftConnection)
	storage := &SwiftStorage{
		NewBaseStorage(location, baseURL),
		&Connection,
		ContainerName,
	}

	return storage
}

func (f *SwiftStorageFile) ReadAll() ([]byte, error) {
	return ioutil.ReadAll(f)
}

func (f *SwiftStorageFile) Size() int64 {
	length, err := f.ObjectOpenFile.Length()
	if err != nil {
		return 0
	}
	return length
}

func (f *SwiftStorage) Delete(filepath string) error {
	return f.Connection.ObjectDelete(f.ContainerName, f.Location + filepath)
}

func (f *SwiftStorage) SaveWithContentType(filepath string, file File, contentType string) error {
	bytes, err := file.ReadAll();
	if err != nil{
		return err
	}
	err = f.Connection.ObjectPutBytes(f.ContainerName, f.Location + filepath, bytes, contentType);
	if err != nil{
		return err
	}
	return nil
}

func (f *SwiftStorage) Save(filepath string, file File) error {
	return f.SaveWithContentType(filepath, file, mime.TypeByExtension(filepath))
}

func (f *SwiftStorage) Exists(filepath string) bool {
	_, err := f.ObjectMetadata(filepath);	
	if err != nil {
		return false
	}
	return true
}

func (f *SwiftStorage) Open(filepath string) (File, error) {
	
	file, _, err := f.Connection.ObjectOpen(f.ContainerName, f.Location + filepath, false, swift.Headers{})

	if err != nil {
		return nil, err
	}

	return &SwiftStorageFile{
		file,
		}, nil
}

func (f *SwiftStorage) ObjectMetadata(filepath string) (swift.Object, error) {
	info, _, err := f.Connection.Object(f.ContainerName, f.Location + filepath)
	if err != nil {
		return swift.Object{}, err
	}
	
	return info, nil
}


func (f *SwiftStorage) ModifiedTime(filepath string) (time.Time, error) {
	metadata, err := f.ObjectMetadata(filepath);	
	if err != nil {
		return time.Time{}, err
	}
	return metadata.LastModified, nil
}

func (f *SwiftStorage) Size(filepath string) int64 {
	metadata, _	 := f.ObjectMetadata(filepath);	
	return metadata.Bytes;
}