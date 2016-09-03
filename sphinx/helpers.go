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

func (s Strings) B() [][]byte {
	strs := s.S()
	results := make([][]byte, 0, len(strs))
	for _, s := range strs {
		results = append(results, []byte(s))
	}
	return results
}

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
var endChar byte = '\x00'

func (s String) S() string {
	if len(s) == 0 {
		return end
	}
	if s[len(s)-1] != endChar {
		return string(s) + end
	}
	return string(s)
}

func optToCommandLn(prevLn *pocketsphinx.CommandLn, opt map[String]interface{}) *pocketsphinx.CommandLn {
	if prevLn == nil {
		prevLn = pocketsphinx.NewCommandLn()
	}
	for name, v := range opt {
		switch x := v.(type) {
		case String:
			pocketsphinx.CommandLnSetStrR(prevLn, name.S(), x.S())
		case string:
			pocketsphinx.CommandLnSetStrR(prevLn, name.S(), String(x).S())
		case float32:
			pocketsphinx.CommandLnSetFloatR(prevLn, name.S(), float64(x))
		case float64:
			pocketsphinx.CommandLnSetFloatR(prevLn, name.S(), x)
		case bool:
			if x {
				pocketsphinx.CommandLnSetIntR(prevLn, name.S(), 1)
			} else {
				pocketsphinx.CommandLnSetIntR(prevLn, name.S(), 0)
			}
		case int8:
			pocketsphinx.CommandLnSetIntR(prevLn, name.S(), int(x))
		case int16:
			pocketsphinx.CommandLnSetIntR(prevLn, name.S(), int(x))
		case int32:
			pocketsphinx.CommandLnSetIntR(prevLn, name.S(), int(x))
		case int64:
			pocketsphinx.CommandLnSetIntR(prevLn, name.S(), int(x))
		case int:
			pocketsphinx.CommandLnSetIntR(prevLn, name.S(), int(x))
		case uint8:
			pocketsphinx.CommandLnSetIntR(prevLn, name.S(), int(x))
		case uint16:
			pocketsphinx.CommandLnSetIntR(prevLn, name.S(), int(x))
		case uint32:
			pocketsphinx.CommandLnSetIntR(prevLn, name.S(), int(x))
		case uint64:
			pocketsphinx.CommandLnSetIntR(prevLn, name.S(), int(x))
		case uint:
			pocketsphinx.CommandLnSetIntR(prevLn, name.S(), int(x))
		}
	}
	return prevLn
}
