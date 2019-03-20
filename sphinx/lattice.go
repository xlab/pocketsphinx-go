package sphinx

import (
	"fmt"

	"github.com/xlab/pocketsphinx-go/pocketsphinx"
)

// Lattice word graph structure used in bestpath/nbest search.
type Lattice struct {
	lat *pocketsphinx.Lattice
}

// LatticeLink represents links between DAG nodes.
//
// A link corresponds to a single hypothesized instance of a word with
// a given start and end point.
type LatticeLink pocketsphinx.Latlink

// LatticeLinkIter iterator over DAG links.
type LatticeLinkIter pocketsphinx.LatlinkIter

// LatticeNode represents DAG nodes.
//
// A node corresponds to a number of hypothesized instances of a word
// which all share the same starting point.
type LatticeNode pocketsphinx.Latnode

// LatticeNodeIter iterator over DAG nodes.
type LatticeNodeIter pocketsphinx.LatnodeIter

// NewLattice reads a lattice from a file on disk.
func (d *Decoder) NewLattice(filename String) (*Lattice, error) {
	lat := pocketsphinx.LatticeRead(d.dec, filename.S())
	if lat == nil {
		err := fmt.Errorf("sphinx: failed to load lattice from %s", filename)
		return nil, err
	}
	l := &Lattice{
		lat: lat,
	}
	return l, nil
}

// Lattice returns a retained copy of underlying reference to pocketsphinx.Lattice.
func (l *Lattice) Lattice() *pocketsphinx.Lattice {
	return pocketsphinx.LatticeRetain(l.lat)
}

func (l *Lattice) Retain() {
	l.lat = pocketsphinx.LatticeRetain(l.lat)
}

func (l *Lattice) Destroy() bool {
	if l.lat != nil {
		ret := pocketsphinx.LatticeFree(l.lat)
		l.lat = nil
		return ret == 0
	}
	return true
}

// WriteTo writes a lattice to disk.
func (l *Lattice) WriteTo(filename String) bool {
	ret := pocketsphinx.LatticeWrite(l.lat, filename.S())
	return ret == 0
}

// WriteToHTK writes a lattice to disk in HTK format.
func (l *Lattice) WriteToHTK(filename String) bool {
	ret := pocketsphinx.LatticeWriteHtk(l.lat, filename.S())
	return ret == 0
}

// LogMath gets the log-math computation object for this lattice.
//
// The lattice retains ownership of this pointer, so you should not attempt to
// free it manually. Use LogMath.Retain() if you wish to
// reuse it elsewhere.
func (l *Lattice) LogMath() *LogMath {
	m := pocketsphinx.LatticeGetLogmath(l.lat)
	return &LogMath{
		m: m,
	}
}

// Iter starts iterating over nodes in the lattice.
//
// No particular order of traversal is guaranteed, and you
// should not depend on this.
func (l Lattice) Iter() *LatticeNodeIter {
	iter := pocketsphinx.GetLatnodeIter(l.lat)
	return (*LatticeNodeIter)(iter)
}

// Next moves to next node in iteration.
func (l *LatticeNodeIter) Next() *LatticeNodeIter {
	iter := pocketsphinx.LatnodeIterNext((*pocketsphinx.LatnodeIter)(l))
	return (*LatticeNodeIter)(iter)
}

// Close stops iterating over nodes.
func (l *LatticeNodeIter) Close() {
	it := (*pocketsphinx.LatnodeIter)(l)
	pocketsphinx.LatnodeIterFree(it)
}

// Node gets node from iterator.
func (l *LatticeNodeIter) Node() *LatticeNode {
	it := (*pocketsphinx.LatnodeIter)(l)
	node := pocketsphinx.LatnodeIterNode(it)
	return (*LatticeNode)(node)
}

// Times gets start and end time range for a node.
//
// first — end frame of first exit from this node.
// last — end frame of last exit from this node.
// start — start frame for all edges exiting this node.
func (l *LatticeNode) Times() (start int32, first, last int16) {
	start = pocketsphinx.LatnodeTimes((*pocketsphinx.Latnode)(l), &first, &last)
	return
}

// Word gets word string for this node.
func (l *Lattice) Word(node *LatticeNode) string {
	return pocketsphinx.LatnodeWord(l.lat, (*pocketsphinx.Latnode)(node))
}

// BaseWord gets base word string for this node.
func (l *Lattice) BaseWord(node *LatticeNode) string {
	return pocketsphinx.LatnodeBaseWord(l.lat, (*pocketsphinx.Latnode)(node))
}

// Exits returns an iterator over exits from this node.
func (l *LatticeNode) Exits() *LatticeLinkIter {
	iter := pocketsphinx.LatnodeExits((*pocketsphinx.Latnode)(l))
	return (*LatticeLinkIter)(iter)
}

// Entries returns an iterator over entries to this node.
func (l *LatticeNode) Entries() *LatticeLinkIter {
	iter := pocketsphinx.LatnodeEntries((*pocketsphinx.Latnode)(l))
	return (*LatticeLinkIter)(iter)
}

// Close stops the iteration over links.
func (l *LatticeLinkIter) Close() {
	pocketsphinx.LatlinkIterFree((*pocketsphinx.LatlinkIter)(l))
}

// Next gets next link from a lattice link iterator.
func (l *LatticeLinkIter) Next() *LatticeLinkIter {
	iter := pocketsphinx.LatlinkIterNext((*pocketsphinx.LatlinkIter)(l))
	return (*LatticeLinkIter)(iter)
}

// ProbabilityOf node gets the best posterior probability and associated acoustic score from a lattice node.
// Returns exit link with highest posterior probability and the probability of this link.
//
// Log is expressed in the log-base used in the decoder. To convert to linear floating-point, use
// Lattice.LogMath().Exp(prob).
func (l *Lattice) ProbabilityOf(node *LatticeNode) (*LatticeLink, int32) {
	var outLink *pocketsphinx.Latlink
	prob := pocketsphinx.LatnodeProb(l.lat, (*pocketsphinx.Latnode)(node), &outLink)
	return (*LatticeLink)(outLink), prob
}

// Link gets link from iterator.
func (l *LatticeLinkIter) Link() *LatticeLink {
	link := pocketsphinx.LatlinkIterLink((*pocketsphinx.LatlinkIter)(l))
	return (*LatticeLink)(link)
}

// Times gets start and end times from a lattice link.
//
// start - start frame of this link.
// end - end frame of this link.
//
// These are inclusive, i.e. the last frame of
// this word is end, not end-1.
func (l *LatticeLink) Times() (start int32, end int16) {
	start = pocketsphinx.LatlinkTimes((*pocketsphinx.Latlink)(l), &end)
	return
}

// Nodes gets destination and source nodes from a lattice link
func (l *LatticeLink) Nodes() (source, dest *LatticeNode) {
	var s *pocketsphinx.Latnode
	d := pocketsphinx.LatlinkNodes((*pocketsphinx.Latlink)(l), &s)
	return (*LatticeNode)(s), (*LatticeNode)(d)
}

// Prev gets predecessor link in best path.
func (l *LatticeLink) Prev() *LatticeLink {
	link := pocketsphinx.LatlinkPred((*pocketsphinx.Latlink)(l))
	return (*LatticeLink)(link)
}

// LinkWord gets word string from a lattice link (possibly a pronunciation variant).
func (l *Lattice) LinkWord(link *LatticeLink) string {
	return pocketsphinx.LatlinkWord(l.lat, (*pocketsphinx.Latlink)(link))
}

// LinkBaseWord gets base word string from a lattice link.
func (l *Lattice) LinkBaseWord(link *LatticeLink) string {
	return pocketsphinx.LatlinkBaseword(l.lat, (*pocketsphinx.Latlink)(link))
}

// LinkProbability gets acoustic score and posterior probability from a lattice link.
//
// Posterior probability for this link. Log is expressed in
// the log-base used in the decoder. To convert to linear
// floating-point, use Lattice.LogMath().Exp(prob).
func (l *Lattice) LinkProbability(link *LatticeLink) (score, prob int32) {
	prob = pocketsphinx.LatlinkProb(l.lat, (*pocketsphinx.Latlink)(link), &score)
	return
}

// NewLink creates a directed link between from and to nodes, but if a link already exists,
// chooses one with the best score.
func (l *Lattice) NewLink(from, to *LatticeNode, score, endFrame int32) {
	pocketsphinx.LatticeLink(l.lat, (*pocketsphinx.Latnode)(from),
		(*pocketsphinx.Latnode)(to), score, endFrame)
}

// TraverseEdges starts a forward traversal of edges in a word graph.
//
// A keen eye will notice an inconsistency in this API versus
// other types of iterators in PocketSphinx. The reason for this is
// that the traversal algorithm is much more efficient when it is able
// to modify the lattice structure. Therefore, to avoid giving the
// impression that multiple traversals are possible at once, no
// separate iterator structure is provided.
func (l *Lattice) TraverseEdges(start, end *LatticeNode) *LatticeLink {
	link := pocketsphinx.LatticeTraverseEdges(l.lat, (*pocketsphinx.Latnode)(start), (*pocketsphinx.Latnode)(end))
	return (*LatticeLink)(link)
}

// TraverseNext gets the next link in forward traversal.
func (l *Lattice) TraverseNext(end *LatticeNode) *LatticeLink {
	link := pocketsphinx.LatticeTraverseNext(l.lat, (*pocketsphinx.Latnode)(end))
	return (*LatticeLink)(link)
}

// ReverseEdges starts a reverse traversal of edges in a word graph.
//
// See Lattice.TraverseEdges() for why this API is the way it is.
func (l *Lattice) ReverseEdges(start, end *LatticeNode) *LatticeLink {
	link := pocketsphinx.LatticeReverseEdges(l.lat, (*pocketsphinx.Latnode)(start), (*pocketsphinx.Latnode)(end))
	return (*LatticeLink)(link)
}

// ReverseNext gets the next link in reverse traversal.
func (l *Lattice) ReverseNext(start *LatticeNode) *LatticeLink {
	link := pocketsphinx.LatticeTraverseNext(l.lat, (*pocketsphinx.Latnode)(start))
	return (*LatticeLink)(link)
}

// BestPath does N-Gram based best-path search on a word graph using A*.
// Returns the final link in best path or nil upon error.
//
// This function calculates both the best path as well as the forward
// probability used in confidence estimation.
func (l *Lattice) BestPath(model *NGramModel, lwf, ascale float32) *LatticeLink {
	link := pocketsphinx.LatticeBestpath(l.lat, model.n, lwf, ascale)
	return (*LatticeLink)(link)
}

// Calculate link posterior probabilities on a word graph. Returns
// posterior probability of the utterance as a whole.
//
// WARN: This function assumes that Lattice.BestPath() search has already been done.
func (l *Lattice) Posterior(model *NGramModel, ascale float32) int32 {
	return pocketsphinx.LatticePosterior(l.lat, model.n, ascale)
}

// PosteriorPrune prunes all links (and associated nodes) below a certain posterior probability,
// return number of arcs removed.
//
// beam represents the minimum posterior probability for links. This is
// expressed in the log-base used in the decoder. To convert
// from linear floating-point, use Lattice.LogMath().Log(prob).
//
// WARN: This function assumes that Lattice.Posterior() has already been called.
func (l *Lattice) PosteriorPrune(model *NGramModel, ascale float32) int32 {
	return pocketsphinx.LatticePosterior(l.lat, model.n, ascale)
}

// NGramExpand expands lattice using an N-gram language model.
//
// This function expands the lattice such that each node represents a
// unique N-gram history, and adds language model scores to the links.
// func (l *Lattice) NGramExpand(model *NGramModel) bool {
// 	ret := pocketsphinx.LatticeNgramExpand(l.lat, model.n)
// 	return ret == 0
// }

// Frames gets the number of frames in the lattice.
func (l *Lattice) Frames() int32 {
	return pocketsphinx.LatticeNFrames(l.lat)
}
