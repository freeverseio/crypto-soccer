package postgres

import (
	"github.com/freeverseio/crypto-soccer/go/notary/storage"
)

type PlaystoreOrderHistoryService struct {
	service storage.PlaystoreOrderService
}

func NewPlaystoreOrderHistoryService(service storage.PlaystoreOrderService) *PlaystoreOrderHistoryService {
	return &PlaystoreOrderHistoryService{
		service: service,
	}
}

func (b PlaystoreOrderHistoryService) UpdateState(order storage.PlaystoreOrder) error {
	return b.service.UpdateState(order)
}

func (b PlaystoreOrderHistoryService) PendingOrders() ([]storage.PlaystoreOrder, error) {
	return b.service.PendingOrders()
}

func (b PlaystoreOrderHistoryService) Order(orderId string) (*storage.PlaystoreOrder, error) {
	return b.service.Order(orderId)
}

func (b PlaystoreOrderHistoryService) Insert(order *storage.PlaystoreOrder) error {
	return b.service.Insert(order)
}
