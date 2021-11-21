package sample

import (
	"github.com/arcbjorn/store-management-system/pb/pb"
	"github.com/golang/protobuf/ptypes"
)

func NewKeyboard() *pb.Keyboard {
	keyboard := &pb.Keyboard {
		Layout: randomKeyboadLayout(),
		Backlit: randomBool(),
	}

	return keyboard
}

func NewCPU() *pb.CPU {
	brand := randomCPUBrand()
	name := randomCPUName(brand)

	coreNumber := randomInt(2, 8)
	threadNumber := randomInt(coreNumber, 12)

	minGhz := randomFloat64(2.0, 3.5)
	maxGhz := randomFloat64(minGhz, 5.0)

	cpu := &pb.CPU{
		Brand: brand,
		Name: name,
		CoreNumber: uint32(coreNumber),
		ThreadNumber: uint32(threadNumber),
		MinGhz: minGhz,
		MaxGhz: maxGhz,
	}

	return cpu
}

func NewGPU() *pb.GPU {
	brand := randomGPUBrand()
	name := randomCPUName(brand)

	minGhz := randomFloat64(1.0, 1.5)
	maxGhz := randomFloat64(minGhz, 2.0)

	memory := &pb.Memory{
		Value: uint64(randomInt(2, 6)),
		Unit: pb.Memory_MEGABYTE,
	}

	gpu := &pb.GPU{
		Brand: brand,
		Name: name,
		MinGhz: minGhz,
		MaxGhz: maxGhz,
		Memory: memory,
	}

	return gpu
}

func NewRam() *pb.Memory {
	ram := &pb.Memory{
		Value: uint64(randomInt(4, 64)),
		Unit: pb.Memory_GIGABYTE,
	}

	return ram
}

func NewSSD() *pb.Storage {
	ssd := &pb.Storage{
		Driver: pb.Storage_SSD,
		Memory: &pb.Memory{
			Value: uint64(randomInt(128, 1024)),
			Unit: pb.Memory_GIGABYTE,
		},
	}

	return ssd
}

func NewHDD() *pb.Storage {
	ssd := &pb.Storage{
		Driver: pb.Storage_HDD,
		Memory: &pb.Memory{
			Value: uint64(randomInt(1, 6)),
			Unit: pb.Memory_TERABYTE,
		},
	}

	return ssd
}

func NewScreen() *pb.Screen {
	screen := &pb.Screen{
		SizeInch: randomFloat32(13, 17),
		Resolution: randomScreenResolution(),
		Panel: randomScreenPanel(),
		Multitouch: randomBool(),
	}

	return screen
}

func newLaptop() *pb.Laptop {
	brand := randomLaptopBrand()
	name := randomLaptopName(brand)

	laptop := &pb.Laptop{
		Id: randomID(),
		Brand: brand,
		Name: name,
		Cpu: NewCPU(),
		Gpus: []*pb.GPU{NewGPU()},
		Storages: []*pb.Storage{NewSSD(), NewHDD()},
		Screen: NewScreen(),
		Keyboard: NewKeyboard(),
		Weight: &pb.Laptop_WeightKg{
			WeightKg: randomFloat64(1.0, 3.0),
		},
		PriceUsd: randomFloat64(1500, 3000),
		ReleaseYear: uint32(randomInt(2015, 2019)),
		UpdatedAt: ptypes.TimestampNow(),
	}

	return laptop
}