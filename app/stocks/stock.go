package stocks

type Coin struct {
	Price    float32
	Quantity float32
}

type Wallet struct {
	Amount float32
}

type CoinResponse struct {
	WalletBalance float32
	Quantity      float32
}

func (w *Wallet) Add(amount float32) {
	w.Amount += amount
}

func (w *Wallet) Subtract(amount float32) {
	w.Amount -= amount
}

func (c *Coin) AddQuantity(quantity float32) {
	c.Quantity += quantity
}
