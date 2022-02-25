package stocks

type Stock struct {
	Name     string
	Price    float32
	Quantity float32
}

type Wallet struct {
	Amount float32
}

func (w *Wallet) Add(amount float32) {
	w.Amount += amount
}

func (w *Wallet) Subtract(amount float32) {
	w.Amount -= amount
}
