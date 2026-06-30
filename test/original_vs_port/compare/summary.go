package main

type summary struct {
	files      int
	decisions  int
	matches    int
	mismatches int
	errors     int
}

func (s *summary) add(other summary) {
	s.decisions += other.decisions
	s.matches += other.matches
	s.mismatches += other.mismatches
	s.errors += other.errors
}
