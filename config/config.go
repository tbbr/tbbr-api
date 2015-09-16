package config

// HashIDConfig contains properties of hashId
type HashIDConfig struct {
	Salt      string
	MinLength int
}

var (
	// HashID Config singleton
	HashID = HashIDConfig{"mSwyDdV6Ml4BNvmsM9TK", 11}
)
