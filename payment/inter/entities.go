package inter

// DirectChargeResponse mirrors the usual BACEN-style cob response fields used with payment/pix.
type DirectChargeResponse struct {
	Txid        string         `json:"txid"`
	CopyPast    string         `json:"pixCopiaECola"`
	Location    ChargeLocation `json:"loc"`
	Status      string         `json:"status"`
	ChargeValue Value          `json:"valor"`
}

// ChargeLocation describes the cob location payload.
type ChargeLocation struct {
	ID       int64  `json:"id"`
	Location string `json:"location"`
	TipoCob  string `json:"tipoCob"`
}

// Value holds the charge amount string.
type Value struct {
	Original string `json:"original"`
}

// BuildDirectChargeBody builds a BACEN-style immediate charge body (same shape as payment/pix.BuildDirectChargeBody).
func BuildDirectChargeBody(dueDate int, cpf, name, value, pixKey string) map[string]interface{} {
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
		"chave": pixKey,
	}
}

// BuildDueChargeBody builds the main JSON body for cobrança com vencimento (cobv), for use with CreateDueCharge(txid, body).
// dataVencimento must be the due date in YYYY-MM-DD (formato de data do arranjo Pix).
// Informe cpf (11 dígitos) ou cnpj (14 dígitos); deixe vazio o documento que não usar.
func BuildDueChargeBody(dataVencimento, nome, valorOriginal, chavePix, cpf, cnpj string) map[string]interface{} {
	devedor := map[string]interface{}{
		"nome": nome,
	}

	if cnpj != "" {
		devedor["cnpj"] = cnpj
	} else {
		devedor["cpf"] = cpf
	}

	return map[string]interface{}{
		"calendario": map[string]interface{}{
			"dataDeVencimento": dataVencimento,
		},
		"devedor": devedor,
		"valor": map[string]interface{}{
			"original": valorOriginal,
		},
		"chave": chavePix,
	}
}
