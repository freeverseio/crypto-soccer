package storage

import "database/sql"

type PlaystoreOrderState string

const (
	PlaystoreOrderPending  PlaystoreOrderState = "pending"
	PlaystoreOrderComplete PlaystoreOrderState = "complete"
	PlaystoreOrderFailed   PlaystoreOrderState = "failed"
)

type PlaystoreOrder struct {
	OrderId    string
	State      PlaystoreOrderState
	StateExtra string
}

func NewPlaystoreOrder(orderId string) *PlaystoreOrder {
	order := PlaystoreOrder{}
	order.OrderId = orderId
	order.State = PlaystoreOrderPending
	return &order
}

func PlaystoreOrderByOrderId(tx *sql.Tx, orderId string) (*PlaystoreOrder, error) {
	rows, err := tx.Query("SELECT state, state_extra FROM playstore_orders WHERE order_id=$1;", orderId)
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
		&order.State,
		&order.StateExtra,
	)
	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (b PlaystoreOrder) Insert(tx *sql.Tx) error {
	_, err := tx.Exec("INSERT INTO playstore_orders (order_id, state, state_extra) VALUES ($1, $2, $3);",
		b.OrderId,
		b.State,
		b.StateExtra,
	)
	return err
}
