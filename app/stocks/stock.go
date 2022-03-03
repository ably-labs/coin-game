package stocks

type Coin struct {
	Price    float32
	Quantity float32
}
type Wallet struct {
	Name   string
	Amount float32
}

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

func (p *Player) AddBalance(amount float32) {
	p.Wallet += amount
}

func (p *Player) SubtractBalance(amount float32) {
	p.Wallet -= amount
}

func (p *Player) AddQuantity(quantity float32) {
	p.CoinQuantity += quantity
}

func (p *Player) SubtractQuantity(quantity float32) {
	p.CoinQuantity -= quantity
}
