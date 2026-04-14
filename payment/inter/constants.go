package inter

// Version is sent on API requests for traceability.
var Version = "1.0.0"

const (
	// URLSandbox is the Inter Empresas CDPJ API base (sandbox).
	URLSandbox = "https://cdpj-sandbox.partners.uatinter.co"
	// URLProduction is the Inter Empresas CDPJ API base (production).
	URLProduction = "https://cdpj.partners.bancointer.com.br"
)

// DefaultOAuthScope covers Pix cob/cobv and webhooks; override with config key "scope" if your integration uses a subset.
const DefaultOAuthScope = "cob.write cob.read cobv.write cobv.read pix.write pix.read webhook.write webhook.read"
