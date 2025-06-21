package usecase

import "github.com/JeanGrijp/Full-Cycle-Desafio-Clean-Architecture/internal/domain"

type OrderRepository interface {
	List() ([]domain.Order, error)
}

type OrderUseCaseInterface interface {
	Execute() ([]domain.Order, error)
}

type ListOrdersUseCase struct {
	Repo OrderRepository
}

func (uc *ListOrdersUseCase) Execute() ([]domain.Order, error) {
	return uc.Repo.List()
}
