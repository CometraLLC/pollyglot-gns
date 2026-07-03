// Package cloze parses Anki-style cloze deletions: {{c1::hidden text}}.
// Pure functions, mirrored by frontend/src/lib/cloze.ts — keep in sync.
package cloze

import "regexp"

var marker = regexp.MustCompile(`\{\{c\d+::(.+?)\}\}`)

// Deletions returns the hidden texts of every well-formed marker, in order.
func Deletions(front string) []string {
	matches := marker.FindAllStringSubmatch(front, -1)
	if len(matches) == 0 {
		return nil
	}
	deletions := make([]string, 0, len(matches))
	for _, m := range matches {
		deletions = append(deletions, m[1])
	}
	return deletions
}

// Blank replaces each deletion with […] for the study-card front.
func Blank(front string) string {
	return marker.ReplaceAllString(front, "[…]")
}

// Reveal strips the markers, leaving the full text for the card back/flip.
func Reveal(front string) string {
	return marker.ReplaceAllString(front, "$1")
}
