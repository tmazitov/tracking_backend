package bl

type OrderPriceList struct {
	BigCarPrice  uint `json:"bigCarPrice"`
	BigCarTime   uint `json:"bigCarTime"`
	HelperPrice  uint `json:"helperPrice"`
	HelperTime   uint `json:"helperTime"`
	FragilePrice uint `json:"fragilePrice"`
	KM           uint `json:"kmPrice"`
}
