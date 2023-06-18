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

