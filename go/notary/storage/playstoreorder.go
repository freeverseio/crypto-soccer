package storage

import "database/sql"

type PlaystoreOrderState string

const (
	PlaystoreOrderPending  PlaystoreOrderState = "pending"
	PlaystoreOrderComplete PlaystoreOrderState = "complete"
	PlaystoreOrderFailed   PlaystoreOrderState = "failed"
)

type PlaystoreOrder struct {
	OrderId       string
	PackageName   string
	ProductId     string
	PurchaseToken string
	State         PlaystoreOrderState
	StateExtra    string
}

func NewPlaystoreOrder() *PlaystoreOrder {
	order := PlaystoreOrder{}
	order.State = PlaystoreOrderPending
	return &order
}

func PlaystoreOrderByOrderId(tx *sql.Tx, orderId string) (*PlaystoreOrder, error) {
	rows, err := tx.Query(`SELECT 
	package_name,
	product_id,
	purchase_token,
	state, 
	state_extra 
	FROM playstore_orders WHERE order_id=$1;`, orderId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	order := PlaystoreOrder{}
	order.OrderId = orderId

	err = rows.Scan(
		&order.PackageName,
		&order.ProductId,
		&order.PurchaseToken,
		&order.State,
		&order.StateExtra,
	)
	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (b PlaystoreOrder) Insert(tx *sql.Tx) error {
	_, err := tx.Exec(`INSERT INTO playstore_orders (
		order_id, 
		package_name,
		product_id,
		purchase_token,
		state, 
		state_extra
		) VALUES ($1, $2, $3, $4, $5, $6);`,
		b.OrderId,
		b.PackageName,
		b.ProductId,
		b.PurchaseToken,
		b.State,
		b.StateExtra,
	)
	return err
}
