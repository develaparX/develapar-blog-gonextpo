package repository_test

import (
	"database/sql"
	"errors"
	"regexp"
	"testing"

	"develapar-server/model"
	"develapar-server/repository"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type CategoryRepositoryTestSuite struct {
	suite.Suite
	db   *sql.DB
	mock sqlmock.Sqlmock
	repo repository.CategoryRepository
}

func (suite *CategoryRepositoryTestSuite) SetupTest() {
	var err error
	suite.db, suite.mock, err = sqlmock.New()
	assert.NoError(suite.T(), err)
	suite.repo = repository.NewCategoryRepository(suite.db)
}

func (suite *CategoryRepositoryTestSuite) TearDownTest() {
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

func (suite *CategoryRepositoryTestSuite) TestGetAll() {
	query := `SELECT id, name FROM categories`

	rows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(1, "Category 1").
		AddRow(2, "Category 2")

	suite.mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(rows)

	categories, err := suite.repo.GetAll()
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), categories, 2)
	assert.Equal(suite.T(), "Category 1", categories[0].Name)
}

func (suite *CategoryRepositoryTestSuite) TestGetAll_Error() {
	query := `SELECT id, name FROM categories`

	suite.mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnError(errors.New("db error"))

	categories, err := suite.repo.GetAll()
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), categories)
}

func (suite *CategoryRepositoryTestSuite) TestCreateCategory() {
	category := model.Category{Name: "New Category"}
	query := `INSERT INTO categories (name) VALUES($1) RETURNING id, name`

	rows := sqlmock.NewRows([]string{"id", "name"}).AddRow(1, category.Name)

	suite.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(category.Name).
		WillReturnRows(rows)

	createdCategory, err := suite.repo.CreateCategory(category)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), category.Name, createdCategory.Name)
}

func (suite *CategoryRepositoryTestSuite) TestCreateCategory_Error() {
	category := model.Category{Name: "New Category"}
	query := `INSERT INTO categories (name) VALUES($1) RETURNING id, name`

	suite.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(category.Name).
		WillReturnError(errors.New("db error"))

	createdCategory, err := suite.repo.CreateCategory(category)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), model.Category{}, createdCategory)
}

func (suite *CategoryRepositoryTestSuite) TestGetCategoryById() {
	query := `SELECT id, name FROM categories WHERE id = $1`

	rows := sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "Category 1")

	suite.mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(1).WillReturnRows(rows)

	category, err := suite.repo.GetCategoryById(1)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "Category 1", category.Name)
}

func (suite *CategoryRepositoryTestSuite) TestGetCategoryById_Error() {
	query := `SELECT id, name FROM categories WHERE id = $1`

	suite.mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(1).WillReturnError(errors.New("db error"))

	category, err := suite.repo.GetCategoryById(1)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), model.Category{}, category)
}

func (suite *CategoryRepositoryTestSuite) TestUpdateCategory() {
	category := model.Category{Id: 1, Name: "Updated Category"}
	query := `UPDATE categories SET name = $1 WHERE id = $2 RETURNING id, name`

	rows := sqlmock.NewRows([]string{"id", "name"}).AddRow(category.Id, category.Name)

	suite.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(category.Name, category.Id).
		WillReturnRows(rows)

	updatedCategory, err := suite.repo.UpdateCategory(category)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), category.Name, updatedCategory.Name)
}

func (suite *CategoryRepositoryTestSuite) TestUpdateCategory_Error() {
	category := model.Category{Id: 1, Name: "Updated Category"}
	query := `UPDATE categories SET name = $1 WHERE id = $2 RETURNING id, name`

	suite.mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(category.Name, category.Id).
		WillReturnError(errors.New("db error"))

	updatedCategory, err := suite.repo.UpdateCategory(category)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), model.Category{}, updatedCategory)
}

func (suite *CategoryRepositoryTestSuite) TestDeleteCategory() {
	query := `DELETE FROM categories WHERE id = $1`

	suite.mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))

	err := suite.repo.DeleteCategory(1)
	assert.NoError(suite.T(), err)
}

func (suite *CategoryRepositoryTestSuite) TestDeleteCategory_Error() {
	query := `DELETE FROM categories WHERE id = $1`

	suite.mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(1).WillReturnError(errors.New("db error"))

	err := suite.repo.DeleteCategory(1)
	assert.Error(suite.T(), err)
}

func TestCategoryRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(CategoryRepositoryTestSuite))
}
