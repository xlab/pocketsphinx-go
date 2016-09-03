package sphinx

import "github.com/xlab/pocketsphinx-go/pocketsphinx"

func b(b bool) int32 {
	if b {
		return 1
	}
	return 0
}

type String string

type Strings []string

func (s Strings) S() []string {
	for i := range s {
		if len(s[i]) == 0 {
			s[i] = end
			continue
		}
		if s[i][len(s[i])-1] != endChar {
			s[i] = s[i] + end
			continue
		}
	}
	return []string(s)
}

var end = "\x00"
var endChar = '\x00'

func (s String) S() string {
	if len(s) == 0 {
		return end
	}
	if s[len(s)-1] != endChar {
		return s + end
	}
	return s
}

func optToCommandLn(prevLn *pocketsphinx.CommandLn, opt map[String]interface{}) *pocketsphinx.CommandLn {
	if prevLn == nil {
		prevLn = pocketsphinx.NewCommandLn()
	}
	for name, v := range n.opt {
		switch x := v.(type) {
		case String:
			pocketsphinx.CommandLnSetStrR(ln, name.S(), x.S())
		case string:
			pocketsphinx.CommandLnSetStrR(ln, name.S(), String(x).S())
		case float32:
			pocketsphinx.CommandLnSetFloatR(ln, name.S(), float64(x))
		case float64:
			pocketsphinx.CommandLnSetFloatR(ln, name.S(), x)
		case bool:
			if x {
				pocketsphinx.CommandLnSetIntR(ln, name.S(), 1)
			} else {
				pocketsphinx.CommandLnSetIntR(ln, name.S(), 0)
			}
		case int8:
			pocketsphinx.CommandLnSetIntR(ln, name.S(), int(x))
		case int16:
			pocketsphinx.CommandLnSetIntR(ln, name.S(), int(x))
		case int32:
			pocketsphinx.CommandLnSetIntR(ln, name.S(), int(x))
		case int64:
			pocketsphinx.CommandLnSetIntR(ln, name.S(), int(x))
		case int:
			pocketsphinx.CommandLnSetIntR(ln, name.S(), int(x))
		case uint8:
			pocketsphinx.CommandLnSetIntR(ln, name.S(), int(x))
		case uint16:
			pocketsphinx.CommandLnSetIntR(ln, name.S(), int(x))
		case uint32:
			pocketsphinx.CommandLnSetIntR(ln, name.S(), int(x))
		case uint64:
			pocketsphinx.CommandLnSetIntR(ln, name.S(), int(x))
		case uint:
			pocketsphinx.CommandLnSetIntR(ln, name.S(), int(x))
		}
	}
	return ln
}
