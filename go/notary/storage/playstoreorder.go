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
	PlayerId      string
	TeamId        string
	Signature     string
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
	player_id,
	team_id,
	signature,
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
		&order.PlayerId,
		&order.TeamId,
		&order.Signature,
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
		player_id,
		team_id,
		signature,
		state, 
		state_extra
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);`,
		b.OrderId,
		b.PackageName,
		b.ProductId,
		b.PurchaseToken,
		b.PlayerId,
		b.TeamId,
		b.Signature,
		b.State,
		b.StateExtra,
	)
	return err
}
