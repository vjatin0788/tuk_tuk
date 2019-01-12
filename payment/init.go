package payment

import (
	"net/http"
	"time"

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
			Timeout: 10 * time.Second,
		},
		Cfg: cfg,
	}
}
