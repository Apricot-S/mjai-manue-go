package game

import (
	"fmt"
	"log"
	"slices"
	"sort"
	"strings"
)

type Pai struct {
	id     uint8
	typ    rune
	number uint8
	isRed  bool
}

const (
	NumIDs = 9*3 + 7

	minPaiID   = 0
	maxSuitID  = 9*3 - 1
	minHonorID = maxSuitID + 1
	maxHonorID = minHonorID + 7 - 1
	minRedID   = maxHonorID + 1
	maxRedID   = minRedID + 3 - 1
	unknownID  = maxRedID + 1
	maxPaiID   = unknownID

	tsupaiType = 't'

	unknownStr    = "?"
	unknownType   = '?'
	unknownNumber = 10
)

var (
	Unknown = newPaiWithIDUnchecked(unknownID)

	types   = [...]rune{'m', 'p', 's', 't'}
	paiStrs = [...]string{
		"1m", "2m", "3m", "4m", "5m", "6m", "7m", "8m", "9m",
		"1p", "2p", "3p", "4p", "5p", "6p", "7p", "8p", "9p",
		"1s", "2s", "3s", "4s", "5s", "6s", "7s", "8s", "9s",
		"E", "S", "W", "N", "P", "F", "C",
		"5mr", "5pr", "5sr",
		"?",
	}
)

func (p *Pai) assertInitialized() {
	if p.number == 0 {
		log.Panic("Pai is not initialized")
	}
}

func (p *Pai) ID() uint8 {
	p.assertInitialized()
	return p.id
}

func (p *Pai) Type() rune {
	p.assertInitialized()
	return p.typ
}

func (p *Pai) Number() uint8 {
	p.assertInitialized()
	return p.number
}

func (p *Pai) IsRed() bool {
	p.assertInitialized()
	return p.isRed
}

func NewPaiWithID(id uint8) (*Pai, error) {
	if id > maxPaiID {
		return nil, fmt.Errorf("id out of range: %d", id)
	}
	return newPaiWithIDUnchecked(id), nil
}

func newPaiWithIDUnchecked(id uint8) *Pai {
	if id == unknownID {
		return &Pai{
			id:     unknownID,
			typ:    unknownType,
			number: unknownNumber,
			isRed:  false,
		}
	}

	isRed := id >= minRedID
	typ := toType(id, isRed)
	number := toNumber(id, isRed)

	return &Pai{
		id:     id,
		typ:    typ,
		number: number,
		isRed:  isRed,
	}
}

func NewPaiWithName(name string) (*Pai, error) {
	id := slices.Index(paiStrs[:], name)
	if id == -1 {
		return nil, fmt.Errorf("unknown pai string: %s", name)
	}
	return newPaiWithIDUnchecked(uint8(id)), nil
}

func NewPaiWithDetail(typ rune, number uint8, isRed bool) (*Pai, error) {
	if !slices.Contains(types[:], typ) {
		return nil, fmt.Errorf("bad type: %c", typ)
	}

	if number < 1 || 9 < number {
		return nil, fmt.Errorf("number out of range: %d", number)
	}
	if typ == tsupaiType && number > 7 {
		return nil, fmt.Errorf("number out of range for tsupai: %d", number)
	}
	if isRed && (number != 5 || typ == tsupaiType) {
		return nil, fmt.Errorf("no reds other than 5: %d", number)
	}

	return newPaiWithDetailUnchecked(typ, number, isRed), nil
}

func newPaiWithDetailUnchecked(typ rune, number uint8, isRed bool) *Pai {
	id := toId(typ, number, isRed)
	return &Pai{
		id:     id,
		typ:    typ,
		number: number,
		isRed:  isRed,
	}
}

func toType(id uint8, isRed bool) rune {
	var typeIndex uint8
	if isRed {
		typeIndex = 2 + id - maxRedID
	} else {
		typeIndex = id / 9
	}
	return types[typeIndex]
}

func toNumber(id uint8, isRed bool) uint8 {
	if isRed {
		return 5
	}
	return id%9 + 1
}

func toId(typ rune, number uint8, isRed bool) uint8 {
	typeIndex := uint8(slices.Index(types[:], typ))
	if isRed {
		return minRedID + typeIndex
	}
	return typeIndex*9 + number - 1
}

func (p *Pai) IsUnknown() bool {
	p.assertInitialized()
	return p.id == unknownID
}

func (p *Pai) HasSameSymbol(other *Pai) bool {
	p.assertInitialized()
	return p.typ == other.typ && p.number == other.number
}

func (p *Pai) NextForDora() *Pai {
	if p.IsUnknown() {
		return newPaiWithIDUnchecked(unknownID)
	}

	number := p.number
	var nextNumber uint8 = 0

	if p.typ == tsupaiType {
		// honors
		switch number {
		case 4:
			// N -> E
			nextNumber = 1
		case 7:
			// C -> P
			nextNumber = 5
		default:
			nextNumber = number + 1
		}
	} else {
		// suits
		switch number {
		case 9:
			nextNumber = 1
		default:
			nextNumber = number + 1
		}
	}

	return newPaiWithDetailUnchecked(p.typ, nextNumber, false)
}

func (p *Pai) IsTsupai() bool {
	p.assertInitialized()
	return p.typ == tsupaiType
}

func (p *Pai) IsYaochu() bool {
	p.assertInitialized()
	return p.typ == tsupaiType || p.number == 1 || p.number == 9
}

func (p *Pai) AddRed() *Pai {
	if p.IsUnknown() {
		return newPaiWithIDUnchecked(unknownID)
	}
	if p.IsTsupai() || p.Number() != 5 {
		return newPaiWithDetailUnchecked(p.typ, p.number, false)
	}
	return newPaiWithDetailUnchecked(p.typ, p.number, true)
}

func (p *Pai) RemoveRed() *Pai {
	if p.IsUnknown() {
		return newPaiWithIDUnchecked(unknownID)
	}
	return newPaiWithDetailUnchecked(p.typ, p.number, false)
}

func (p *Pai) Next(n int8) *Pai {
	if p.IsUnknown() || p.typ == tsupaiType {
		return nil
	}

	nextNumber := int8(p.number) + n
	if 1 <= nextNumber && nextNumber <= 9 {
		return newPaiWithDetailUnchecked(p.typ, uint8(nextNumber), false)
	}
	return nil
}

func (p *Pai) ToString() string {
	p.assertInitialized()
	return paiStrs[p.id]
}

func PaisToStr(pais []Pai) string {
	var strs []string
	for _, p := range pais {
		strs = append(strs, p.ToString())
	}
	return strings.Join(strs, " ")
}

func StrToPais(str string) ([]Pai, error) {
	pais := []Pai{}
	for _, f := range strings.Fields(str) {
		pai, err := NewPaiWithName(f)
		if err != nil {
			return nil, err
		}
		pais = append(pais, *pai)
	}
	return pais, nil
}

func (p *Pai) convertToCompare() int {
	id := int(p.id)
	switch {
	case id <= 4:
		// 1m - 5m
		return id
	case id == 34:
		// 5mr
		return 5
	case 5 <= id && id <= 13:
		// 6m - 5p
		return id + 1
	case id == 35:
		// 5pr
		return 15
	case 14 <= id && id <= 22:
		// 6p - 5s
		return id + 2
	case id == 36:
		// 5sr
		return 25
	case 23 <= id && id <= 33:
		// 6s - C
		return id + 3
	case id == unknownID:
		// ? (Unknown)
		return unknownID + 3
	default:
		log.Panicf("invalid pai id: %d", id)
		return -1
	}
}

func (p *Pai) compareTo(other *Pai) int {
	if other == nil {
		log.Panic("Other pai is nil")
	}
	p.assertInitialized()
	return p.convertToCompare() - other.convertToCompare()
}

type Pais []Pai

func (ps Pais) Len() int {
	return len(ps)
}

func (ps Pais) Less(i, j int) bool {
	return ps[i].compareTo(&ps[j]) < 0
}

func (ps Pais) Swap(i, j int) {
	ps[i], ps[j] = ps[j], ps[i]
}

func GetUniquePais(ps Pais, del func(Pai) bool) []Pai {
	unique := slices.Clone(ps)
	sort.Sort(unique)
	unique = slices.CompactFunc(unique, func(a, b Pai) bool {
		return a.ID() == b.ID()
	})
	if del == nil {
		return unique
	}
	unique = slices.DeleteFunc(unique, func(p Pai) bool {
		return del(p)
	})
	return unique
}
