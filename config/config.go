package config

// HashIDConfig contains properties of hashId
type HashIDConfig struct {
	Salt      string
	MinLength int
}

var (
	// HashID Config singleton
	HashID = HashIDConfig{"mSwyDdV6Ml4BNvmsM9TK", 11}

	// FBMessengerBotToken holds the verify token needed to verify facebook's
	// messenger bot
	FBMessengerBotToken = "zu4klu2QcPRw64ausbf4"
)
