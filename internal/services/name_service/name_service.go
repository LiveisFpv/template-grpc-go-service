package name_service

import (
	"context"
	"template-grpc-go-service/internal/domain/models"
	"time"

	"github.com/sirupsen/logrus"
)

// All methods
type NameStorage interface {
	GetNamebyID(ctx context.Context, name_id int) (name *models.Name, err error)
	GetAllName(ctx context.Context, pagination *models.Pagination, filter []*models.Filter, orderby []*models.Sort) ([]*models.Name, *models.Pagination, error)
	CreateName(ctx context.Context, name_title, name_capital, name_area string) (name *models.Name, err error)
	UpdateNamebyID(ctx context.Context, name *models.Name) (err error)
	DeleteNamebyID(ctx context.Context, name_id int) (name *models.Name, err error)
}

type NameService struct {
	log         *logrus.Logger
	nameStorage NameStorage
	tokenTTL    time.Duration
}

// Constructor service of Name
func New(
	log *logrus.Logger,
	nameStorage NameStorage,
	tokenTTL time.Duration,
) *NameService {
	return &NameService{
		log:         log,
		nameStorage: nameStorage,
		tokenTTL:    tokenTTL,
	}
}

// Add_Name implements namegrpc.Name.
func (c *NameService) Add_Name(ctx context.Context, name_title, name_capital, name_area string) (name *models.Name, err error) {
	const op = "Name.Create"
	log := c.log.WithFields(
		logrus.Fields{
			"op":      op,
			"title":   name_title,
			"capital": name_capital,
			"area":    name_area,
		},
	)
	log.Info("Start Create Name")

	name, err = c.nameStorage.CreateName(ctx, name_title, name_capital, name_area)
	if err != nil {
		c.log.Error("failed to create name", err)
		return nil, err
	}

	return name, nil
}

// Delete_NamebyID implements namegrpc.Name.
func (c *NameService) Delete_NamebyID(ctx context.Context, name_id int) (*models.Name, error) {
	const op = "Name.Delete"
	log := c.log.WithFields(
		logrus.Fields{
			"op": op,
			"id": name_id,
		},
	)
	log.Info("Start Delete Name")
	res, err := c.nameStorage.DeleteNamebyID(ctx, name_id)
	if err != nil {
		c.log.Error("failed to delete name", err)
		return nil, err
	}
	return res, nil
}

// Get_All_Name implements namegrpc.Name.
func (c *NameService) Get_All_Name(ctx context.Context, pagination *models.Pagination, filter []*models.Filter, orderby []*models.Sort) ([]*models.Name, *models.Pagination, error) {
	const op = "Name.GetAll"
	log := c.log.WithFields(
		logrus.Fields{
			"op": op,
		},
	)
	log.Info("Start Get ALL Name")

	countries, new_pagination, err := c.nameStorage.GetAllName(ctx, pagination, filter, orderby)
	if err != nil {
		c.log.Error("failed to get all countries", err)
		return nil, nil, err
	}

	return countries, new_pagination, nil
}

// Get_NamebyID implements namegrpc.Name.
func (c *NameService) Get_NamebyID(ctx context.Context, name_id int) (name *models.Name, err error) {
	const op = "Name.GetbyID"
	log := c.log.WithFields(
		logrus.Fields{
			"op": op,
			"id": name_id,
		},
	)
	log.Info("Start Get by ID Name")
	name, err = c.nameStorage.GetNamebyID(ctx, name_id)
	if err != nil {
		c.log.Error("failed to get name by id", err)
		return nil, err
	}
	return name, nil
}

// Update_NamebyID implements namegrpc.Name.
func (c *NameService) Update_NamebyID(ctx context.Context, name *models.Name) (err error) {
	const op = "Name.Update"
	log := c.log.WithFields(
		logrus.Fields{
			"op":      op,
			"id":      name.Name_id,
			"title":   name.Name_title,
			"capital": name.Name_capital,
			"area":    name.Name_area,
		},
	)
	log.Info("Start Update name")
	err = c.nameStorage.UpdateNamebyID(ctx, name)
	if err != nil {
		c.log.Error("failed to update name", err)
		return err
	}
	return nil
}
