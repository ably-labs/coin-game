package stocks

type Player struct {
	Name         string
	Wallet       float32
	CoinQuantity float32
}
type CoinResponse struct {
	Player        string
	WalletBalance float32
	CoinQuantity  float32
}

type BuySellRequest struct {
	Player           string
	CurrentCoinPrice float32
	Quantity         float32
}
