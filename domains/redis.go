package domains

type RedisRepository interface {
	Get(key string) (result string, err error)
	Set(key string, value interface{}) (err error)
}
