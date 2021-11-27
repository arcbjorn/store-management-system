package services

import (
	"errors"
	"fmt"
	"sync"

	"github.com/arcbjorn/store-management-system/pb/laptop"
	"github.com/jinzhu/copier"
)

var ErrAlreadyExists = errors.New("record already exists")

type LaptopStore interface {
	Save(laptop *laptop.Laptop) error
	Find(id string) (*laptop.Laptop, error)
}

type InMemoryLaptopStore struct {
	mutex sync.RWMutex
	data  map[string]*laptop.Laptop
}

func NewInMemoryLaptopStore() *InMemoryLaptopStore {
	return &InMemoryLaptopStore{
		data: make(map[string]*laptop.Laptop),
	}
}

func (store *InMemoryLaptopStore) Save(laptopDto *laptop.Laptop) error {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	if store.data[laptopDto.Id] != nil {
		return ErrAlreadyExists
	}

	// deep copy
	other := &laptop.Laptop{}
	err := copier.Copy(other, laptopDto)
	if err != nil {
		return fmt.Errorf("cannot copy laptop data: %w", err)
	}

	store.data[other.Id] = other
	return nil
}

func (store *InMemoryLaptopStore) Find(id string) (*laptop.Laptop, error) {
	store.mutex.RLock()
	defer store.mutex.RLock()

	foundLaptop := store.data[id]
	if foundLaptop == nil {
		return nil, nil
	}

	other := &laptop.Laptop{}
	err := copier.Copy(other, foundLaptop)
	if err != nil {
		return nil, fmt.Errorf("cannot copy laptop data: %w", err)
	}

	return other, nil
}
