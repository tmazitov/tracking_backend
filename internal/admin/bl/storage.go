package bl

type UserStorage interface {
	OfferAccept(offerId int) (*DB_User, error)
	OfferReject(offerId int) error
	UserOfferList() ([]DB_UserOffer, error)
	StaffRemove(userId int64) error
	StaffWorkTimeUpdate(workTime *StaffWorkTime) error
}

type OrderStorage interface {
	OrderPricelistUpdate(priceList *OrderPriceList) error
}

type Storage interface {
	UserStorage() UserStorage
	OrderStorage() OrderStorage
}
