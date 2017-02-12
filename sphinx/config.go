package sphinx

import "github.com/xlab/pocketsphinx-go/pocketsphinx"

type Config struct {
	opt       map[String]interface{}
	evaluated *pocketsphinx.CommandLn
}

// NewConfig creates a new command-line argument set based on the provided config options.
func NewConfig(opts ...Option) *Config {
	cfg := &Config{
		opt: make(map[String]interface{}, 32),
	}
	for i := range opts {
		opts[i](cfg)
	}
	return cfg
}

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
	prevLn := pocketsphinx.CommandLnParseR(nil, defn, 0, nil, 0)
	c.evaluated = optToCommandLn(prevLn, c.opt)
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

type Option func(c *Config)

// Options for debugging and logging.

// LogFileOption sets file to write log messages in.
func LogFileOption(filename string) Option {
	return func(c *Config) {
		c.opt[String("-logfn")] = String(filename)
	}
}

// DebugOption sets verbosity level for debugging messages.
func DebugOption(level int) Option {
	return func(c *Config) {
		c.opt[String("-debug")] = level
	}
}

// MFCLogDirOption sets directory to log feature files to.
func MFCLogDirOption(dir string) Option {
	return func(c *Config) {
		c.opt[String("-mfclogdir")] = String(dir)
	}
}

// RawLogDirOption sets directory to log raw audio files to.
func RawLogDirOption(dir string) Option {
	return func(c *Config) {
		c.opt[String("-rawlogdir")] = String(dir)
	}
}

// SenLogDirOption sets directory to log senone score files to.
func SenLogDirOption(dir string) Option {
	return func(c *Config) {
		c.opt[String("-senlogdir")] = String(dir)
	}
}

// Options defining beam width parameters for tuning the search.

// BeamOption sets beam width applied to every frame in Viterbi search (smaller values mean wider beam).
//
// Default: 1e-48
func BeamOption(width float64) Option {
	return func(c *Config) {
		c.opt[String("-beam")] = width
	}
}

// WBeamOption sets beam width applied to word exits.
//
// Default: 7e-29
func WBeamOption(width float64) Option {
	return func(c *Config) {
		c.opt[String("-wbeam")] = width
	}
}

// PBeamOption sets beam width applied to phone transitions.
//
// Default: 1e-48
func PBeamOption(width float64) Option {
	return func(c *Config) {
		c.opt[String("-pbeam")] = width
	}
}

// LPBeamOption sets beam width applied to last phone in words.
//
// Default: 1e-40
func LPBeamOption(width float64) Option {
	return func(c *Config) {
		c.opt[String("-lpbeam")] = width
	}
}

// LPOnlyBeamOption sets beam width applied to last phone in single-phone words.
//
// Default: 7e-29
func LPOnlyBeamOption(width float64) Option {
	return func(c *Config) {
		c.opt[String("-lponlybeam")] = width
	}
}

// FwdFlatBeamOption sets beam width applied to every frame in second-pass flat search.
//
// Default: 1e-64
func FwdFlatBeamOption(width float64) Option {
	return func(c *Config) {
		c.opt[String("-fwdflatbeam")] = width
	}
}

// FwdFlatWBeamOption sets beam width applied to word exits in second-pass flat search.
//
// Default: 7e-29
func FwdFlatWBeamOption(width float64) Option {
	return func(c *Config) {
		c.opt[String("-fwdflatwbeam")] = width
	}
}

// PLWindowOption sets phoneme lookahead window size, in frames.
//
// Default: 5
func PLWindowOption(frames int) Option {
	return func(c *Config) {
		c.opt[String("-pl_window")] = frames
	}
}

// PLBeamOption sets beam width applied to phone loop search for lookahead.
//
// Default: 1e-10
func PLBeamOption(width float64) Option {
	return func(c *Config) {
		c.opt[String("-pl_beam")] = width
	}
}

// PLPBeamOption sets beam width applied to phone loop transitions for lookahead.
//
// Default: 1e-10
func PLPBeamOption(width float64) Option {
	return func(c *Config) {
		c.opt[String("-pl_pbeam")] = width
	}
}

// PipOption sets phone insertion penalty for phone loop.
//
// Default: 1.0
func PipOption(penalty float32) Option {
	return func(c *Config) {
		c.opt[String("-pip")] = penalty
	}
}

// PLWeightOption sets weight for phoneme lookahead penalties.
//
// Default: 3.0
func PLWeightOption(weight float64) Option {
	return func(c *Config) {
		c.opt[String("-pl_weight")] = weight
	}
}

// Options defining other parameters for tuning the search.

// CompAllSenOption enables compute all senone scores in every frame (can be faster when there are many senones).
//
// Default: false
func CompAllSenOption(compallsen bool) Option {
	return func(c *Config) {
		c.opt[String("-compallsen")] = compallsen
	}
}

// FwdTreeOption enables run forward lexicon-tree search (1st pass).
//
// Default: true
func FwdTreeOption(fwdtree bool) Option {
	return func(c *Config) {
		c.opt[String("-fwdtree")] = fwdtree
	}
}

// FwdFlatOption enables run forward flat-lexicon search over word lattice (2nd pass).
//
// Default: true
func FwdFlatOption(fwdflat bool) Option {
	return func(c *Config) {
		c.opt[String("-fwdflat")] = fwdflat
	}
}

// BestpathOption enables run bestpath (Dijkstra) search over word lattice (3rd pass).
//
// Default: true
func BestpathOption(bestpath bool) Option {
	return func(c *Config) {
		c.opt[String("-bestpath")] = bestpath
	}
}

// BacktraceOption enables printing results and backtraces to log.
//
// Default: false
func BacktraceOption(backtrace bool) Option {
	return func(c *Config) {
		c.opt[String("-backtrace")] = backtrace
	}
}

// LatsizeOption sets initial backpointer table size.
//
// Default: 5000
func LatsizeOption(size int) Option {
	return func(c *Config) {
		c.opt[String("-latsize")] = size
	}
}

// MaxWPFOption sets maximum number of distinct word exits at each frame (or -1 for no pruning).
//
// Default: -1
func MaxWPFOption(max int) Option {
	return func(c *Config) {
		c.opt[String("-maxwpf")] = max
	}
}

// MaxHMMPFOption sets maximum number of active HMMs to maintain at each frame (or -1 for no pruning).
//
// Default: 30000
func MaxHMMPFOption(max int) Option {
	return func(c *Config) {
		c.opt[String("-maxhmmpf")] = max
	}
}

// MinEndFrOption sets nodes ignored in lattice construction if they persist for fewer than N frames.
//
// Default: 0
func MinEndFrOption(n int) Option {
	return func(c *Config) {
		c.opt[String("-min_endfr")] = n
	}
}

// FwdFlateFWidOption sets minimum number of end frames for a word to be searched in fwdflat search.
//
// Default: 4
func FwdFlateFWidOption(frames int) Option {
	return func(c *Config) {
		c.opt[String("-fwdflatefwid")] = frames
	}
}

// FwdFlatSfWinOption sets window of frames in lattice to search for successor words in fwdflat search.
//
// Default: 25
func FwdFlatSfWinOption(frames int) Option {
	return func(c *Config) {
		c.opt[String("-fwdflatsfwin")] = frames
	}
}

// Options for keyphrase spotting.

// KeyphraseOption sets keyphrase to spot.
func KeyphraseOption(keyphrase string) Option {
	return func(c *Config) {
		c.opt[String("-keyphrase")] = String(keyphrase)
	}
}

// KeywordsFileOption sets a file with keyphrases to spot, one per line.
func KeywordsFileOption(filename string) Option {
	return func(c *Config) {
		c.opt[String("-kws")] = String(filename)
	}
}

// KeywordsPLPOption sets phone loop probability for keyphrase spotting.
//
// Default: 1e-1
func KeywordsPLPOption(prob float64) Option {
	return func(c *Config) {
		c.opt[String("-kws_plp")] = prob
	}
}

// KeywordsDelayOption sets delay to wait for best detection score.
//
// Default: 10
func KeywordsDelayOption(delay int) Option {
	return func(c *Config) {
		c.opt[String("-kws_delay")] = delay
	}
}

// KeywordsThresholdOption threshold for p(hyp)/p(alternatives) ratio.
//
// Default: 1.0
func KeywordsThresholdOption(threshold float64) Option {
	return func(c *Config) {
		c.opt[String("-kws_threshold")] = threshold
	}
}

// FiniteStateGrammarsOption for finite state grammars.
func FiniteStateGrammarsOption(filepath string) Option {
	return func(c *Config) {
		c.opt[String("-fsg")] = String(filepath)
	}
}

// Options for statistical language models (N-Gram).

// AllPhoneFileOption sets filepath for phoneme decoding with phonetic lm.
func AllPhoneFileOption(filepath string) Option {
	return func(c *Config) {
		c.opt[String("-allphone")] = String(filepath)
	}
}

// AllPhoneCIOption enables perform phoneme decoding with phonetic lm and context-independent units only.
//
// Default: false
func AllPhoneCIOption(ciUnitsOnly bool) Option {
	return func(c *Config) {
		c.opt[String("-allphone_ci")] = ciUnitsOnly
	}
}

// LMFileOption sets word trigram language model input file.
func LMFileOption(filename string) Option {
	return func(c *Config) {
		c.opt[String("-lm")] = String(filename)
	}
}

// LMSetOption specifies a set of language model.
func LMSetOption(set string) Option {
	return func(c *Config) {
		c.opt[String("-lmctl")] = String(set)
	}
}

// LMNameOption sets which language model in LMSetOption to use by default.
func LMNameOption(name string) Option {
	return func(c *Config) {
		c.opt[String("-lmname")] = String(name)
	}
}

// LWeightOption sets language model probability weight.
//
// Default: 6.5
func LWeightOption(weight float32) Option {
	return func(c *Config) {
		c.opt[String("-lw")] = weight
	}
}

// FwdFlatLWeightOption sets language model probability weight for flat lexicon (2nd pass) decoding.
//
// Default: 8.5
func FwdFlatLWeightOption(weight float32) Option {
	return func(c *Config) {
		c.opt[String("-fwdflatlw")] = weight
	}
}

// BestPathLWeightOption sets language model probability weight for bestpath search.
//
// Default: 9.5
func BestPathLWeightOption(weight float32) Option {
	return func(c *Config) {
		c.opt[String("-bestpathlw")] = weight
	}
}

// AScaleOption sets inverse of acoustic model scale for confidence score calculation.
//
// Default: 20.0
func AScaleOption(scale float32) Option {
	return func(c *Config) {
		c.opt[String("-ascale")] = scale
	}
}

// WIPenaltyOption sets word insertion penalty.
//
// Default: 0.65
func WIPenaltyOption(penalty float32) Option {
	return func(c *Config) {
		c.opt[String("-wip")] = penalty
	}
}

// NewWordPenaltyOption sets new word transition penalty.
//
// Default: 1.0
func NewWordPenaltyOption(penalty float32) Option {
	return func(c *Config) {
		c.opt[String("-nwpen")] = penalty
	}
}

// PIPenaltyOption sets phone insertion penalty.
//
// Default: 1.0
func PIPenaltyOption(penalty float32) Option {
	return func(c *Config) {
		c.opt[String("-pip")] = penalty
	}
}

// UnigramWeightOption sets unigram weight.
//
// Default: 1.0
func UnigramWeightOption(weight float32) Option {
	return func(c *Config) {
		c.opt[String("-uw")] = weight
	}
}

// SilProbOption sets silence word transition probability.
//
// Default: 0.005
func SilProbOption(prob float32) Option {
	return func(c *Config) {
		c.opt[String("-silprob")] = prob
	}
}

// FillProbOption sets filler word transition probability.
//
// Default: 1e-8
func FillProbOption(prob float32) Option {
	return func(c *Config) {
		c.opt[String("-fillprob")] = prob
	}
}

// Options for dictionaries.

// DictFileOption sets main pronunciation dictionary (lexicon) input file.
func DictFileOption(filename string) Option {
	return func(c *Config) {
		c.opt[String("-dict")] = String(filename)
	}
}

// FillerDictFileOption sets noise word pronunciation dictionary input file.
func FillerDictFileOption(filename string) Option {
	return func(c *Config) {
		c.opt[String("-fdict")] = String(filename)
	}
}

// DictCaseOption enables if dictionary is case sensitive (NOTE: case insensitivity applies to ASCII characters only).
//
// Default: false
func DictCaseOption(sens bool) Option {
	return func(c *Config) {
		c.opt[String("-dictcase")] = sens
	}
}

// Options for acoustic modeling.

// HMMDirOption sets directory containing acoustic model files.
func HMMDirOption(dir string) Option {
	return func(c *Config) {
		c.opt[String("-hmm")] = String(dir)
	}
}

// FeatParamsFileOption sets file containing feature extraction parameters.
func FeatParamsFileOption(filename string) Option {
	return func(c *Config) {
		c.opt[String("-featparams")] = String(filename)
	}
}

// MDefFileOption sets model definition input file.
func MDefFileOption(filename string) Option {
	return func(c *Config) {
		c.opt[String("-mdef")] = String(filename)
	}
}

// TMatFileOption sets HMM state transition matrix input file.
func TMatFileOption(filename string) Option {
	return func(c *Config) {
		c.opt[String("-tmat")] = String(filename)
	}
}

// TMatFloorOption sets HMM state transition probability floor (applied to TMatFileOption file).
func TMatFloorOption(floor float32) Option {
	return func(c *Config) {
		c.opt[String("-tmatfloor")] = floor
	}
}

// MeansFileOption sets mixture gaussian means input file.
func MeansFileOption(filename string) Option {
	return func(c *Config) {
		c.opt[String("-mean")] = String(filename)
	}
}

// VarFileOption sets mixture gaussian variances input file.
func VarFileOption(filename int) Option {
	return func(c *Config) {
		c.opt[String("-var")] = String(filename)
	}
}

// VarFloorOption sets mixture gaussian variance floor (applied to data from VarFileOption file).
//
// Default: 0.0001
func VarFloorOption(floor float32) Option {
	return func(c *Config) {
		c.opt[String("-varfloor")] = floor
	}
}

// MixWFileOption sets senone mixture weights input file (uncompressed).
func MixWFileOption(filename string) Option {
	return func(c *Config) {
		c.opt[String("-mixw")] = String(filename)
	}
}

// MixWFloorOption sets senone mixture weights floor (applied to data from MixWFileOption file).
//
// Default: 0.0000001
func MixWFloorOption(floor float32) Option {
	return func(c *Config) {
		c.opt[String("-mixwfloor")] = floor
	}
}

// AWeightOption sets inverse weight applied to acoustic scores.
//
// Default: 1
func AWeightOption(weight int) Option {
	return func(c *Config) {
		c.opt[String("-aw")] = weight
	}
}

// SenDumpFileOption sets senone dump (compressed mixture weights) input file.
func SenDumpFileOption(filename string) Option {
	return func(c *Config) {
		c.opt[String("-sendump")] = String(filename)
	}
}

// MLLRFileOption sets MLLR transformation to apply to means and variances.
func MLLRFileOption(filename string) Option {
	return func(c *Config) {
		c.opt[String("-mllr")] = String(filename)
	}
}

// MMapOption enables use of memory-mapped I/O (if possible) for model files.
//
// Default: true
func MMapOption(mmap bool) Option {
	return func(c *Config) {
		c.opt[String("-mmap")] = mmap
	}
}

// DsRatioOption sets frame GMM computation downsampling ratio.
//
// Default: 1
func DsRatioOption(ratio int) Option {
	return func(c *Config) {
		c.opt[String("-ds")] = ratio
	}
}

// TopNOption sets maximum number of top Gaussians to use in scoring.
//
// Default: 4
func TopNOption(max int) Option {
	return func(c *Config) {
		c.opt[String("-topn")] = max
	}
}

// TopNBeamOption sets beam width used to determine top-N Gaussians (or a list, per-feature).
//
// Default: "0"
func TopNBeamOption(width string) Option {
	return func(c *Config) {
		c.opt[String("-topn_beam")] = String(width)
	}
}

// LogBaseOption sets base in which all log-likelihoods calculated.
//
// Default: 1.0001
func LogBaseOption(base float32) Option {
	return func(c *Config) {
		c.opt[String("-logbase")] = base
	}
}

// Misc options.

// SampleRateOption sets sample rate.
//
// Default: 16000.0
func SampleRateOption(rate float32) Option {
	return func(c *Config) {
		c.opt[String("-samprate")] = rate
	}
}

// InputEndianOption sets endianess of the input.
//
// Default: "little"
func InputEndianOption(endian string) Option {
	return func(c *Config) {
		c.opt[String("-input_endian")] = String(endian)
	}
}

// UserOption sets a user specified option to a custom value.
func UserOption(name string, v interface{}) Option {
	return func(c *Config) {
		c.opt[String(name)] = v
	}
}
