package output

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"golang.org/x/term"
)

var spinnerFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

// Progress — строка "⠹ Pulling 3/5", которая обновляется на месте.
// Пишет в stderr и только в терминал, чтобы не мусорить в pipe, файлах и CI.
type Progress struct {
	mu      sync.Mutex
	w       io.Writer
	enabled bool
	label   string
	total   int
	done    int
	frame   int

	stop    chan struct{}
	stopped chan struct{}
}

func NewProgress(label string, total int) *Progress {
	return newProgress(os.Stderr, term.IsTerminal(int(os.Stderr.Fd())), label, total, 100*time.Millisecond)
}

// newProgress с tick = 0 не крутит спиннер сам. Удобно в тестах.
func newProgress(w io.Writer, enabled bool, label string, total int, tick time.Duration) *Progress {
	p := &Progress{w: w, enabled: enabled, label: label, total: total}

	if !enabled {
		return p
	}

	p.render()

	if tick > 0 {
		p.stop = make(chan struct{})
		p.stopped = make(chan struct{})

		go p.spin(tick)
	}

	return p
}

func (p *Progress) spin(tick time.Duration) {
	defer close(p.stopped)

	ticker := time.NewTicker(tick)
	defer ticker.Stop()

	for {
		select {
		case <-p.stop:
			return
		case <-ticker.C:
			p.mu.Lock()
			p.frame++
			p.render()
			p.mu.Unlock()
		}
	}
}

// Advance отмечает ещё один готовый элемент. Можно вызывать из разных горутин.
func (p *Progress) Advance() {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.done++
	p.render()
}

// Finish останавливает спиннер и стирает строку перед выводом результата.
func (p *Progress) Finish() {
	if p.stop != nil {
		close(p.stop)
		<-p.stopped
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	if p.enabled {
		fmt.Fprint(p.w, "\r\033[K")
	}
}

// render перерисовывает строку: \r — в начало строки, \033[K — стереть её.
func (p *Progress) render() {
	if p.enabled {
		fmt.Fprintf(p.w, "\r\033[K%s %s %d/%d", spinnerFrames[p.frame%len(spinnerFrames)], p.label, p.done, p.total)
	}
}
