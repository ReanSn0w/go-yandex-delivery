package delivery_test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/ReanSn0w/go-yandex-delivery/pkg/api"
	"github.com/ReanSn0w/go-yandex-delivery/pkg/delivery"
	"github.com/ReanSn0w/go-yandex-delivery/pkg/utils"
	"github.com/ReanSn0w/gokit/pkg/app"
	"github.com/ReanSn0w/gokit/pkg/tool"
	"github.com/ReanSn0w/gokit/pkg/web/httpdebug"
	"github.com/go-pkgz/lgr"
	"github.com/stretchr/testify/assert"
)

var (
	l lgr.L
	d *delivery.Delivery

	opts = struct {
		app.Debug

		Token string `long:"token" env:"TOKEN" description:"API Token"`
	}{}
)

func init() {
	log, err := app.LoadConfiguration("Delivery Package", "test", &opts)
	if err != nil {
		panic(err)
	}

	l = log

	httpdebug := httpdebug.New(log, http.DefaultClient, true)
	d = api.New(utils.Development, httpdebug, opts.Token).Delivery()
}

func TestDelivery_GetPredictedPrice(t *testing.T) {
	cases := []struct {
		Name        string
		IsOversized bool
		Request     delivery.PredictPriceRequest
		HasError    bool
	}{
		{
			Name:        "Yandex Delivery Example",
			HasError:    false,
			IsOversized: false,
			Request: delivery.PredictPriceRequest{
				TotalWeight:        240,
				TotalAssessedPrice: 0,
				ClientPrice:        0,
				Tariff:             delivery.LMP_TimeInterval,
				PaymentMethod:      delivery.PM_AlreadyPaid,
				Source: delivery.Source{
					PlatformStationID: "fbed3aa1-2cc6-4370-ab4d-59c5cc9bb924",
				},
				Destination: delivery.Destination{
					Address: "Москва, Вернадского пр-кт, 91к2, кв. 171",
				},
				Places: []delivery.Place{
					{
						PhysicalDims: delivery.PhysicalDims{
							WeightGross:      240,
							Dx:               5,
							Dy:               10,
							Dz:               20,
							PredefinedVolume: 1000,
						},
					},
				},
			},
		},
		{
			Name:        "Yandex Delivery Example",
			HasError:    true,
			IsOversized: false,
			Request: delivery.PredictPriceRequest{
				TotalWeight:        240,
				TotalAssessedPrice: 0,
				ClientPrice:        0,
				Tariff:             delivery.LMP_TimeInterval,
				PaymentMethod:      delivery.PM_AlreadyPaid,
				Source: delivery.Source{
					PlatformStationID: "fbed3aa1-2cc6-4370-ab4d-59c5cc9bb924",
				},
				Destination: delivery.Destination{
					Address: "Санкт-Петербург, Большая Монетная улица, 1к1А",
					IntervalUTC: &delivery.IntervalUTC{
						From: time.Now().AddDate(0, 0, 2).Add(time.Hour),
						To:   time.Now().AddDate(0, 0, 7).Add(8 * time.Hour),
					},
				},
				Places: []delivery.Place{
					{
						PhysicalDims: delivery.PhysicalDims{
							WeightGross:      240,
							Dx:               5,
							Dy:               10,
							Dz:               20,
							PredefinedVolume: 1000,
						},
					},
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			resp, err := d.GetPredictedPrice(c.IsOversized, c.Request)
			if c.HasError {
				assert.NotNil(t, err)
				assert.Nil(t, resp)
				return
			}

			assert.Nil(t, err, "%s", err)
			assert.False(t, resp.Error)
			assert.NotEmpty(t, resp.PricingTotal)

			if c.Request.PaymentMethod != delivery.PM_AlreadyPaid {
				assert.NotEmpty(t, resp.PricingCommissionOnDeliveryPaymentAmount)
				assert.NotEmpty(t, resp.PricingCommissionOnDeliveryPayment)
				assert.NotEmpty(t, resp.Pricing)
			}
		})
	}
}

func TestDelivery_GetDeliveryIntervals(t *testing.T) {
	cases := []struct {
		Name           string
		IsOversized    bool
		LastMilePolicy delivery.LastMilePolicy
		Request        delivery.DeliveryIntervalsRequest
		HasError       bool
	}{
		{
			Name:           "Yandex Delivery Time Interval",
			IsOversized:    false,
			LastMilePolicy: delivery.LMP_TimeInterval,
			Request: delivery.DeliveryIntervalsRequest{
				Source: delivery.Source{PlatformStationID: "fbed3aa1-2cc6-4370-ab4d-59c5cc9bb924"},
				Destination: delivery.Destination{
					Address: "Москва, Вернадского пр-кт, 91к2",
				},
				Places: []delivery.Place{
					{
						PhysicalDims: delivery.PhysicalDims{
							WeightGross:      240,
							Dx:               5,
							Dy:               10,
							Dz:               20,
							PredefinedVolume: 1000,
						},
					},
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			resp, err := d.GetDeliveryIntervals(
				c.IsOversized,
				c.LastMilePolicy,
				c.Request)
			if c.HasError {
				assert.NotNil(t, err)
				assert.Nil(t, resp)
				return
			}

			assert.Nil(t, err, "%s", err)
			assert.GreaterOrEqual(t, len(resp.Offers), 1)

			for _, offer := range resp.Offers {
				assert.False(t, offer.From.IsZero())
				assert.False(t, offer.To.IsZero())
			}
		})
	}
}

func TestDelvivery_GetDeliveryPoints(t *testing.T) {
	cases := []struct {
		Name      string
		Address   string
		Request   delivery.DeliveryPointsRequest
		HasErrors bool
	}{
		{
			Name:    "Успешный запрос списка пунктов выдачи",
			Address: "г. Москва",
			Request: delivery.DeliveryPointsRequest{
				Type:                       delivery.PST_PickupPoint,
				PaymentMethod:              delivery.PM_AlreadyPaid,
				IsNotBrandedPartnerStation: true,
				IsPostOffice:               true,
				PaymentMethods:             []delivery.PaymentMethod{delivery.PM_AlreadyPaid},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			resp, err := d.GetLocationID(c.Address)
			if assert.Nil(t, err) && assert.NotNil(t, resp) {
				for _, variant := range resp.Variants {
					assert.NotZero(t, variant.GeoID)
					assert.NotEmpty(t, variant.Address)
				}

				if assert.NotEmpty(t, resp.Variants) {
					c.Request.GeoID = resp.Variants[0].GeoID

					resp, err := d.GetDeliveryPoints(c.Request)
					if assert.Nil(t, err) && assert.NotNil(t, resp) {
						for _, point := range resp.Points {
							assert.NotEmpty(t, point.ID)
							assert.NotEmpty(t, point.OperatorStationID)
							assert.NotEmpty(t, point.Name)
							assert.NotEmpty(t, point.Address)
							assert.NotNil(t, point.PaymentMethods)
						}

						assert.NotEmpty(t, resp.Points)
					}
				}
			}
		})
	}
}

func TestDelivery_CreateOfferWOEditing(t *testing.T) {
	barcode := tool.NewID()

	cases := []struct {
		Name     string
		Region   string
		Request  delivery.CreateOfferRequest
		HasError bool
	}{
		{
			Name:     "Успешное создание заявки на доставку",
			HasError: false,
			Region:   "Москва",
			Request: delivery.CreateOfferRequest{
				Info:           delivery.Info{OperatorRequestID: tool.NewID()},
				Source:         delivery.Source{PlatformStation: &delivery.PlatformStation{PlatformID: "fbed3aa1-2cc6-4370-ab4d-59c5cc9bb924"}},
				LastMilePolicy: delivery.LMP_TimeInterval,
				RecipientInfo:  delivery.Contact{FirstName: "Иван", Phone: "+79261234567"},
				BillingInfo:    delivery.BillingInfo{PaymentMethod: delivery.PM_AlreadyPaid},
				Items: func() []delivery.Item {
					id := tool.NewID()

					return []delivery.Item{
						{
							Count:        1,
							Name:         "Чехол для iPhone 16 Pro Max",
							Article:      id,
							PlaceBarcode: barcode,
							BillingDetails: delivery.BillingDetails{
								UnitPrice:         10000,
								AssessedUnitPrice: 10000,
							},
						},
					}
				}(),
				Places: []delivery.Place{
					{
						Barcode: barcode,
						PhysicalDims: delivery.PhysicalDims{
							WeightGross:      240,
							Dx:               5,
							Dy:               10,
							Dz:               20,
							PredefinedVolume: 1000,
						},
					},
				},
			},
		},
		{
			Name:     "Успешное создание заявки на самовывоз",
			HasError: false,
			Region:   "Москва",
			Request: delivery.CreateOfferRequest{
				Info:           delivery.Info{OperatorRequestID: tool.NewID()},
				Source:         delivery.Source{PlatformStation: &delivery.PlatformStation{PlatformID: "fbed3aa1-2cc6-4370-ab4d-59c5cc9bb924"}},
				LastMilePolicy: delivery.LMP_TimeInterval,
				RecipientInfo:  delivery.Contact{FirstName: "Иван", Phone: "+79261234567"},
				BillingInfo:    delivery.BillingInfo{PaymentMethod: delivery.PM_AlreadyPaid},
				Items: func() []delivery.Item {
					id := tool.NewID()

					return []delivery.Item{
						{
							Count:        1,
							Name:         "Чехол для iPhone 16 Pro Max",
							Article:      id,
							PlaceBarcode: barcode,
							BillingDetails: delivery.BillingDetails{
								UnitPrice:         10000,
								AssessedUnitPrice: 10000,
							},
						},
					}
				}(),
				Places: []delivery.Place{
					{
						Barcode: barcode,
						PhysicalDims: delivery.PhysicalDims{
							WeightGross:      240,
							Dx:               5,
							Dy:               10,
							Dz:               20,
							PredefinedVolume: 1000,
						},
					},
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			// Шаг 1: Получаем Location ID для указанного адреса.
			locResp, locErr := d.GetLocationID(c.Region)
			assert.Nil(t, locErr)
			if !assert.NotNil(t, locResp) {
				return
			}

			// Шаг 2: Корректируем запрос на создание заказа в зависимости от типа его получения
			switch c.Request.LastMilePolicy {
			case delivery.LMP_SelfPickup:
				pointsResp, pointsErr := d.GetDeliveryPoints(delivery.DeliveryPointsRequest{
					GeoID:                      locResp.Variants[0].GeoID,
					Type:                       delivery.PST_PickupPoint,
					PaymentMethod:              delivery.PM_AlreadyPaid,
					IsNotBrandedPartnerStation: true,
					IsPostOffice:               true,
				})
				assert.Nil(t, pointsErr)
				assert.NotNil(t, pointsResp)
				if !assert.NotEmpty(t, pointsResp.Points) {
					return
				}

				c.Request.Destination = delivery.Destination{
					Type:            c.Request.LastMilePolicy.DestinationType(),
					PlatformStation: &delivery.PlatformStation{PlatformID: pointsResp.Points[0].ID},
				}
			case delivery.LMP_TimeInterval:
				intervalsResp, intervalsErr := d.GetDeliveryIntervals(false, c.Request.LastMilePolicy, delivery.DeliveryIntervalsRequest{
					Source:      delivery.Source{PlatformStationID: c.Request.Source.PlatformStation.PlatformID},
					Destination: delivery.Destination{Address: "Москва, пр-кт Вернадского 91к2"},
					Places:      c.Request.Places,
				})
				assert.Nil(t, intervalsErr)
				assert.NotNil(t, intervalsResp)
				if !assert.NotEmpty(t, intervalsResp.Offers) {
					return
				}

				c.Request.Destination = delivery.Destination{
					Type: c.Request.LastMilePolicy.DestinationType(),
					IntervalUTC: &delivery.IntervalUTC{
						From: intervalsResp.Offers[0].From,
						To:   intervalsResp.Offers[0].To,
					},
					CustomLocation: &delivery.CustomLocation{
						Details: delivery.Address{
							FullAddress: c.Request.Destination.Address,
						},
					},
				}
			default:
				t.Error("не корректный вариант доставки")
				t.FailNow()
				return
			}

			// Шаг 3: Создаем Offer.
			offerResp, offerErr := d.CreateOffer(c.Request)
			if c.HasError {
				assert.NotNil(t, offerErr)
				assert.Nil(t, offerResp)
				return
			}
			assert.Nil(t, offerErr)
			assert.NotNil(t, offerResp)
			if !assert.NotEmpty(t, offerResp.Offers) {
				return
			}

			// Шаг 4: Подтверждаем созданный Offer.
			confirmResp, confirmErr := d.ConfirmOffer(offerResp.Offers[0].OfferID)
			assert.Nil(t, confirmErr)
			if !assert.NotNil(t, confirmResp) {
				return
			}

			// Шаг 5: Получаем информацию о заявке на доставку.
			reqInfoResp, reqInfoErr := d.GetRequestInfo(confirmResp.RequestID, false)
			assert.Nil(t, reqInfoErr)
			assert.NotNil(t, reqInfoResp)

			// Шаг 6: Отменяем заявку на доставку.
			cancelResp, cancelErr := d.CancelRequest(confirmResp.RequestID)
			assert.Nil(t, cancelErr)
			assert.NotNil(t, cancelResp)
		})
	}
}

func TestDelivery_CreateOrderWOEditing(t *testing.T) {
	cases := []struct {
		Name             string
		Address          string
		CreateOfferReq   delivery.CreateOfferRequest
		CreateRequestReq delivery.CreateRequestRequest
		HasError         bool
		RequestID        string
		EditingTaskID    string
	}{ /* заполняемые тестовые случаи */ }

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			// Шаг 1: Получаем Location ID для указанного адреса.
			locResp, locErr := d.GetLocationID(c.Address)
			assert.Nil(t, locErr)
			assert.NotNil(t, locResp)

			// Шаг 2: Получаем Delivery Points на основе Location ID.
			c.CreateOfferReq.Source.PlatformStationID = fmt.Sprintf("%d", locResp.Variants[0].GeoID)
			pointsResp, pointsErr := d.GetDeliveryPoints(delivery.DeliveryPointsRequest{
				GeoID: locResp.Variants[0].GeoID,
			})
			assert.Nil(t, pointsErr)
			assert.NotNil(t, pointsResp)

			// Шаг 3: Создаем Offer.
			offerResp, offerErr := d.CreateOffer(c.CreateOfferReq)
			if c.HasError {
				assert.NotNil(t, offerErr)
				assert.Nil(t, offerResp)
				return
			}
			assert.Nil(t, offerErr)
			assert.NotNil(t, offerResp)

			// Шаг 4: Подтверждаем созданный Offer.
			confirmResp, confirmErr := d.ConfirmOffer(offerResp.Offers[0].OfferID)
			assert.Nil(t, confirmErr)
			assert.NotNil(t, confirmResp)

			// Шаг 5: Создаем Request на основе подтвержденного оффера.
			c.CreateRequestReq.Source.PlatformStationID = c.CreateOfferReq.Source.PlatformStationID
			createReqResp, createReqErr := d.CreateRequest(c.CreateRequestReq)
			assert.Nil(t, createReqErr)
			assert.NotNil(t, createReqResp)

			// Шаг 6: Генерация ярлыков и получения акта приема/передачи.
			// Эти методы возвращают ошибки "не реализовано", поэтому просто проверим их правильную обработку.
			labels, labelsErr := d.GenerateRequestLabels(delivery.GenerateRequestLabelsRequest{
				RequestIDS: []string{createReqResp.RequestID},
			})
			assert.NotNil(t, labelsErr)
			assert.Nil(t, labels)

			handoverAct, handoverActErr := d.GetRequestHandoverAct(createReqResp.RequestID)
			assert.NotNil(t, handoverActErr)
			assert.Nil(t, handoverAct)
		})
	}
}
