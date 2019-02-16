package config

// C is the package config
var C Config

// Config is the server configurtion
type Config struct {
	LionelAddress   string
	StackersAddress string

	Keystore struct {
		Account string
		Path    string
		Passwd  string
	}

	DB struct {
		Path string
	}

	Web3 struct {
		MaxGasPrice uint64
		RPCURL      string
	}

	API struct {
		Port int
	}
}
