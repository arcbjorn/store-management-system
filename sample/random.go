package sample

import (
	"math/rand"
	"time"

	"github.com/arcbjorn/store-management-system/pb/laptop"
	"github.com/google/uuid"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func randomKeyboadLayout() laptop.Keyboard_Layout {
	switch rand.Intn(3) {
		case 1:
			return laptop.Keyboard_QWERTY
		case 2:
			return laptop.Keyboard_QWERTZ
		default:
			return laptop.Keyboard_AZERTY
	}
}

func randomCPUBrand() string {
	return randomStringFromSet("Intel", "AMD")
}

func randomCPUName(brand string) string {
	if brand == "Intel" {
		return randomStringFromSet(
			"Xeon E-2286M",
			"Core i9-9980HK",
			"Core i7-9750H",
			"Core i5-9400F",
			"Core i3-1005G1",
		)
	}

	return randomStringFromSet(
			"Ryzen 7 PRO 2700U",
			"Ryzen 7 PRO 3500U",
			"Ryzen 7 PRO 3200GE",
		)
}

func randomGPUBrand() string {
	return randomStringFromSet(
			"NVIDIA",
			"AMD",
		)
}

func randomGPUName(brand string) string {
	if brand == "NVIDIA" {
		return randomStringFromSet(
			"RTX 2060",
			"RTX 2070",
			"RTX 1660-Ti",
			"RTX 1070",
		)
	}

	return randomStringFromSet(
			"RX 590",
			"RX 580",
			"RX 5700-XT",
			"RX Vega-56",
		)
}

func randomLaptopBrand() string {
	return randomStringFromSet(
			"Apple",
			"Dell",
			"Lenovo",
		)
}

func randomLaptopName(brand string) string {
	switch brand {
	case "Apple":
		return randomStringFromSet("Macbook Air", "Macbook Pro")
	case "Dell":
		return randomStringFromSet("Latitude", "Vostro", "XPS", "Alienware")
	default:
		return randomStringFromSet("Thinkpad X1", "Thinkpad P1", "Thinkpad P53" )
	}
}

func randomScreenResolution() *laptop.Screen_Resolution {
	height := randomInt(1080, 4320)
	width := height * 16 / 9

	resolution := &laptop.Screen_Resolution{
		Height: uint32(height),
		Width: uint32(width),
	}

	return resolution
}

func randomScreenPanel() laptop.Screen_Panel {
	if rand.Intn(2) == 1 {
		return laptop.Screen_IPS
	}

	return laptop.Screen_OLED
}

func randomStringFromSet(a ...string) string {
	n := len(a)
	if n == 0 {
		return ""
	}
	return a[rand.Intn(n)]
}

func randomBool() bool {
	return rand.Intn(2) == 1
}

func randomInt(min int, max int) int {
	return min + rand.Intn(max-min+1)
}

func randomFloat64(min float64, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func randomFloat32(min float32, max float32) float32 {
	return min + rand.Float32()*(max-min)
}

func randomID() string {
	return uuid.New().String()
}