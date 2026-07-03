package conversation

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/base-go/backend/internal/shared/factory"
	"github.com/base-go/backend/internal/shared/models"
)

func TestSocraticTutorGreeting(t *testing.T) {
	tutor := NewSocraticTutor()

	greeting := tutor.Greeting("Japanese")

	assert.Contains(t, greeting, "Japanese")
	assert.True(t, strings.HasSuffix(strings.TrimSpace(greeting), "?"),
		"the Socratic tutor always ends with a question")
}

func TestSocraticTutorReplyIsDeterministic(t *testing.T) {
	tutor := NewSocraticTutor()
	history := []models.ConversationMessage{
		factory.Message().FromTutor().WithContent("greeting").Build(),
	}

	a, err := tutor.Reply(context.Background(), "Japanese", history, "こんにちは")
	require.NoError(t, err)
	b, err := tutor.Reply(context.Background(), "Japanese", history, "こんにちは")
	require.NoError(t, err)

	assert.Equal(t, a, b, "same history + message must produce the same reply")
}

func TestSocraticTutorAlwaysEndsWithQuestion(t *testing.T) {
	// The persona rule from the original repo's socratic-tutor.mdc:
	// "Always end your response with a question".
	tutor := NewSocraticTutor()

	var history []models.ConversationMessage
	for i, msg := range []string{"こんにちは", "ねこ", "みず", "いぬ", "ありがとう", "はい"} {
		reply, err := tutor.Reply(context.Background(), "Japanese", history, msg)
		require.NoError(t, err)
		assert.True(t, strings.HasSuffix(strings.TrimSpace(reply), "?"),
			"reply %d must end with a question, got %q", i, reply)

		history = append(history,
			factory.Message().WithContent(msg).Build(),
			factory.Message().FromTutor().WithContent(reply).Build(),
		)
	}
}

func TestSocraticTutorRepliesVaryAcrossTurns(t *testing.T) {
	tutor := NewSocraticTutor()

	var history []models.ConversationMessage
	seen := map[string]bool{}
	distinct := 0
	for range 4 {
		reply, err := tutor.Reply(context.Background(), "Spanish", history, "hola")
		require.NoError(t, err)
		if !seen[reply] {
			distinct++
		}
		seen[reply] = true
		history = append(history,
			factory.Message().WithContent("hola").Build(),
			factory.Message().FromTutor().WithContent(reply).Build(),
		)
	}

	assert.GreaterOrEqual(t, distinct, 3, "consecutive turns should not repeat the same script")
}

func TestSocraticTutorEchoesTheLearnersWords(t *testing.T) {
	tutor := NewSocraticTutor()

	reply, err := tutor.Reply(context.Background(), "Japanese", nil, "生き甲斐")
	require.NoError(t, err)

	assert.Contains(t, reply, "生き甲斐", "the first probe quotes what the learner said")
}
