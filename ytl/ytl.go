package ytl

// Youtube Live model definitions and helpers

import "time"

type Event interface {
	Chat | SuperChat | SuperSticker | Membership | Milestone | Ban | Deletion
}

// Regular live chat message
type Chat struct {
	ID string `faker:"actionId" json:"id"` // unique ID identifying a chat or other similar event
	Author
	Origin
	Message     string `faker:"ParagraphWithSentenceCount(1)" json:"msg"`
	IsModerator bool   `faker:"bool" json:"mod"`
	IsOwner     bool   `faker:"bool" json:"own"`
	IsVerified  bool   `faker:"bool" json:"vrf"`
	Membership  string `faker:"string" json:"mem"`
	YTTimestamp
}

type SuperChat struct {
	ID string `faker:"actionId" json:"id"`
	Author
	Origin
	Message      string  `faker:"ParagraphWithSentenceCount(1)" json:"msg"`
	Currency     string  `faker:"currencycode" json:"cur"`
	Amount       float64 `faker:"float64" json:"amo"`
	Significance int     `faker:"IntInRange(1,7)" json:"sig"`
	YTTimestamp
}

type SuperSticker struct {
	ID string `faker:"actionId" json:"id"`
	Author
	Origin
	Text string `faker:"string" json:"txt"`
	YTTimestamp
	Currency string  `faker:"currencycode" json:"cur"`
	Amount   float64 `faker:"Float64InRange(2, 20)" json:"amo"`
}

// Membership join chat message
type Membership struct {
	ID string `faker:"actionId" json:"id"`
	Author
	Origin
	YTTimestamp
	Level string `faker:"-" json:"lv"`
	Since string `faker:"membershipSince" json:"s"`
}

// Member milestone message
type Milestone struct {
	ID string `faker:"actionId" json:"id"`
	Author
	Origin
	Message string `faker:"ParagraphWithSentenceCount(1)" json:"msg"`
	Level   string `faker:"string" json:"lv"`
	Since   string `faker:"membershipSince" json:"s"`
	YTTimestamp
}

type Ban struct {
	ID        string `faker:"actionId" json:"oid"`
	ChannelId string `faker:"stringWithSize(24)" json:"cid"`
	Origin
	YTTimestamp
}

type Deletion struct {
	ID       string `faker:"actionId" json:"oid"`
	TargetId string `faker:"actionId" json:"tid"`
	Origin
	YTTimestamp
	Retracted bool `faker:"bool" json:"r"`
}

// Embedded models

type Author struct {
	ChannelId string `faker:"stringWithSize(24)" json:"cid"`
	Name      string `faker:"username" json:"aut"`
}

type Origin struct {
	ChannelId string `faker:"holomemChannelId" json:"ocid"`
	VideoId   string `faker:"stringWithSize(11)" json:"ovid"`
}

type YTTimestamp struct {
	Timestamp time.Time `faker:"customTimeStamp" json:"ts"`
}
