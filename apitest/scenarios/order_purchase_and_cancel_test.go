package scenarios

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Reensef/go-microservices-course/apitest/internal/clients"
	orderApi "github.com/Reensef/go-microservices-course/shared/pkg/openapi/order/v1"
	inventoryGrpc "github.com/Reensef/go-microservices-course/shared/pkg/proto/inventory/v1"
)

// TestOrderPurchaseAndCancel — основной пользовательский путь: логин, просмотр
// каталога, покупка одной детали, оплата, и отдельно — отмена второго,
// неоплаченного заказа.
//
// Это ЕДИНЫЙ сквозной сценарий, а не набор независимых проверок: каждый шаг
// использует состояние (session_uuid, part_uuid, order_uuid), полученное на
// предыдущем, поэтому шаги идут как обычные вызовы внутри одного теста, а не
// как t.Run-подтесты — при падении любого шага сценарий сразу останавливается
// (require.*), т.к. без этого состояния дальнейшие шаги не имеют смысла.
//
// Сценарий бьёт по уже поднятому dev-стеку (iam, inventory, order, payment) —
// сам он ничего не поднимает и не останавливает, см. `task compose:*:up`.
// Адреса — localhost:50051-50054 по умолчанию, переопределяются через
// APITEST_*_ADDR / APITEST_ORDER_BASE_URL (см. internal/env).
//
// Шаги:
//  1. Регистрация и логин нового пользователя через iam — получаем session_uuid.
//  2. Список деталей в inventory (ListParts) — требует валидной сессии.
//  3. Получение конкретной детали по UUID (GetPart).
//  4. Создание заказа с этой деталью через order (REST); user_uuid заказа
//     order берёт из сессии сам, в теле запроса он не передаётся.
//  5. Заказ сразу после создания должен быть в статусе PENDING_PAYMENT.
//  6. Оплата заказа (REST) — order вызывает payment; payment независимо
//     проверяет сессию и сверяет её с user_uuid платежа.
//  7. После оплаты статус должен стать PAID.
//  8. Создание второго заказа — отдельный путь для проверки отмены.
//  9. Отмена второго заказа и проверка, что статус стал CANCELED.
func TestOrderPurchaseAndCancel(t *testing.T) {
	ctx := context.Background()
	c := clients.New(t)

	// Шаг 1: логинимся под новым пользователем.
	sessionUUID, login := clients.RegisterAndLogin(t, ctx, c)
	ctx = clients.WithSession(ctx, sessionUUID)
	t.Logf("пользователь %s залогинен, session_uuid=%s", login, sessionUUID)

	// Шаг 2: получаем список деталей.
	partsResp, err := c.Inventory.ListParts(ctx, &inventoryGrpc.ListPartsRequest{
		Filter: &inventoryGrpc.PartsFilter{},
	})
	require.NoError(t, err, "ListParts")
	require.NotEmpty(t, partsResp.GetParts(), "inventory должен содержать хотя бы одну деталь (посевные данные)")
	partUUID := partsResp.GetParts()[0].GetId()
	t.Logf("выбрана деталь %s", partUUID)

	// Шаг 3: получаем деталь по UUID и проверяем, что вернулась именно она.
	partResp, err := c.Inventory.GetPart(ctx, &inventoryGrpc.GetPartRequest{Id: partUUID})
	require.NoError(t, err, "GetPart")
	require.Equal(t, partUUID, partResp.GetPart().GetId())
	t.Logf("деталь получена: %s", partResp.GetPart().GetName())

	// Шаг 4: создаём заказ с этой деталью.
	createRes, err := c.Order.CreateOrder(ctx, &orderApi.CreateOrderRequest{
		PartIds: []string{partUUID},
	})
	require.NoError(t, err, "CreateOrder: transport error")
	createResp, ok := createRes.(*orderApi.CreateOrderResponse)
	require.Truef(t, ok, "CreateOrder: ожидался успешный ответ, получено %T", createRes)
	orderUUID := createResp.OrderUUID
	t.Logf("создан заказ %s на сумму %.2f", orderUUID, createResp.TotalPrice)

	// Шаг 5: заказ сразу после создания должен быть в PENDING_PAYMENT.
	orderDto := getOrder(t, ctx, c, orderUUID)
	require.Equal(t, orderApi.OrderStatusPENDINGPAYMENT, orderDto.Status)

	// Шаг 6: оплачиваем заказ.
	payRes, err := c.Order.PayOrder(ctx,
		&orderApi.PayOrderRequest{PaymentMethod: orderApi.PaymentMethodCARD},
		orderApi.PayOrderParams{OrderUUID: orderUUID},
	)
	require.NoError(t, err, "PayOrder: transport error")
	_, ok = payRes.(*orderApi.PayOrderResponse)
	require.Truef(t, ok, "PayOrder: ожидался успешный ответ, получено %T", payRes)
	t.Logf("заказ %s оплачен", orderUUID)

	// Шаг 7: после оплаты статус должен стать PAID.
	orderDto = getOrder(t, ctx, c, orderUUID)
	require.Equal(t, orderApi.OrderStatusPAID, orderDto.Status)

	// Шаг 8: создаём второй заказ — отдельный путь для проверки отмены.
	createRes2, err := c.Order.CreateOrder(ctx, &orderApi.CreateOrderRequest{
		PartIds: []string{partUUID},
	})
	require.NoError(t, err, "CreateOrder (2): transport error")
	createResp2, ok := createRes2.(*orderApi.CreateOrderResponse)
	require.Truef(t, ok, "CreateOrder (2): ожидался успешный ответ, получено %T", createRes2)
	order2UUID := createResp2.OrderUUID
	t.Logf("создан второй заказ %s", order2UUID)

	order2Dto := getOrder(t, ctx, c, order2UUID)
	require.Equal(t, orderApi.OrderStatusPENDINGPAYMENT, order2Dto.Status)

	// Шаг 9: отменяем второй заказ и проверяем итоговый статус.
	cancelRes, err := c.Order.CancelOrder(ctx, orderApi.CancelOrderParams{OrderUUID: order2UUID})
	require.NoError(t, err, "CancelOrder: transport error")
	_, ok = cancelRes.(*orderApi.CancelOrderNoContent)
	require.Truef(t, ok, "CancelOrder: ожидался 204 No Content, получено %T", cancelRes)

	order2Dto = getOrder(t, ctx, c, order2UUID)
	require.Equal(t, orderApi.OrderStatusCANCELED, order2Dto.Status)
	t.Logf("заказ %s отменён", order2UUID)
}

// getOrder — небольшой хелпер поверх GetOrderByUUID: сценарий проверяет статус
// заказа несколько раз, не хочется дублировать разбор union-типа ответа
// на каждом шаге.
func getOrder(t *testing.T, ctx context.Context, c *clients.Clients, orderUUID string) *orderApi.OrderDto {
	t.Helper()

	res, err := c.Order.GetOrderByUUID(ctx, orderApi.GetOrderByUUIDParams{OrderUUID: orderUUID})
	require.NoError(t, err, "GetOrderByUUID: transport error")
	dto, ok := res.(*orderApi.OrderDto)
	require.Truef(t, ok, "GetOrderByUUID: ожидался успешный ответ, получено %T", res)

	return dto
}
