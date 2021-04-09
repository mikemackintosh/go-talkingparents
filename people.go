package tp

import (
  "fmt"
)

var (
  // instance variables
  people = &People{}
)

// People is a collection of persons in the thread.
// This will always have 2 count.
type People []string

// AddPersonIfNotExists will always populate, but skip duplicates.
func (p *People) AddPersonIfNotExists(person string) {
	for _, name := range *p {
		if name == person {
			return
		}
	}

	*p = append(*p, person)
}

// Not will return the other person in the people slice, otherwise throw an error.
func (p *People) Not(person string) (string, error) {
	for _, name := range *p {
		if name != person {
			return name, nil
		}
	}

	return "", fmt.Errorf("could not find an alternate to '%s'", person)
}
