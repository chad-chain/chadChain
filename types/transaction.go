package types

type Transaction struct {
	To     string
	From   string
	Amount int
}

func (t *Transaction) String() string {
	res := "Transaction{\n"
	res += "  To: " + t.To + "\n"
	res += "  From: " + t.From + "\n"
	res += "  Amount: " + string(rune(t.Amount)) + "\n"
	res += "}"
	return res
}
