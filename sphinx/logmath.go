package sphinx

import "github.com/xlab/pocketsphinx-go/pocketsphinx"

/*
 * Fast integer logarithmic addition operations.
 *
 * In evaluating HMM models, probability values are often kept in log
 * domain, to avoid overflow.  To enable these logprob values to be
 * held in int32 variables without significant loss of precision, a
 * logbase of (1+epsilon) (where epsilon < 0.01 or so) is used.  This
 * module maintains this logbase (B).
 *
 * However, maintaining probabilities in log domain creates a problem
 * when adding two probability values.  This problem can be solved by
 * table lookup.  Note that:
 *
 *  - \f$ b^z = b^x + b^y \f$
 *  - \f$ b^z = b^x(1 + b^{y-x})     = b^y(1 + e^{x-y}) \f$
 *  - \f$ z   = x + log_b(1 + b^{y-x}) = y + log_b(1 + b^{x-y}) \f$
 *
 * So:
 *
 *  - when \f$ y > x, z = y + logadd\_table[-(x-y)] \f$
 *  - when \f$ x > y, z = x + logadd\_table[-(y-x)] \f$
 *  - where \f$ logadd\_table[n] = log_b(1 + b^{-n}) \f$
 *
 * The first entry in <i>logadd_table</i> is
 * simply \f$ log_b(2.0) \f$, for
 * the case where \f$ y = x \f$ and thus
 * \f$ z = log_b(2x) = log_b(2) + x \f$.  The last entry is zero,
 * where \f$ log_b(x+y) = x = y \f$ due to loss of precision.
 *
 * Since this table can be quite large particularly for small
 * logbases, an option is provided to compress it by dropping the
 * least significant bits of the table.
 */

// LogMath integer log math computation class.
type LogMath struct {
	m *pocketsphinx.Logmath
}

// LogMath returns a retained copy of underlying reference to pocketsphinx.Logmath.
func (l *LogMath) LogMath() *pocketsphinx.Logmath {
	return pocketsphinx.LogmathRetain(l.m)
}

// WriteTo writes a log table to a file.
func (l LogMath) WriteTo(filename String) bool {
	ret := pocketsphinx.LogmathWrite(l.m, filename.S())
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
