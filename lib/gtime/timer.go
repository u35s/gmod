package gtime

// 以 UnixNano 为单位
type Timer struct {
	lastNano uint
	interval uint
}

func (t *Timer) Init(interval uint) {
	t.interval = interval
	t.Reset()
}

func (t *Timer) Reset()       { t.lastNano = TimeNano() }
func (t *Timer) Elapse() uint { return TimeNano() - t.lastNano }

func (t *Timer) TimeUp(nano uint) bool {
	v := t.lastNano
	for nano > t.lastNano+t.interval {
		t.lastNano += t.interval
	}
	return v != t.lastNano
}
