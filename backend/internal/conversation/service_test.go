package conversation

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/base-go/backend/internal/shared/factory"
	"github.com/base-go/backend/internal/shared/models"
)

// --- fakes ---

type fakeRepo struct {
	conversations map[uuid.UUID]*models.Conversation
	messages      map[uuid.UUID][]models.ConversationMessage
	forceErr      error
}

func newFakeRepo() *fakeRepo {
	return &fakeRepo{
		conversations: make(map[uuid.UUID]*models.Conversation),
		messages:      make(map[uuid.UUID][]models.ConversationMessage),
	}
}

func (f *fakeRepo) CreateConversation(_ context.Context, c *models.Conversation) error {
	if f.forceErr != nil {
		return f.forceErr
	}
	c.ID = uuid.New()
	now := time.Now()
	c.CreatedAt, c.UpdatedAt = now, now
	f.conversations[c.ID] = c
	return nil
}

func (f *fakeRepo) GetConversationsByUser(_ context.Context, userID uuid.UUID) ([]models.Conversation, error) {
	if f.forceErr != nil {
		return nil, f.forceErr
	}
	var result []models.Conversation
	for _, c := range f.conversations {
		if c.UserID == userID && c.DeletedAt == nil {
			result = append(result, *c)
		}
	}
	return result, nil
}

func (f *fakeRepo) GetConversationByID(_ context.Context, id uuid.UUID) (*models.Conversation, error) {
	if f.forceErr != nil {
		return nil, f.forceErr
	}
	c, ok := f.conversations[id]
	if !ok || c.DeletedAt != nil {
		return nil, gorm.ErrRecordNotFound
	}
	return c, nil
}

func (f *fakeRepo) CreateMessage(_ context.Context, m *models.ConversationMessage) error {
	if f.forceErr != nil {
		return f.forceErr
	}
	m.ID = uuid.New()
	m.CreatedAt = time.Now()
	f.messages[m.ConversationID] = append(f.messages[m.ConversationID], *m)
	return nil
}

func (f *fakeRepo) GetMessagesByConversation(_ context.Context, conversationID uuid.UUID) ([]models.ConversationMessage, error) {
	if f.forceErr != nil {
		return nil, f.forceErr
	}
	return f.messages[conversationID], nil
}

type fakeTutor struct {
	greeting string
	reply    string
	err      error

	gotLanguage string
	gotHistory  []models.ConversationMessage
	gotMessage  string
}

func (f *fakeTutor) Greeting(language string) string {
	return f.greeting
}

func (f *fakeTutor) Reply(_ context.Context, language string, history []models.ConversationMessage, message string) (string, error) {
	f.gotLanguage, f.gotHistory, f.gotMessage = language, history, message
	return f.reply, f.err
}

func seedConversation(repo *fakeRepo, userID uuid.UUID) *models.Conversation {
	conv := factory.Conversation().WithUserID(userID).Build()
	repo.conversations[conv.ID] = &conv
	return &conv
}

// --- tests ---

func TestCreateConversation(t *testing.T) {
	userID := uuid.New()

	t.Run("creates the conversation with the tutor's greeting", func(t *testing.T) {
		repo := newFakeRepo()
		tutor := &fakeTutor{greeting: "Welcome! What do you know?"}
		svc := NewService(repo, tutor)

		resp, status, err := svc.CreateConversation(context.Background(), userID, CreateConversationRequest{
			Language: "Japanese",
		})

		require.NoError(t, err)
		assert.Equal(t, http.StatusCreated, status)
		assert.Equal(t, "Practice Japanese", resp.Title, "default title from language")
		assert.Equal(t, "Japanese", resp.Language)

		msgs := repo.messages[resp.ID]
		require.Len(t, msgs, 1, "greeting stored as the first message")
		assert.Equal(t, models.RoleTutor, msgs[0].Role)
		assert.Equal(t, "Welcome! What do you know?", msgs[0].Content)
	})

	t.Run("honors a custom title", func(t *testing.T) {
		repo := newFakeRepo()
		svc := NewService(repo, &fakeTutor{})

		resp, _, err := svc.CreateConversation(context.Background(), userID, CreateConversationRequest{
			Title: "Ordering food", Language: "Spanish",
		})

		require.NoError(t, err)
		assert.Equal(t, "Ordering food", resp.Title)
	})

	t.Run("requires a language", func(t *testing.T) {
		svc := NewService(newFakeRepo(), &fakeTutor{})

		_, status, err := svc.CreateConversation(context.Background(), userID, CreateConversationRequest{})

		require.Error(t, err)
		assert.Equal(t, http.StatusBadRequest, status)
	})

	t.Run("maps repo failure to 500", func(t *testing.T) {
		repo := newFakeRepo()
		repo.forceErr = errors.New("db down")
		svc := NewService(repo, &fakeTutor{})

		_, status, err := svc.CreateConversation(context.Background(), userID, CreateConversationRequest{
			Language: "Japanese",
		})

		require.Error(t, err)
		assert.Equal(t, http.StatusInternalServerError, status)
	})
}

func TestListConversations(t *testing.T) {
	userID := uuid.New()

	t.Run("returns only the user's conversations", func(t *testing.T) {
		repo := newFakeRepo()
		mine := seedConversation(repo, userID)
		seedConversation(repo, uuid.New())
		svc := NewService(repo, &fakeTutor{})

		resp, status, err := svc.ListConversations(context.Background(), userID)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, status)
		require.Len(t, resp, 1)
		assert.Equal(t, mine.ID, resp[0].ID)
	})

	t.Run("empty list serializes as [] not null", func(t *testing.T) {
		svc := NewService(newFakeRepo(), &fakeTutor{})

		resp, _, err := svc.ListConversations(context.Background(), userID)

		require.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Empty(t, resp)
	})
}

func TestGetMessages(t *testing.T) {
	userID := uuid.New()

	t.Run("returns the conversation's messages", func(t *testing.T) {
		repo := newFakeRepo()
		conv := seedConversation(repo, userID)
		repo.messages[conv.ID] = []models.ConversationMessage{
			factory.Message().WithConversationID(conv.ID).FromTutor().WithContent("hi").Build(),
			factory.Message().WithConversationID(conv.ID).WithContent("こんにちは").Build(),
		}
		svc := NewService(repo, &fakeTutor{})

		resp, status, err := svc.GetMessages(context.Background(), userID, conv.ID)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, status)
		require.Len(t, resp, 2)
		assert.Equal(t, models.RoleTutor, resp[0].Role)
		assert.Equal(t, "こんにちは", resp[1].Content)
	})

	t.Run("404 on another user's conversation", func(t *testing.T) {
		repo := newFakeRepo()
		conv := seedConversation(repo, uuid.New())
		svc := NewService(repo, &fakeTutor{})

		_, status, err := svc.GetMessages(context.Background(), userID, conv.ID)

		require.Error(t, err)
		assert.Equal(t, http.StatusNotFound, status)
	})

	t.Run("empty conversation serializes as [] not null", func(t *testing.T) {
		repo := newFakeRepo()
		conv := seedConversation(repo, userID)
		svc := NewService(repo, &fakeTutor{})

		resp, _, err := svc.GetMessages(context.Background(), userID, conv.ID)

		require.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Empty(t, resp)
	})
}

func TestSendMessage(t *testing.T) {
	userID := uuid.New()

	t.Run("stores the exchange and returns both messages", func(t *testing.T) {
		repo := newFakeRepo()
		conv := seedConversation(repo, userID)
		existing := factory.Message().WithConversationID(conv.ID).FromTutor().WithContent("greeting").Build()
		repo.messages[conv.ID] = []models.ConversationMessage{existing}
		tutor := &fakeTutor{reply: "What do you think it means?"}
		svc := NewService(repo, tutor)

		resp, status, err := svc.SendMessage(context.Background(), userID, conv.ID, SendMessageRequest{
			Content: "こんにちは",
		})

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, status)
		assert.Equal(t, models.RoleUser, resp.UserMessage.Role)
		assert.Equal(t, "こんにちは", resp.UserMessage.Content)
		assert.Equal(t, models.RoleTutor, resp.TutorMessage.Role)
		assert.Equal(t, "What do you think it means?", resp.TutorMessage.Content)

		// provider saw the language, prior history, and the new message
		assert.Equal(t, conv.Language, tutor.gotLanguage)
		require.Len(t, tutor.gotHistory, 1)
		assert.Equal(t, "greeting", tutor.gotHistory[0].Content)
		assert.Equal(t, "こんにちは", tutor.gotMessage)

		// both messages persisted after the greeting
		assert.Len(t, repo.messages[conv.ID], 3)
	})

	t.Run("provider failure persists nothing and maps to 502", func(t *testing.T) {
		repo := newFakeRepo()
		conv := seedConversation(repo, userID)
		svc := NewService(repo, &fakeTutor{err: errors.New("llm down")})

		_, status, err := svc.SendMessage(context.Background(), userID, conv.ID, SendMessageRequest{
			Content: "こんにちは",
		})

		require.Error(t, err)
		assert.Equal(t, http.StatusBadGateway, status)
		assert.Empty(t, repo.messages[conv.ID], "a failed exchange must not half-persist")
	})

	t.Run("validates the content", func(t *testing.T) {
		repo := newFakeRepo()
		conv := seedConversation(repo, userID)
		svc := NewService(repo, &fakeTutor{})

		_, status, err := svc.SendMessage(context.Background(), userID, conv.ID, SendMessageRequest{})

		require.Error(t, err)
		assert.Equal(t, http.StatusBadRequest, status)
	})

	t.Run("404 on another user's conversation", func(t *testing.T) {
		repo := newFakeRepo()
		conv := seedConversation(repo, uuid.New())
		svc := NewService(repo, &fakeTutor{reply: "?"})

		_, status, err := svc.SendMessage(context.Background(), userID, conv.ID, SendMessageRequest{
			Content: "hi",
		})

		require.Error(t, err)
		assert.Equal(t, http.StatusNotFound, status)
		assert.Empty(t, repo.messages[conv.ID])
	})
}
