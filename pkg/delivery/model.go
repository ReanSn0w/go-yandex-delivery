package delivery

import "time"

const (
	PM_AlreadyPaid    PaymentMethod = "already_paid"
	PM_CardOnReceiot  PaymentMethod = "card_on_receipt"
	PM_CashOnDelivery PaymentMethod = "cash_on_delivery"

	LMP_TimeInterval LastMilePolicy = "time_interval"
	LMP_SelfPickup   LastMilePolicy = "self_pickup"

	PST_PickupPoint   PickupStationType = "pickup_point"
	PST_Terminal      PickupStationType = "terminal"
	PST_PostOffice    PickupStationType = "post_office"
	PST_SortingCenter PickupStationType = "sorting_center"

	ERS_Pending   EditingRequestStatus = "pending"
	ERS_Execution EditingRequestStatus = "execution"
	ERS_Success   EditingRequestStatus = "success"
	ERS_Failure   EditingRequestStatus = "failure"

	R_Cancel_ShopCanceled               Reason = "SHOP_CANCELLED"                // Отправитель отменил заказ
	R_Cancel_UserChangedMind            Reason = "USER_CHANGED_MIND"             // Покупатель передумал
	R_Cancel_DeliveryProblems           Reason = "DELIVERY_PROBLEMS"             // Проблемы с доставкой
	R_Cancel_DimensionsExceeded         Reason = "DIMENSIONS_EXCEEDED"           // Посылка слишком большая для способа доставки
	R_Cancel_DimensionsExceededLocker   Reason = "DIMENSIONS_EXCEEDED_LOCKER"    // Превышены допустимые габариты постамата
	R_Cancel_NoPassport                 Reason = "NO_PASSPORT"                   // Нет паспорта
	R_Cancel_OrderIsDamaged             Reason = "ORDER_IS_DAMAGED"              // Заказ поврежден
	R_Cancel_ExtraRescheduling          Reason = "EXTRA_RESCHEDULING"            // Заказ отменен из-за частых переносов
	R_Cancel_BrokenItem                 Reason = "BROKEN_ITEM"                   // Товар оказался бракованным
	R_Cancel_OrderItemsQuantityMismatch Reason = "ORDER_ITEMS_QUANTITY_MISMATCH" // Не совпадает количество товаров
	R_Cancel_OrderWasLost               Reason = "ORDER_WAS_LOST"                // Заказ утерян
	R_Cancel_LateContact                Reason = "LATE_CONTACT"                  // С пользователем связались слишком поздно
	R_Cancel_PickupExpired              Reason = "PICKUP_EXPIRED"                // Срок хранения в пункте выдачи истек

	R_Change_ClientRequest                   Reason = "CLIENT_REQUEST"                      // "По просьбе клиента"
	R_Change_CouriesCouldNotContactRecipient Reason = "COURIES_COULD_NOT_CONTACT_RECIPIENT" // "Курьер не смог дозвониться"
	R_Change_DeliveryDateUupdatedByDelivery  Reason = "DELIVERY_DATE_UPDATED_BY_DELIVERY"   // "Задержка обработки заказа партнёром"
	R_Change_DeliveryDateUupdatedByRecipient Reason = "DELIVERY_DATE_UPDATED_BY_RECIPIENT"  // "По запросу от пользователя"
	R_Change_DeliveryDateUupdatedByShop      Reason = "DELIVERY_DATE_UPDATED_BY_SHOP"       // "По запросу от магазина"
	R_Change_LastMileChangedByUser           Reason = "LAST_MILE_CHANGED_BY_USER"           // "Последняя миля изменена по инициативе пользователя"
	R_Change_LockerFull                      Reason = "LOCKER_FULL"                         // "Нет свободных ячеек подходящего размера"
	R_Change_NoPassport                      Reason = "NO_PASSPORT"                         // "Нет паспорта"
	R_Change_PickupPointTechnicalIssues      Reason = "PICKUPPOINT_TECHNICAL_ISSUES"        // "Технические проблемы в ПВЗ"

	R_Unknown Reason = "UNKNOWN" // "Не определён"
	R_Other   Reason = "OTHER"   // "Другая проблема"
)

type (
	PaymentMethod        string // Способ оплаты товаров
	LastMilePolicy       string // Варианты доставки
	PickupStationType    string // Тип точки приема/выдачи заказа.
	EditingRequestStatus string // Статус запроса на редактирование
	Reason               string // Описание причины переноса/отмены
)

func (l LastMilePolicy) DestinationType() string {
	switch l {
	case LMP_TimeInterval:
		return "custom_location"
	case LMP_SelfPickup:
		return "platform_station"
	default:
		return "unknown"
	}
}

type PredictPriceRequest struct {
	Source             Source         `json:"source"`
	Destination        Destination    `json:"destination"`
	PaymentMethod      PaymentMethod  `json:"payment_method"`
	Places             []Place        `json:"places"`
	Tariff             LastMilePolicy `json:"tariff"`               // Тариф доставки
	TotalWeight        int64          `json:"total_weight"`         // Cуммарный вес посылки в граммах
	TotalAssessedPrice int64          `json:"total_assessed_price"` // Суммарная оценочная стоимость посылок в копейках
	ClientPrice        int64          `json:"client_price"`         // Cумма к оплате с получателя в копейках
}

// Информация о точке получения заказа
type Destination struct {
	// Тип целевой точки.
	// Для доставки до двери — custom_location (2),
	// для доставки до ПВЗ — platform_station (1)
	Type string `json:"type,omitempty"`

	// Описание целевой станции в случае, если она зарегистрирована в платформе
	PlatformStation *PlatformStation `json:"platform_station,omitempty"`

	// Полное описание целевого адреса доставки
	CustomLocation *CustomLocation `json:"custom_location,omitempty"`

	// Временной интервал (в UNIX)
	IntervalUTC *IntervalUTC `json:"interval_utc,omitempty"`

	//ID ПВЗ или постамата, зарегистрированного в платформе, в который нужна доставка
	PlatformStationID string `json:"platform_station_id,omitempty"`

	// Адрес получения с указанием города, улицы и номера дома.
	// Номер квартиры, подъезда и этаж указывать не нужно
	Address string `json:"address,omitempty"`
}

// Информация о местах в заказе
type Place struct {
	PhysicalDims PhysicalDims `json:"physical_dims"`         // Физические параметры места
	Barcode      string       `json:"barcode,omitempty"`     // Текущий штрихкод грузоместа
	Description  string       `json:"description,omitempty"` // Описание грузоместа
}

type PhysicalDims struct {
	WeightGross      int64 `json:"weight_gross"`      // Вес брутто, граммы
	Dx               int64 `json:"dx"`                // Длина, сантиметры
	Dy               int64 `json:"dy"`                // Высота, сантиметры
	Dz               int64 `json:"dz"`                // Ширина, сантиметры
	PredefinedVolume int64 `json:"predefined_volume"` // Объем (в см3)
}

// Информация о точке отправления заказа
type Source struct {
	// ID склада отправки, зарегистрированного в платформе
	PlatformStationID string           `json:"platform_station_id,omitempty"`
	PlatformStation   *PlatformStation `json:"platform_station,omitempty"`
	IntervalUTC       *IntervalUTC     `json:"interval_utc,omitempty"`
}

type PredictPriceResponse struct {
	Error bool `json:"error"`
	// Суммарная стоимость доставки с учетом дополнительных услуг (с НДС)
	PricingTotal string `json:"pricing_total"`

	// Размер комисиии за прием наложенного платежа в руб (с НДС).
	// Заполняется в случае указания способа оплаты отлиного от already_paid
	PricingCommissionOnDeliveryPaymentAmount string `json:"pricing_commission_on_delivery_payment_amount"`

	// Размер комисиии за прием наложенного платежа в %.
	// Заполняется в случае указания способа оплаты отлиного от already_paid
	PricingCommissionOnDeliveryPayment string `json:"pricing_commission_on_delivery_payment"`

	// Стоимость за услугу доставки и страхование посылки в руб (с НДС).
	// В случае указания способа оплаты already_paid не заполняется
	Pricing string `json:"pricing"`
}

type DeliveryIntervalsRequest struct {
	Source      Source      `json:"source"`
	Destination Destination `json:"destination"`
	Places      []Place     `json:"places"` // Информация о местах в заказе
}

type DeliveryIntervalsResponse struct {
	Offers []Offer `json:"offers"`
}

type Offer struct {
	// UNIX или UTC время предлагаемого начала доставки
	From time.Time `json:"from"`

	// UNIX или UTC время предлагаемого окончания доставки
	To time.Time `json:"to"`

	// Будет ли заказ доставляться почтой россии
	DeliveredByPost bool `json:"delivered_by_post"`
}

type LocationIDResponse struct {
	Variants []LocationDetectedVariant `json:"variants"`
}

type LocationDetectedVariant struct {
	GeoID   int64  `json:"geo_id"`  // Идентификатор населенного пункта (geo_id)
	Address string `json:"address"` // Вариант адреса
}

type DeliveryPointsRequest struct {
	PickupPointIDS []string `json:"pickup_point_ids"`

	// Идентификатор населенного пункта (geo_id)
	GeoID                      int64               `json:"geo_id"`
	Longitude                  *CoordinateInterval `json:"longitude,omitempty"`                      // Интервал для выбора всех объектов в отрезке по долготе.
	Latitude                   *CoordinateInterval `json:"latitude,omitempty"`                       // Интервал для выбора всех объектов в отрезке по широте.
	Type                       PickupStationType   `json:"type,omitempty"`                           // Тип точки приема/выдачи заказа.
	PaymentMethod              PaymentMethod       `json:"payment_method,omitempty"`                 // Тип оплаты в точке самостоятельного получения заказа.
	AvailableForDropoff        bool                `json:"available_for_dropoff,omitempty"`          // Возможность отгрузки заказов в точку (самопривоз).
	IsYandexBranded            bool                `json:"is_yandex_branded,omitempty"`              // Признак брендированные ли ПВЗ.
	IsNotBrandedPartnerStation bool                `json:"is_not_branded_partner_station,omitempty"` // Признак добавляющий партнерские ПВЗ.
	IsPostOffice               bool                `json:"is_post_office,omitempty"`                 // Признак добавляющий ПВЗ почты россии.
	PaymentMethods             []PaymentMethod     `json:"payment_methods,omitempty"`                // Набор типов оплаты, которы должны быть доступны в самостоятельного точке получения заказа.
}

type CoordinateInterval struct {
	From int64 `json:"from"` // Нижняя граница интервала
	To   int64 `json:"to"`   // Верхняя граница интервала
}

type DeliveryPointsResponse struct {
	Points []Point `json:"points"`
}

type Point struct {
	// Идентификатор точки забора заказа.
	// Должен использоваться при получении вариантов доставки в качестве конечной точки
	ID                string   `json:"ID"`
	OperatorStationID string   `json:"operator_station_id"` // Идентификатор точки забора заказа в системе оператора
	Name              string   `json:"name"`                // Название точки забора заказа
	Type              string   `json:"type"`                // Тип точки забора заказа
	Position          Position `json:"position"`            // Точка, в которой находится точка забора заказа
	Address           Address  `json:"address"`             // Полный адрес точки забора заказа
	Instruction       string   `json:"instruction"`         // Дополнительные указания по тому, как добраться до точки получения заказа
	PaymentMethods    []string `json:"payment_methods"`     // Возможные методы оплаты заказа при получении
	Contact           Contact  `json:"contact"`             // Данные для связи с точкой забора заказа
	Schedule          Schedule `json:"schedule"`            // Расписание работы точки
	IsYandexBranded   bool     `json:"is_yandex_branded"`   // Признак брендированные ли ПВЗ
	IsMarketPartner   bool     `json:"is_market_partner"`   // Признак партнерского ПВЗ
	IsDarkStore       bool     `json:"is_dark_store"`       // Признак даркстора
	IsPostOffice      bool     `json:"is_post_office"`      // Признак Почты России
	Dayoffs           []Dayoff `json:"dayoffs"`             // Нерабочие дни ПВЗ
}

type Address struct {
	// Комментарий
	Comment string `json:"comment"`

	// Полный адрес с указанием города, улицы и номера дома.
	// Номер квартиры, подъезда и этаж указывать не нужно
	FullAddress string `json:"full_address"`

	// Номер квартиры или офиса, обязателен при наличии
	Room string `json:"room"`
}

type Contact struct {
	FirstName  string `json:"first_name"`           // Имя
	LastName   string `json:"last_name,omitempty"`  // Фамилия
	Partonymic string `json:"partonymic,omitempty"` // Отчество
	Phone      string `json:"phone"`                // Телефон
	Email      string `json:"email,omitempty"`      // Электронная почта
}

type Dayoff struct {
	Date string `json:"date_utc"` // Дата в формате UTC
}

type Position struct {
	Latitude  float64 `json:"latitude"`  // Широта
	Longitude float64 `json:"longitude"` // Долгота
}

type Schedule struct {
	TimeZone     int64         `json:"time_zone"`
	Restrictions []Restriction `json:"restrictions"`
}

type Restriction struct {
	// Номера дней недели, к которым применяется правило.
	// 1 - понедельник,
	// 2 - вторник,
	// ...,
	// 7 - воскресенье.
	Days     []int64 `json:"days"`
	TimeFrom Time    `json:"time_from"` // Время начала работы.
	TimeTo   Time    `json:"time_to"`   // Время окончания работы.
}

type Time struct {
	Hours   int64 `json:"hours"`   // Часы
	Minutes int64 `json:"minutes"` // Минуты
}

type CreateOfferRequest struct {
	Info           Info           `json:"info"`             // Базовый набор метаданных по запросу
	Source         Source         `json:"source"`           // Информация о точке отправления заказа
	Destination    Destination    `json:"destination"`      // Информация о точке получения заказа
	Items          []Item         `json:"items"`            // Информация о предметах в заказе
	Places         []Place        `json:"places"`           // Информация о местах в заказе
	BillingInfo    BillingInfo    `json:"billing_info"`     // Данные для биллинга
	RecipientInfo  Contact        `json:"recipient_info"`   // Данные о получателе
	LastMilePolicy LastMilePolicy `json:"last_mile_policy"` // Требуемый способ доставки

	// Разрешен ли частичный выкуп
	// true — разрешен частичный выкуп заказа
	// false — частичный выкуп заказа недоступен
	// Значение по умолчанию: false
	ParticularItemsRefuse bool `json:"particular_items_refuse"`
}

type BillingInfo struct {
	PaymentMethod PaymentMethod `json:"payment_method"` // Метод оплаты

	// Сумма, которую нужно взять с получателя за доставку. Актуально только для заказов с постоплатой (типы оплаты cash_on_receipt и card_on_receipt)
	DeliveryCost int64 `json:"delivery_cost,omitempty"`
}

type CustomLocation struct {
	Latitude  float64 `json:"latitude,omitempty"`  // Широта
	Longitude float64 `json:"longitude,omitempty"` // Долгота
	Details   Address `json:"details"`             // Детали
}

type BillingDetails struct {
	// Значение НДС.
	// Допустимые значения — 0, 5, 7, 10, 20.
	// Если заказ без НДС, передавайте значение -1
	NDS               int64  `json:"nds,omitempty"`
	Inn               string `json:"inn,omitempty"`       // ИНН
	UnitPrice         int64  `json:"unit_price"`          // Цена за единицу товара (передается в копейках)
	AssessedUnitPrice int64  `json:"assessed_unit_price"` // Оценочная цена за единицу товара (передается в копейках)
}

type CreateOfferResponse struct {
	Offers []OfferItem `json:"offers"`
}

type OfferItem struct {
	OfferID      string       `json:"offer_id"`      // Идентификатор предложения маршрутного листа (оффера).
	ExpiresAt    time.Time    `json:"expires_at"`    // Timestamp окончания действия предложения маршрутного листа в Гринвиче.
	OfferDetails OfferDetails `json:"offer_details"` // Подробности оффера.
}

type OfferDetails struct {
	DeliveryInterval                         DeliveryInterval `json:"delivery_interval"`
	PickupInterval                           PickupInterval   `json:"pickup_interval"`
	Pricing                                  string           `json:"pricing"`
	PricingCommissionOnDeliveryPayment       string           `json:"pricing_commission_on_delivery_payment"`
	PricingCommissionOnDeliveryPaymentAmount string           `json:"pricing_commission_on_delivery_payment_amount"`
	PricingTotal                             string           `json:"pricing_total"`
}

type DeliveryInterval struct {
	Min    time.Time `json:"min"`    // Нижняя граница временного интервала доставки, UTC Timestamp в Гринвиче.
	Max    time.Time `json:"max"`    // Верхняя граница временного интервала доставки, UTC Timestamp в Гринвиче.
	Policy string    `json:"policy"` // Политика доставки последней мили в предложенном варианте.
}

type PickupInterval struct {
	Min time.Time `json:"min"` // Нижняя граница временного интервала забора, UTC Timestamp в Гринвиче.
	Max time.Time `json:"max"` // Верхняя граница временного интервала забора, UTC Timestamp в Гринвиче.
}

type ConfirmOfferResponse struct {
	RequestID string `json:"request_id"` // Идентификатор только что созданного заказа
}

type GetRequestInfoResponse RequestElement

type ActualDeliveryInterval struct {
	From string `json:"from"` // "10:00+03:00"
	To   string `json:"to"`   // "11:00+03:00"
}

type AvailableActions struct {
	UpdateDatesAvailable           bool `json:"update_dates_available"`
	UpdateAddressAvailable         bool `json:"update_address_available"`
	UpdateCourierToPickupAvailable bool `json:"update_courier_to_pickup_available"`
	UpdatePickupToCourierAvailable bool `json:"update_pickup_to_courier_available"`
	UpdatePickupToPickupAvailable  bool `json:"update_pickup_to_pickup_available"`
	UpdateItems                    bool `json:"update_items"`
	UpdateRecipient                bool `json:"update_recipient"`
	UpdatePlaces                   bool `json:"update_places"`
}

type IntervalUTC struct {
	From time.Time `json:"from"`
	To   time.Time `json:"to"`
}

type PlatformStation struct {
	PlatformID string `json:"platform_id"`
}

type Info struct {
	OperatorRequestID string `json:"operator_request_id"` // Идентификатор заказа у отправителя
	Comment           string `json:"comment"`             // Комментарий к заказу
}

type Item struct {
	Count          int64          `json:"count"`                   // Количество товара
	Name           string         `json:"name"`                    // Название товара
	Article        string         `json:"article"`                 // Артикул товара
	MarkingCode    string         `json:"marking_code,omitempty"`  // Код маркировки товара
	Uin            string         `json:"uin,omitempty"`           // Уникальный идентификатор товара
	BillingDetails BillingDetails `json:"billing_details"`         // Детали оплаты товара
	PhysicalDims   *PhysicalDims  `json:"physical_dims,omitempty"` // Физические размеры товара
	PlaceBarcode   string         `json:"place_barcode"`           // Штрихкод места хранения товара
}

type State struct {
	Status       string    `json:"status"`        // Статус, описывающий текущее состояние заказа
	Description  string    `json:"description"`   // Описание статуса
	TimestampUTC time.Time `json:"timestamp_utc"` // Временная метка в формате UTC
	Reason       Reason    `json:"reason"`
}

type GetRequestsInfoResponse struct {
	Requests []RequestElement `json:"requests"`
}

type RequestElement struct {
	RequestID      string      `json:"request_id"` // ID заказа в логистической платформе
	Request        RequestInfo `json:"request"`
	State          State       `json:"state"`            // Текущий статус заказа
	FullItemsPrice int64       `json:"full_items_price"` // Общая стоимость всех предметов в заказе
	SharingURL     string      `json:"sharing_url"`      // Ссылка на страницу с трекингом заказа для получателя
	CourierOrderID string      `json:"courier_order_id"` // Номер заказа в системе оператора
}

type RequestInfo struct {
	Info                  Info             `json:"info"`
	Source                Source           `json:"source"`
	Destination           Destination      `json:"destination"`
	Items                 []Item           `json:"items"`
	Places                []Place          `json:"places"`
	BillingInfo           BillingInfo      `json:"billing_info"`
	RecipientInfo         Contact          `json:"recipient_info"`
	LastMilePolicy        LastMilePolicy   `json:"last_mile_policy"`
	ParticularItemsRefuse bool             `json:"particular_items_refuse"`
	AvailableActions      AvailableActions `json:"available_actions"`
}

type GetRequestActualInfoResponse struct {
	DeliveryDate     string           `json:"delivery_date"`     // Дата доставки
	DeliveryInterval DeliveryInterval `json:"delivery_interval"` // Интервал доставки
}

type EditRequestInfoRequest struct {
	RequestID      string         `json:"request_id"`       // ID заказа
	RecipientInfo  Contact        `json:"recipient_info"`   // Информация о получателе
	Destination    Destination    `json:"destination"`      // Информация о пункте назначения
	LastMilePolicy string         `json:"last_mile_policy"` // Политика доставки
	Places         []PlaceElement `json:"places"`           // Данные о грузоместах
}

type PlaceElement struct {
	Barcode string `json:"barcode"` // Текущий штрихкод грузоместа
	Place   Place  `json:"place"`   // Информация о грузоместе
}

type EditRequestInfoResponse struct {
	CompletedUpdates []Update `json:"completed_updates"` // Выполненные изменения заказа
	ActiveUpdates    []Update `json:"active_updates"`    // Выполняющиеся изменения заказа
	IgnoredUpdates   []Update `json:"ignored_updates"`   // Невыполненные изменения заказа
	EditID           string   `json:"edit_id"`           // ID операции редактирования
}

type Update struct {
	Reason       string   `json:"reason"`        // Причина формирования статуса
	Status       string   `json:"status"`        // Статус обновления
	Type         string   `json:"type"`          // Тип изменения
	ErrorDetails []string `json:"error_details"` // Описание ошибок
	Code         string   `json:"code"`          // Кол статуса
}

type GetRequestRedeliveryOptionsRequest struct {
	RequestID   string      `json:"request_id"`  // ID заказа
	Destination Destination `json:"destination"` // Информация о пункте назначения
}

type GetRequestRedeliveryOptionsResponse struct {
	// Возможные интервалы доставки
	// Интервал времени в формате UTC.
	Options []IntervalUTC `json:"options"`
}

type GetRequestHistoryResponse struct {
	StateHistory []StateHistory `json:"state_history"` // История изменения статусов заказа
}

type StateHistory struct {
	Status       string `json:"status"`        // Статус, описывающий текущее состояние заказа
	Description  string `json:"description"`   // Описание статуса
	TimestampUTC string `json:"timestamp_utc"` // Временная метка в формате UTC
	Reason       Reason `json:"reason"`        // Причина изменения статуса
}

type CancelRequestResponse struct {
	// Статус отмены заявки. Может принимать только значения из enum.
	// CREATED: отмена заявки инициирована в платформе
	// SUCCESS: отмена заявки успешно выполнена
	// ERROR: запрос завершился с ошибкой
	Status      string `json:"status"`
	Description string `json:"description"` // Комментарий к результату выполнения запроса
}

type CreateRequestRequest struct {
	Info           Info           `json:"info"`             // Базовый набор метаданных по запросу
	Source         Source         `json:"source"`           // Информация о точке отправления заказа
	Destination    Destination    `json:"destination"`      // Информация о точке получения заказа
	Items          []Item         `json:"items"`            // Информация о предметах в заказе
	Places         []Place        `json:"places"`           // Информация о местах в заказе
	BillingInfo    BillingInfo    `json:"billing_info"`     // Данные для биллинга
	RecipientInfo  Contact        `json:"recipient_info"`   // Данные о получателе
	LastMilePolicy LastMilePolicy `json:"last_mile_policy"` // Требуемый способ доставки

	// Разрешен ли частичный выкуп
	// true — разрешен частичный выкуп заказа
	// false — частичный выкуп заказа недоступен
	// Значение по умолчанию: false
	ParticularItemsRefuse bool `json:"particular_items_refuse"`
}

type CreateRequestResponse struct {
	RequestID string `json:"request_id"` // Идентификатор только что созданного заказа
}

type EditRequestPlacesRequest struct {
	RequestID string `json:"request_id"` // Идентификатор заказа в системе
	Places    Places `json:"places"`     // Информация о местах в заказе
}

type Places struct {
	Dimensions Dimensions   `json:"dimensions"` // Размеры коробки
	Barcode    string       `json:"barcode"`    // Штрихкод коробки
	Items      []PlacesItem `json:"items"`      // Список товаров в коробке
}

type Dimensions struct {
	WeightGross int64 `json:"weight_gross"` // Вес коробки в граммах
	Dx          int64 `json:"dx"`           // Ширина коробки в сантиметрах
	Dy          int64 `json:"dy"`           // Высота коробки в сантиметрах
	Dz          int64 `json:"dz"`           // Длина коробки в сантиметрах
}

type PlacesItem struct {
	Count       int64  `json:"count"`        // Количество товара в коробке
	ItemBarcode string `json:"item_barcode"` // Штрихкод товара, который находится в этой коробке
}

type EditRequestPlacesResponse struct {
	EditingTaskID string `json:"editing_task_id"` // Идентификатор созданного запроса на редактирование для уточнения его статуса
}

type GetEditRequestStatusResponse struct {
	// Статус запроса на редактирование
	// pending: ожидает выполнения
	// execution: выполняется
	// success: успешно выполнен
	// failure: в процессе выполнения произошла ошибка
	Status EditingRequestStatus `json:"status"`
}

type EditRequestItemsRequest struct {
	// ID запроса на редактирование
	RequestID string `json:"request_id"`

	// Список товаров
	// Указываются маркировки для редактирования
	ItemsInstances []ItemsInstance `json:"items_instances"`
}

type ItemsInstance struct {
	ItemBarcode string `json:"item_barcode"` // Штрихкод товара
	Article     string `json:"article"`      // Артикул товара
	MarkingCode string `json:"marking_code"` // Код маркировки товара
}

type EditRequestItemsResponse struct {
	EditingTaskID string `json:"editing_task_id"` // Идентификатор созданного запроса на редактирование для уточнения его статуса
}

type GenerateRequestLabelsRequest struct {
	RequestIDS   []string `json:"request_ids"`   // Список ID заказов. Количество заказов не должно превышать предельно допустимого.
	GenerateType string   `json:"generate_type"` // Формат генерации ярлыков. one - один ярлык на страницу. many - максимум ярлыков на страницу.
	Language     string   `json:"language"`      // Язык надписей на этикетке
}
