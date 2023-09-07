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

func BuildDirectChargeBody(dueDate int, cpf string, name string, value string) map[string]interface{} {
	return map[string]interface{}{

		"calendario": map[string]interface{}{
			"expiracao": dueDate,
		},
		"devedor": map[string]interface{}{

			"cpf":  cpf,
			"nome": name,
		},
		"valor": map[string]interface{}{

			"original": value,
		},
		"chave": "235d0898-d5e0-419e-97f2-ceb3017751f7",
	}
}
