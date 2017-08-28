package sphinx

import (
	"errors"
	"fmt"
	"runtime"

	"github.com/xlab/pocketsphinx-go/pocketsphinx"
)

/**
 * @file jsgf.go JSGF grammar compiler
 *
 * This file defines the data structures for parsing JSGF grammars
 * into Sphinx finite-state grammars.
 **/

type JSGF struct {
	j *pocketsphinx.JSGF
}

// NewJSGFGrammar creates a new JSGF grammar. Parent is optional parent
// grammar for this one (nil, usually). Rturns new JSGF grammar object, or nil on failure.
func NewJSGFGrammar(parent *JSGF) (*JSGF, error) {
	var p *pocketsphinx.JSGF
	if parent != nil {
		p = parent.j
	}

	grammar := pocketsphinx.JSGFGrammarNew(p)
	if grammar != nil {
		runtime.SetFinalizer(grammar, func(j *pocketsphinx.JSGF) {
			pocketsphinx.JSGFGrammarFree(j)
		})
		return &JSGF{
			j: grammar,
		}, nil
	}
	err := errors.New("pocketsphinx: failed to create JSGF grammar")
	return nil, err
}

func JSGFParseFile(filename String, parent *JSGF) (*JSGF, error) {
	var p *pocketsphinx.JSGF
	if parent != nil {
		p = parent.j
	}

	grammar := pocketsphinx.JSGFParseFile(filename.S(), p)
	if grammar != nil {
		runtime.SetFinalizer(grammar, func(j *pocketsphinx.JSGF) {
			pocketsphinx.JSGFGrammarFree(j)
		})
		return &JSGF{
			j: grammar,
		}, nil
	}
	err := fmt.Errorf("pocketsphinx: failed to parse JSGF grammar from file %s", filename)
	return nil, err
}

func JSGFParseString(data String, parent *JSGF) (*JSGF, error) {
	var p *pocketsphinx.JSGF
	if parent != nil {
		p = parent.j
	}

	grammar := pocketsphinx.JSGFParseString(data.S(), p)
	if grammar != nil {
		runtime.SetFinalizer(grammar, func(j *pocketsphinx.JSGF) {
			pocketsphinx.JSGFGrammarFree(j)
		})
		return &JSGF{
			j: grammar,
		}, nil
	}
	err := errors.New("pocketsphinx: failed to parse JSGF grammar")
	return nil, err
}

func (j *JSGF) GrammarName() string {
	if j == nil || j.j == nil {
		return ""
	}
	return pocketsphinx.JSGFGrammarName(j.j)
}

// func (j *JSGF) Free() {
// 	if j == nil || j.j == nil {
// 		return
// 	}
// 	pocketsphinx.JSGFGrammarFree(j.j)
// }

type JSGFRuleIter pocketsphinx.JSGFRuleIter
