package newrelic

import (
	"github.com/newrelic/go-agent/v3/integrations/logcontext-v2/nrlogrus"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/sirupsen/logrus"
)

type ConfigOptions struct {
	Licence    string
	AppName    string
	LogEnabled bool
}

func NewApplication(options ConfigOptions) (*newrelic.Application, error) {
	return newrelic.NewApplication(
		newrelic.ConfigAppName(options.AppName),
		newrelic.ConfigLicense(options.Licence),
		newrelic.ConfigAppLogForwardingEnabled(options.LogEnabled),
	)
}

func Logger(app *newrelic.Application) *logrus.Logger {
	nrlogrusFormatter := nrlogrus.NewFormatter(app, &logrus.TextFormatter{})
	logger := logrus.New()
	logger.SetFormatter(nrlogrusFormatter)
	return logger
}
