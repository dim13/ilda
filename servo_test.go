package ilda

import "testing"

func TestServoStepResponse(t *testing.T) {
	s := DefaultServo()
	const sp = 20000.0
	const dt = 1.0 / 30000 // 30k pps

	var max, last float64
	for range 2000 {
		last = s.Step(sp, dt)
		if last > max {
			max = last
		}
	}

	// critically damped: no overshoot past the setpoint, monotonic settle
	if max > sp+1 {
		t.Errorf("expected no overshoot past %v, got peak %v", sp, max)
	}
	if d := last - sp; d > 50 || d < -50 {
		t.Errorf("expected settling near %v, got %v", sp, last)
	}
}
