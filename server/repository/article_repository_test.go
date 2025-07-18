package repository

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"testing"
	"time"

	"develapar-server/model"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ArticleRepositoryTestSuite struct {
	suite.Suite
	db   *sql.DB
	mock sqlmock.Sqlmock
	repo ArticleRepository
}

func (suite *ArticleRepositoryTestSuite) SetupTest() {
	var err error
	suite.db, suite.mock, err = sqlmock.New()
	assert.NoError(suite.T(), err)
	suite.repo = NewArticleRepository(suite.db)
}

func (suite *ArticleRepositoryTestSuite) TearDownTest() {
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

func (suite *ArticleRepositoryTestSuite) TestGetAll() {
	query := `SELECT a.id, a.title, a.slug, a.content, a.views, a.created_at, a.updated_at, u.id, u.name, u.email, c.id, c.name FROM articles a JOIN users u ON a.user_id = u.id JOIN categories c ON a.category_id = c.id ORDER BY a.created_at DESC;`

	rows := sqlmock.NewRows([]string{"id", "title", "slug", "content", "views", "created_at", "updated_at", "user_id", "user_name", "user_email", "category_id", "category_name"}).
		AddRow(1, "Title 1", "slug-1", "Content 1", 10, time.Now(), time.Now(), 1, "User 1", "user1@example.com", 1, "Category 1").
		AddRow(2, "Title 2", "slug-2", "Content 2", 20, time.Now(), time.Now(), 2, "User 2", "user2@example.com", 2, "Category 2")

	suite.mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(rows)

	articles, err := suite.repo.GetAll(context.Background())
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), articles, 2)
	assert.Equal(suite.T(), "Title 1", articles[0].Title)
	assert.Equal(suite.T(), "User 1", articles[0].User.Name)
}

func (suite *ArticleRepositoryTestSuite) TestGetAll_Error() {
	query := `SELECT a.id, a.title, a.slug, a.content, a.views, a.created_at, a.updated_at, u.id, u.name, u.email, c.id, c.name FROM articles a JOIN users u ON a.user_id = u.id JOIN categories c ON a.category_id = c.id ORDER BY a.created_at DESC;`

	suite.mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnError(errors.New("db error"))

	articles, err := suite.repo.GetAll(context.Background())
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), articles)
}

func (suite *ArticleRepositoryTestSuite) TestCreateArticle() {
	article := model.Article{
		Title:    "New Article",
		Slug:     "new-article",
		Content:  "New Content",
		User:     model.User{Id: 1},
		Category: model.Category{Id: 1},
	}

	query := `INSERT INTO articles (title, content, slug, user_id, category_id) VALUES ($1, $2, $3, $4, $5) RETURNING id, title, slug, content, user_id, category_id, views, created_at, updated_at`

	rows := sqlmock.NewRows([]string{"id", "title", "slug", "content", "user_id", "category_id", "views", "created_at", "updated_at"}).
		AddRow(1, article.Title, article.Slug, article.Content, article.User.Id, article.Category.Id, 0, time.Now(), time.Now())

	suite.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(article.Title, article.Content, article.Slug, article.User.Id, article.Category.Id).
		WillReturnRows(rows)

	createdArticle, err := suite.repo.CreateArticle(context.Background(), article)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), article.Title, createdArticle.Title)
}

func (suite *ArticleRepositoryTestSuite) TestCreateArticle_Error() {
	article := model.Article{
		Title:    "New Article",
		Slug:     "new-article",
		Content:  "New Content",
		User:     model.User{Id: 1},
		Category: model.Category{Id: 1},
	}

	query := `INSERT INTO articles (title, content, slug, user_id, category_id) VALUES ($1, $2, $3, $4, $5) RETURNING id, title, slug, content, user_id, category_id, views, created_at, updated_at`

	suite.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(article.Title, article.Content, article.Slug, article.User.Id, article.Category.Id).
		WillReturnError(errors.New("db error"))

	createdArticle, err := suite.repo.CreateArticle(context.Background(), article)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), model.Article{}, createdArticle)
}

func (suite *ArticleRepositoryTestSuite) TestUpdateArticle() {
	article := model.Article{
		Id:       1,
		Title:    "Updated Title",
		Slug:     "updated-slug",
		Content:  "Updated Content",
		Category: model.Category{Id: 2},
	}

	query := `UPDATE articles SET title = $1, slug = $2, content = $3, category_id=$4, updated_at = NOW() WHERE id = $5 RETURNING id, title, slug, content, category_id, created_at, updated_at`

	rows := sqlmock.NewRows([]string{"id", "title", "slug", "content", "category_id", "created_at", "updated_at"}).
		AddRow(article.Id, article.Title, article.Slug, article.Content, article.Category.Id, time.Now(), time.Now())

	suite.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(article.Title, article.Slug, article.Content, article.Category.Id, article.Id).
		WillReturnRows(rows)

	updatedArticle, err := suite.repo.UpdateArticle(context.Background(), article)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), article.Title, updatedArticle.Title)
}

func (suite *ArticleRepositoryTestSuite) TestUpdateArticle_Error() {
	article := model.Article{
		Id:       1,
		Title:    "Updated Title",
		Slug:     "updated-slug",
		Content:  "Updated Content",
		Category: model.Category{Id: 2},
	}

	query := `UPDATE articles SET title = $1, slug = $2, content = $3, category_id=$4, updated_at = NOW() WHERE id = $5 RETURNING id, title, slug, content, category_id, created_at, updated_at`

	suite.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(article.Title, article.Slug, article.Content, article.Category.Id, article.Id).
		WillReturnError(errors.New("db error"))

	updatedArticle, err := suite.repo.UpdateArticle(context.Background(), article)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), model.Article{}, updatedArticle)
}

func (suite *ArticleRepositoryTestSuite) TestGetArticleById() {
	query := `SELECT id, title, slug, content, views, user_id, category_id, created_at, updated_at FROM articles WHERE id = $1`

	rows := sqlmock.NewRows([]string{"id", "title", "slug", "content", "views", "user_id", "category_id", "created_at", "updated_at"}).
		AddRow(1, "Title 1", "slug-1", "Content 1", 10, 1, 1, time.Now(), time.Now())

	suite.mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(1).WillReturnRows(rows)

	article, err := suite.repo.GetArticleById(context.Background(), 1)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "Title 1", article.Title)
}

func (suite *ArticleRepositoryTestSuite) TestGetArticleById_Error() {
	query := `SELECT id, title, slug, content, views, user_id, category_id, created_at, updated_at FROM articles WHERE id = $1`

	suite.mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(1).WillReturnError(errors.New("db error"))

	article, err := suite.repo.GetArticleById(context.Background(), 1)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), model.Article{}, article)
}

func (suite *ArticleRepositoryTestSuite) TestGetArticleByUserId() {
	query := `SELECT a.id, a.title, a.slug, a.content, a.views, a.created_at, a.updated_at, u.id, u.name, u.email, c.id, c.name FROM articles a JOIN users u ON a.user_id = u.id JOIN categories c ON a.category_id = c.id WHERE a.user_id = $1 ORDER BY a.created_at DESC;`

	rows := sqlmock.NewRows([]string{"id", "title", "slug", "content", "views", "created_at", "updated_at", "user_id", "user_name", "user_email", "category_id", "category_name"}).
		AddRow(1, "Title 1", "slug-1", "Content 1", 10, time.Now(), time.Now(), 1, "User 1", "user1@example.com", 1, "Category 1")

	suite.mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(1).WillReturnRows(rows)

	articles, err := suite.repo.GetArticleByUserId(context.Background(), 1)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), articles, 1)
	assert.Equal(suite.T(), "Title 1", articles[0].Title)
}

func (suite *ArticleRepositoryTestSuite) TestGetArticleByUserId_Error() {
	query := `SELECT a.id, a.title, a.slug, a.content, a.views, a.created_at, a.updated_at, u.id, u.name, u.email, c.id, c.name FROM articles a JOIN users u ON a.user_id = u.id JOIN categories c ON a.category_id = c.id WHERE a.user_id = $1 ORDER BY a.created_at DESC;`

	suite.mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(1).WillReturnError(errors.New("db error"))

	articles, err := suite.repo.GetArticleByUserId(context.Background(), 1)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), articles)
}

func (suite *ArticleRepositoryTestSuite) TestGetArticleBySlug() {
	query := `SELECT a.id, a.title, a.slug, a.content, a.views, a.created_at, a.updated_at, u.id, u.name, u.email, c.id, c.name FROM articles a JOIN users u ON a.user_id = u.id JOIN categories c ON a.category_id = c.id WHERE a.slug = $1;`

	rows := sqlmock.NewRows([]string{"id", "title", "slug", "content", "views", "created_at", "updated_at", "user_id", "user_name", "user_email", "category_id", "category_name"}).
		AddRow(1, "Title 1", "slug-1", "Content 1", 10, time.Now(), time.Now(), 1, "User 1", "user1@example.com", 1, "Category 1")

	suite.mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs("slug-1").WillReturnRows(rows)

	article, err := suite.repo.GetArticleBySlug(context.Background(), "slug-1")
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "Title 1", article.Title)
}

func (suite *ArticleRepositoryTestSuite) TestGetArticleBySlug_Error() {
	query := `SELECT a.id, a.title, a.slug, a.content, a.views, a.created_at, a.updated_at, u.id, u.name, u.email, c.id, c.name FROM articles a JOIN users u ON a.user_id = u.id JOIN categories c ON a.category_id = c.id WHERE a.slug = $1;`

	suite.mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs("slug-1").WillReturnError(errors.New("db error"))

	article, err := suite.repo.GetArticleBySlug(context.Background(), "slug-1")
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), model.Article{}, article)
}

func (suite *ArticleRepositoryTestSuite) TestGetArticleByCategory() {
	query := `SELECT a.id, a.title, a.slug, a.content, a.views, a.created_at, a.updated_at, u.id, u.name, u.email, c.id, c.name FROM articles a JOIN users u ON a.user_id = u.id JOIN categories c ON a.category_id = c.id WHERE c.name = $1;`

	rows := sqlmock.NewRows([]string{"id", "title", "slug", "content", "views", "created_at", "updated_at", "user_id", "user_name", "user_email", "category_id", "category_name"}).
		AddRow(1, "Title 1", "slug-1", "Content 1", 10, time.Now(), time.Now(), 1, "User 1", "user1@example.com", 1, "Category 1")

	suite.mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs("Category 1").WillReturnRows(rows)

	articles, err := suite.repo.GetArticleByCategory(context.Background(), "Category 1")
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), articles, 1)
	assert.Equal(suite.T(), "Title 1", articles[0].Title)
}

func (suite *ArticleRepositoryTestSuite) TestGetArticleByCategory_Error() {
	query := `SELECT a.id, a.title, a.slug, a.content, a.views, a.created_at, a.updated_at, u.id, u.name, u.email, c.id, c.name FROM articles a JOIN users u ON a.user_id = u.id JOIN categories c ON a.category_id = c.id WHERE c.name = $1;`

	suite.mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs("Category 1").WillReturnError(errors.New("db error"))

	articles, err := suite.repo.GetArticleByCategory(context.Background(), "Category 1")
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), articles)
}

func (suite *ArticleRepositoryTestSuite) TestDeleteArticle() {
	query := `DELETE FROM articles WHERE id = $1`

	suite.mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))

	err := suite.repo.DeleteArticle(context.Background(), 1)
	assert.NoError(suite.T(), err)
}

func (suite *ArticleRepositoryTestSuite) TestDeleteArticle_Error() {
	query := `DELETE FROM articles WHERE id = $1`

	suite.mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(1).WillReturnError(errors.New("db error"))

	err := suite.repo.DeleteArticle(context.Background(), 1)
	assert.Error(suite.T(), err)
}

func TestArticleRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(ArticleRepositoryTestSuite))
}
