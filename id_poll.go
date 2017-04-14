package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"sync/atomic"
	"time"
)

var IDPollInited bool

var addInt int64
var previousTimeStamp int

type BaseConfig struct {
	NodeNumberFmt string //node or base device sign
	NodeNumber    string //base ID channel volume
	FmtNumber     int
	IDSerialMax   int64 //the max numbers of bass-id
}

const (
	UserID      int = 10
	CommodityID int = 20
	PartnerNO   int = 30
	OrderID     int = 40
	TransID     int = 50
	SerialID    int = 60
	OptSerialID int = 70
	PingID      int = 90
)

var (
	baseConfigure    BaseConfig
	NodeAddTimestemp int
)

func init() {
	baseConfigure.NodeNumberFmt = "%d1"
	baseConfigure.IDSerialMax = 9999
	baseConfigure.FmtNumber = 0 //01
	IDPollInited = false
}

// BaseIDChannel BaseID pool container
var BaseIDChannel = make(chan string, 1000)

// IDPoolStart start the id pool service
// you can init IDPool in main package with
// utils.IDPoolStart()
func IDPoolStart() {
	go BaseIDsPool()
	IDPollInited = true
}

//BaseIDsPool main IDPOOL channel handle
func BaseIDsPool() {
	for {
		BaseIDChannel <- getBaseID()
	}
}

func NewID(usageSign int) string {
	var usageSignStr string
	switch usageSign {
	case UserID:
		usageSignStr = "10"
	case CommodityID:
		usageSignStr = "20"
	case PartnerNO:
		usageSignStr = "30"
	case OrderID:
		usageSignStr = "40"
	case TransID:
		usageSignStr = "50" //Using for password creating
	case SerialID:
		usageSignStr = "60" //Serial id
	case OptSerialID:
		usageSignStr = "70"
	case PingID:
		usageSignStr = "90"
	}

	newID := usageSignStr + <-BaseIDChannel
	return newID
}

func NewPingID() string {
	return NewID(PingID)
}

func NewOptSerialID() string {
	return NewID(OptSerialID)
}

func NewPaymentSerialID() string {
	return NewID(SerialID)
}

func NewUserID() string {
	return NewID(UserID)
}

func NewCommodityID() string {
	return NewID(CommodityID)
}

func NewOrderID() string {
	return NewID(OrderID)
}

func NewPartnerID() string {
	return NewUserID()
}

func NewPartnerNO() string {
	return NewID(PartnerNO)
}

/* package static functions */

// create a timeStamp using for bass-ID-pools(10bit)
func getBaseID() string {
	timeStampInt := int(time.Now().Unix())
	if previousTimeStamp != timeStampInt {
		addIntInit()
	} else if addInt == baseConfigure.IDSerialMax+1 {
		if baseConfigure.FmtNumber == 8 {
			time.Sleep(time.Second * 1)
		}
		baseConfigure.FmtNumber++
		baseConfigure.NodeNumber = fmt.Sprintf(baseConfigure.NodeNumberFmt,
			baseConfigure.FmtNumber)
		addInt = 0
	}

	previousTimeStamp = timeStampInt
	timeStampInt -= 912008000
	timeStampStr := fmt.Sprintf("%010d", timeStampInt)

	atomic.AddInt64(&addInt, 1)
	newBassID := baseConfigure.NodeNumber +
		timeStampStr +
		fmt.Sprintf("%04d", addInt%baseConfigure.IDSerialMax)

	return newBassID
}

//addIntInit return the ID counter(in single second) to 0
func addIntInit() {
	addInt = 0
	baseConfigure.FmtNumber = 0
	baseConfigure.NodeNumber = fmt.Sprintf(baseConfigure.NodeNumberFmt,
		baseConfigure.FmtNumber)
}

// ShortDigitalID produces a "unique" 8 bytes long digital string.
// Do not use as a reliable way to get unique IDs, instead use for things like logging.
// Standerd: "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
// OurDigit: "0123456789012345678901234567890123456789012345678901234567890123"
func ShortDigitalID() string {
	encodeDigital := "0123456789012345678901234567890123456789012345678901234567890123"
	b := make([]byte, 6) // 8 bytes len after BASE64 encoded
	io.ReadFull(rand.Reader, b)
	digitalEncode := base64.NewEncoding(encodeDigital)
	return digitalEncode.EncodeToString(b)
}
