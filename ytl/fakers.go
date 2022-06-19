package ytl

// Custom Faker definitions

import (
	"github.com/pioz/faker"
)

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

var customTimeStamp = func(params ...string) (interface{}, error) {
	// twoYearsAgo := now.AddDate(-2, 0, 0)
	return faker.TimeNow(), nil
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
		{"customTimeStamp", "time.Time", customTimeStamp},
	}

	var err error

	for _, builder := range builders {
		err = faker.RegisterBuilder(builder.name, builder.typ, builder.def)
		if err != nil {
			panic(err)
		}
	}
}

var ChatFactory = func() interface{} {
	var c Chat
	faker.Build(&c)
	return c
}

var SuperChatFactory = func() interface{} {
	var sc SuperChat
	faker.Build(&sc)
	return sc
}

var SuperStickerFactory = func() interface{} {
	var ss SuperSticker
	faker.Build(&ss)
	return ss
}

var MembershipFactory = func() interface{} {
	var membership Membership
	faker.Build(&membership)
	return membership
}

var MilestoneFactory = func() interface{} {
	var milestone Milestone
	faker.Build(&milestone)
	return milestone
}

var BanFactory = func() interface{} {
	var ban Ban
	faker.Build(&ban)
	return ban
}

var DeletionFactory = func() interface{} {
	var deletion Deletion
	faker.Build(&deletion)
	return deletion
}
