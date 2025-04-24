package api

import (
	"fmt"

	"github.com/ReanSn0w/go-yandex-delivery/pkg/delivery"
	"github.com/ReanSn0w/go-yandex-delivery/pkg/utils"
	"github.com/ReanSn0w/gokit/pkg/web"
)

func New(environment utils.Environment, client web.HTTPClient, token string) *API {
	return &API{
		environment: environment,
		client:      client,
		token:       token,
	}
}

type API struct {
	environment utils.Environment
	client      web.HTTPClient
	token       string
}

func (a *API) Delivery() *delivery.Delivery {
	return delivery.New(a, a.environment)
}

func (a *API) Request(base, path string) *web.JsonRequest {
	return web.NewJsonRequest(a.client, fmt.Sprintf("%v%v", base, path)).
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", fmt.Sprintf("Bearer %v", a.token))
}
