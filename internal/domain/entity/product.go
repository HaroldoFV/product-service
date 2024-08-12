package entity

import (
	"errors"
	"github.com/google/uuid"
)

const (
	DISABLED = "disabled"
	ENABLED  = "enabled"
)

type Product struct {
	id          string
	name        string
	description string
	price       float64
	status      string
}

func NewProduct(name, description string, price float64) (*Product, error) {
	product := &Product{
		id:          uuid.New().String(),
		name:        name,
		description: description,
		price:       price,
		status:      DISABLED,
	}
	err := product.IsValid()
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (p *Product) Update(name, description string) error {
	p.name = name
	p.description = description
	return p.IsValid()
}

func (p *Product) IsValid() error {
	if p.id == "" {
		return errors.New("invalid id")
	}
	if p.name == "" {
		return errors.New("name cannot be empty")
	}
	if len(p.name) > 100 {
		return errors.New("name cannot be longer than 100 characters")
	}
	if len(p.description) > 500 {
		return errors.New("description cannot be longer than 500 characters")
	}
	if p.status == "" {
		p.status = DISABLED
	}
	if p.status != ENABLED && p.status != DISABLED {
		return errors.New("status must be enabled or disabled")
	}
	if p.price < 0 {
		return errors.New("price must be greater or equal zero")
	}
	return nil
}

func (p *Product) Enable() error {
	if p.price > 0 {
		p.status = ENABLED
		return nil
	}
	err := p.IsValid()
	if err != nil {
		return err
	}
	return nil
}

func (p *Product) Disable() error {
	p.status = DISABLED
	err := p.IsValid()
	if err != nil {
		return err
	}
	return nil
}

func (p *Product) ChangePrice(price float64) error {
	p.price = price
	err := p.IsValid()
	if err != nil {
		return err
	}
	return nil
}

func (p *Product) GetID() string {
	return p.id
}

func (p *Product) GetName() string {
	return p.name
}

func (p *Product) GetDescription() string {
	return p.description
}

func (p *Product) GetStatus() string {
	return p.status
}

func (p *Product) GetPrice() float64 {
	return p.price
}

func (p *Product) SetID(id string) {
	p.id = id
}
