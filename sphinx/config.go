package sphinx

import "github.com/xlab/pocketsphinx-go/pocketsphinx"

type Config struct {
	options   map[string]interface{}
	evaluated *pocketsphinx.CommandLn
}

// NewConfig creates a new command-line argument set based on the provided config options.
func NewConfig(opts ...Option) *Config {
	return &Config{
		options: make(map[string]interface{}, 32),
	}
}

type Option func(c *Config)

// NewConfigRetain gets a new config while retaining ownership of a command-line argument set.
func NewConfigRetain(ln *pocketsphinx.CommandLn) *Config {
	return &Config{
		evaluated: pocketsphinx.CommandLnRetain(ln),
	}
}

// Retain retains ownership of a command-line argument set.
func (c *Config) Retain() {
	c.evaluated = pocketsphinx.CommandLnRetain(c.evaluated)
}

func (c *Config) CommandLn() *pocketsphinx.CommandLn {
	if c.evaluated != nil {
		return c.evaluated
	}
	defn := pocketsphinx.Args()
	argv := make([][]byte, 0, 32)
	// for name, opt := range c.options {
	//	 TODO: real config opts
	// }
	c.evaluated = pocketsphinx.CommandLnParseR(nil, defn, int32(len(argv)), argv, 0)
	return c.evaluated
}

func (c *Config) Destroy() bool {
	if c.evaluated != nil {
		ret := pocketsphinx.CommandLnFreeR(c.evaluated)
		c.evaluated = nil
		return ret == 0
	}
	return true
}
