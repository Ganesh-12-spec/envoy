package crypto

type Secret struct {
	Ciphertext string `json:"ciphertext"`
	Nonce      string `json:"nonce"`
}
type Vault struct {
	Secrets map[string]Secret `json:"..."`
}
