package repository

type Bucket string

// Bucket - is name
const (
	AccessTokens  Bucket = "access_tocken"
	RequestTokens Bucket = "request_tocken"
)

type TokenRepositorier interface {
	Save(chatID int64, token string, bucket Bucket) error
	Get(chatID int64, bucket Bucket) (string, error)
}
