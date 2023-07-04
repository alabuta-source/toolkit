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

func uploadFile(c *gin.Context) {
	waitress, err := gcpKey(c.Request, "go-tool")
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
    
	fileUrl, wErr := waitress.UploadFile(file, fileHeader.Filename, "user")
	if wErr != nil {
		//handle error
		return
	}
	c.JSON(http.StatusOK, fileUrl)
}

func main() {
	router := gin.Default()
	router.POST("/upload", uploadFile)

	err := router.Run(":8080")
	if err != nil {
		//handle error
		return
	}
}
```


