package efipay

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
