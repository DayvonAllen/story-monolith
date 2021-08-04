package services

import (
	"story-app-monolith/domain"
	"story-app-monolith/repo"
)

type ConversationService interface {
	FindConversation(owner, to string) (*domain.Conversation, error)
	GetConversationPreviews(owner string) (*[]domain.ConversationPreview, error)
}

type DefaultConversationService struct {
	repo repo.ConversationRepo
}

func (c DefaultConversationService) FindConversation(owner, to string) (*domain.Conversation, error) {
	conversation, err := c.repo.FindConversation(owner, to)
	if err != nil {
		return nil, err
	}
	return conversation, nil
}

func (c DefaultConversationService) GetConversationPreviews(owner string) (*[]domain.ConversationPreview, error) {
	conversation, err := c.repo.GetConversationPreviews(owner)
	if err != nil {
		return nil, err
	}
	return conversation, nil
}

func NewConversationService(repository repo.ConversationRepo) DefaultConversationService {
	return DefaultConversationService{repository}
}