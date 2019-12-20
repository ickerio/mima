package services

type Service interface {
	Start() error
	Stop() error
	Create() error
}
