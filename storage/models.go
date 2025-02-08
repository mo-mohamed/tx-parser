package store

type Transaction struct {
	// Hash is the unique identifier for this transaction.
	Hash string `json:"hash"`
	// From is the Ethereum address that initiated the transaction.
	From string `json:"from"`
	// To is the Ethereum address of the account that is the recipient of the transaction.
	To string `json:"to"`
	// Value is the transaction amount.
	Value string `json:"value"`
	// BlockNumber is the number of the transaction.
	BlockNumber string `json:"blockNumber"`
}
