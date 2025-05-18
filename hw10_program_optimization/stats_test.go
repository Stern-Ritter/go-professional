//go:build !bench
// +build !bench

package hw10programoptimization

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetDomainStat(t *testing.T) {
	data := `{"Id":1,"Name":"Howard Mendoza","Username":"0Oliver","Email":"aliquid_qui_ea@Browsedrive.gov","Phone":"6-866-899-36-79","Password":"InAQJvsq","Address":"Blackbird Place 25"}
{"Id":2,"Name":"Jesse Vasquez","Username":"qRichardson","Email":"mLynch@broWsecat.com","Phone":"9-373-949-64-00","Password":"SiZLeNSGn","Address":"Fulton Hill 80"}
{"Id":3,"Name":"Clarence Olson","Username":"RachelAdams","Email":"RoseSmith@Browsecat.com","Phone":"988-48-97","Password":"71kuz3gA5w","Address":"Monterey Park 39"}
{"Id":4,"Name":"Gregory Reid","Username":"tButler","Email":"5Moore@Teklist.net","Phone":"520-04-16","Password":"r639qLNu","Address":"Sunfield Park 20"}
{"Id":5,"Name":"Janice Rose","Username":"KeithHart","Email":"nulla@Linktype.com","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}`

	t.Run("find 'com'", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(data), "com")
		require.NoError(t, err)
		require.Equal(t, DomainStat{
			"browsecat.com": 2,
			"linktype.com":  1,
		}, result)
	})

	t.Run("find 'gov'", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(data), "gov")
		require.NoError(t, err)
		require.Equal(t, DomainStat{"browsedrive.gov": 1}, result)
	})

	t.Run("find 'unknown'", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(data), "unknown")
		require.NoError(t, err)
		require.Equal(t, DomainStat{}, result)
	})
}

func TestGetDomainStat_EmptyInput(t *testing.T) {
	r := bytes.NewReader([]byte(""))
	domainStat, err := GetDomainStat(r, "gilead")

	require.NoError(t, err, "unexpected error")
	assert.Equal(t, DomainStat{}, domainStat, "should return empty domain stat, byt got: %v", domainStat)
}

func TestGetDomainStat_ValidInput(t *testing.T) {
	data := `{"ID": 1, "Name": "Roland Deschain", "Username": "roland", "Email": "roland@dark.gilead"}
{"ID": 2, "Name": "Eddie Dean", "Username": "eddie", "Email": "eddie@dark.newyork"}
{"ID": 3, "Name": "Susannah Dean", "Username": "susannah", "Email": "susannah@dark.gilead"}
{"ID": 4, "Name": "Jake Chambers", "Username": "jake", "Email": "jake@dark.gilead"}
{"ID": 5, "Name": "Walter o'Dim", "Username": "walter", "Email": "walter@dark.dis"}`
	r := bytes.NewReader([]byte(data))
	domainStat, err := GetDomainStat(r, "gilead")

	require.NoError(t, err, "unexpected error")
	expected := DomainStat{
		"dark.gilead": 3,
	}
	assert.Equal(t, expected, domainStat, "should return domain stat: %v, but got: %v", expected, domainStat)
}

func TestGetDomainStat_InvalidInput(t *testing.T) {
	r := bytes.NewReader([]byte(`{"ID": 1, "Name": "Roland Deschain", "Username": "roland", "Email": "roland@dark.gilead"`))
	_, err := GetDomainStat(r, "gilead")

	require.Error(t, err, "should return error")
}
