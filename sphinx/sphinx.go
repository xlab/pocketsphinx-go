package sphinx

import (
	"errors"
	"time"

	"github.com/xlab/pocketsphinx-go/pocketsphinx"
)

type Decoder struct {
	cfg *Config
	dec *pocketsphinx.Decoder

	maxRawdataSize int32
	rawdataBuf     [][]int16
}

// Config gets the configuration object for this decoder.
func (d *Decoder) Config() *Config {
	return d.cfg
}

// NewDecoder initializes the decoder from a configuration object.
func NewDecoder(cfg *Config) (*Decoder, error) {
	if cfg == nil {
		cfg = NewConfig()
	}
	dec := &Decoder{
		cfg: cfg,
		dec: pocketsphinx.Init(cfg.CommandLn()),
	}
	if dec.dec == nil {
		cfg.Destroy()
		err := errors.New("pocketsphinx.Init failed")
		return nil, err
	}
	dec.SetRawDataSize(0)
	return dec, nil
}

// Reconfigure reinitializes the decoder with updated configuration.
//
// This function allows you to switch the acoustic model, dictionary,
// or other configuration without creating an entirely new decoding
// object.
//
// An optional new configuration to use. If cfg is
// nil, the previous configuration will be reloaded,
// with any changes applied.
func (d *Decoder) Reconfigure(cfg *Config) {
	pocketsphinx.Reinit(d.dec, cfg.CommandLn())
}

func (d *Decoder) Destroy() bool {
	if d.dec != nil {
		ret := pocketsphinx.Free(d.dec)
		d.dec = nil
		return ret == 0
	}
	d.cfg.Destroy()
	return true
}

// LogMath gets the log-math computation object for this decoder.
//
// The decoder retains ownership of this pointer, so you should not attempt to
// free it manually. Use LogMath.Retain() if you wish to
// reuse it elsewhere.
func (d *Decoder) LogMath() *LogMath {
	return &LogMath{
		m: pocketsphinx.GetLogmath(d.dec),
	}
}

// Decoder returns a retained copy of underlying reference to pocketsphinx.Decoder.
func (d *Decoder) Decoder() *pocketsphinx.Decoder {
	return pocketsphinx.Retain(d.dec)
}

// UpdateMLLR adapts current acoustic model using a linear transform (Maximum Likelihood Linear Regression).
//
// mllr is the new transform to use, or nil to update the existing
// transform. The decoder retains ownership of this pointer,
// so you should not attempt to free it manually. Use
// MLLR.Retain() if you wish to reuse it elsewhere.
//
// Returns the updated transform object for this decoder, or
// nil on failure.
func (d *Decoder) UpdateMLLR(mllr *MLLR) *MLLR {
	var m *pocketsphinx.Mllr
	if mllr == nil {
		m = pocketsphinx.UpdateMllr(d.dec, nil)
	} else {
		m = pocketsphinx.UpdateMllr(d.dec, mllr.m)
	}
	if m != nil {
		return &MLLR{
			m: m,
		}
	}
	return nil
}

// ReadDict reloads the pronunciation dictionary from a file.
//
// This function replaces the current pronunciation dictionary with
// the one stored in dictFile. This also causes the active search
// module(s) to be reinitialized, in the same manner as calling
// Decoder.AddWord() with update=true.
//
// dictFile is the path to dictionary file to load.
// fillerDictFile is the path to filler dictionary to load,
// or empty string to keep the existing filler dictionary.
func (d *Decoder) ReadDict(dictFile, fillerDictFile String) bool {
	ret := pocketsphinx.LoadDict(d.dec, dictFile.S(), fillerDictFile.S(), end)
	return ret == 0
}

// WriteDict writes the current pronunciation dictionary to a file.
func (d *Decoder) WriteDict(dictFile String) bool {
	ret := pocketsphinx.SaveDict(d.dec, dictFile.S(), end)
	return ret == 0
}

// AddWord adds a word to the pronunciation dictionary.
//
// This function adds a word to the pronunciation dictionary and the
// current language model (but not to the current FSG if
// FSG mode is enabled). If the word is already present in one or the
// other, it does whatever is necessary to ensure that the word can be
// recognized.
//
// word is a word string to add, e.g. "hello". phones is a whitespace-separated list of phoneme strings
// describing pronunciation of the word, e.g. "HH AH L OW". If update is true, updates the
// search module (whichever one is currently active) to recognize the newly added word.
// If adding multiple words, it is more efficient to pass false here in all but the last word.
//
// Returns the internal ID (>= 0) of the newly added word.
func (d *Decoder) AddWord(word, phones String, update bool) (id int32, ok bool) {
	ret := pocketsphinx.AddWord(d.dec, word.S(), phones.S(), b(update))
	if ret < 0 {
		return 0, false
	}
	return ret, true
}

// LookupWord lookups for the word in the dictionary and returns phone transcription for it.
//
// Returns whitespace-spearated phone string describing the pronunciation of the word,
// or empty string if word is not present in the dictionary.
func (d *Decoder) LookupWord(word String) (string, bool) {
	phones := pocketsphinx.LookupWord(d.dec, word.S())
	if phones != nil {
		return string(phones), true
	}
	return "", false
}

// StartStream starts processing of the stream of speech. Channel parameters like
// noise-level are maintained for the stream and reused among utterances.
// Times returned in segment iterators are also stream-wide.
func (d *Decoder) StartStream() bool {
	ret := pocketsphinx.StartStream(d.dec)
	return ret == 0
}

// StartUtt starts utterance processing.
// This function should be called before any utterance data is passed
// to the decoder. It marks the start of a new utterance and
// reinitializes internal data structures.
func (d *Decoder) StartUtt() bool {
	ret := pocketsphinx.StartUtt(d.dec)
	return ret == 0
}

// EndUtt ends utterance processing.
func (d *Decoder) EndUtt() bool {
	ret := pocketsphinx.EndUtt(d.dec)
	return ret == 0
}

// ProcessRaw decodes a raw audio stream.
//
// No headers are recognized in this files. The configuration
// parameters SampleRateOption and InputEndianOption are used
// to determine the sampling rate and endianness of the stream,
// respectively. Audio is always assumed to be 16-bit signed PCM.
//
// If noSearch is enabled, performs feature extraction but does no
// any recognition yet. This may be necessary if your processor has
// trouble doing recognition in real-time.
//
// fullUtterance shows that this block of data is a full utterance
// worth of data. This may allow the recognizer to
// produce more accurate results.
//
// Returns number of frames of data searched.
func (d *Decoder) ProcessRaw(data []int16, noSearch, fullUtterance bool) (frames int32, ok bool) {
	frames = pocketsphinx.ProcessRaw(d.dec, data, uint(len(data)), b(noSearch), b(fullUtterance))
	if frames < 0 {
		return 0, false
	}
	return frames, true
}

// ProcessCep decodes acoustic feature data.
//
// If noSearch is enabled, performs feature extraction but does no
// any recognition yet. This may be necessary if your processor has
// trouble doing recognition in real-time.
//
// fullUtterance shows that this block of data is a full utterance
// worth of data. This may allow the recognizer to
// produce more accurate results.
//
// Returns number of frames of data searched.
func (d *Decoder) ProcessCep(data [][]float32, noSearch, fullUtterance bool) (frames int32, ok bool) {
	frames = pocketsphinx.ProcessCep(d.dec, data, int32(len(data)), b(noSearch), b(fullUtterance))
	if frames < 0 {
		return 0, false
	}
	return frames, true
}

// FramesSearched gets the number of frames of data searched.
//
// Note that there is a delay between this and the number of frames of
// audio which have been input to the system. This is due to the fact
// that acoustic features are computed using a sliding window of
// audio, and dynamic features are computed over a sliding window of
// acoustic features.
//
// Returns number of frames of speech data which have been recognized so far.
func (d *Decoder) FramesSearched() int32 {
	return pocketsphinx.GetNFrames(d.dec)
}

// Hypothesis gets hypothesis string and path score.
//
// Returns string containing best hypothesis at this point in
// decoding. Empty if no hypothesis is available. And path score of that string.
func (d *Decoder) Hypothesis() (hyp string, score int32) {
	hyp = pocketsphinx.GetHyp(d.dec, &score)
	return
}

// Probability gets posterior probability of the best hypothesis.
//
// Unless the BestpathOption option is enabled, this function will
// always return zero (corresponding to a posterior probability of
// 1.0). Even if BestpathOption is enabled, it will also return zero when
// called on a partial result. Ongoing research into effective
// confidence annotation for partial hypotheses may result in these
// restrictions being lifted in future versions.
func (d *Decoder) Probability() int32 {
	return pocketsphinx.GetProb(d.dec)
}

// WordLattice gets the word lattice object containing all hypotheses so far.
//
// The pointer is owned by the decoder and you should not attempt to free it manually.
// It is only valid until the next utterance, unless you use
// Lattice.Retain() to retain it.
func (d *Decoder) WordLattice() *Lattice {
	lat := pocketsphinx.GetLattice(d.dec)
	return &Lattice{
		lat: lat,
	}
}

// UttDuration gets performance information for the current utterance.
//
// speech — number of seconds of speech.
// cpu — number of seconds of CPU time used.
// wall — number of seconds of wall time used.
func (d *Decoder) UttDuration() (speech, cpu, wall time.Duration) {
	var speechSeconds float64
	var cpuSeconds float64
	var wallSeconds float64
	pocketsphinx.GetUttTime(d.dec, &speechSeconds, &cpuSeconds, &wallSeconds)
	speech = time.Duration(speechSeconds * float64(time.Second))
	cpu = time.Duration(cpuSeconds * float64(time.Second))
	wall = time.Duration(wallSeconds * float64(time.Second))
	return
}

// AllDuration gets overall performance information.
//
// speech — number of seconds of speech.
// cpu — number of seconds of CPU time used.
// wall — number of seconds of wall time used.
func (d *Decoder) AllDuration() (speech, cpu, wall time.Duration) {
	var speechSeconds float64
	var cpuSeconds float64
	var wallSeconds float64
	pocketsphinx.GetAllTime(d.dec, &speechSeconds, &cpuSeconds, &wallSeconds)
	speech = time.Duration(speechSeconds * float64(time.Second))
	cpu = time.Duration(cpuSeconds * float64(time.Second))
	wall = time.Duration(wallSeconds * float64(time.Second))
	return
}

// IsInSpeech checks if the last feed audio buffer contained speech.
func (d *Decoder) IsInSpeech() bool {
	v := pocketsphinx.GetInSpeech(d.dec)
	return v == 1
}

// SetRawDataSize sets the limit of the raw audio data to store in decoder
// to retrieve it later with Decoder.RawData().
func (d *Decoder) SetRawDataSize(frames int32) {
	d.maxRawdataSize = frames
	d.rawdataBuf = [][]int16{
		make([]int16, frames),
	}
	pocketsphinx.SetRawdataSize(d.dec, frames*2)
}

// RawData retrieves the raw data collected during utterance decoding.
func (d *Decoder) RawData() []int16 {
	var size int32
	pocketsphinx.GetRawdata(d.dec, d.rawdataBuf, &size)
	return d.rawdataBuf[0][:size/2]
}
