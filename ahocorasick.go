package goahocorasick

import (
	"fmt"
)

import (
	"github.com/anknown/godarts"
	"unicode/utf8"
)

const FAIL_STATE = -1
const ROOT_STATE = 1

type Machine struct {
	trie    *godarts.DoubleArrayTrie
	failure []int
	output  [][]string // slice for speed
}

type Term struct {
	Pos  int
	Word string
}

func (m *Machine) Build(keywords [][]rune) (err error) {
	if len(keywords) == 0 {
		return fmt.Errorf("empty keywords")
	}

	d := new(godarts.Darts)

	trie := new(godarts.LinkedListTrie)
	m.trie, trie, err = d.Build(keywords)
	if err != nil {
		return err
	}

	output := make(map[int]([]string), 0)
	for idx, val := range d.Output {
		output[idx] = append(output[idx], string(val))
	}

	queue := make([](*godarts.LinkedListTrieNode), 0)
	m.failure = make([]int, len(m.trie.Base))
	for _, c := range trie.Root.Children {
		m.failure[c.Base] = godarts.ROOT_NODE_BASE
	}
	queue = append(queue, trie.Root.Children...)

	for {
		if len(queue) == 0 {
			break
		}

		node := queue[0]
		for _, n := range node.Children {
			if n.Base == godarts.END_NODE_BASE {
				continue
			}
			inState := m.f(node.Base)
		set_state:
			outState := m.g(inState, n.Code-godarts.ROOT_NODE_BASE)
			if outState == FAIL_STATE {
				inState = m.f(inState)
				goto set_state
			}
			if _, ok := output[outState]; ok != false {
				output[n.Base] = append(output[outState], output[n.Base]...)
			}
			m.setF(n.Base, outState)
		}
		queue = append(queue, node.Children...)
		queue = queue[1:]
	}

	// converting output to slice for speed
	maxOut := 0
	for key := range output{

		if key > maxOut{
			maxOut = key
		}
	}
	m.output = make([][]string, maxOut+1)
	for key, val := range output{
		m.output[key] = val
	}

	return nil
}

func (m *Machine) g(inState int, input rune) int {

	t := inState + int(input) + godarts.ROOT_NODE_BASE

	if t >= len(m.trie.Base) {
		if inState == ROOT_STATE {
			return ROOT_STATE
		}
		return FAIL_STATE
	}
	if inState == m.trie.Check[t] {
		return m.trie.Base[t]
	}

	if inState == ROOT_STATE {
		return ROOT_STATE
	}

	return FAIL_STATE
}

func (m *Machine) f(index int) (state int) {
	return m.failure[index]
}

func (m *Machine) setF(inState, outState int) {
	m.failure[inState] = outState
}

func (m *Machine) MultiPatternSearch(content string, returnImmediately bool) []Term {
	terms := make([]Term, 0, 16)

	state := ROOT_STATE
	for pos, c := range content {
	start:
		newState := m.g(state, c)
		if newState == FAIL_STATE {
			state = m.f(state)
			goto start
		} else {
			state = newState
			if state >= len(m.output){
				continue
			}
			for _, word := range m.output[state] {
				term := Term{
					Pos: pos + utf8.RuneLen(c) - len(word),
					Word: word,
				}
				terms = append(terms, term)
				if returnImmediately{
					return terms
				}
			}
		}
	}

	return terms
}
