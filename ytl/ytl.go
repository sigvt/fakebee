package ytl

// Youtube Live model definitions and helpers

import "time"

// Regular live chat message
type Chat struct {
	ID string `faker:"uuid;unique" json:"id"`
	Author
	Origin
	Message     string `faker:"ParagraphWithSentenceCount(1)" json:"msg"`
	IsModerator bool   `faker:"bool" json:"mod"`
	IsOwner     bool   `faker:"bool" json:"own"`
	IsVerified  bool   `faker:"bool" json:"vrf"`
	Membership  string `bson:"membership" json:"mem"`
	YTTimestamp
}

type SuperChat struct {
	ID string `bson:"id" json:"id"`
	Author
	Origin
	Message      string  `bson:"message" json:"msg"`
	Currency     string  `bson:"currency" json:"cur"`
	Amount       float64 `bson:"purchaseAmount" json:"amo"`
	Significance int     `bson:"significance" json:"sig"`
	YTTimestamp
}

type SuperSticker struct {
	ID string `bson:"id" json:"id"`
	Author
	Origin
	Text      string    `bson:"text" json:"txt"`
	Timestamp time.Time `bson:"timestamp" json:"ts"`
	Currency  string    `bson:"currency" json:"cur"`
	Amount    float64   `bson:"purchaseAmount" json:"amo"`
}

// Membership join chat message
type Membership struct {
	ID string `bson:"id" json:"id"`
	Author
	Origin
	Timestamp time.Time `bson:"timestamp" json:"ts"`
	Level     string    `bson:"level" json:"lv"`
	Since     string    `bson:"since" json:"s"`
}

// Member milestone message
type Milestone struct {
	ID string `bson:"id" json:"id"`
	Author
	Origin
	Message string `bson:"message" json:"msg"`
	Level   string `bson:"level" json:"lv"`
	Since   string `bson:"since" json:"s"`
	YTTimestamp
}

type Ban struct {
	ID        string `bson:"_id" json:"oid"`
	ChannelId string `bson:"channelId" json:"cid"`
	Origin
	YTTimestamp
}

type Deletion struct {
	ID       string `bson:"_id" json:"oid"`
	TargetId string `bson:"targetId" json:"tid"`
	Origin
	YTTimestamp
	Retracted bool `bson:"retracted" json:"r"`
}

type Author struct {
	ChannelId string `bson:"authorChannelId" json:"cid"`
	Name      string `faker:"username" json:"aut"`
}

type Origin struct {
	ChannelId string `bson:"originChannelId" json:"ocid"`
	VideoId   string `bson:"originVideoId" json:"ovid"`
}

type YTTimestamp struct {
	Timestamp time.Time `faker:"time" json:"ts"`
}
