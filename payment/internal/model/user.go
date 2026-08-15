package model

type User struct {
	Uuid  string
	Login string
	Email string
}

type TransferData struct {
	Method    string
	Args      []string
	Nonce     string
	PublicKey string
	Signature string
}

type Transfer struct {
	RequestID  string
	TransferID string
	Channel    string
	Chaincode  string
	Data       TransferData
}
