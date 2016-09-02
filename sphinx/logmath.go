package sphinx

import "github.com/xlab/pocketsphinx-go/pocketsphinx"

type LogMath struct {
	m *pocketsphinx.Logmath
}

// DumpTo writes a log table to a file.
func (l LogMath) DumpTo(filename string) bool {
	ret := pocketsphinx.LogmathWrite(l.m, filename+"\x00")
	return ret == 0
}

// GetTableShape gets the log table size and dimensions.
func (l LogMath) GetTableShape() (size, width, shift uint32, ok bool) {
	ret := pocketsphinx.LogmathGetTableShape(l.m, &size, &width, &shift)
	ok = ret == 0
	return
}

// GetBase gets the log base.
func (l LogMath) GetBase() float64 {
	return pocketsphinx.LogmathGetBase(l.m)
}

// GetZero gets the smallest possible value represented in this base.
func (l LogMath) GetZero() int32 {
	return pocketsphinx.LogmathGetZero(l.m)
}

// GetWidth gets the width of the values in a log table.
func (l LogMath) GetWidth() int32 {
	return pocketsphinx.LogmathGetWidth(l.m)
}

// GetShift gets the shift of the values in a log table.
func (l LogMath) GetShift() int32 {
	return pocketsphinx.LogmathGetShift(l.m)
}

// AddExact adds two values in log space exactly and slowly (without using add table).
func (l LogMath) AddExact(p, q int32) int32 {
	return pocketsphinx.LogmathAddExact(l.m, p, q)
}

// Add two values in log space (i.e. return log(exp(p)+exp(q)))
func (l LogMath) Add(p, q int32) int32 {
	return pocketsphinx.LogmathAdd(l.m, p, q)
}

// Log converts linear floating point number to integer log in base B.
func (l LogMath) Log(p float64) int32 {
	return pocketsphinx.LogmathLog(l.m, p)
}

// Exp converts integer log in base B to linear floating point.
func (l LogMath) Exp(p int32) float64 {
	return pocketsphinx.LogmathExp(l.m, p)
}

// LnToLog converts natural log (in floating point) to integer log in base B.
func (l LogMath) LnToLog(p float64) int32 {
	return pocketsphinx.LogmathLnToLog(l.m, p)
}

// LogToLn converts integer log in base B to natural log (in floating point).
func (l LogMath) LogToLn(p int32) float64 {
	return pocketsphinx.LogmathLogToLn(l.m, p)
}

// Log10ToLog converts base 10 log (in floating point) to integer log in base B.
func (l LogMath) Log10ToLog(p float64) int32 {
	return pocketsphinx.LogmathLog10ToLog(l.m, p)
}

// LogToLog10 converts integer log in base B to base 10 log (in floating point).
func (l LogMath) LogToLog10(p int32) float64 {
	return pocketsphinx.LogmathLogToLog10(l.m, p)
}

// Log10ToLogFloat converts base 10 log (in floating point) to float log in base B.
func (l LogMath) Log10ToLogFloat(p float64) float32 {
	return pocketsphinx.LogmathLog10ToLogFloat(l.m, p)
}

// LogFloatToLog10 converts float log in base B to base 10 log.
func (l LogMath) LogFloatToLog10(p float32) float64 {
	return pocketsphinx.LogmathLogFloatToLog10(l.m, p)
}

func (l *LogMath) Destroy() bool {
	if l.m != nil {
		ret := pocketsphinx.LogmathFree(l.m)
		l.m = nil
		return ret == 0
	}
	return true
}

func (l *LogMath) Retain() {
	l.m = pocketsphinx.LogmathRetain(l.m)
}
