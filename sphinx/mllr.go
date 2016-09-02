package sphinx

import (
	"fmt"

	"github.com/xlab/pocketsphinx-go/pocketsphinx"
)

type MLLR struct {
	m *pocketsphinx.Mllr
}

func (m *MLLR) Retain() {
	m.m = pocketsphinx.MllrRetain(m.m)
}

func (m *MLLR) Destroy() bool {
	if m.m != nil {
		ret := pocketsphinx.MllrFree(m.m)
		m.m = nil
		return ret == 0
	}
	return true
}

// NewMLLR reads a speaker-adaptive linear transform from a file (mllr_matrix).
// See http://cmusphinx.sourceforge.net/wiki/tutorialadapt for details.
func NewMLLR(filename string) (*MLLR, error) {
	m := pocketsphinx.MllrRead(filename + "\x00")
	if m == nil {
		err := fmt.Errorf("sphinx: failed to load MLLR transform matrix from %s", filename)
		return nil, err
	}
	mllr := &MLLR{
		m: m,
	}
	return mllr, nil
}
