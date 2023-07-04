# toolkit

### install
```bash
go get -u github.com/alabuta-source/toolkit@latest
```
### Using token maker to generate token
```go
package main
import "github.com/alabuta-source/toolkit"

func GenerateToken(tokenID string, username string, expirationTime time.Time) (string, error) {
	symmetricKey := "12345678912345678912345678912345"
	tokenMaker := toolkit.NewTokenMaker(symmetricKey)
	
    token, err := tokenMaker.CreateToken(tokenID, username, expirationTime)
    if err != nil {
        return "", err
    }
    return token, nil
}

```

### Using token maker to validate token
```go  
package main
import "github.com/alabuta-source/toolkit"

func ValidateToken(token string) (*toolkit.TokenPayload, error) {
    symmetricKey := "12345678912345678912345678912345"
    tokenMaker := toolkit.NewTokenMaker(symmetricKey)
    
    payload, err := tokenMaker.ValidateToken(token)
    if err != nil {
        return nil, err
    }
    return payload, nil
}
```

### Using Email sender
```go
package main
import "github.com/alabuta-source/toolkit"

func SendEmail(to string, subject string, body string) error {
	config := &toolkit.EmailSenderConfig{
		Password:     "test123",
		From:         "some@test.com",
		ServerConfig: "smtp.gmail.com",
		Port:         123,
    }
	
    emailSender := toolkit.NewEmailSender(config)
   return emailSender.SendEmail(to, subject, body)
}
```

### Using Email sender with template
```go  
package main
import "github.com/alabuta-source/toolkit"

func SendWellComeEmailWithTemplate(to string, subject string, templatePath string, username string, url string) error {
	config := &toolkit.EmailSenderConfig{
		Password:     "test123",
		From:         "some@test.com",
		ServerConfig: "smtp.gmail.com",
		Port:         123,
	}
	
	body := &toolkit.EmailTemplateBody{
		Name: username,
		URL:  url,
	}
    
    emailSender := toolkit.NewEmailSender(config)
  return emailSender.SendEmailWithSimpleTemplate(to, subject, templatePath, body)
}
```

### ### Upload file using GCP Waitress
```go
package main
import (
	"encoding/json"
	"github.com/alabuta-source/toolkit"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
)

//init a GCP waitress manager
func gcpKey(request *http.Request, bucket string) (toolkit.GCPWaitressManager, error) {
	key := new(toolkit.GCPBucketAuthJson)
	bytes, err := os.ReadFile("key.json")
	if err != nil {
		return nil, err
	}
	_ = json.Unmarshal(bytes, key)
	return toolkit.NewGCPWaitress(bucket, request, key)
}

//uploadFile example, you can send any file extension
func uploadFile(c *gin.Context) {
	waitress, err := gcpKey(c.Request, "your bucket name")
	if err != nil {
		//handle error
		return
	}

	fileHeader, er := c.FormFile("file")
	if er != nil {
		//handle error
		return
	}

	file, e := fileHeader.Open()
	if e != nil {
		//handle error
	}
    
	// this method will return the object url
	// example: https://storage.cloud.google.com/your-bucket/prefix/file.png
	fileUrl, wErr := waitress.UploadFile(file, fileHeader.Filename, "the prefix")
	if wErr != nil {
		//handle error
		return
	}
	c.JSON(http.StatusOK, fileUrl)
}

// List all bucket files example.
// Just send the prefix to get all files
// example: users, and you'll get all users images
func getFiles(c *gin.Context) {
	waitress, err := gcpKey(c.Request, "your bucket name")
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	//response example
	/*
		[
			"https://storage.cloud.google.com/your-bucket/prefix/file.png"
			"https://storage.cloud.google.com/your-bucket/prefix/file2.png"
		]
	*/
	images, er := waitress.ListFiles("the prefix")
	if er != nil {
		c.JSON(http.StatusBadRequest, er)
		return
	}
	c.JSON(http.StatusOK, images)
}

// This method is to delete a single file
func delete(c *gin.Context) {
	waitress, err := gcpKey(c.Request, "your bucket name")
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	file := "https://storage.cloud.google.com/your-bucket/prefix/file2.png"
	if er := waitress.DeleteFile(file); er != nil {
		c.JSON(http.StatusBadRequest, er.Error())
		return
	}
	c.JSON(http.StatusOK, "done")
}

func main() {
	router := gin.Default()
	router.POST("/file", uploadFile)
	router.GET("/files", getFiles)
	router.DELETE("/file", delete)

	err := router.Run(":8080")
	if err != nil {
		//handle error
		return
	}
}
```


