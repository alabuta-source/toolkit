package main

import (
	"encoding/json"
	"fmt"
	"github.com/alabuta-source/toolkit/payment/src/efipay/pix"
	"log"
	"path"
	"runtime"
)

var Credentials = map[string]interface{}{
	"client_id":     "",
	"client_secret": "",
	"sandbox":       true,
	"timeout":       4,
	"CA":            fmt.Sprintf("%s/efi-payment/src/certs/sand.crt.pem", getRootDir()),
	"Key":           fmt.Sprintf("%s/efi-payment/src/certs/sand.key.pem", getRootDir()),
}

func main() {

	body := pix.BuildDirectChargeBody(3600, "12345678000", "user test", "00.01")

	client := pix.NewEfiPay(Credentials)
	resp, err := client.CreateImmediateCharge(body)

	chargeResponse := pix.DirectChargeResponse{}
	_ = json.Unmarshal(resp, &chargeResponse)
	if err != nil {
		log.Println(err)
	}

	qrcodeBytes, er := client.PixGenerateQRCode(chargeResponse.WithQrCodeParam())
	code := pix.QrCodePix{}
	_ = json.Unmarshal(qrcodeBytes, &code)
	if er != nil {
		log.Println(er)
	}
	fmt.Println(code)
}

func getRootDir() string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Panic("root dir")
	}

	pwd := path.Dir(filename)
	return path.Join(pwd, "..", "..")
}
