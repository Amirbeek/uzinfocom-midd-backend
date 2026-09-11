package order

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Amirbeek/uzinfocom-midd-backend/internal/database"
	"github.com/Amirbeek/uzinfocom-midd-backend/pkg/models"
)

var ErrInsufficientStock = errors.New("stockda qolmadi")

type Repo interface {
	CreateOrder(ctx context.Context, order models.Order, idempotency string, userId int64) error
	GetOrder(ctx context.Context, orderID int64, userID int64) (*models.Order, error)
	CancelOrder(ctx context.Context, orderID int64, userID int64) error
	CancelExpiredOrders(ctx context.Context) error
}

type repo struct{ db *sql.DB }

func NewRepo(db database.Service) Repo {
	return &repo{db: db.DB()}
}

func (r *repo) CreateOrder(ctx context.Context, order models.Order, idempotency string, userId int64) error {
	//Mashq: POST /orders — bir nechta item'li buyurtma, Idempotency-Key header majburiy (bir xil key bilan qayta yuborilsa, stock ikkinchi marta kamaymasligi kerak)

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	var existingOrderID int64

	err = tx.QueryRowContext(ctx, `
        SELECT order_id
        FROM idempotency_keys
        WHERE key = $1 AND user_id = $2
    `, idempotency, userId).Scan(&existingOrderID)

	if err == nil {
		// agar shu stage da kelsa, bu degani key bilan order oldin yarailgan
		return nil
	}

	if !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	// Order yaratish
	err = tx.QueryRowContext(ctx, `
        INSERT INTO orders (user_id, status)
        VALUES ($1, 'pending')
        RETURNING id
    `, userId).Scan(&order.ID)

	if err != nil {
		return err
	}

	// Har itemni qoshish
	for _, item := range order.Items {

		_, err = tx.ExecContext(ctx, `
            INSERT INTO order_items (
                order_id,
                product_id,
                quantity
            )
            VALUES ($1, $2, $3)
        `, order.ID, item.ProductID, item.Quantity)
		if err != nil {
			return err
		}

		// 4. Stockni kamaytirish
		result, err := tx.ExecContext(ctx, `
            UPDATE products
            SET stock_quantity = stock_quantity - $1
            WHERE id = $2
              AND stock_quantity >= $1
        `, item.Quantity, item.ProductID)

		if err != nil {
			return err
		}

		rows, err := result.RowsAffected()
		if err != nil {
			return err
		}

		if rows == 0 {
			return ErrInsufficientStock
		}
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO idempotency_keys (key, user_id, order_id) VALUES ($1, $2, $3)
	`, idempotency, userId, order.ID)
	if err != nil {
		return err
	}

	// 6. Hammasi hatosiz otsa COMMIT
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *repo) GetOrder(ctx context.Context, orderID int64, userID int64) (*models.Order, error) {
	// MASHQ:  GET /orders/{id} — status: pending → confirmed / cancelled
	var order models.Order

	err := r.db.QueryRowContext(ctx, `
		SELECT id, user_id, status, created_at
		FROM orders
		WHERE id = $1 AND user_id = $2
	`, orderID, userID).Scan(
		&order.ID,
		&order.UserID,
		&order.Status,
		&order.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (r *repo) CancelOrder(ctx context.Context, orderID int64, userID int64) error {
	// POST /orders/{id}/cancel — reserved stock qaytarilishi kerak
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	// Birinchi orderni statusini tekshiramiz chunki agar pedning bolmasa approve/rejectted oldin bolsa error qayataramiz,
	var status string
	err = tx.QueryRowContext(ctx, `
		SELECT status
		FROM orders
		WHERE id = $1 AND user_id = $2
	`, orderID, userID).Scan(&status)

	if err != nil {
		return err
	}

	if status != "pending" {
		return errors.New("order cancel qilib bolmadi")
	}

	// Stockdaki mahsulotlarni sonini qaytarish querysi. Bu yerda  "order_items.order_id = $1" topamiz va productga tegishli orderlarni topib stock quantityga qoshib qoyish querysi
	_, err = tx.ExecContext(ctx, `
		UPDATE products
		SET stock_quantity = products.stock_quantity + order_items.quantity
		FROM order_items
		WHERE order_items.order_id = $1 AND products.id = order_items.product_id
	`, orderID)

	if err != nil {
		return err
	}

	// endi orderni cancel qilamiz
	_, err = tx.ExecContext(ctx, `
		UPDATE orders
		SET status = 'cancelled'
		WHERE id = $1 AND user_id = $2
	`, orderID, userID)

	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *repo) CancelExpiredOrders(ctx context.Context) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `
		UPDATE products
		set stock_quantity = products.stock_quantity + o_items.quantity
		FROM order_items o_items
		WHERE o_items.order_id IN (
			SELECT id FROM orders WHERE status = 'pending'
			  AND created_at < NOW() - INTERVAL '15 minutes'
		)
		AND products.id = o_items.product_id
	`)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `
		UPDATE orders
		SET status = 'cancelled'
		WHERE status = 'pending' AND created_at < NOW() - INTERVAL '15 minutes'
	`)
	if err != nil {
		return err
	}

	return tx.Commit()
}
