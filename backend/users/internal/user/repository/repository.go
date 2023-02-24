package repository

type Repository interface {
	Put(key, val string) error
	Get(key string) (res, err error)
}
