package cryptor

import "testing"

type ksKDFParams struct {
	// Shared parameters
	Salt  string `json:"salt"`
	DKLen int    `json:"dklen"`
	// Scrypt-specific parameters
	N int `json:"n,omitempty"`
	P int `json:"p,omitempty"`
	R int `json:"r,omitempty"`
	// PBKDF2-specific parameters
	C   int    `json:"c,omitempty"`
	PRF string `json:"prf,omitempty"`
}
type ksKDF struct {
	Function string       `json:"function"`
	Params   *ksKDFParams `json:"params"`
	Message  string       `json:"message"`
}
type ksChecksum struct {
	Function string                 `json:"function"`
	Params   map[string]interface{} `json:"params"`
	Message  string                 `json:"message"`
}
type ksCipherParams struct {
	// AES-128-CTR-specific parameters
	IV string `json:"iv,omitempty"`
}
type ksCipher struct {
	Function string          `json:"function"`
	Params   *ksCipherParams `json:"params"`
	Message  string          `json:"message"`
}
type keystoreV4 struct {
	KDF      *ksKDF      `json:"kdf"`
	Checksum *ksChecksum `json:"checksum"`
	Cipher   *ksCipher   `json:"cipher"`
}

const (
	name    = "keystore"
	version = 4
)

type options struct {
	cipher    string
	costPower uint
}

type optionFunc func(*options)

func (of optionFunc) apply(opt *options) {
	of(opt)
}

type Option interface {
	apply(*options)
}

func WithCipher(cipher string) Option {
	return optionFunc(func(opt *options) {
		opt.cipher = cipher
	})
}

// WithCost sets the cipher key cost for the encryptor to 2^power overriding
// the default value of 18 (ie. 2^18=262144). Higher values increases the
// cost of an exhaustive search but makes encoding and decoding proportionally slower.
// This should only be in testing as it affects security. It panics if t is nil.
func WithCost(t *testing.T, costPower uint) Option {
	if t == nil {
		panic("nil testing parameter")
	}
	return optionFunc(func(o *options) {
		o.costPower = costPower
	})
}

type Cryptor struct {
	cipher string
	cost   int
}

// NewEncryptor creates a new keystore V4 encryptor.
// This takes the following options:
// - cipher: the cipher to use when encrypting the secret, can be either "pbkdf2" (default) or "scrypt"
// - costPower: the cipher key cost to use as power of 2, default is 18 (ie. 2^18).
func NewCryptor(opts ...Option) *Cryptor {
	defaultOpt := options{
		cipher:    "pbkdf2",
		costPower: 18,
	}
	for _, op := range opts {
		op.apply(&defaultOpt)
	}
	return &Cryptor{
		cipher: defaultOpt.cipher,
		cost:   1 << defaultOpt.costPower,
	}
}

// Name returns the name of this encryptor
func (c *Cryptor) Name() string {
	return name
}

// Version returns the version of this encryptor
func (c *Cryptor) Version() uint {
	return version
}
