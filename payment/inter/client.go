package inter

import (
	"fmt"
)

// Client exposes Pix Cobrança methods for Banco Inter (CDPJ API).
type Client struct {
	endpoints
}

// NewInter builds a Client from the same style of config map used by payment/pix.NewEfiPay.
//
// Required keys: client_id, client_secret, CA (path to PEM certificate), Key (path to PEM private key), sandbox (bool), timeout (int seconds).
//
// Optional keys: scope (string OAuth scope, default DefaultOAuthScope), conta_corrente (string, digits only, sent as x-conta-corrente when set).
func NewInter(configs map[string]interface{}) (*Client, error) {
	clientID, _ := configs["client_id"].(string)
	clientSecret, _ := configs["client_secret"].(string)
	ca, _ := configs["CA"].(string)
	key, _ := configs["Key"].(string)
	sandbox, _ := configs["sandbox"].(bool)

	if clientID == "" || clientSecret == "" {
		return nil, fmt.Errorf("inter: client_id and client_secret are required")
	}

	if ca == "" || key == "" {
		return nil, fmt.Errorf("inter: CA and Key (certificate paths) are required")
	}

	timeout, ok := configs["timeout"].(int)
	if !ok {
		return nil, fmt.Errorf("inter: config \"timeout\" must be an int (seconds)")
	}

	scope, _ := configs["scope"].(string)
	contaCorrente := stringFromConfig(configs, "conta_corrente", "contaCorrente", "x_conta_corrente", "x-conta-corrente")

	req, err := newRequester(clientID, clientSecret, ca, key, sandbox, timeout, scope, contaCorrente, nil)
	if err != nil {
		return nil, err
	}

	return &Client{endpoints{requester: req}}, nil
}

func stringFromConfig(m map[string]interface{}, keys ...string) string {
	for _, k := range keys {
		if v, ok := m[k].(string); ok && v != "" {
			return v
		}
	}

	return ""
}
