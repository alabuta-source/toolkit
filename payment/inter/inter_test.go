package inter

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestBuildDueChargeBody_CPFVsCNPJ(t *testing.T) {
	t.Parallel()

	cpfBody := BuildDueChargeBody("2025-12-31", "Fulano", "10.00", "chave-pix", "12345678901", "")
	raw, _ := json.Marshal(cpfBody)
	if !strings.Contains(string(raw), `"cpf"`) || strings.Contains(string(raw), `"cnpj"`) {
		t.Fatalf("expected cpf only: %s", raw)
	}

	cnpjBody := BuildDueChargeBody("2025-12-31", "Empresa SA", "100.00", "chave-pix", "", "12345678000195")
	raw2, _ := json.Marshal(cnpjBody)
	if !strings.Contains(string(raw2), `"cnpj"`) {
		t.Fatalf("expected cnpj: %s", raw2)
	}
}

func TestDigitsOnly(t *testing.T) {
	t.Parallel()

	if got := digitsOnly("12-34"); got != "1234" {
		t.Fatalf("got %q", got)
	}

	if got := digitsOnly(""); got != "" {
		t.Fatalf("got %q", got)
	}
}

func TestAppendQuery(t *testing.T) {
	t.Parallel()

	out := appendQuery("/pix/v2/cob", map[string]string{
		"inicio": "2024-01-01",
		"fim":    "2024-01-02",
	})

	if !strings.Contains(out, "inicio=") || !strings.Contains(out, "fim=") {
		t.Fatalf("unexpected query %q", out)
	}
}

func TestGetRouteSubstitutesPathParams(t *testing.T) {
	t.Parallel()

	params := map[string]string{"txid": "abc-1", "extra": "z"}
	path := getRoute("/pix/v2/cob/:txid", params)
	if path != "/pix/v2/cob/abc-1" {
		t.Fatalf("path %q", path)
	}

	if params["extra"] != "z" {
		t.Fatalf("query param should remain: %#v", params)
	}

	if _, ok := params["txid"]; ok {
		t.Fatal("txid should be removed after substitution")
	}
}

func TestOAuthAndPixRequest(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	mux.HandleFunc("/oauth/v2/token", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("token: want POST, got %s", r.Method)
			w.WriteHeader(http.StatusMethodNotAllowed)

			return
		}

		if err := r.ParseForm(); err != nil {
			t.Errorf("parse form: %v", err)
			w.WriteHeader(http.StatusBadRequest)

			return
		}

		if r.Form.Get("grant_type") != "client_credentials" {
			t.Errorf("grant_type: %q", r.Form.Get("grant_type"))
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"access_token":"tok1","expires_in":3600,"token_type":"Bearer"}`))
	})

	var sawAuth string
	mux.HandleFunc("/pix/v2/cob/tx-1", func(w http.ResponseWriter, r *http.Request) {
		sawAuth = r.Header.Get("Authorization")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ok":true}`))
	})

	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)

	hc := srv.Client()
	a := &auth{
		clientID:     "id",
		clientSecret: "secret",
		scope:        DefaultOAuthScope,
		baseURL:      srv.URL,
		httpClient:   hc,
	}

	r := &requester{
		auth:           a,
		baseURL:        srv.URL,
		timeoutSeconds: 30,
		httpClient:     hc,
	}

	c := &Client{endpoints: endpoints{requester: r}}
	out, err := c.DetailCharge("tx-1")
	if err != nil {
		t.Fatal(err)
	}

	if string(out) != `{"ok":true}` {
		t.Fatalf("body %s", out)
	}

	if sawAuth != "Bearer tok1" {
		t.Fatalf("Authorization %q", sawAuth)
	}
}
