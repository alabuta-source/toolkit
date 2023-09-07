package main

import (
	"encoding/json"
	"fmt"
	"github.com/alabuta-source/toolkit/payment/src/efipay"
	"github.com/alabuta-source/toolkit/payment/src/efipay/pix"
	"log"
	"path"
	"runtime"
)

var Credentials = map[string]interface{}{
	"client_id":     "Client_Id_551a23cc9ef9d95008161b1cc3272d666ffc7d2a",
	"client_secret": "Client_Secret_01ecb0b788e104bdad842c07c1e9b67e7c526243",
	"sandbox":       true,
	"timeout":       20,
	"CA":            fmt.Sprintf("%s/efi-payment/src/certs/sand.crt.pem", getRootDir()), //caminho da chave publica
	"Key":           fmt.Sprintf("%s/efi-payment/src/certs/sand.key.pem", getRootDir()), //caminho da chave privada
}

func main() {

	body := map[string]interface{}{

		"calendario": map[string]interface{}{
			"expiracao": 3600,
		},
		"devedor": map[string]interface{}{

			"cnpj": "12345678000195",
			"nome": "Empresa de Serviços SA",
		},
		"valor": map[string]interface{}{

			"original": "00.01",
		},
		"chave": "235d0898-d5e0-419e-97f2-ceb3017751f7",
	}

	client := pix.NewEfiPay(Credentials)
	resp, err := client.CreateImmediateCharge(body)

	chargeResponse := efipay.DirectChargeResponse{}
	_ = json.Unmarshal(resp, &chargeResponse)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(resp)
}

func getRootDir() string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Panic("root dir")
	}

	pwd := path.Dir(filename)
	return path.Join(pwd, "..", "..")
}
