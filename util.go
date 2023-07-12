package toolkit

import (
	"fmt"
	"github.com/google/uuid"
	"math/rand"
	"strings"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func randomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func randomOwner() string {
	return randomString(6)
}

func generateUUID() string {
	return uuid.New().String()
}

func formatErr(msg string, args ...interface{}) string {
	return fmt.Sprintf(msg, args...)
}

func cutSpaces(value string) string {
	return strings.Replace(value, " ", "", -1)
}

func removeBucketName(path, bucket string) string {
	return strings.Replace(path, fmt.Sprintf("/%s/", bucket), "", -1)
}

func welcomeTemplate() string {
	return ""
}

func verifyEmailTemplate() string {
	return ""
}

func resetPassTemplate() string {
	return `<!DOCTYPE html>
<html>

<head>
    <title>test</title>

    <link rel="preconnect" href="https://fonts.googleapis.com">
    <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
    <link href="https://fonts.googleapis.com/css2?family=Inter:wght@100;400;700&display=swap" rel="stylesheet">

    <style type="text/css">

        main {
            display: flex;
            flex-direction: column;
            align-items: center;
        }

        .wrapper {
          background: #8CAA8E;
          min-height: 400px;
          padding: 100px 30px;
        }

        h1 {
            color: #1A2229;
            text-align: center;
            font-family: Inter;
            font-size: 24px;
            font-style: normal;
            font-weight: 700;
            line-height: normal;
        }

        p {
            color: #1A2229;
            text-align: center;
            font-family: Inter;
            font-size: 14px;
            font-style: normal;
            font-weight: 400;
            line-height: normal;
            margin-top: 80px;
        }

        .destaque {
            color: #6D8D6F;
        }

        .name {
          margin-bottom: 20px;
        }

    </style>
</head>

<body>
    <main>
      <div class="wrapper">

        <h1 class="name">Olá, {{ .Name }}</h1>
        <h1>Use o código <span class="destaque">{{ .Message }}</span> para finalizar seu cadastro Bionicblack.</h1>

        <p>
            Copie o código, retorne ao site e cole no espaço indicado.
        </p>
      </div>

    </main>
</body>

</html>`
}

//        <div class="code-div">
//            <span class="code-span">{{ .Message }}</span>
//        </div>

//        .code-span {
//            color: rgba(0, 0, 0, 0.70);
//            text-align: center;
//            font-family: Inter;
//            font-size: 16px;
//            font-style: normal;
//            font-weight: 700;
//            line-height: normal;
//        }

//
