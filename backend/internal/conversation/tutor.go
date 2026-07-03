package conversation

import (
	"context"
	"fmt"

	"github.com/base-go/backend/internal/shared/models"
)

// TutorProvider is the pluggable conversation partner (see D-007). The
// scripted Socratic tutor is the deterministic default; an LLM-backed
// provider plugs in behind this same interface.
type TutorProvider interface {
	// Greeting opens a new conversation in the given language.
	Greeting(language string) string
	// Reply answers the learner's message given the conversation so far.
	Reply(ctx context.Context, language string, history []models.ConversationMessage, message string) (string, error)
}

// SocraticTutor is a deterministic, scripted implementation of the
// Socratic persona from the original Pollyglot repo: it never answers
// directly, guides through questions, and always ends with a question.
type SocraticTutor struct{}

func NewSocraticTutor() *SocraticTutor {
	return &SocraticTutor{}
}

func (t *SocraticTutor) Greeting(language string) string {
	return fmt.Sprintf(
		"Hello! I'm your %s practice partner. I won't hand you answers — I'll ask questions until you find them yourself. To start: what is a %s word or phrase you already know?",
		language, language,
	)
}

// probes cycle by turn; {0} is the learner's message, {1} the language.
var probes = []string{
	"「%[1]s」— interesting. What do you think %[1]s means, in your own words?",
	"Good. How would you use %[1]s in a full %[2]s sentence?",
	"Nice try. What other %[2]s words do you know that relate to %[1]s?",
	"Let's flip it: if a friend used %[1]s, how would you reply in %[2]s?",
	"Almost there. Can you recall a situation today where %[1]s would fit — what would you say?",
}

func (t *SocraticTutor) Reply(_ context.Context, language string, history []models.ConversationMessage, message string) (string, error) {
	userTurns := 0
	for _, m := range history {
		if m.Role == models.RoleUser {
			userTurns++
		}
	}

	probe := probes[userTurns%len(probes)]
	return fmt.Sprintf(probe, message, language), nil
}
