package payment

import (
	"net/http"

	"github.com/TukTuk/core"
)

var PayClient *PaymentClient

type PaymentClient struct {
	Client *http.Client
	Cfg    *core.Config
}

func InitPayment(cfg *core.Config) {
	PayClient = &PaymentClient{
		Client: &http.Client{
			Timeout: cfg.Payment.Timeout,
		},
		Cfg: cfg,
	}
}
