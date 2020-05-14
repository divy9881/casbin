package effect

// DefaultEffectorStream is the default implementation of interface EffectorStream.
type DefaultEffectorStream struct {
	done bool
	res bool
	expr string
	idx uint
	cap uint
	expl []uint
}

func (s *DefaultEffectorStream) Current() bool {
	if !s.done {
		panic()
	}
	return s.res
}

func (s *DefaultEffectorStream) Explain() []uint {
	if !s.done {
		panic()
	}
	return s.expl
}

func (s *DefaultEffectorStream) PushEffect(eft Effect) bool {
	hasPolicy := s.cap > 1

	if s.expr == "some(where (p_eft == allow))" {
		if eft == Allow {
			s.done = true
			s.res = true

			if hasPolicy {
				s.expl = append(s.expl, s.idx)
			}
		}
	} else if s.expr == "some(where (p_eft == allow)) && !some(where (p_eft == deny))" {
		if eft == Allow {
			s.res = true

			if hasPolicy {
				s.expl = append(s.expl, s.idx)
			}
		} else if eft == Deny {
			s.done = true
			s.res = false

			if hasPolicy {
				s.expl = append(s.expl, s.idx)
			}
		}
	} else if s.expr == "!some(where (p_eft == deny))" {
		if eft == Deny {
			s.done = true
			s.res = false

			if hasPolicy {
				s.expl = append(s.expl, s.idx)
			}
		}
	} else if s.expr == "priority(p_eft) || deny" && eft != Indeterminate {
		if eft == Allow {
			s.res = true
		} else {
			s.res = false
		}
		s.done = true

		if hasPolicy {
			s.expl = append(s.expl, s.idx)
		}
	}

	if s.idx + 1 == s.cap {
		s.done = true
	}
	s.idx++
	
	return s.done
}
