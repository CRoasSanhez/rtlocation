package services

import (
	"rtlocation/repositories"
	"rtlocation/models"
)

// PaymentService ...
type PaymentService interface{
	AddPayment(userIDHex string, payment models.Payment)(models.Payment)
	DeletePayment(userIDHex string, payment models.Payment)(error)
}

// NewPaymentService ...
func NewPaymentService(repo repositories.PaymentRepository)PaymentService{
	return &mPaymentService{
		repo: repo,
	}
}

type mPaymentService struct{
	repo repositories.PaymentRepository
}

// AddPayment ...
func(s *mPaymentService) AddPayment(userIDHex string, payment models.Payment)models.Payment{

	s.repo.AddPayment(userIDHex, payment)

	return models.Payment{}
}

// DeletePayment ...
func(s *mPaymentService)DeletePayment(userIDHex string, payment models.Payment)(error){
	return nil
}