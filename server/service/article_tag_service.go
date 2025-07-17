package service

import (
	"context"
	"develapar-server/model"
	"develapar-server/repository"
)

type ArticleTagService interface {
	AssignTags(articleId int, tagId []int) error
	AsignTagsByName(articleId int, tagNames []string) error
	FindTagByArticleId(articleId int) ([]model.Tags, error)
	FindArticleByTagId(tagId int) ([]model.Article, error)
	RemoveTagFromArticle(articleId, tagId int) error
}

type articleTagService struct {
	articleTagRepo repository.ArticleTagRepository
	tagRepo        repository.TagRepository
}

// RemoveTagFromArticle implements ArticleTagService.
func (a *articleTagService) RemoveTagFromArticle(articleId int, tagId int) error {
	return a.articleTagRepo.RemoveTagFromArticle(context.Background(), articleId, tagId)
}

// AsignTagsByName implements ArticleTagService.
func (a *articleTagService) AsignTagsByName(articleId int, tagNames []string) error {

	var tagIds []int

	for _, tagName := range tagNames {
		tag, err := a.tagRepo.GetTagByName(context.Background(), tagName)
		if err != nil {
			newTag, err := a.tagRepo.CreateTag(context.Background(), model.Tags{Name: tagName})
			if err != nil {
				return err
			}
			tagIds = append(tagIds, newTag.Id)
		} else {
			tagIds = append(tagIds, tag.Id)
		}
	}

	return a.articleTagRepo.AssignTags(context.Background(), articleId, tagIds)
}

// AssignTags implements ArticleTagService.
func (a *articleTagService) AssignTags(articleId int, tagId []int) error {
	return a.articleTagRepo.AssignTags(context.Background(), articleId, tagId)
}

// FindArticleByTagId implements ArticleTagService.
func (a *articleTagService) FindArticleByTagId(tagId int) ([]model.Article, error) {
	return a.articleTagRepo.GetArticleByTagId(context.Background(), tagId)
}

// FindTagByArticleId implements ArticleTagService.
func (a *articleTagService) FindTagByArticleId(articleId int) ([]model.Tags, error) {
	return a.articleTagRepo.GetTagsByArticleId(context.Background(), articleId)
}

func NewArticleTagService(tagRepo repository.TagRepository, articleTagRepo repository.ArticleTagRepository) ArticleTagService {
	return &articleTagService{
		tagRepo:        tagRepo,
		articleTagRepo: articleTagRepo,
	}
}
