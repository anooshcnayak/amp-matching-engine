package types

import (
	"encoding/json"
	"math/big"
	"time"

	"github.com/Proofsuite/amp-matching-engine/utils/math"
	"github.com/ethereum/go-ethereum/common"

	validation "github.com/go-ozzo/ozzo-validation"
	"gopkg.in/mgo.v2/bson"
)

// Pair struct is used to model the pair data in the system and DB
type Pair struct {
	ID                 bson.ObjectId  `json:"-" bson:"_id"`
	BaseTokenSymbol    string         `json:"baseTokenSymbol,omitempty" bson:"baseTokenSymbol"`
	BaseTokenAddress   common.Address `json:"baseTokenAddress,omitempty" bson:"baseTokenAddress"`
	BaseTokenDecimals  int            `json:"baseTokenDecimals,omitempty" bson:"baseTokenDecimals"`
	QuoteTokenSymbol   string         `json:"quoteTokenSymbol,omitempty" bson:"quoteTokenSymbol"`
	QuoteTokenAddress  common.Address `json:"quoteTokenAddress,omitempty" bson:"quoteTokenAddress"`
	QuoteTokenDecimals int            `json:"quoteTokenDecimals,omitempty" bson:"quoteTokenDecimals"`
	Active             bool           `json:"active,omitempty" bson:"active"`
	MakeFee            *big.Int       `json:"makeFee,omitempty" bson:"makeFee"`
	TakeFee            *big.Int       `json:"takeFee,omitempty" bson:"takeFee"`
	CreatedAt          time.Time      `json:"-" bson:"createdAt"`
	UpdatedAt          time.Time      `json:"-" bson:"updatedAt"`
}

func (p *Pair) MarshalJSON() ([]byte, error) {
	pair := map[string]interface{}{
		"baseTokenSymbol":    p.BaseTokenSymbol,
		"baseTokenDecimals":  p.BaseTokenDecimals,
		"quoteTokenSymbol":   p.QuoteTokenSymbol,
		"quoteTokenDecimals": p.QuoteTokenDecimals,
		"baseTokenAddress":   p.BaseTokenAddress,
		"quoteTokenAddress":  p.QuoteTokenAddress,
		"active":             p.Active,
	}

	if p.MakeFee != nil {
		pair["makeFee"] = p.MakeFee.String()
	}

	if p.TakeFee != nil {
		pair["takeFee"] = p.TakeFee.String()
	}

	return json.Marshal(pair)
}

type PairAddresses struct {
	Name       string         `json:"name" bson:"name"`
	BaseToken  common.Address `json:"baseToken" bson:"baseToken"`
	QuoteToken common.Address `json:"quoteToken" bson:"quoteToken"`
}

type PairAddressesRecord struct {
	Name       string `json:"name" bson:"name"`
	BaseToken  string `json:"baseToken" bson:"baseToken"`
	QuoteToken string `json:"quoteToken" bson:"quoteToken"`
}

type PairRecord struct {
	ID bson.ObjectId `json:"id" bson:"_id"`

	BaseTokenSymbol    string    `json:"baseTokenSymbol" bson:"baseTokenSymbol"`
	BaseTokenAddress   string    `json:"baseTokenAddress" bson:"baseTokenAddress"`
	BaseTokenDecimals  int       `json:"baseTokenDecimals" bson:"baseTokenDecimals"`
	QuoteTokenSymbol   string    `json:"quoteTokenSymbol" bson:"quoteTokenSymbol"`
	QuoteTokenAddress  string    `json:"quoteTokenAddress" bson:"quoteTokenAddress"`
	QuoteTokenDecimals int       `json:"quoteTokenDecimals" bson:"quoteTokenDecimals"`
	Active             bool      `json:"active" bson:"active"`
	MakeFee            string    `json:"makeFee" bson:"makeFee"`
	TakeFee            string    `json:"takeFee" bson:"takeFee"`
	CreatedAt          time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt          time.Time `json:"updatedAt" bson:"updatedAt"`
}

func (p *Pair) BaseTokenMultiplier() *big.Int {
	return math.Exp(big.NewInt(10), big.NewInt(int64(p.BaseTokenDecimals)))
}

func (p *Pair) QuoteTokenMultiplier() *big.Int {
	return math.Exp(big.NewInt(10), big.NewInt(int64(p.QuoteTokenDecimals)))
}

func (p *Pair) PairMultiplier() *big.Int {
	defaultMultiplier := math.Exp(big.NewInt(10), big.NewInt(18))
	baseTokenMultiplier := math.Exp(big.NewInt(10), big.NewInt(int64(p.BaseTokenDecimals)))

	return math.Mul(defaultMultiplier, baseTokenMultiplier)
}

func (p *Pair) Code() string {
	code := p.BaseTokenSymbol + "/" + p.QuoteTokenSymbol + "::" + p.BaseTokenAddress.Hex() + "::" + p.QuoteTokenAddress.Hex()
	return code
}

func (p *Pair) AddressCode() string {
	code := p.BaseTokenAddress.Hex() + "::" + p.QuoteTokenAddress.Hex()
	return code
}

func (p *Pair) Name() string {
	name := p.BaseTokenSymbol + "/" + p.QuoteTokenSymbol
	return name
}

func (p *Pair) MinQuoteAmount() *big.Int {
	return math.Add(math.Mul(big.NewInt(2), p.MakeFee), math.Mul(big.NewInt(2), p.TakeFee))
}

func (p *Pair) SetBSON(raw bson.Raw) error {
	decoded := &PairRecord{}

	err := raw.Unmarshal(decoded)
	if err != nil {
		return err
	}

	makeFee := big.NewInt(0)
	makeFee, _ = makeFee.SetString(decoded.MakeFee, 10)

	takeFee := big.NewInt(0)
	takeFee, _ = takeFee.SetString(decoded.TakeFee, 10)

	p.ID = decoded.ID
	p.BaseTokenSymbol = decoded.BaseTokenSymbol
	p.BaseTokenAddress = common.HexToAddress(decoded.BaseTokenAddress)
	p.BaseTokenDecimals = decoded.BaseTokenDecimals
	p.QuoteTokenSymbol = decoded.QuoteTokenSymbol
	p.QuoteTokenAddress = common.HexToAddress(decoded.QuoteTokenAddress)
	p.QuoteTokenDecimals = decoded.QuoteTokenDecimals
	p.Active = decoded.Active
	p.MakeFee = makeFee
	p.TakeFee = takeFee

	p.CreatedAt = decoded.CreatedAt
	p.UpdatedAt = decoded.UpdatedAt
	return nil
}

func (p *Pair) GetBSON() (interface{}, error) {
	return &PairRecord{
		ID: p.ID,

		BaseTokenSymbol:    p.BaseTokenSymbol,
		BaseTokenAddress:   p.BaseTokenAddress.Hex(),
		BaseTokenDecimals:  p.BaseTokenDecimals,
		QuoteTokenSymbol:   p.QuoteTokenSymbol,
		QuoteTokenAddress:  p.QuoteTokenAddress.Hex(),
		QuoteTokenDecimals: p.QuoteTokenDecimals,
		Active:             p.Active,
		MakeFee:            p.MakeFee.String(),
		TakeFee:            p.TakeFee.String(),
		CreatedAt:          p.CreatedAt,
		UpdatedAt:          p.UpdatedAt,
	}, nil
}

// Validate function is used to verify if an instance of
// struct satisfies all the conditions for a valid instance
func (p Pair) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.BaseTokenAddress, validation.Required),
		validation.Field(&p.QuoteTokenAddress, validation.Required),
		validation.Field(&p.BaseTokenSymbol, validation.Required),
		validation.Field(&p.QuoteTokenSymbol, validation.Required),
	)
}

// GetOrderBookKeys
func (p *Pair) GetOrderBookKeys() (sell, buy string) {
	return p.GetKVPrefix() + "::SELL", p.GetKVPrefix() + "::BUY"
}

func (p *Pair) GetKVPrefix() string {
	return p.BaseTokenAddress.Hex() + "::" + p.QuoteTokenAddress.Hex()
}

type PairData struct {
	Pair        PairID   `json:"id,omitempty" bson:"_id"`
	Close       *big.Int `json:"close,omitempty" bson:"close"`
	Count       *big.Int `json:"count,omitempty" bson:"count"`
	High        *big.Int `json:"high,omitempty" bson:"high"`
	Low         *big.Int `json:"low,omitempty" bson:"low"`
	Open        *big.Int `json:"open,omitempty" bson:"open"`
	Volume      *big.Int `json:"volume,omitempty" bson:"volume"`
	Timestamp   int64    `json:"timestamp,omitempty" bson:"timestamp"`
	OrderVolume *big.Int `json:"orderVolume,omitempty" bson:"orderVolume"`
	OrderCount  *big.Int `json:"orderCount,omitempty" bson:"orderCount"`
}

func (p *PairData) MarshalJSON() ([]byte, error) {
	pairData := map[string]interface{}{
		"pair": map[string]interface{}{
			"pairName":   p.Pair.PairName,
			"baseToken":  p.Pair.BaseToken.Hex(),
			"quoteToken": p.Pair.QuoteToken.Hex(),
		},
		"timestamp": p.Timestamp,
	}

	if p.Open != nil {
		pairData["open"] = p.Open.String()
	}

	if p.High != nil {
		pairData["high"] = p.High.String()
	}

	if p.Low != nil {
		pairData["low"] = p.Low.String()
	}

	if p.Volume != nil {
		pairData["volume"] = p.Volume.String()
	}

	if p.Close != nil {
		pairData["close"] = p.Close.String()
	}

	if p.Count != nil {
		pairData["count"] = p.Count.String()
	}

	if p.OrderVolume != nil {
		pairData["orderVolume"] = p.OrderVolume.String()
	}

	if p.OrderCount != nil {
		pairData["orderCount"] = p.OrderCount.String()
	}

	bytes, err := json.Marshal(pairData)
	return bytes, err
}

func (p *PairData) AddressCode() string {
	code := p.Pair.BaseToken.Hex() + "::" + p.Pair.QuoteToken.Hex()
	return code
}
