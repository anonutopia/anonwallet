package main

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
)

const satInBtc = uint64(100000000)

type Transaction struct {
	gorm.Model
	TxId      string `sql:"size:255"`
	Processed bool   `sql:"DEFAULT:false"`
}

type User struct {
	gorm.Model
	Nickname                string `sql:"size:255;unique_index"`
	Email                   string `sql:"size:255;unique_index"`
	Country                 string `sql:"size:255"`
	City                    string `sql:"size:255"`
	Address                 string `sql:"size:255;unique_index"`
	PasswordHash            string `sql:"size:255"`
	Referral                string `sql:"size:255"`
	BitcoinAddr             string `sql:"size:255"`
	BitcoinBalanceNew       int    `sql:"DEFAULT:0"`
	BitcoinBalanceProcessed int    `sql:"DEFAULT:0"`
	EtherAddr               string `sql:"size:255"`
	EtherBalanceNew         int    `sql:"DEFAULT:0"`
	EtherBalanceProcessed   int    `sql:"DEFAULT:0"`
	NextFacebookAward       int    `sql:"DEFAULT:100000000"`
	LastFacebookAwardTime   *time.Time
	ProfitEth               uint64
	ProfitWav               uint64
	ProfitBtc               uint64
	ProfitEthTotal          uint64
	ProfitWavTotal          uint64
	ProfitBtcTotal          uint64
	ReferralProfitEth       uint64
	ReferralProfitWav       uint64
	ReferralProfitBtc       uint64
	ReferralProfitEthTotal  uint64
	ReferralProfitWavTotal  uint64
	ReferralProfitBtcTotal  uint64
	ReceivedFreeAnote       bool    `sql:"DEFAULT:false"`
	EmailVerified           bool    `sql:"DEFAULT:false"`
	TelegramId              int     `sql:"DEFAULT:0"`
	Badges                  []Badge `gorm:"many2many:user_badges;"`
}

func (u *User) ProfitWavString() string {
	return fmt.Sprintf("%.8f", float64(u.ProfitWav)/float64(100000000))
}

func (u *User) ProfitBtcString() string {
	return fmt.Sprintf("%.8f", float64(u.ProfitBtc)/float64(100000000))
}

func (u *User) ProfitEthString() string {
	return fmt.Sprintf("%.8f", float64(u.ProfitEth)/float64(100000000))
}

func (u *User) ReferralProfitWavString() string {
	return fmt.Sprintf("%.8f", float64(u.ReferralProfitWav)/float64(100000000))
}

func (u *User) ReferralProfitBtcString() string {
	return fmt.Sprintf("%.8f", float64(u.ReferralProfitBtc)/float64(100000000))
}

func (u *User) ReferralProfitEthString() string {
	return fmt.Sprintf("%.8f", float64(u.ReferralProfitEth)/float64(100000000))
}

func (u *User) HasBadges() bool {
	var badges []*Badge
	db.Model(u).Association("Badges").Find(&badges)
	return len(badges) > 0
}

func (u *User) IsFounder() bool {
	return u.HasBadge("founder")
}

func (u *User) FacebookShareEnabled() bool {
	if u.LastFacebookAwardTime == nil && u.NextFacebookAward > 0 {
		return true
	} else if u.LastFacebookAwardTime.Add(24*time.Hour).Before(time.Now()) && u.NextFacebookAward > 0 {
		return true
	}

	return false
}

func (u *User) HasBadge(badge string) bool {
	var badges []*Badge
	db.Model(u).Association("Badges").Find(&badges)
	for _, b := range badges {
		if b.Name == badge {
			return true
		}
	}
	return false
}

func (u *User) ReferredUsers() []*User {
	var users []*User
	db.Where("referral = ?", u.Address).Find(&users)
	return users
}

func (u *User) ReferredUsersVerifiedCount() int {
	count := 0
	users := u.ReferredUsers()

	for _, user := range users {
		if user.EmailVerified {
			count++
		}
	}

	return count
}

type Badge struct {
	gorm.Model
	Name  string `sql:"size:255;unique_index"`
	AUBIX uint16
}

type KeyValue struct {
	gorm.Model
	Key   string `sql:"size:255;unique_index"`
	Value uint64 `sql:"type:bigint"`
}

type UsedAddress struct {
	gorm.Model
	Address string `sql:"size:255;unique_index"`
	Type    uint8
	UserID  int
	User    User
	Balance uint64 `sql:"type:bigint"`
}

type SharePost struct {
	gorm.Model
	URI       string `sql:"size:255;unique_index"`
	Processed bool   `sql:"DEFAULT:false"`
}
