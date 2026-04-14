package inter

// Pix API paths use the /pix/v2 prefix on the CDPJ host (see Banco Inter developer references).

type endpoints struct {
	requester *requester
}

func (e endpoints) CreateImmediateCharge(body map[string]interface{}) ([]byte, error) {
	return e.requester.request("/pix/v2/cob", httpMethodPost, nil, body)
}

func (e endpoints) CreateCharge(txid string, body map[string]interface{}) ([]byte, error) {
	params := map[string]string{"txid": txid}

	return e.requester.request("/pix/v2/cob/:txid", httpMethodPut, params, body)
}

func (e endpoints) UpdateCharge(txid string) ([]byte, error) {
	params := map[string]string{"txid": txid}

	return e.requester.request("/pix/v2/cob/:txid", httpMethodPatch, params, nil)
}

func (e endpoints) DetailCharge(txid string) ([]byte, error) {
	params := map[string]string{"txid": txid}

	return e.requester.request("/pix/v2/cob/:txid", httpMethodGet, params, nil)
}

func (e endpoints) ListCharges(inicio, fim string) ([]byte, error) {
	params := map[string]string{
		"inicio": inicio,
		"fim":    fim,
	}

	return e.requester.request("/pix/v2/cob?inicio=:inicio&fim=:fim", httpMethodGet, params, nil)
}

func (e endpoints) PixCreateLocation(body map[string]interface{}) ([]byte, error) {
	return e.requester.request("/pix/v2/loc", httpMethodPost, nil, body)
}

func (e endpoints) PixUnlinkTxidLocation(id string, body map[string]interface{}) ([]byte, error) {
	params := map[string]string{"id": id}

	return e.requester.request("/pix/v2/loc/:id/txid", httpMethodDelete, params, body)
}

func (e endpoints) PixDetailLocation(id string) ([]byte, error) {
	params := map[string]string{"id": id}

	return e.requester.request("/pix/v2/loc/:id", httpMethodGet, params, nil)
}

func (e endpoints) PixLocationList(inicio, fim string) ([]byte, error) {
	params := map[string]string{
		"inicio": inicio,
		"fim":    fim,
	}

	return e.requester.request("/pix/v2/loc?inicio=:inicio&fim=:fim", httpMethodGet, params, nil)
}

func (e endpoints) PixGenerateQRCode(id string) ([]byte, error) {
	params := map[string]string{"id": id}

	return e.requester.request("/pix/v2/loc/:id/qrcode", httpMethodGet, params, nil)
}

func (e endpoints) CreateDueCharge(txid string, body map[string]interface{}) ([]byte, error) {
	params := map[string]string{"txid": txid}

	return e.requester.request("/pix/v2/cobv/:txid", httpMethodPut, params, body)
}

func (e endpoints) PixUpdateDueCharge(txid string, body map[string]interface{}) ([]byte, error) {
	params := map[string]string{"txid": txid}

	return e.requester.request("/pix/v2/cobv/:txid", httpMethodPatch, params, body)
}

func (e endpoints) DetailDueCharge(txid string) ([]byte, error) {
	params := map[string]string{"txid": txid}

	return e.requester.request("/pix/v2/cobv/:txid", httpMethodGet, params, nil)
}

func (e endpoints) PixListDueCharges(params map[string]string) ([]byte, error) {
	return e.requester.request("/pix/v2/cobv", httpMethodGet, params, nil)
}

func (e endpoints) PixDevolution(e2eid, id string, body map[string]interface{}) ([]byte, error) {
	p := map[string]string{
		"e2eid": e2eid,
		"id":    id,
	}

	return e.requester.request("/pix/v2/pix/:e2eid/devolucao/:id", httpMethodPut, p, body)
}

func (e endpoints) PixDetailDevolution(e2eid, id string) ([]byte, error) {
	p := map[string]string{
		"e2eid": e2eid,
		"id":    id,
	}

	return e.requester.request("/pix/v2/pix/:e2eid/devolucao/:id", httpMethodGet, p, nil)
}

func (e endpoints) PixReceivedList(params map[string]string) ([]byte, error) {
	return e.requester.request("/pix/v2/pix", httpMethodGet, params, nil)
}

func (e endpoints) PixDetailReceived(e2eid string) ([]byte, error) {
	params := map[string]string{"e2eid": e2eid}

	return e.requester.request("/pix/v2/pix/:e2eid", httpMethodGet, params, nil)
}

func (e endpoints) PixConfigWebhook(chave string, body map[string]interface{}) ([]byte, error) {
	params := map[string]string{"chave": chave}

	return e.requester.request("/pix/v2/webhook/:chave", httpMethodPut, params, body)
}

func (e endpoints) PixDeleteWebhook(chave string, body map[string]interface{}) ([]byte, error) {
	params := map[string]string{"chave": chave}

	return e.requester.request("/pix/v2/webhook/:chave", httpMethodDelete, params, body)
}

func (e endpoints) PixDetailWebhook(chave string) ([]byte, error) {
	params := map[string]string{"chave": chave}

	return e.requester.request("/pix/v2/webhook/:chave", httpMethodGet, params, nil)
}

func (e endpoints) PixListWebhooks(inicio, fim string) ([]byte, error) {
	params := map[string]string{
		"inicio": inicio,
		"fim":    fim,
	}

	return e.requester.request("/pix/v2/webhook?inicio=:inicio&fim=:fim", httpMethodGet, params, nil)
}

const (
	httpMethodGet    = "GET"
	httpMethodPost   = "POST"
	httpMethodPut    = "PUT"
	httpMethodPatch  = "PATCH"
	httpMethodDelete = "DELETE"
)
