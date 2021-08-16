package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type Product struct {
	ID            uint64    `gorm:"primary_key;auto_increment" json:"id"`
	Code          string    `gorm:"index:;size:255;not null;unique" json:"code"`
	Name          string    `gorm:"size:255;not null;" json:"name"`
	User          User      `json:"user"`
	UserID        uint32    `sql:"type:int REFERENCES users(id)" json:"user_id"`
	Description   string    `gorm:"size:255;not null;" json:"description"`
	DefaultPrice  uint32    `gorm:"type:int;" json:"default_price"`
	Image		  string	`gorm:"size:255;not null;" json:"pic_name"`
	CreatedAt     time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt     time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (p *Product) Prepare() {
	p.ID = 0
	p.Code = html.EscapeString(strings.TrimSpace(p.Code))
	p.Name = html.EscapeString(strings.TrimSpace(p.Name))
	p.User = User{}
	p.Description = html.EscapeString(strings.TrimSpace(p.Description))
	p.DefaultPrice = p.DefaultPrice
	//p.Image = html.EscapeString(strings.TrimSpace(p.Image))
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
}

func (p *Product) Validate() error {

	if p.Code == "" {
		return errors.New("Required Product Code")
	}
	if p.Name == "" {
		return errors.New("Required Product Name")
	}
	if p.UserID < 1 {
		return errors.New("Required User to Assign")
	}
	return nil
}

func (p *Product) SavePost(db *gorm.DB) (*Product, error) {
	var err error
	err = db.Debug().Model(&Product{}).Create(&p).Error
	if err != nil {
		return &Product{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", p.UserID).Take(&p.User).Error
		if err != nil {
			return &Product{}, err
		}
	}
	return p, nil
}

func (p *Product) FindAllProducts(db *gorm.DB) (*[]Product, error) {
	var err error
	products := []Product{}
	err = db.Debug().Model(&Product{}).Limit(100).Find(&products).Error
	if err != nil {
		return &[]Product{}, err
	}
	if len(products) > 0 {
		for i, _ := range products {
			err := db.Debug().Model(&User{}).Where("id = ?", products[i].UserID).Take(&products[i].User).Error
			if err != nil {
				return &[]Product{}, err
			}
		}
	}
	return &products, nil
}

func (p *Product) FindPostByID(db *gorm.DB, pid uint64) (*Product, error) {
	var err error
	err = db.Debug().Model(&Product{}).Where("id = ?", pid).Take(&p).Error
	if err != nil {
		return &Product{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", p.UserID).Take(&p.User).Error
		if err != nil {
			return &Product{}, err
		}
	}
	return p, nil
}

func (p *Product) UpdateAPost(db *gorm.DB) (*Product, error) {

	var err error
	err = db.Debug().Model(&Product{}).Where("id = ?", p.ID).Updates(Product{Code: p.Code, Name: p.Name, Description: p.Description, DefaultPrice: p.DefaultPrice, Image: p.Image, UpdatedAt: time.Now()}).Error
	if err != nil {
		return &Product{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", p.UserID).Take(&p.User).Error
		if err != nil {
			return &Product{}, err
		}
	}
	return p, nil
}

func (p *Product) DeleteAPost(db *gorm.DB, pid uint64, uid uint32) (int64, error) {

	db = db.Debug().Model(&Product{}).Where("id = ? and user_id = ?", pid, uid).Take(&Product{}).Delete(&Product{})

	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("Post not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}