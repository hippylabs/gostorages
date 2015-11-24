package gostorages

import (
	"testing"
	"io/ioutil"
	"github.com/stretchr/testify/assert"
	"time"
	"os"
)

var swiftConnection = &SwiftConnection{
	UserName: os.Getenv("SWIFT_API_USER"),
	ApiKey: os.Getenv("SWIFT_API_KEY"),
	AuthUrl: os.Getenv("SWIFT_AUTH_URL"),
}

var containerName = "test_container";
var locationName = "test_folder";
var basePath = "/basepath"

var swiftStorage = NewSwiftStorage(
	swiftConnection, 
	containerName, 
	locationName, 
	basePath,
)

var CONTENT = "(╯°□°）╯︵ ┻━┻";


func TestSwiftSaveAndFileOps(t *testing.T) {
	err := swiftStorage.Save("folder3/file1.txt", NewContentFile([]byte(CONTENT)))
	if err != nil {
		t.Fatal(err)
	}

	file, err := swiftStorage.Open("folder3/file1.txt")
	if err != nil {
		t.Fatal(err)
	}

	bytes, err := file.ReadAll();
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, string(bytes), CONTENT)
	assert.Equal(t, int(file.Size()), len(CONTENT))
}

func TestSwiftSize(t *testing.T) {
	size := swiftStorage.Size("folder3/file1.txt")
	assert.Equal(t, int(size), len(CONTENT))
}


func TestSwiftOpen(t *testing.T) {
	file, err := swiftStorage.Open("folder3/file1.txt")
	if err != nil {
		t.Fatal(err)
	}
	
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}
	
	assert.Equal(t, string(bytes), CONTENT)
}

func TestSwiftPath(t *testing.T) {
	path := swiftStorage.Path("folder3/file1.txt")
	assert.Equal(t, path, locationName +"/folder3/file1.txt")
}

func TestSwiftExists(t *testing.T) {
	exist := swiftStorage.Exists("folder3/file1.txt")
	assert.Equal(t, exist, true)
	
	no_exist := swiftStorage.Exists("folder3/file.txt.1")
	assert.Equal(t, no_exist, false)
}

func TestSwiftModifiedTime(t *testing.T) {
	modified, err := swiftStorage.ModifiedTime("folder3/file1.txt")
	if err != nil{
		t.Fatal(err)
	}
	assert.Equal(t, modified.Format(time.RFC822), time.Now().UTC().Format(time.RFC822))

}

func TestSwiftURL(t *testing.T) {
	url := swiftStorage.URL("folder3/file1.txt")
	assert.Equal(t, url, basePath + "/" +locationName + "/folder3/file1.txt")
}

func TestSwiftHasBaseURL(t *testing.T) {
	assert.Equal(t, swiftStorage.HasBaseURL(), true)
}

func TestSwiftDelete(t *testing.T) {
	err := swiftStorage.Delete("folder3/file1.txt")
	if err != nil{
		t.Fatal(err)
	}
	
	exist := swiftStorage.Exists("folder3/file.txt")
	assert.Equal(t, exist, false)
}
