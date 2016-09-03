package sphinx

import (
	"fmt"

	"github.com/xlab/pocketsphinx-go/pocketsphinx"
)

// NGramModel is a type representing an N-Gram based language model.
type NGramModel struct {
	n *pocketsphinx.NgramModel
}

// NGramModel returns a retained copy of underlying reference to pocketsphinx.NgramModel.
func (n *NGramModel) NGramModel() *pocketsphinx.NgramModel {
	return pocketsphinx.Retain(n.n)
}

// NGramOptions is an utility to construct CommandLn for NewNGramModel.
type NGramOptions struct {
	opt       map[String]interface{}
	evaluated *pocketsphinx.CommandLn
}

func (n *NGramOptions) CommandLn() *pocketsphinx.CommandLn {
	if n.evaluated != nil {
		return n.evaluated
	}
	n.evaluated = optToCommandLn(nil, n.opt)
	return n.evaluated
}

func (n *NGramOptions) Destroy() bool {
	// TODO: not sure if the CommandLn is retained by us
	// after passing to the C land.
	//
	// if n.evaluated != nil {
	// 	n.evaluated.Free()
	// 	n.evaluated = nil
	// }
	// return true

	if n.evaluated != nil {
		ret := pocketsphinx.CommandLnFreeR(n.evaluated)
		n.evaluated = nil
		return ret == 0
	}
	return true
}

// MMap options sets whether to use memory-mapped I/O.
func (n *NGramOptions) MMap(v bool) {
	n.opt[String("-mmap")] = v
}

// LanguageWeight to apply to the model.
func (n *NGramOptions) LanguageWeight(v float32) {
	n.opt[String("-lw")] = v
}

// WordInsertionPenalty to apply to the model.
func (n *NGramOptions) WordInsertionPenalty(v float32) {
	n.opt[String("-wip")] = v
}

// NewNGramModel reads an N-Gram model from a file on disk.
//
// lmath carries log-math parameters to use for probability
// calculations. Ownership of this object is assumed by
// the newly created NGramModel, and you should not
// attempt to free it manually. If you wish to reuse it
// elsewhere, you must retain it with LogMath.Retain().
func NewNGramModel(fileName String, fileType NGramFileType,
	lmath *LogMath, opt ...NGramOptions) (*NGramModel, error) {
	if len(opt) > 0 {
		m := pocketsphinx.NgramModelRead(opt[0].CommandLn(), fileName.S(), fileType, lmath.m)
		if m == nil {
			err := fmt.Errorf("sphinx: failed to load n-gram model from %s", filename)
			return nil, err
		}
		return m, nil
	}
	m := pocketsphinx.NgramModelRead(nil, fileName.S(), fileType, lmath.m)
	if m == nil {
		err := fmt.Errorf("sphinx: failed to load n-gram model from %s", filename)
		return nil, err
	}
	return m, nil
}

func (n *NGramModel) Destroy() bool {
	if n.n != nil {
		ret := pocketsphinx.NgramModelFree(n.n)
		n.n = nil
		return ret == 0
	}
	return true
}

func (n *NGramModel) Retain() {
	n.n = pocketsphinx.NgramModelRetain(n.n)
}

// NgramFileType as declared in sphinxbase/ngram_model.h:81
type NGramFileType int32

// File types for N-Gram files.
const (
	// NGramInvalid is not a valid file type.
	NGramInvalid NGramFileType = pocketsphinx.NgramInvalid
	// NGramAuto to determine file type automatically.
	NGramAuto NGramFileType = pocketsphinx.NgramAuto
	// NGramArpa is for ARPABO text format (the standard).
	NGramArpa NGramFileType = pocketsphinx.NgramArpa
	// NGramBin is for sphinx .DMP format.
	NGramBin NGramFileType = pocketsphinx.NgramBin
)

// NgramCase as declared in sphinxbase/ngram_model.h:166
type NGramCase int32

// Constants for case folding.
const (
	NGramUpper NGramCase = pocketsphinx.NgramUpper
	NGramLower NGramCase = pocketsphinx.NgramLower
)

// CaseFold word strings in an N-Gram model.
//
// WARNING: This is not Unicode aware, so any non-ASCII characters
// will not be converted.
func (n *NGramModel) CaseFold(c NGramCase) bool {
	ret := pocketsphinx.NgramModelCasefold(n.n, c)
	return ret == 0
}

// WriteTo writes an N-Gram model to disk.
func (n *NGramModel) WriteTo(filename String, format NGramFileType) bool {
	ret := pocketsphinx.NgramModelWrite(n.n, filename.S(), format)
	return ret == 0
}

// ApplyWeights applies a language weight and insertion penalty weight to a
// language model.
//
// This will change the values output by NGramModel.Score() and friends.
// This is done for efficiency since in decoding, these are the only
// values we actually need. Call NGramModel.Probability() if you want the "raw"
// N-Gram probability estimate.
//
// To remove all weighting, call NGramModel.ApplyWeights(1.0, 1.0).
func (n *NGramModel) ApplyWeights(langWeight, insertionPenalty float32) bool {
	ret := pocketsphinx.NgramModelApplyWeights(n.n, langWeight, insertionPenalty)
	return ret == 0
}

// Weights gets the current language weight from a language model
// and logarithm of word insertion penalty.
func (n *NGramModel) Weights() (langWeight float32, wipLog int32) {
	langWeight = pocketsphinx.NgramModelGetWeights(n.n, &wipLog)
	return
}

// TrigramScore does quick trigram score lookup.
func (n *NGramModel) TrigramScore(w3, w2, w1 int32) (score, nUsed int32) {
	score = pocketsphinx.NgramTgScore(n.n, w3, w2, w1, &nUsed)
	return
}

// TrigramScore does quick bigram score lookup.
func (n *NGramModel) BigramScore(w2, w1 int32) (score, nUsed int32) {
	score = pocketsphinx.NgramBgScore(n.n, w2, w1, &nUsed)
	return
}

// Score get the score (scaled, interpolated log-probability) for a general
// N-Gram. See TrigramScore and BigramScore for particular cases.
//
// If one of the words is not in the LM's vocabulary, the result will
// depend on whether this is an open or closed vocabulary language
// model.  For an open-vocabulary model, unknown words are all mapped
// to the unigram "UNK" which has a non-zero probability and also
// participates in higher-order N-Grams. Therefore, you will get a
// score of some sort in this case.
//
// For a closed-vocabulary model, unknown words are impossible and
// thus have zero probability.  Therefore, if wordID is
// unknown, this function will return a "zero" log-probability, i.e. a
// large negative number. To obtain this number for comparison, call
// NGramModel.Zero().
func (n *NGramModel) Score(wordID int32, history []int32) (score, nUsed int32) {
	score = pocketsphinx.NgramNgScore(n.n, wordID, history, int32(len(history)), &nUsed)
	return
}

// Probability gets the "raw" log-probability for a general N-Gram.
//
// This returns the log-probability of an N-Gram, as defined in the
// language model file, before any language weighting, interpolation,
// or insertion penalty has been applied.
//
// When backing off to a unigram from a bigram or trigram, the
// unigram weight (interpolation with uniform) is not removed.
func (n *NGramModel) Probability(words Strings) int32 {
	return pocketsphinx.NgramProb(n.n, words.S(), int32(len(words)))
}

// Quick "raw" probability lookup for a general N-Gram.
//
// See documentation for NGramModel.Score() and NGramModel.ApplyWeights()
// for an explanation of this.
func (n *NGramModel) QuickProbability(wordID int32, history []int32) (prob, nUsed int32) {
	prob = pocketsphinx.NgramNgProb(n.n, wordID, history, int32(len(history)), &nUsed)
	return
}

// ScoreToProbability converts score to "raw" log-probability.
//
// The unigram weight (interpolation with uniform) is not
// removed, since there is no way to know which order of N-Gram
// generated score.
func (n *NGramModel) ScoreToProbability(score int32) int32 {
	return pocketsphinx.NgramScoreToProb(n.n, score)
}

// WordID lookups numerical word ID.
func (n *NGramModel) WordID(word String) int32 {
	return pocketsphinx.NgramWid(n.n, word.S())
}

// Word lookups word string for numerical wordID.
func (n *NGramModel) Word(wordID int32) string {
	return pocketsphinx.NgramWord(n.n, wordID)
}

// UnknownWordID gets the unknown word ID for a language model.
//
// Language models can be either "open vocabulary" or "closed
// vocabulary". The difference is that the former assigns a fixed
// non-zero unigram probability to unknown words, while the latter
// does not allow unknown words (or, equivalently, it assigns them
// zero probability). If this is a closed vocabulary model, this
// function will return NGgramInvalidWordID.
//
// The ID for the unknown word, or NGgramInvalidWordID if none
// exists.
func (n *NGramModel) UnknownWordID() int32 {
	return pocketsphinx.NgramUnknownWid(n.n)
}

var NGgramInvalidWordID int32 = pocketsphinx.NgramInvalidWid

// Zero gets the "zero" log-probability value for a language model.
func (n *NGramModel) Zero() int32 {
	return pocketsphinx.NgramZero(n.n)
}

// Size gets the order of the N-gram model (i.e. the "N" in "N-gram")
func (n *NGramModel) Size() int32 {
	return pocketsphinx.NgramModelGetSize(n.n)
}

// Counts gets the counts of the various N-grams in the model.
func (n *NGramModel) Counts() []uint32 {
	return pocketsphinx.NgramModelGetCounts(n.n)
}

// AddWord adds a word (unigram) to the language model and returns
// the word ID for the new word.
//
// The semantics of this are not particularly well-defined for
// model sets, and may be subject to change. Currently this will add
// the word to all of the submodels
func (n *NGramModel) AddWord(word String, weight int32) int32 {
	id := pocketsphinx.NgramModelAddWord(n.n, word.S(), weight)
	return id
}

// ReadClassDef reads a class definition file and add classes to a language model.
//
// This function assumes that the class tags have already been defined
// as unigrams in the language model. All words in the class
// definition will be added to the vocabulary as special in-class words.
// For this reason is is necessary that they not have the same names
// as any words in the general unigram distribution. The convention
// is to suffix them with ":class_tag", where class_tag is the class
// tag minus the enclosing square brackets.
func (n *NGramModel) ReadClassDef(filename String) bool {
	ret := pocketsphinx.NgramModelReadClassdef(n.n, filename.S())
	return ret == 0
}

// Add a new class to a language model.
//
// If className already exists in the unigram set for NGramModel,
// then it will be converted to a class tag, and weight will be ignored.
// Otherwise, a new unigram will be created as in NGramModel.AddWord().
func (n *NGramModel) AddClass(className String, weight float32, words Strings, weights []float32) bool {
	ret := pocketsphinx.NgramModelAddClass(n.n, className.S(), weight, words.S(), weights, int32(len(words)))
	return ret == 0
}

// AddClassWord adds a word to a class in a language model with a weight
// of this word relative to the within-class uniform distribution. Returns word ID.
func (n *NGramModel) AddClassWord(className, word String, weight float32) int32 {
	id := pocketsphinx.NgramModelAddClassWord(n.n, className.S(), word.S(), weight)
	return id
}

// TODO: implement model sets

// Flush any cached N-Gram information.
func (n *NGramModel) Flush() {
	pocketsphinx.NgramModelFlush(n.n)
}
