package db

import (
	"fmt"

	domain "disfactory/imgur-backup/domain"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
)

type PGParameters struct {
	Account  string
	Password string
	Host     string
	Port     string
	DBName   string
}

type DBFactoryImageRepository struct {
	db     *sqlx.DB
	params PGParameters
}

func (pg *PGParameters) Dialect() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", pg.Account, pg.Password, pg.Host, pg.Port, pg.DBName)
}

func NewDBFactoryImageRepository(params PGParameters) (*DBFactoryImageRepository, error) {
	repo := DBFactoryImageRepository{}

	db, err := sqlx.Connect("postgres", params.Dialect())
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}
	repo.db = db
	repo.params = params

	return &repo, nil
}

func (repo *DBFactoryImageRepository) Close() {
	repo.db.Close()
}

func (repo *DBFactoryImageRepository) GetImages(size int, offset int) ([]domain.FactoryImage, error) {
	var images []domain.FactoryImage
	err := repo.db.Select(&images, "SELECT id, image_path FROM factory_image LIMIT $1 OFFSET $2;", size, offset)
	if err != nil {
		return nil, err
	}
	return images, nil
}
