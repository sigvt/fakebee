package ytl

import (
	"github.com/pioz/faker"
)

var HolomemChannels = []string{"UCDqI2jOz0weumE8s7paEk6g",
	"UC1CfXB_kRs3C-zaeTG3oGyg",
	"UChAnqc_AY5_I3Px5dig3X1Q",
	"UCyl1z3jo3XHR1riLFKG5UAg"}

var HolomemChannelIdBuilder = func(params ...string) (interface{}, error) {
	return faker.Pick(HolomemChannels...), nil
}
