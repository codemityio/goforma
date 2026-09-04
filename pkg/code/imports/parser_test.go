package imports

import (
	_ "embed"
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	//go:embed testdata/pkg.json
	pkg string

	//go:embed testdata/pkg-app.json
	pkgApp string

	//go:embed testdata/pkg-integration.json
	pkgIntegration string

	//go:embed testdata/pkg-recursive.json
	pkgRecursive string

	//go:embed testdata/pkg-recursive-depth-1.json
	pkgRecursiveDepth1 string

	//go:embed testdata/pkg-recursive-exclude-standard.json
	pkgRecursiveExcludeStandard string

	//go:embed testdata/pkg-recursive-exclude-vendor.json
	pkgRecursiveExcludeVendor string

	//go:embed testdata/pkg-recursive-exclude-internal.json
	pkgRecursiveExcludeInternal string
)

func TestDefaultParser_Parse(t *testing.T) {
	wd, err := os.Getwd()
	require.NoError(t, err)

	toMakeTheTestPassOnDifferentOS := []string{
		"crypto/fips140",
		"crypto/internal/alias",
		"crypto/internal/bigmod",
		"crypto/internal/boring",
		"crypto/internal/boring/bbig",
		"crypto/internal/boring/sig",
		"crypto/internal/edwards25519",
		"crypto/internal/edwards25519/field",
		"crypto/internal/entropy",
		"crypto/internal/fips140",
		"crypto/internal/fips140/aes",
		"crypto/internal/fips140/aes/gcm",
		"crypto/internal/fips140/alias",
		"crypto/internal/fips140/bigmod",
		"crypto/internal/fips140/check",
		"crypto/internal/fips140/drbg",
		"crypto/internal/fips140/ecdh",
		"crypto/internal/fips140/ecdsa",
		"crypto/internal/fips140/ed25519",
		"crypto/internal/fips140/edwards25519",
		"crypto/internal/fips140/edwards25519/field",
		"crypto/internal/fips140/hkdf",
		"crypto/internal/fips140/hmac",
		"crypto/internal/fips140/mlkem",
		"crypto/internal/fips140/nistec",
		"crypto/internal/fips140/nistec/fiat",
		"crypto/internal/fips140/rsa",
		"crypto/internal/fips140/sha256",
		"crypto/internal/fips140/sha3",
		"crypto/internal/fips140/sha512",
		"crypto/internal/fips140/subtle",
		"crypto/internal/fips140/tls12",
		"crypto/internal/fips140/tls13",
		"crypto/internal/fips140deps/byteorder",
		"crypto/internal/fips140deps/cpu",
		"crypto/internal/fips140deps/godebug",
		"crypto/internal/fips140hash",
		"crypto/internal/fips140only",
		"crypto/internal/hpke",
		"crypto/internal/impl",
		"crypto/internal/mlkem768",
		"crypto/internal/nistec",
		"crypto/internal/nistec/fiat",
		"crypto/internal/randutil",
		"crypto/internal/sysrand",
		"crypto/sha3",
		"crypto/subtle",
		"crypto/tls/internal/fips140tls",
		"crypto/x509/internal/macos",
		"embed",
		"golang.org/x/crypto/hkdf",
		"golang.org/x/crypto/internal/poly1305",
		"golang.org/x/crypto/sha3",
		"golang.org/x/net/route",
		"golang.org/x/sys/cpu",
		"internal/abi",
		"internal/byteorder",
		"internal/chacha8rand",
		"internal/concurrent",
		"internal/cpu",
		"internal/goarch",
		"internal/goos",
		"internal/routebsd",
		"internal/runtime/maps",
		"internal/runtime/math",
		"internal/runtime/sys",
		"internal/runtime/syscall",
		"internal/strconv",
		"internal/stringslite",
		"internal/sync",
		"internal/weak",
		"io/fs",
		"math/bits",
		"math/rand/v2",
		"os",
		"path/filepath",
		"reflect",
		"runtime",
		"slices",
		"sync",
		"sync/atomic",
		"unique",
		"unsafe",
	}

	// Define test cases
	tests := []struct {
		name                                            string
		inputPath                                       string
		depth                                           int
		owned                                           []string
		excludePaths                                    []string
		expectedResult                                  string
		excludeStandard, excludeVendor, excludeInternal bool
		writePath                                       string
	}{
		{
			name:           "pkg",
			inputPath:      "./pkg",
			depth:          DefaultDepth,
			owned:          []string{"github.com/codemityio"},
			excludePaths:   []string{},
			expectedResult: pkg,
			writePath:      "testdata/pkg.json",
		},
		{
			name:           "pkg-app",
			inputPath:      "./pkg/app",
			depth:          DefaultDepth,
			owned:          []string{"github.com/codemityio"},
			excludePaths:   toMakeTheTestPassOnDifferentOS,
			expectedResult: pkgApp,
			writePath:      "testdata/pkg-app.json",
		},
		{
			name:           "pkg-integration",
			inputPath:      "./pkg/integration",
			depth:          DefaultDepth,
			owned:          []string{"github.com/codemityio"},
			excludePaths:   toMakeTheTestPassOnDifferentOS,
			expectedResult: pkgIntegration,
			writePath:      "testdata/pkg-integration.json",
		},
		{
			name:           "pkg-recursive",
			inputPath:      "./pkg/...",
			depth:          DefaultDepth,
			owned:          []string{"github.com/codemityio"},
			excludePaths:   toMakeTheTestPassOnDifferentOS,
			expectedResult: pkgRecursive,
			writePath:      "testdata/pkg-recursive.json",
		},
		{
			name:           "pkg-recursive-depth-1",
			inputPath:      "./pkg/...",
			depth:          1,
			owned:          []string{"github.com/codemityio"},
			excludePaths:   toMakeTheTestPassOnDifferentOS,
			expectedResult: pkgRecursiveDepth1,
			writePath:      "testdata/pkg-recursive-depth-1.json",
		},
		{
			name:            "pkg-recursive-exclude-standard",
			inputPath:       "./pkg/...",
			depth:           DefaultDepth,
			owned:           []string{"github.com/codemityio"},
			excludePaths:    toMakeTheTestPassOnDifferentOS,
			expectedResult:  pkgRecursiveExcludeStandard,
			excludeStandard: true,
			writePath:       "testdata/pkg-recursive-exclude-standard.json",
		},
		{
			name:           "pkg-recursive-exclude-vendor",
			inputPath:      "./pkg/...",
			depth:          DefaultDepth,
			owned:          []string{"github.com/codemityio"},
			excludePaths:   toMakeTheTestPassOnDifferentOS,
			expectedResult: pkgRecursiveExcludeVendor,
			excludeVendor:  true,
			writePath:      "testdata/pkg-recursive-exclude-vendor.json",
		},
		{
			name:            "pkg-recursive-exclude-internal",
			inputPath:       "./pkg/...",
			depth:           DefaultDepth,
			owned:           []string{"github.com/codemityio"},
			excludePaths:    toMakeTheTestPassOnDifferentOS,
			expectedResult:  pkgRecursiveExcludeInternal,
			excludeInternal: true,
			writePath:       "testdata/pkg-recursive-exclude-internal.json",
		},
	}

	// Iterate over test cases
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			parser := New(
				WithRootPath(wd+"/testdata/code"),
				WithDepth(test.depth),
				WithOwned(test.owned),
				WithExcludePaths(test.excludePaths),
				WithExcludeStandard(test.excludeStandard),
				WithExcludeVendor(test.excludeVendor),
				WithExcludeInternal(test.excludeInternal),
			)

			output, err := parser.Parse(test.inputPath)
			require.NoError(t, err)

			result, err := json.MarshalIndent(output, "", "  ")
			require.NoError(t, err)

			assert.JSONEq(t, test.expectedResult, string(result))

			require.NoError(t, os.WriteFile(test.writePath, result, 0o644)) // #nosec G306
		})
	}
}
