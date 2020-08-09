package markov

import (
	"fmt"
)

// Flexible moving window with callback for printout
func (m *MarkovChain) processCorpusInternal(
	order, startOn int, msgCallback func(int)) {
	// # Moving window.
	for i := startOn; i < len(m.corpus)-order+1; i++ {
		window := m.corpus[i : i+order]

		// # Feedback: send i to caller for printout of progress.
		msgCallback(i)

		// # Update db with word pair relationship info.
		current, others := window[0], window[1:] // [tag #1]
		for dst, other := range others {
			m.db.IncrementPair(current, other, dst+1)
		}
	}
}

// ProcessCorpusByOrder uses 'order' to generate a markov chain.
// 'order' is the size of a moving window in this context.
func (m *MarkovChain) ProcessCorpusByOrder(order int, verbose bool) {
	size := len(m.corpus)
	// # Option for max window.
	if order < 0 || order > size {
		order = size
	}
	// # slice bounds on [tag #1] in processCorpusInternal.
	if order == 0 {
		order = 1
	}
	m.processCorpusInternal(order, 0, func(i int) {
		if verbose {
			fmt.Printf("\r Chunks remaining:%d", size-i)
		}
	})

}

// ProcessCorpusComplete avoids 'order' when generating a markov chain.
// This approach is more comprehensive but naturally requires more time.
func (m *MarkovChain) ProcessCorpusComplete(verbose bool) {
	size := len(m.corpus)
	for i := 0; i < size; i++ {

		m.processCorpusInternal(size-i, i, func(j int) {
			// @ callback not necessary; evaluating direction of changes.
			if verbose {
				s := size - i // scans remaining
				fmt.Printf("\r Scans Remaining:%d", s-1)
			}
		})

	}
}
