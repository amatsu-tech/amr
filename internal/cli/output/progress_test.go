package output

import (
	"bytes"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestProgress(t *testing.T) {
	var buf bytes.Buffer

	p := newProgress(&buf, true, "Pulling", 2, 0)
	p.Advance()
	p.Advance()
	p.Finish()

	want := "\r\033[K⠋ Pulling 0/2" + "\r\033[K⠋ Pulling 1/2" + "\r\033[K⠋ Pulling 2/2" + "\r\033[K"
	if got := buf.String(); got != want {
		t.Errorf("output = %q, want %q", got, want)
	}
}

func TestProgressDisabled(t *testing.T) {
	var buf bytes.Buffer

	p := newProgress(&buf, false, "Pulling", 2, time.Millisecond)
	p.Advance()
	p.Finish()

	if buf.Len() != 0 {
		t.Errorf("disabled progress wrote %q, want nothing", buf.String())
	}
}

func TestProgressSpinsAndStops(t *testing.T) {
	var buf safeBuffer

	p := newProgress(&buf, true, "Pulling", 1, time.Millisecond)
	time.Sleep(20 * time.Millisecond)
	p.Finish()

	output := buf.String()

	// Спиннер успел смениться хотя бы раз.
	if !strings.Contains(output, "⠙") {
		t.Errorf("spinner did not advance: %q", output)
	}

	// После Finish спиннер больше ничего не пишет.
	time.Sleep(10 * time.Millisecond)

	if buf.String() != output {
		t.Error("spinner kept writing after Finish")
	}
}

func TestProgressConcurrent(t *testing.T) {
	var buf safeBuffer

	const total = 50

	p := newProgress(&buf, true, "Cloning", total, time.Millisecond)

	var wg sync.WaitGroup
	for range total {
		wg.Go(p.Advance)
	}
	wg.Wait()
	p.Finish()

	if !strings.Contains(buf.String(), "Cloning 50/50") {
		t.Error("progress did not reach 50/50")
	}
}

// safeBuffer — bytes.Buffer, который можно читать, пока в него пишет спиннер.
type safeBuffer struct {
	mu  sync.Mutex
	buf bytes.Buffer
}

func (b *safeBuffer) Write(p []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	return b.buf.Write(p)
}

func (b *safeBuffer) String() string {
	b.mu.Lock()
	defer b.mu.Unlock()

	return b.buf.String()
}
