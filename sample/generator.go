package sample

import (
	"github.com/arcbjorn/store-management-system/pb/laptop"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func NewKeyboard() *laptop.Keyboard {
	keyboard := &laptop.Keyboard {
		Layout: randomKeyboadLayout(),
		Backlit: randomBool(),
	}

	return keyboard
}

func NewCPU() *laptop.CPU {
	brand := randomCPUBrand()
	name := randomCPUName(brand)

	coreNumber := randomInt(2, 8)
	threadNumber := randomInt(coreNumber, 12)

	minGhz := randomFloat64(2.0, 3.5)
	maxGhz := randomFloat64(minGhz, 5.0)

	cpu := &laptop.CPU{
		Brand: brand,
		Name: name,
		CoreNumber: uint32(coreNumber),
		ThreadNumber: uint32(threadNumber),
		MinGhz: minGhz,
		MaxGhz: maxGhz,
	}

	return cpu
}

func NewGPU() *laptop.GPU {
	brand := randomGPUBrand()
	name := randomGPUName(brand)

	minGhz := randomFloat64(1.0, 1.5)
	maxGhz := randomFloat64(minGhz, 2.0)

	memory := &laptop.Memory{
		Value: uint64(randomInt(2, 6)),
		Unit: laptop.Memory_MEGABYTE,
	}

	gpu := &laptop.GPU{
		Brand: brand,
		Name: name,
		MinGhz: minGhz,
		MaxGhz: maxGhz,
		Memory: memory,
	}

	return gpu
}

func NewRam() *laptop.Memory {
	ram := &laptop.Memory{
		Value: uint64(randomInt(4, 64)),
		Unit: laptop.Memory_GIGABYTE,
	}

	return ram
}

func NewSSD() *laptop.Storage {
	ssd := &laptop.Storage{
		Driver: laptop.Storage_SSD,
		Memory: &laptop.Memory{
			Value: uint64(randomInt(128, 1024)),
			Unit: laptop.Memory_GIGABYTE,
		},
	}

	return ssd
}

func NewHDD() *laptop.Storage {
	ssd := &laptop.Storage{
		Driver: laptop.Storage_HDD,
		Memory: &laptop.Memory{
			Value: uint64(randomInt(1, 6)),
			Unit: laptop.Memory_TERABYTE,
		},
	}

	return ssd
}

func NewScreen() *laptop.Screen {
	screen := &laptop.Screen{
		SizeInch: randomFloat32(13, 17),
		Resolution: randomScreenResolution(),
		Panel: randomScreenPanel(),
		Multitouch: randomBool(),
	}

	return screen
}

func NewLaptop() *laptop.Laptop {
	brand := randomLaptopBrand()
	name := randomLaptopName(brand)

	laptop := &laptop.Laptop{
		Id: randomID(),
		Brand: brand,
		Name: name,
		Cpu: NewCPU(),
		Gpus: []*laptop.GPU{NewGPU()},
		Storages: []*laptop.Storage{NewSSD(), NewHDD()},
		Screen: NewScreen(),
		Keyboard: NewKeyboard(),
		Weight: &laptop.Laptop_WeightKg{
			WeightKg: randomFloat64(1.0, 3.0),
		},
		PriceUsd: randomFloat64(1500, 3000),
		ReleaseYear: uint32(randomInt(2015, 2019)),
		UpdatedAt: timestamppb.Now(),
	}

	return laptop
}