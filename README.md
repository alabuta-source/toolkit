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

### Init a new GCP Waitress 
```go
package main
import "github.com/alabuta-source/toolkit"

func GCPWaitress(c *gin.Context) (toolkit.GCPWaitressManager, error) {
    bytes, err := os.ReadFile("foo/key.json")
    if err != nil {
        return nil, err
    }
    return toolkit.NewGCPWaitress("bucket-name", c.Request, bytes)
}
```

### Upload file using GCP Waitress
```go
func UploadFile(c *gin.Context) {
    file, fileHeader, err := c.Request.FormFile("file")
    if err != nil {
       return nil, err
    }
   
    waitress, err := GCPWaitress(c)
    if err != nil {
       c.JSON(status, err)
       return
    }
   
   // params: file, fileHeader, prefix and cacheControl
   //if the cacheControl is empty, the default value will be used
   filename, err := waitress.SaveFile(file, fileHeader, "folder/", "")
   if err != nil {
      c.JSON(status, err)
      return
   }
   c.JSON(http.Status.ok, fmt.Sprintf("%s uploaded", filename))
}
```


