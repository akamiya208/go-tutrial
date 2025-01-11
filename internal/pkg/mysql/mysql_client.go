package mysql

import (
	"github.com/akamiya208/go-tutrial/internal/pkg/models"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type IClient interface {
	GetTask(ID uint) (models.Task, error)
	GetTasksByName(name string) ([]models.Task, error)
	CreateTask(task *models.Task) error
	UpdateTask(task *models.Task) error
	DeleteTask(task *models.Task) error
	DB() *gorm.DB
}

var _ IClient = &Client{}

type Client struct {
	db *gorm.DB
}

func NewMySQLClient() (IClient, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")
	if err := viper.ReadInConfig(); err != nil {
		panic(err.Error())
	}

	dsn := viper.GetString("mysql.dsn")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	return &Client{db: db}, nil
}

func (c *Client) GetTask(ID uint) (models.Task, error) {
	var task models.Task
	if err := c.db.First(&task, ID).Error; err != nil {
		return models.Task{}, err
	}
	return task, nil
}

func (c *Client) GetTasksByName(name string) ([]models.Task, error) {
	var tasks []models.Task
	if err := c.db.Where("name = ?", name).Find(&tasks).Error; err != nil {
		return []models.Task{}, err
	}
	return tasks, nil
}

func (c *Client) CreateTask(task *models.Task) error {
	return c.db.Create(&task).Error
}

func (c *Client) UpdateTask(task *models.Task) error {
	return c.db.Save(&task).Error
}

func (c *Client) DeleteTask(task *models.Task) error {
	return c.db.Delete(&task).Error
}

func (c *Client) DB() *gorm.DB {
	return c.db
}
