package hw10programoptimization

import (
	"archive/zip"
	"testing"

	"github.com/stretchr/testify/require"
)

func BenchmarkGetDomainStat(b *testing.B) {
	b.Helper()
	r, err := zip.OpenReader("testdata/users.dat.zip")
	require.NoError(b, err)
	defer r.Close()

	for i := 0; i < b.N; i++ {
		data, err := r.File[0].Open()
		require.NoError(b, err)
		defer data.Close()

		_, err = GetDomainStat(data, "biz")
		require.NoError(b, err)
	}
}
