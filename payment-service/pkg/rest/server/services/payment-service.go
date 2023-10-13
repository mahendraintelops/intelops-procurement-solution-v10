package services

import (
	"github.com/mahendraintelops/intelops-procurement-solution-v10/payment-service/pkg/rest/server/daos"
	"github.com/mahendraintelops/intelops-procurement-solution-v10/payment-service/pkg/rest/server/models"
)

type PaymentService struct {
	paymentDao *daos.PaymentDao
}

func NewPaymentService() (*PaymentService, error) {
	paymentDao, err := daos.NewPaymentDao()
	if err != nil {
		return nil, err
	}
	return &PaymentService{
		paymentDao: paymentDao,
	}, nil
}

func (paymentService *PaymentService) CreatePayment(payment *models.Payment) (*models.Payment, error) {
	return paymentService.paymentDao.CreatePayment(payment)
}

func (paymentService *PaymentService) UpdatePayment(id int64, payment *models.Payment) (*models.Payment, error) {
	return paymentService.paymentDao.UpdatePayment(id, payment)
}

func (paymentService *PaymentService) DeletePayment(id int64) error {
	return paymentService.paymentDao.DeletePayment(id)
}

func (paymentService *PaymentService) ListPayments() ([]*models.Payment, error) {
	return paymentService.paymentDao.ListPayments()
}

func (paymentService *PaymentService) GetPayment(id int64) (*models.Payment, error) {
	return paymentService.paymentDao.GetPayment(id)
}
