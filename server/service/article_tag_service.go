package service

import (
	"develapar-server/model"
	"develapar-server/repository"
)

type ArticleTagService interface {
	AssignTags(articleId int, tagId []int) error
	AsignTagsByName(articleId int, tagNames []string) error
	FindTagByArticleId(articleId int) ([]model.Tags, error)
	FindArticleByTagId(tagId int) ([]model.Article, error)
}

type articleTagService struct {
	articleTagRepo repository.ArticleTagRepository
	tagRepo        repository.TagRepository
}

// AsignTagsByName implements ArticleTagService.
func (a *articleTagService) AsignTagsByName(articleId int, tagNames []string) error {

	var tagIds []int

	for _, tagName := range tagNames {
		tag, err := a.tagRepo.GetTagByName(tagName)
		if err != nil {
			newTag, err := a.tagRepo.CreateTag(model.Tags{Name: tagName})
			if err != nil {
				return err
			}
			tagIds = append(tagIds, newTag.Id)
		} else {
			tagIds = append(tagIds, tag.Id)
		}
	}

	return a.articleTagRepo.AssignTags(articleId, tagIds)
}

// AssignTags implements ArticleTagService.
func (a *articleTagService) AssignTags(articleId int, tagId []int) error {
	return a.articleTagRepo.AssignTags(articleId, tagId)
}

// FindArticleByTagId implements ArticleTagService.
func (a *articleTagService) FindArticleByTagId(tagId int) ([]model.Article, error) {
	return a.articleTagRepo.GetArticleByTagId(tagId)
}

// FindTagByArticleId implements ArticleTagService.
func (a *articleTagService) FindTagByArticleId(articleId int) ([]model.Tags, error) {
	return a.articleTagRepo.GetTagsByArticleId(articleId)
}

func NewArticleTagService(tagRepo repository.TagRepository, articleTagRepo repository.ArticleTagRepository) ArticleTagService {
	return &articleTagService{
		tagRepo:        tagRepo,
		articleTagRepo: articleTagRepo,
	}
}
