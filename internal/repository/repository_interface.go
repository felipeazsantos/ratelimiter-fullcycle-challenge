package repository

type IRepositoryLimiter interface {
	Save() error
	GetByIp() error
	GetByToken() error
}
