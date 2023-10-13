package daos

import (
	"database/sql"
	"errors"
	"github.com/mahendraintelops/intelops-procurement-solution-v10/payment-service/pkg/rest/server/daos/clients/sqls"
	"github.com/mahendraintelops/intelops-procurement-solution-v10/payment-service/pkg/rest/server/models"
	log "github.com/sirupsen/logrus"
)

type PaymentDao struct {
	sqlClient *sqls.SQLiteClient
}

func migratePayments(r *sqls.SQLiteClient) error {
	query := `
	CREATE TABLE IF NOT EXISTS payments(
		Id INTEGER PRIMARY KEY AUTOINCREMENT,
        
		Amount REAL NOT NULL,
		InvoiceNumber TEXT NOT NULL,
        CONSTRAINT id_unique_key UNIQUE (Id)
	)
	`
	_, err1 := r.DB.Exec(query)
	return err1
}

func NewPaymentDao() (*PaymentDao, error) {
	sqlClient, err := sqls.InitSqliteDB()
	if err != nil {
		return nil, err
	}
	err = migratePayments(sqlClient)
	if err != nil {
		return nil, err
	}
	return &PaymentDao{
		sqlClient,
	}, nil
}

func (paymentDao *PaymentDao) CreatePayment(m *models.Payment) (*models.Payment, error) {
	insertQuery := "INSERT INTO payments(Amount, InvoiceNumber)values(?, ?)"
	res, err := paymentDao.sqlClient.DB.Exec(insertQuery, m.Amount, m.InvoiceNumber)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	m.Id = id

	log.Debugf("payment created")
	return m, nil
}

func (paymentDao *PaymentDao) UpdatePayment(id int64, m *models.Payment) (*models.Payment, error) {
	if id == 0 {
		return nil, errors.New("invalid payment ID")
	}
	if id != m.Id {
		return nil, errors.New("id and payload don't match")
	}

	payment, err := paymentDao.GetPayment(id)
	if err != nil {
		return nil, err
	}
	if payment == nil {
		return nil, sql.ErrNoRows
	}

	updateQuery := "UPDATE payments SET Amount = ?, InvoiceNumber = ? WHERE Id = ?"
	res, err := paymentDao.sqlClient.DB.Exec(updateQuery, m.Amount, m.InvoiceNumber, id)
	if err != nil {
		return nil, err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rowsAffected == 0 {
		return nil, sqls.ErrUpdateFailed
	}

	log.Debugf("payment updated")
	return m, nil
}

func (paymentDao *PaymentDao) DeletePayment(id int64) error {
	deleteQuery := "DELETE FROM payments WHERE Id = ?"
	res, err := paymentDao.sqlClient.DB.Exec(deleteQuery, id)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sqls.ErrDeleteFailed
	}

	log.Debugf("payment deleted")
	return nil
}

func (paymentDao *PaymentDao) ListPayments() ([]*models.Payment, error) {
	selectQuery := "SELECT * FROM payments"
	rows, err := paymentDao.sqlClient.DB.Query(selectQuery)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)
	var payments []*models.Payment
	for rows.Next() {
		m := models.Payment{}
		if err = rows.Scan(&m.Id, &m.Amount, &m.InvoiceNumber); err != nil {
			return nil, err
		}
		payments = append(payments, &m)
	}
	if payments == nil {
		payments = []*models.Payment{}
	}

	log.Debugf("payment listed")
	return payments, nil
}

func (paymentDao *PaymentDao) GetPayment(id int64) (*models.Payment, error) {
	selectQuery := "SELECT * FROM payments WHERE Id = ?"
	row := paymentDao.sqlClient.DB.QueryRow(selectQuery, id)
	m := models.Payment{}
	if err := row.Scan(&m.Id, &m.Amount, &m.InvoiceNumber); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sqls.ErrNotExists
		}
		return nil, err
	}

	log.Debugf("payment retrieved")
	return &m, nil
}
