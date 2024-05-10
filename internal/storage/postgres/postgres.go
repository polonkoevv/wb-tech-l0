package postgres

import (
	"context"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/polonkoevv/wb-tech/internal/models"
)

type Postgres struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Postgres {
	return &Postgres{
		db: db,
	}
}

func (s *Postgres) getAllItems(ctx context.Context, uid string) ([]models.Items, error) {
	op := "storage.postgres.GetAllItems"

	query := `SELECT
--     json_agg(
--             json_build_object(
                    chrt_id,
					order_uid,
					track_number,
					price,
					rid,
					name,
					sale,
					size,
					total_price,
					nm_id,
					brand,
					status
--             )
--     ) AS items
FROM order_items
WHERE order_uid = $1; 
	`

	rows, err := s.db.QueryxContext(ctx, query, uid)
	if err != nil {
		return nil, fmt.Errorf("%s : %w", op, err)
	}

	var cache []models.Items
	var c struct{}
	for rows.Next() {
		var item models.Items
		if err := rows.StructScan(&item); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		cache = append(cache, item)
	}

	fmt.Println(c)

	return cache, nil
}

func (s *Postgres) Save(ctx context.Context, order models.Order) error {
	op := "storage.postgres.Save"
	fmt.Println("Before")
	tx, err := s.db.BeginTxx(ctx, nil)
	fmt.Println("After")

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	_, err = tx.NamedExecContext(ctx,
		`
		INSERT INTO order_info (order_uid, track_number, entry, customer_id, delivery_service, date_created,
			shardkey, sm_id, oof_shard, locale, internal_signature)
		VALUES (:order_uid, :track_number, :entry, :customer_id, :delivery_service, :date_created,
	:shardkey, :sm_id, :oof_shard, :locale, :internal_signature)`,
		order,
	)
	if err != nil {
		return fmt.Errorf("%s : %w", op, err)
	}

	_, err = tx.NamedExecContext(ctx,
		`INSERT INTO delivery_info (order_uid, name, phone, zip, city, address, region, email)
			 VALUES (:order_uid, :name, :phone, :zip, :city, :address, :region, :email)`, order)
	if err != nil {
		return fmt.Errorf("%s : %w", op, err)
	}

	_, err = tx.NamedExecContext(ctx,
		`INSERT INTO payment_info (order_uid, transaction, request_id, currency, provider, amount, 
                          payment_dt, bank, delivery_cost, goods_total, custom_fee) 
			 VALUES (:order_uid, :transaction, :request_id, :currency, :provider, :amount, :payment_dt, :bank, 
			         :delivery_cost, :goods_total, :custom_fee)`,
		order)
	if err != nil {
		return fmt.Errorf("%s : %w", op, err)
	}

	stmt, err := tx.PrepareNamedContext(ctx,
		`INSERT INTO order_items (chrt_id, order_uid, track_number, price, rid, name, sale, size, total_price,
	                    nm_id, brand, status)
			   VALUES (:chrt_id, :order_uid, :track_number, :price, :rid, :name, :sale, :size, :total_price,
			           :nm_id, :brand, :status)`)

	for _, item := range order.Items {
		_, err = stmt.ExecContext(ctx, item)
		if err != nil {
			return fmt.Errorf("%s : %w", op, err)
		}
	}

	return nil
}

func (s *Postgres) LoadCache(ctx context.Context) (map[string]models.Order, error) {
	op := "storage.postgres.LoadCache"

	query := `
		SELECT
				oi.*,
				di.name, di.phone, di.zip,
				di.city, di.address, di.region, di.email,
				pi.transaction, pi.request_id, pi.currency, pi.provider, pi.amount, pi.payment_dt, pi.bank,
				pi.delivery_cost, pi.goods_total, pi.custom_fee
			FROM order_info oi
			JOIN delivery_info di ON oi.order_uid = di.order_uid
			JOIN payment_info pi ON oi.order_uid = pi.order_uid
	`

	rows, err := s.db.QueryxContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s : %w", op, err)
	}
	c, _ := rows.Rows.Columns()
	fmt.Println(c)

	defer rows.Close()

	cache := make(map[string]models.Order)

	for rows.Next() {
		var order models.Order

		if err := rows.StructScan(&order); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		cache[order.OrderUid] = order
	}

	for k, v := range cache {
		item, err := s.getAllItems(ctx, k)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		v.Items = item
		cache[k] = v
	}
	return cache, nil
}
