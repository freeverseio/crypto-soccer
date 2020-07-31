package postgres

// type PlaystoreOrderHistoryService struct {
// 	PlaystoreOrderService
// }

// func NewPlaystoreOrderHistoryService(tx *sql.Tx) *PlaystoreOrderHistoryService {
// 	return &PlaystoreOrderHistoryService{*NewPlaystoreOrderService(tx)}
// }

// func (b PlaystoreOrderHistoryService) UpdateState(order storage.PlaystoreOrder) error {
// 	if err := b.PlaystoreOrderService.UpdateState(order); err != nil {
// 		return err
// 	}
// 	return b.insertHistory(order)
// }

// func (b PlaystoreOrderHistoryService) Insert(order storage.PlaystoreOrder) error {
// 	if err := b.PlaystoreOrderService.Insert(order); err != nil {
// 		return err
// 	}
// 	return b.insertHistory(order)
// }

// func (b PlaystoreOrderHistoryService) insertHistory(order storage.PlaystoreOrder) error {
// 	_, err := b.tx.Exec(`INSERT INTO playstore_orders_histories (
// 		order_id,
// 		package_name,
// 		product_id,
// 		purchase_token,
// 		player_id,
// 		team_id,
// 		signature,
// 		state,
// 		state_extra
// 		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);`,
// 		order.OrderId,
// 		order.PackageName,
// 		order.ProductId,
// 		order.PurchaseToken,
// 		order.PlayerId,
// 		order.TeamId,
// 		order.Signature,
// 		order.State,
// 		order.StateExtra,
// 	)
// 	return err
// }
