package main

import (
	"encoding/json"
	"fmt"
	"github.com/alabuta-source/toolkit/payment/pix"
	"log"
	"path"
	"runtime"
)

var Credentials = map[string]interface{}{
	"client_id":     "",
	"client_secret": "",
	"sandbox":       true,
	"timeout":       4,
	"CA":            fmt.Sprintf("%s/toolkit/payment/certs/sand.crt.pem", getRootDir()),
	"Key":           fmt.Sprintf("%s/toolkit/payment/certs/sand.key.pem", getRootDir()),
}

func main() {
	client := pix.NewEfiPay(Credentials)

	body := pix.BuildDirectChargeBody(3600, "12345678000", "user test", "00.01")
	resp, err := client.CreateImmediateCharge(body)

	chargeResponse := pix.DirectChargeResponse{}
	_ = json.Unmarshal(resp, &chargeResponse)
	if err != nil {
		log.Println(err)
	}

	qrcodeBytes, er := client.PixGenerateQRCode("3")
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
