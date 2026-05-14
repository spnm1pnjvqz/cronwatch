package alert

// Notifier is the interface implemented by all alert delivery mechanisms.
type Notifier interface {
	Send(a Alert) error
}

// MultiNotifier fans out a single alert to multiple Notifier implementations.
type MultiNotifier struct {
	notifiers []Notifier
}

// NewMultiNotifier creates a MultiNotifier that sends alerts to all provided notifiers.
func NewMultiNotifier(notifiers ...Notifier) *MultiNotifier {
	return &MultiNotifier{notifiers: notifiers}
}

// Send delivers the alert to every registered notifier.
// It collects all errors and returns the first one encountered, but still
// attempts delivery to all remaining notifiers.
func (m *MultiNotifier) Send(a Alert) error {
	var firstErr error
	for _, n := range m.notifiers {
		if err := n.Send(a); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}
