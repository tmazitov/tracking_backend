package order

import (
	"errors"
	"fmt"

	"github.com/tmazitov/tracking_backend.git/internal/tms/bl"
)

func (s *Storage) OrderBillUpdatePrice(orderId int64, bill bl.R_OrderBill) error {

	conn, err := s.repo.Conn()
	if err != nil {
		return errors.New("DB conn error: " + err.Error())
	}

	defer s.repo.Close()

	execString := `
	INSERT INTO public.order_bills 
		( 
			order_id, 
			car_type_id,
			car_price,
			car_hours,
			helper_price,
			helper_hours,
			helper_count,
			km,
			price,
			is_fragile_cargo	
		)
		VALUES ( $1, $2, $3, $4, $5, $6, $7, $8, $9, $10 )
	ON CONFLICT (order_id) DO
	UPDATE
		SET 
			car_type_id = $2,
			car_price = $3,
			car_hours = $4,
			helper_price = $5,
			helper_hours = $6,
			helper_count = $7,
			km = $8,
			price = $9,
			is_fragile_cargo = $10
			
		WHERE order_bills.order_id=$1
	`

	fmt.Println(
		execString,
		orderId,
		bill.CarTypeID,
		bill.CarPrice,
		bill.CarHours,
		bill.HelperPrice,
		bill.HelperHours,
		bill.HelperCount,
		bill.KM,
		bill.Total,
		bill.IsFragileCargo,
	)

	return conn.QueryRow(
		execString,
		orderId,
		bill.CarTypeID,
		bill.CarPrice,
		bill.CarHours,
		bill.HelperPrice,
		bill.HelperHours,
		bill.HelperCount,
		bill.KM,
		bill.Total,
		bill.IsFragileCargo,
	).Err()
}
