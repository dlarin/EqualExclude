package EqualExclude

import (
	"fmt"
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/stretchr/testify/require"
)

type obj struct {
	SubObj  subObj
	Id      int
	Name    string
	Column1 string
	Column2 int
	Column3 *string
	Column4 *int
	column1 string
	column2 int
	column3 *string
	column4 *int
}

type subObj struct {
	Column1 string
	Column2 string
	Column3 []string
}

func TestEqualExcludeObject(t *testing.T) {
	column31 := "test3"
	column32 := "test4"
	column41 := 100
	column42 := 200

	var q1, q2 obj
	require.NoError(t, faker.FakeData(&q1))
	require.NoError(t, faker.FakeData(&q2))
	q1.Name, q2.Name = "test name", "test name"
	q1.column1 = "test1"
	q2.column1 = "test2"
	q1.column2 = 100
	q2.column2 = 200
	q1.column3 = &column31
	q2.column3 = &column32
	q1.column4 = &column41
	q2.column4 = &column42
	q1.SubObj.Column2 = "test sub 1"
	q2.SubObj.Column2 = "test sub 1"
	q1.SubObj.Column3 = []string{"1", "2", "-1"}
	q2.SubObj.Column3 = []string{"2", "2", "3"}

	EqualExclude(t, &q1, &q2, "Id", "Column1", "Column2", "Column3",
		"Column4",
		"column1",
		"column2",
		"column3",
		"column4",
		"SubObj.Column1",
		"SubObj.Column2",
		"SubObj.Column3.0",
		"SubObj.Column3.2",
	)

}

func TestEqualExcludeSplit(t *testing.T) {
	tests := []struct {
		q1 interface{}
		q2 interface{}
		s  []string
	}{
		{
			[]string{"1", "2", "3"},
			[]string{"1", "5", "3"},
			[]string{"1"},
		},
		{
			[]int{1, 2, 6},
			[]int{1, 2, 3},
			[]string{"2"},
		},
		{
			[]interface{}{1, 2, 6, "123", "47556"},
			[]interface{}{1, 2, 3, "123", "456"},
			[]string{"2", "4"},
		},
		{
			[]interface{}{1, 2, []int{1, 3, 3}, "123", "456"},
			[]interface{}{1, 2, []int{1, 2, 3}, "123", "456123"},
			[]string{"2.1", "4"},
		},
	}
	for i, tt := range tests {
		tt := tt
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			EqualExclude(t, tt.q1, tt.q2, tt.s...)
		})
	}
}

func TestEqualExcludeSplitMask(t *testing.T) {
	tests := []struct {
		q1 interface{}
		q2 interface{}
		s  []string
	}{
		{
			[]interface{}{1, 2, []int{1, 3, 3}, "123", "456"},
			[]interface{}{1, 2, []int{1, 2, 3}, "123", "456123"},
			[]string{"*.1", "4"},
		},
		{
			[]obj{
				{Id: 1, Name: "test1"},
				{Id: 2, Name: "test2"},
			},
			[]obj{
				{Id: 1, Name: "test3"},
				{Id: 2, Name: "test4"},
			},
			[]string{"*.Name"},
		},
	}
	for i, tt := range tests {
		tt := tt
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			EqualExclude(t, tt.q1, tt.q2, tt.s...)
		})
	}
}

func TestEqualExcludeMap(t *testing.T) {
	tests := []struct {
		q1 interface{}
		q2 interface{}
		s  []string
	}{
		{
			map[string]int{"1": 1, "2": 7, "3": 3},
			map[string]int{"1": 1, "2": 5, "3": 3},
			[]string{"2"},
		},
		{
			map[string]string{"1": "1", "2": "21", "3": "3"},
			map[string]string{"1": "1", "2": "22", "3": "3"},
			[]string{"2"},
		},
		{
			map[string]obj{
				"1": {Id: 1, Name: "test1"},
			},
			map[string]obj{
				"1": {Id: 1, Name: "test3"},
			},
			[]string{"1.Name"},
		},
	}
	for i, tt := range tests {
		tt := tt
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			EqualExclude(t, tt.q1, tt.q2, tt.s...)
		})
	}
}
