package sphinx

import "github.com/xlab/pocketsphinx-go/pocketsphinx"

type Lattice struct {
	lat *pocketsphinx.Lattice
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

type LatticeIter pocketsphinx.LatnodeIter

func (l Lattice) Iter() *LatticeIter {
	iter := pocketsphinx.GetLatnodeIter(l.lat)
	return (*LatticeIter)(iter)
}

func (l *LatticeIter) Next() *LatticeIter {
	it := (*pocketsphinx.LatnodeIter)(l)
	next := pocketsphinx.LatnodeIterNext(it)
	return (*LatticeIter)(next)
}

func (l *LatticeIter) Close() {
	it := (*pocketsphinx.LatnodeIter)(l)
	pocketsphinx.LatnodeIterFree(it)
}

type LatticeNode pocketsphinx.Latnode

func (l *LatticeIter) Node() *LatticeNode {
	it := (*pocketsphinx.LatnodeIter)(l)
	node := pocketsphinx.LatnodeIterNode(it)
	return (*LatticeNode)(node)
}

func (l *Lattice) Word(node *LatticeNode) string {
	return pocketsphinx.LatnodeWord(l.lat, (*pocketsphinx.Latnode)(node))
}

func (l *Lattice) BaseWord(node *LatticeNode) string {
	return pocketsphinx.LatnodeWord(l.lat, (*pocketsphinx.Latnode)(node))
}
