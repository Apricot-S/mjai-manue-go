package game

import (
	"fmt"
	"reflect"
	"testing"
)

func TestNewPaiSetWithPais(t *testing.T) {
	type args struct {
		pais []Pai
	}
	type testCase struct {
		name    string
		args    args
		want    *PaiSet
		wantErr bool
	}
	tests := []testCase{}

	pais0 := []Pai{}
	array0 := [NumIDs]int{}
	ps0 := PaiSet(array0)
	tests = append(tests, testCase{"empty", args{pais0}, &ps0, false})

	pais1 := []Pai{}
	for i := uint8(0); i < NumIDs; i++ {
		p, _ := NewPaiWithID(i)
		pais1 = append(pais1, *p)
	}
	array1 := [NumIDs]int{}
	for i := range array1 {
		array1[i] = 1
	}
	ps1 := PaiSet(array1)
	tests = append(tests, testCase{"all1", args{pais1}, &ps1, false})

	pais2 := []Pai{}
	for i := uint8(0); i < NumIDs; i++ {
		p, _ := NewPaiWithID(i)
		pais2 = append(pais2, *p, *p)
	}
	array2 := [NumIDs]int{}
	for i := range array2 {
		array2[i] = 2
	}
	ps2 := PaiSet(array2)
	tests = append(tests, testCase{"all2", args{pais2}, &ps2, false})

	redPais := []Pai{}
	for _, n := range []string{"5mr", "5pr", "5sr"} {
		r, _ := NewPaiWithName(n)
		redPais = append(redPais, *r)
	}
	redArray := [NumIDs]int{4: 1, 4 + 9: 1, 4 + 18: 1}
	redPs := PaiSet(redArray)
	tests = append(tests, testCase{"red", args{redPais}, &redPs, false})

	unknowns := []Pai{}
	u, _ := NewPaiWithName("?")
	unknowns = append(unknowns, *u)
	tests = append(tests, testCase{"cannot add unknown", args{unknowns}, nil, true})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewPaiSetWithPais(tt.args.pais)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPaiSetWithPais() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPaiSetWithPais() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetAll(t *testing.T) {
	type args struct {
		array [NumIDs]int
	}
	type testCase struct {
		name string
		args args
		want *PaiSet
	}
	tests := []testCase{}

	array4 := [NumIDs]int{}
	for i := range array4 {
		array4[i] = 4
	}
	ps := PaiSet(array4)
	tests = append(tests, testCase{"all4", args{array4}, &ps})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetAll(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPaiSet_ToPais(t *testing.T) {
	type fields struct {
		array [NumIDs]int
	}
	type testCase struct {
		name   string
		fields fields
		want   []Pai
	}
	tests := []testCase{}

	array0 := [NumIDs]int{}
	pais0 := []Pai{}
	tests = append(tests, testCase{"empty", fields{array0}, pais0})

	array1 := [NumIDs]int{}
	for i := range array1 {
		array1[i] = 1
	}
	pais1 := []Pai{}
	for i := uint8(0); i < NumIDs; i++ {
		p, _ := NewPaiWithID(i)
		pais1 = append(pais1, *p)
	}
	tests = append(tests, testCase{"all1", fields{array1}, pais1})

	array2 := [NumIDs]int{}
	for i := range array2 {
		array2[i] = 2
	}
	pais2 := []Pai{}
	for i := uint8(0); i < NumIDs; i++ {
		p, _ := NewPaiWithID(i)
		pais2 = append(pais2, *p, *p)
	}
	tests = append(tests, testCase{"all2", fields{array2}, pais2})

	arrayMinus1 := [NumIDs]int{}
	for i := range arrayMinus1 {
		arrayMinus1[i] = -1
	}
	tests = append(tests, testCase{"all-1", fields{arrayMinus1}, pais0})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ps := PaiSet(tt.fields.array)
			if got := ps.ToPais(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PaiSet.ToPais() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPaiSet_Count(t *testing.T) {
	type fields struct {
		array [NumIDs]int
	}
	type args struct {
		pai *Pai
	}
	type testCase struct {
		name    string
		fields  fields
		args    args
		want    int
		wantErr bool
	}
	tests := []testCase{}

	array0 := [NumIDs]int{}
	for i := uint8(0); i < NumIDs; i++ {
		p, _ := NewPaiWithID(i)
		tests = append(tests, testCase{p.ToString() + " 0", fields{array0}, args{p}, 0, false})
	}
	for _, n := range []string{"5mr", "5pr", "5sr"} {
		r, _ := NewPaiWithName(n)
		tests = append(tests, testCase{r.ToString() + " 0", fields{array0}, args{r}, 0, false})
	}

	array1 := [NumIDs]int{}
	for i := range array1 {
		array1[i] = 1
	}
	for i := uint8(0); i < NumIDs; i++ {
		p, _ := NewPaiWithID(i)
		tests = append(tests, testCase{p.ToString() + " 1", fields{array1}, args{p}, 1, false})
	}
	for _, n := range []string{"5mr", "5pr", "5sr"} {
		r, _ := NewPaiWithName(n)
		tests = append(tests, testCase{r.ToString() + " 1", fields{array1}, args{r}, 1, false})
	}

	arrayMinus1 := [NumIDs]int{}
	for i := range arrayMinus1 {
		arrayMinus1[i] = -1
	}
	for i := uint8(0); i < NumIDs; i++ {
		p, _ := NewPaiWithID(i)
		tests = append(tests, testCase{p.ToString() + " -1", fields{arrayMinus1}, args{p}, -1, false})
	}
	for _, n := range []string{"5mr", "5pr", "5sr"} {
		r, _ := NewPaiWithName(n)
		tests = append(tests, testCase{r.ToString() + " -1", fields{arrayMinus1}, args{r}, -1, false})
	}

	u, _ := NewPaiWithName("?")
	tests = append(tests, testCase{"?", fields{array1}, args{u}, 0, true})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ps := PaiSet(tt.fields.array)
			got, err := ps.Count(tt.args.pai)
			if (err != nil) != tt.wantErr {
				t.Errorf("PaiSet.Count() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PaiSet.Count() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPaiSet_Has(t *testing.T) {
	type fields struct {
		array [NumIDs]int
	}
	type args struct {
		pai *Pai
	}
	type testCase struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}
	tests := []testCase{}

	array0 := [NumIDs]int{}
	for i := uint8(0); i < NumIDs; i++ {
		p, _ := NewPaiWithID(i)
		tests = append(tests, testCase{p.ToString() + " 0", fields{array0}, args{p}, false, false})
	}
	for _, n := range []string{"5mr", "5pr", "5sr"} {
		r, _ := NewPaiWithName(n)
		tests = append(tests, testCase{r.ToString() + " 0", fields{array0}, args{r}, false, false})
	}

	array1 := [NumIDs]int{}
	for i := range array1 {
		array1[i] = 1
	}
	for i := uint8(0); i < NumIDs; i++ {
		p, _ := NewPaiWithID(i)
		tests = append(tests, testCase{p.ToString() + " 1", fields{array1}, args{p}, true, false})
	}
	for _, n := range []string{"5mr", "5pr", "5sr"} {
		r, _ := NewPaiWithName(n)
		tests = append(tests, testCase{r.ToString() + " 1", fields{array1}, args{r}, true, false})
	}

	arrayMinus1 := [NumIDs]int{}
	for i := range arrayMinus1 {
		arrayMinus1[i] = -1
	}
	for i := uint8(0); i < NumIDs; i++ {
		p, _ := NewPaiWithID(i)
		tests = append(tests, testCase{p.ToString() + " -1", fields{arrayMinus1}, args{p}, false, false})
	}
	for _, n := range []string{"5mr", "5pr", "5sr"} {
		r, _ := NewPaiWithName(n)
		tests = append(tests, testCase{r.ToString() + " -1", fields{arrayMinus1}, args{r}, false, false})
	}

	u, _ := NewPaiWithName("?")
	tests = append(tests, testCase{"?", fields{array1}, args{u}, false, true})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ps := PaiSet(tt.fields.array)
			got, err := ps.Has(tt.args.pai)
			if (err != nil) != tt.wantErr {
				t.Errorf("PaiSet.Has() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PaiSet.Has() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPaiSet_AddPai(t *testing.T) {
	type fields struct {
		array [NumIDs]int
	}
	type args struct {
		pai *Pai
		n   int
	}
	type testCase struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		count   int
	}
	tests := []testCase{}

	for i := uint8(0); i < NumIDs; i++ {
		array := [NumIDs]int{}
		p, _ := NewPaiWithID(i)
		tests = append(tests, testCase{p.ToString() + "+0", fields{array}, args{p, 0}, false, 0})
	}
	for _, n := range []string{"5mr", "5pr", "5sr"} {
		array := [NumIDs]int{}
		r, _ := NewPaiWithName(n)
		tests = append(tests, testCase{r.ToString() + "+0", fields{array}, args{r, 0}, false, 0})
	}

	for i := uint8(0); i < NumIDs; i++ {
		array := [NumIDs]int{}
		p, _ := NewPaiWithID(i)
		tests = append(tests, testCase{p.ToString() + "+1", fields{array}, args{p, 1}, false, 1})
	}
	for _, n := range []string{"5mr", "5pr", "5sr"} {
		array := [NumIDs]int{}
		r, _ := NewPaiWithName(n)
		tests = append(tests, testCase{r.ToString() + "+1", fields{array}, args{r, 1}, false, 1})
	}

	for i := uint8(0); i < NumIDs; i++ {
		array := [NumIDs]int{}
		p, _ := NewPaiWithID(i)
		tests = append(tests, testCase{p.ToString() + "-1", fields{array}, args{p, -1}, false, -1})
	}
	for _, n := range []string{"5mr", "5pr", "5sr"} {
		array := [NumIDs]int{}
		r, _ := NewPaiWithName(n)
		tests = append(tests, testCase{r.ToString() + "-1", fields{array}, args{r, -1}, false, -1})
	}

	// 1 - 1 = 0
	for i := uint8(0); i < NumIDs; i++ {
		array := [NumIDs]int{}
		array[i] = 1
		p, _ := NewPaiWithID(i)
		tests = append(tests, testCase{p.ToString() + "+1-1", fields{array}, args{p, -1}, false, 0})
	}
	for i, n := range []string{"5mr", "5pr", "5sr"} {
		array := [NumIDs]int{}
		array[4+i*9] = 1
		r, _ := NewPaiWithName(n)
		tests = append(tests, testCase{r.ToString() + "+1-1", fields{array}, args{r, -1}, false, 0})
	}

	{
		array := [NumIDs]int{}
		u, _ := NewPaiWithName("?")
		tests = append(tests, testCase{u.ToString(), fields{array}, args{u, 0}, true, 0})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ps := PaiSet(tt.fields.array)
			if err := ps.AddPai(tt.args.pai, tt.args.n); (err != nil) != tt.wantErr {
				t.Errorf("PaiSet.AddPai() error = %v, wantErr %v", err, tt.wantErr)
			}
			if c, _ := ps.Count(tt.args.pai); c != tt.count {
				t.Errorf("PaiSet.Count() = %v, count %v", c, tt.count)
			}
		})
	}
}

func TestPaiSet_AddPais(t *testing.T) {
	type fields struct {
		array [NumIDs]int
	}
	type args struct {
		pais []Pai
	}
	type testCase struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		pai     *Pai
		count   int
	}
	tests := []testCase{}

	for i := uint8(0); i < NumIDs; i++ {
		array := [NumIDs]int{}
		p, _ := NewPaiWithID(i)
		ps := []Pai{}
		tests = append(tests, testCase{p.ToString() + "+0", fields{array}, args{ps}, false, p, 0})
	}
	for _, n := range []string{"5mr", "5pr", "5sr"} {
		array := [NumIDs]int{}
		r, _ := NewPaiWithName(n)
		ps := []Pai{}
		tests = append(tests, testCase{r.ToString() + "+0", fields{array}, args{ps}, false, r, 0})
	}

	for i := uint8(0); i < NumIDs; i++ {
		array := [NumIDs]int{}
		p, _ := NewPaiWithID(i)
		ps := []Pai{*p}
		tests = append(tests, testCase{p.ToString() + "+1", fields{array}, args{ps}, false, p, 1})
	}
	for _, n := range []string{"5mr", "5pr", "5sr"} {
		array := [NumIDs]int{}
		r, _ := NewPaiWithName(n)
		ps := []Pai{*r}
		tests = append(tests, testCase{r.ToString() + "+1", fields{array}, args{ps}, false, r, 1})
	}

	for i := uint8(0); i < NumIDs; i++ {
		array := [NumIDs]int{}
		p, _ := NewPaiWithID(i)
		ps := []Pai{*p, *p}
		tests = append(tests, testCase{p.ToString() + "+2", fields{array}, args{ps}, false, p, 2})
	}
	for _, n := range []string{"5mr", "5pr", "5sr"} {
		array := [NumIDs]int{}
		r, _ := NewPaiWithName(n)
		ps := []Pai{*r, *r}
		tests = append(tests, testCase{r.ToString() + "+2", fields{array}, args{ps}, false, r, 2})
	}

	{
		array := [NumIDs]int{}
		u, _ := NewPaiWithName("?")
		ps := []Pai{*u}
		tests = append(tests, testCase{u.ToString(), fields{array}, args{ps}, true, u, 0})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ps := PaiSet(tt.fields.array)
			if err := ps.AddPais(tt.args.pais); (err != nil) != tt.wantErr {
				t.Errorf("PaiSet.AddPais() error = %v, wantErr %v", err, tt.wantErr)
			}
			if c, _ := ps.Count(tt.pai); c != tt.count {
				t.Errorf("PaiSet.Count() = %v, count %v", c, tt.count)
			}
		})
	}
}

func TestPaiSet_RemovePaiSet(t *testing.T) {
	type fields struct {
		array [NumIDs]int
	}
	type args struct {
		paiSet *PaiSet
	}
	type testCase struct {
		name   string
		fields fields
		args   args
		want   PaiSet
	}
	tests := []testCase{}

	// 1 - 1 = 0
	for i := uint8(0); i < NumIDs; i++ {
		array := [NumIDs]int{}
		array[i] = 1
		p, _ := NewPaiWithID(i)
		ps, _ := NewPaiSetWithPais([]Pai{*p})
		want := PaiSet{}
		tests = append(tests, testCase{p.ToString() + "+1-1", fields{array}, args{ps}, want})
	}
	for i, n := range []string{"5mr", "5pr", "5sr"} {
		array := [NumIDs]int{}
		array[4+i*9] = 1
		r, _ := NewPaiWithName(n)
		ps, _ := NewPaiSetWithPais([]Pai{*r})
		want := PaiSet{}
		tests = append(tests, testCase{r.ToString() + "+1-1", fields{array}, args{ps}, want})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ps := PaiSet(tt.fields.array)
			ps.RemovePaiSet(tt.args.paiSet)

			if !reflect.DeepEqual(ps, tt.want) {
				t.Errorf("removed = %v, want %v", ps, tt.want)
			}
		})
	}
}

func TestPaiSet_ToString(t *testing.T) {
	type fields struct {
		array [NumIDs]int
	}
	type testCase struct {
		name   string
		fields fields
		want   string
	}
	tests := []testCase{}

	tests = append(tests, testCase{"empty", fields{[NumIDs]int{}}, ""})

	for i := uint8(0); i < NumIDs; i++ {
		array := [NumIDs]int{}
		array[i] = 1
		p, _ := NewPaiWithID(i)
		want := p.ToString()
		tests = append(tests, testCase{want, fields{array}, want})
	}

	for i := uint8(0); i < NumIDs; i++ {
		array := [NumIDs]int{}
		array[i] = 2
		p, _ := NewPaiWithID(i)
		want := fmt.Sprintf("%s %s", p.ToString(), p.ToString())
		tests = append(tests, testCase{want, fields{array}, want})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ps := PaiSet(tt.fields.array)
			if got := ps.ToString(); got != tt.want {
				t.Errorf("PaiSet.ToString() = %v, want %v", got, tt.want)
			}
		})
	}
}
