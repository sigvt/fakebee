package ytl

// Custom Faker definitions

import "github.com/pioz/faker"

var holomemChannels = []string{"UCDqI2jOz0weumE8s7paEk6g",
	"UC1CfXB_kRs3C-zaeTG3oGyg",
	"UChAnqc_AY5_I3Px5dig3X1Q",
	"UCyl1z3jo3XHR1riLFKG5UAg"}

var durations = []string{"1 month", "2 months", "1 year", "2 years"}

var holomemChannelIdBuilder = func(params ...string) (interface{}, error) {
	return faker.Pick(holomemChannels...), nil
}

var actionIdBuilder = func(params ...string) (interface{}, error) {
	return faker.StringWithSize(78) + "%3D", nil
}

var membershipSinceBuilder = func(params ...string) (interface{}, error) {
	return faker.Pick(durations...), nil
}

// Must be called once to register the defined builders
func RegisterBuilders() {
	var builders = []struct {
		name, typ string
		def       func(...string) (interface{}, error)
	}{
		{"holomemChannelId", "string", holomemChannelIdBuilder},
		{"actionId", "string", actionIdBuilder},
		{"membershipSince", "string", membershipSinceBuilder},
	}

	var err error

	for _, builder := range builders {
		err = faker.RegisterBuilder(builder.name, builder.typ, builder.def)
		if err != nil {
			panic(err)
		}
	}
}
