package pix

import "strconv"

type DirectChargeResponse struct {
	Txid        string         `json:"txid"`
	Location    ChargeLocation `json:"loc"`
	Status      string         `json:"status"`
	ChargeValue Value          `json:"valor"`
}

type ChargeLocation struct {
	ID       int64  `json:"id"`
	Location string `json:"location"`
	TipoCob  string `json:"tipoCob"`
}

type Value struct {
	Original string `json:"original"`
}

type QrCodePix struct {
	QrCode       string `json:"qrCode"`
	ImagemQrcode string `json:"imagemQrcode"`
}

func (pixCharge *DirectChargeResponse) WithQrCodeParam() string {
	return strconv.FormatInt(pixCharge.Location.ID, 10)
}
