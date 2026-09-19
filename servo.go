package ilda

// PID is a discrete PID controller.
type PID struct {
	Kp, Ki, Kd float64
	integral   float64
	prevErr    float64
}

// Step advances the controller by dt seconds toward setpoint sp from
// measured value pv, returning the control output.
func (c *PID) Step(sp, pv, dt float64) float64 {
	err := sp - pv
	c.integral += err * dt
	deriv := (err - c.prevErr) / dt
	c.prevErr = err
	return c.Kp*err + c.Ki*c.integral + c.Kd*deriv
}

// Servo emulates one axis of a laser galvanometer: a PID controller driving
// a mirror with mass, fighting inertia and damping. Unlike an ideal point-to-
// point jump, this produces the overshoot and settling time real galvos show
// on fast moves.
type Servo struct {
	PID
	Inertia float64 // mirror+coil moment of inertia
	Damping float64 // mechanical/eddy-current damping

	pos, vel float64
}

// NewServo builds a galvo axis. Kp/Ki/Kd tune the controller, inertia and
// damping the mechanics; larger inertia or smaller damping means slower
// settling and more overshoot/ringing.
func NewServo(kp, ki, kd, inertia, damping float64) *Servo {
	return &Servo{PID: PID{Kp: kp, Ki: ki, Kd: kd}, Inertia: inertia, Damping: damping}
}

// DefaultServo returns a critically damped galvo (no overshoot) whose
// bandwidth is deliberately below a typical ILDA test pattern's fastest
// direction changes, at 30 kpps: dwelled corners (many repeated points, as
// in a square) are reached cleanly, while a sparse fast circle scan is
// visibly rounded and undersized, same as a real galvo's finite bandwidth
// would render it. Ki and Kd are left at zero: the mechanical damping term
// already provides velocity feedback, so adding derivative-of-error on top
// would just reintroduce derivative kick every time the setpoint jumps.
func DefaultServo() *Servo {
	return NewServo(4.6225e8, 0, 0, 1, 43000)
}

// moveSubsteps is the number of internal physics steps per Move call. A
// single Euler step per ILDA point is too coarse for stiff (high-bandwidth)
// gains: it shows spurious numerical overshoot that isn't real servo
// behavior. Sub-stepping fixes the integration without changing the output
// cadence (still one reported position per input point).
const moveSubsteps = 8

// step advances the servo by one physics tick of size h.
func (s *Servo) step(sp, h float64) float64 {
	force := s.PID.Step(sp, s.pos, h)
	accel := (force - s.Damping*s.vel) / s.Inertia
	s.vel += accel * h
	s.pos += s.vel * h
	return s.pos
}

// Step the servo toward setpoint sp by dt seconds and returns the
// resulting position.
func (s *Servo) Step(sp, dt float64) float64 {
	h := dt / moveSubsteps
	var pos float64
	for range moveSubsteps {
		pos = s.step(sp, h)
	}
	return pos
}

func clampInt16(v float64) int16 {
	switch {
	case v < -32768:
		return -32768
	case v > 32767:
		return 32767
	default:
		return int16(v)
	}
}

// Move replaces each point with the path a real galvo pair (x, y) traces
// tracking these points as setpoints at the given scan rate (points per
// second): every physics sub-step between two ILDA points is emitted as its
// own point, so the drawn path curves between setpoints instead of jumping
// straight to each one. Z and color pass through unchanged.
func (f Frame) Move(pps float64, x, y *Servo) Frame {
	h := 1 / pps / moveSubsteps
	out := f
	out.Points = make([]Point, 0, len(f.Points)*moveSubsteps)
	for _, pt := range f.Points {
		for range moveSubsteps {
			out.Points = append(out.Points, Point{
				X:     clampInt16(x.step(float64(pt.X), h)),
				Y:     clampInt16(y.step(float64(pt.Y), h)),
				Z:     pt.Z,
				Color: pt.Color,
			})
		}
	}
	return out
}
