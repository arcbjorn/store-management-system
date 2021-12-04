package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/arcbjorn/store-management-system/pb/laptop"
	"github.com/jinzhu/copier"
)

var ErrAlreadyExists = errors.New("record already exists")

type LaptopStore interface {
	Save(laptop *laptop.Laptop) error
	Find(id string) (*laptop.Laptop, error)
	Search(ctx context.Context, filter *laptop.Filter, found func(laptop *laptop.Laptop) error) error
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

	other, err := deepCopy(laptopDto)
	if err != nil {
		return err
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

	return deepCopy(foundLaptop)
}

func (store *InMemoryLaptopStore) Search(ctx context.Context, filter *laptop.Filter, found func(laptop *laptop.Laptop) error) error {
	store.mutex.RLock()
	defer store.mutex.RUnlock()

	for _, lp := range store.data {
		time.Sleep(time.Second)
		log.Print("checking laptop id: ", lp.GetId())

		if ctx.Err() == context.Canceled || ctx.Err() == context.DeadlineExceeded {
			log.Printf("context is cancelled")
			return errors.New("context is cancelled")
		}

		if isQualified(filter, lp) {
			other, err := deepCopy(lp)
			if err != nil {
				return err
			}

			err = found(other)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func isQualified(filter *laptop.Filter, laptop *laptop.Laptop) bool {
	if laptop.GetPriceUsd() > filter.GetMaxPriceUsd() {
		return false
	}

	if laptop.GetCpu().GetCoreNumber() < filter.MinCpuCores {
		return false
	}

	if laptop.GetCpu().GetMinGhz() < filter.GetMinCpuGhz() {
		return false
	}

	if toBit(laptop.GetRam()) < toBit(filter.GetMinRam()) {
		return false
	}

	return true
}

func toBit(memory *laptop.Memory) uint64 {
	value := memory.GetValue()

	switch memory.GetUnit() {
	case laptop.Memory_BIT:
		return value
	case laptop.Memory_BYTE:
		return value << 3
	case laptop.Memory_KILOBYTE:
		return value << 13
	case laptop.Memory_MEGABYTE:
		return value << 23
	case laptop.Memory_GIGABYTE:
		return value << 33
	case laptop.Memory_TERABYTE:
		return value << 43
	default:
		return 0
	}
}

func deepCopy(lt *laptop.Laptop) (*laptop.Laptop, error) {
	other := &laptop.Laptop{}
	err := copier.Copy(other, lt)
	if err != nil {
		return nil, fmt.Errorf("cannot copy laptop data: %w", err)
	}

	return other, nil
}
