package delivery

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/ReanSn0w/go-yandex-delivery/pkg/utils"
	"github.com/ReanSn0w/gokit/pkg/web"
)

const (
	production  = "https://b2b-authproxy.taxi.yandex.net/api/b2b/platform"
	development = "https://b2b.taxi.tst.yandex.net/api/b2b/platform"
)

// Создает новый экземпляр структуры для работы с API доставки на следующий день
func New(api utils.API, env utils.Environment) *Delivery {
	base := production
	if env == utils.Development {
		base = development
	}

	return &Delivery{api: api, base: base}
}

type Delivery struct {
	api  utils.API
	base string
}

func (d *Delivery) request(path string) *web.JsonRequest {
	return d.api.Request(d.base, path)
}

// GetPredictedPrice возвращает предварительную оценку стоимости доставки
// is_oversized - Флаг КГТ
func (d *Delivery) GetPredictedPrice(isOversized bool, req PredictPriceRequest) (*PredictPriceResponse, error) {
	res := PredictPriceResponse{}
	err := d.request("/pricing-calculator").
		SetQuery("is_oversized", fmt.Sprint(isOversized)).
		SetBody(req).
		Do(&res)
	return &res, err
}

// GetDeliveryIntervals возвращает интервалы доставки
// is_oversized - Флаг КГТ
func (d *Delivery) GetDeliveryIntervals(isOversized bool, lastMilePolicy LastMilePolicy, req DeliveryIntervalsRequest) (*DeliveryIntervalsResponse, error) {
	res := DeliveryIntervalsResponse{}
	err := d.request("/offers/info").
		SetQuery("is_oversized", fmt.Sprint(isOversized)).
		SetQuery("last_mile_policy", string(lastMilePolicy)).
		SetQuery("send_unix", "true").
		SetBody(req).
		Do(&res)
	return &res, err
}

// GetLocationID возвращает идентификатор населенного пункта
func (d *Delivery) GetLocationID(address string) (*LocationIDResponse, error) {
	res := LocationIDResponse{}
	err := d.request("/location/id").
		SetBody(map[string]any{"location": address}).
		Do(&res)
	return &res, err
}

// GetDeliveryPoints возвращает список точек самовывоза и ПВЗ
func (d *Delivery) GetDeliveryPoints(req DeliveryPointsRequest) (*DeliveryPointsResponse, error) {
	res := DeliveryPointsResponse{}
	err := d.request("/ickup-points/list").
		SetBody(req).
		Do(&res)
	return &res, err
}

// CreateOffer создает заявку на доставку
func (d *Delivery) CreateOffer(req CreateOfferRequest) (*CreateOfferResponse, error) {
	res := CreateOfferResponse{}
	err := d.request("/offers/create").
		SetQuery("send_unix", "true").
		SetBody(req).
		Do(&res)
	return &res, err
}

// ConfirmOffer подтверждает заявку на доставку
func (d *Delivery) ConfirmOffer(offerID string) (*ConfirmOfferResponse, error) {
	resp := ConfirmOfferResponse{}
	err := d.request("/offers/confirm").
		SetBody(map[string]any{"offer_id": offerID}).
		Do(&resp)
	return &resp, err
}

// GetRequestInfo возвращает информацию о заявке на доставку
func (d *Delivery) GetRequestInfo(requestID string, slim bool) (*GetRequestInfoResponse, error) {
	resp := GetRequestInfoResponse{}
	err := d.request("/request/info").
		SetQuery("request_id", requestID).
		SetQuery("slim", strconv.FormatBool(slim)).
		Do(&resp)
	return &resp, err
}

// GetRequestsInfo получает информацию о заявках во временном интервале
func (d *Delivery) GetRequestsInfo(from, to time.Time, requestsIds ...string) (*GetRequestsInfoResponse, error) {
	res := GetRequestsInfoResponse{}
	err := d.request("/requests/info").
		SetBody(map[string]any{
			"from":        from.Format(time.RFC3339),
			"to":          to.Format(time.RFC3339),
			"request_ids": requestsIds,
		}).Do(&res)
	return &res, err
}

// GetRequestActualInfo получeние актуальной информации о доставке
func (d *Delivery) GetRequestActualInfo(requestID string) (*GetRequestActualInfoResponse, error) {
	res := GetRequestActualInfoResponse{}
	err := d.request("/request/actual_info").
		SetQuery("request_id", requestID).
		Do(&res)
	return &res, err
}

// EditRequestInfo редактирует информацию о заказе
func (d *Delivery) EditRequestInfo(req EditRequestInfoRequest) (*EditRequestInfoResponse, error) {
	res := EditRequestInfoResponse{}
	err := d.request("/request/edit").
		SetBody(req).
		Do(&res)
	return &res, err
}

// GetRequestRedeliveryOptions получает интервалы доставки для нового места получения заказа
func (d *Delivery) GetRequestRedeliveryOptions(req GetRequestRedeliveryOptionsRequest) (*GetRequestRedeliveryOptionsResponse, error) {
	res := GetRequestRedeliveryOptionsResponse{}
	err := d.request("/request/redelivery_options").
		SetBody(req).
		Do(&res)
	return &res, err
}

// GetRequestHistory получает историю заявки
func (d *Delivery) GetRequestHistory(requestID string) (*GetRequestHistoryResponse, error) {
	resp := GetRequestHistoryResponse{}
	err := d.request("/request/history").
		SetQuery("request_id", requestID).
		Do(&resp)
	return &resp, err
}

// CancelRequest отменяет заявку
func (d *Delivery) CancelRequest(requestID string) (*CancelRequestResponse, error) {
	resp := CancelRequestResponse{}
	err := d.request("/request/cancel").
		SetBody(map[string]any{"request_id": requestID}).
		Do(&resp)
	return &resp, err
}

// CreateRequest создает новый заказ
func (d *Delivery) CreateRequest(req CreateRequestRequest) (*CreateRequestResponse, error) {
	resp := CreateRequestResponse{}
	err := d.request("/request/create").
		SetHeader("Accept-Language", "ru").
		SetQuery("send_unix", "true").
		SetBody(req).
		Do(&resp)
	return &resp, err
}

// EditRequestPlaces редактирование грузомест заказа
func (d *Delivery) EditRequestPlaces(req EditRequestPlacesRequest) (*EditRequestPlacesResponse, error) {
	res := EditRequestPlacesResponse{}
	err := d.request("/request/places/edit").
		SetBody(req).
		Do(&res)
	return &res, err
}

// GetEditRequestStatus получение статуса запроса на редактирование
func (d *Delivery) GetEditRequestStatus(taskID string) (*GetEditRequestStatusResponse, error) {
	resp := GetEditRequestStatusResponse{}
	err := d.request("/request/edit/status").
		SetQuery("editing_task_id", taskID).
		Do(&resp)
	return &resp, err
}

// EditRequestItems редактирование товаров заказа
func (d *Delivery) EditRequestItems(req EditRequestItemsRequest) (*EditRequestItemsResponse, error) {
	res := EditRequestItemsResponse{}
	err := d.request("/request/items-instances/edit").
		SetBody(req).
		Do(&res)
	return &res, err
}

// GenerateRequestLabels генерация транспортных ярлыков
func (d *Delivery) GenerateRequestLabels(req GenerateRequestLabelsRequest) (io.ReadCloser, error) {
	return nil, errors.New("not implemented")
}

// GetRequestHandoverAct получение акта приема/передачи отгрузки
func (d *Delivery) GetRequestHandoverAct(requestIds ...string) (io.ReadCloser, error) {
	return nil, errors.New("not implemented")
}
