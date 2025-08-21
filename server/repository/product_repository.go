package repository

import (
	"context"
	"database/sql"
	"develapar-server/model"
	"develapar-server/utils"
	"time"

	"github.com/google/uuid"
)

type ProductRepository interface {
	// Product Categories
	CreateProductCategory(ctx context.Context, payload model.ProductCategory) (model.ProductCategory, error)
	GetAllProductCategories(ctx context.Context) ([]model.ProductCategory, error)
	GetProductCategoryById(ctx context.Context, id uuid.UUID) (model.ProductCategory, error)
	GetProductCategoryBySlug(ctx context.Context, slug string) (model.ProductCategory, error)
	UpdateProductCategory(ctx context.Context, payload model.ProductCategory) (model.ProductCategory, error)
	DeleteProductCategory(ctx context.Context, id uuid.UUID) error

	// Products
	CreateProduct(ctx context.Context, payload model.Product) (model.Product, error)
	GetAllProducts(ctx context.Context) ([]model.Product, error)
	GetAllProductsWithPagination(ctx context.Context, offset, limit int) ([]model.Product, int, error)
	GetProductById(ctx context.Context, id uuid.UUID) (model.Product, error)
	GetProductsByCategory(ctx context.Context, categoryId uuid.UUID) ([]model.Product, error)
	UpdateProduct(ctx context.Context, payload model.Product) (model.Product, error)
	DeleteProduct(ctx context.Context, id uuid.UUID) error

	// Product Affiliate Links
	CreateProductAffiliateLink(ctx context.Context, payload model.ProductAffiliateLink) (model.ProductAffiliateLink, error)
	GetAffiliateLinksbyProductId(ctx context.Context, productId uuid.UUID) ([]model.ProductAffiliateLink, error)
	UpdateProductAffiliateLink(ctx context.Context, payload model.ProductAffiliateLink) (model.ProductAffiliateLink, error)
	DeleteProductAffiliateLink(ctx context.Context, id uuid.UUID) error

	// Article Product Relations
	AddProductToArticle(ctx context.Context, articleId, productId uuid.UUID) error
	RemoveProductFromArticle(ctx context.Context, articleId, productId uuid.UUID) error
	GetProductsByArticleId(ctx context.Context, articleId uuid.UUID) ([]model.Product, error)
	GetArticlesByProductId(ctx context.Context, productId uuid.UUID) ([]model.Article, error)
}

type productRepository struct {
	db *sql.DB
}

// CreateProductCategory implements ProductRepository
func (r *productRepository) CreateProductCategory(ctx context.Context, payload model.ProductCategory) (model.ProductCategory, error) {
	newId := uuid.Must(uuid.NewV7())
	slug := utils.GenerateSlug(payload.Name)
	var category model.ProductCategory

	query := `INSERT INTO product_categories (id, name, slug, description, created_at, updated_at) 
			  VALUES ($1, $2, $3, $4, $5, $6) 
			  RETURNING id, name, slug, description, created_at, updated_at`

	err := r.db.QueryRowContext(ctx, query, newId, payload.Name, slug, payload.Description, time.Now(), time.Now()).
		Scan(&category.Id, &category.Name, &category.Slug, &category.Description, &category.CreatedAt, &category.UpdatedAt)

	if err != nil {
		if ctx.Err() != nil {
			return model.ProductCategory{}, ctx.Err()
		}
		return model.ProductCategory{}, err
	}

	return category, nil
}

// GetAllProductCategories implements ProductRepository
func (r *productRepository) GetAllProductCategories(ctx context.Context) ([]model.ProductCategory, error) {
	query := `SELECT id, name, slug, description, created_at, updated_at FROM product_categories ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		return nil, err
	}
	defer rows.Close()

	var categories []model.ProductCategory
	for rows.Next() {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		var category model.ProductCategory
		err := rows.Scan(&category.Id, &category.Name, &category.Slug, &category.Description, &category.CreatedAt, &category.UpdatedAt)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

// GetProductCategoryById implements ProductRepository
func (r *productRepository) GetProductCategoryById(ctx context.Context, id uuid.UUID) (model.ProductCategory, error) {
	var category model.ProductCategory
	query := `SELECT id, name, slug, description, created_at, updated_at FROM product_categories WHERE id = $1`

	err := r.db.QueryRowContext(ctx, query, id).
		Scan(&category.Id, &category.Name, &category.Slug, &category.Description, &category.CreatedAt, &category.UpdatedAt)

	if err != nil {
		if ctx.Err() != nil {
			return model.ProductCategory{}, ctx.Err()
		}
		return model.ProductCategory{}, err
	}

	return category, nil
}

// GetProductCategoryBySlug implements ProductRepository
func (r *productRepository) GetProductCategoryBySlug(ctx context.Context, slug string) (model.ProductCategory, error) {
	var category model.ProductCategory
	query := `SELECT id, name, slug, description, created_at, updated_at FROM product_categories WHERE slug = $1`

	err := r.db.QueryRowContext(ctx, query, slug).
		Scan(&category.Id, &category.Name, &category.Slug, &category.Description, &category.CreatedAt, &category.UpdatedAt)

	if err != nil {
		if ctx.Err() != nil {
			return model.ProductCategory{}, ctx.Err()
		}
		return model.ProductCategory{}, err
	}

	return category, nil
}

// UpdateProductCategory implements ProductRepository
func (r *productRepository) UpdateProductCategory(ctx context.Context, payload model.ProductCategory) (model.ProductCategory, error) {
	var category model.ProductCategory
	query := `UPDATE product_categories SET name = $1, slug = $2, description = $3, updated_at = $4 
			  WHERE id = $5 
			  RETURNING id, name, slug, description, created_at, updated_at`

	err := r.db.QueryRowContext(ctx, query, payload.Name, payload.Slug, payload.Description, time.Now(), payload.Id).
		Scan(&category.Id, &category.Name, &category.Slug, &category.Description, &category.CreatedAt, &category.UpdatedAt)

	if err != nil {
		if ctx.Err() != nil {
			return model.ProductCategory{}, ctx.Err()
		}
		return model.ProductCategory{}, err
	}

	return category, nil
}

// DeleteProductCategory implements ProductRepository
func (r *productRepository) DeleteProductCategory(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM product_categories WHERE id = $1`, id)
	if err != nil {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		return err
	}
	return nil
}

func NewProductRepository(database *sql.DB) ProductRepository {
	return &productRepository{db: database}
} // Cr
func (r *productRepository) CreateProduct(ctx context.Context, payload model.Product) (model.Product, error) {
	newId := uuid.Must(uuid.NewV7())
	var product model.Product

	query := `INSERT INTO products (id, product_category_id, name, description, image_url, is_active, created_at, updated_at) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8) 
			  RETURNING id, product_category_id, name, description, image_url, is_active, created_at, updated_at`

	err := r.db.QueryRowContext(ctx, query, newId, payload.ProductCategoryId, payload.Name, payload.Description,
		payload.ImageUrl, payload.IsActive, time.Now(), time.Now()).
		Scan(&product.Id, &product.ProductCategoryId, &product.Name, &product.Description,
			&product.ImageUrl, &product.IsActive, &product.CreatedAt, &product.UpdatedAt)

	if err != nil {
		if ctx.Err() != nil {
			return model.Product{}, ctx.Err()
		}
		return model.Product{}, err
	}

	return product, nil
}

// GetAllProducts implements ProductRepository
func (r *productRepository) GetAllProducts(ctx context.Context) ([]model.Product, error) {
	query := `
		SELECT 
			p.id, p.product_category_id, p.name, p.description, p.image_url, p.is_active, p.created_at, p.updated_at,
			pc.id, pc.name, pc.slug, pc.description
		FROM products p
		LEFT JOIN product_categories pc ON p.product_category_id = pc.id
		ORDER BY p.created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		return nil, err
	}
	defer rows.Close()

	var products []model.Product
	for rows.Next() {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		var product model.Product
		var category model.ProductCategory
		var categoryId sql.NullString
		var categoryName sql.NullString
		var categorySlug sql.NullString
		var categoryDesc sql.NullString

		err := rows.Scan(
			&product.Id, &product.ProductCategoryId, &product.Name, &product.Description,
			&product.ImageUrl, &product.IsActive, &product.CreatedAt, &product.UpdatedAt,
			&categoryId, &categoryName, &categorySlug, &categoryDesc,
		)
		if err != nil {
			return nil, err
		}

		if categoryId.Valid {
			category.Id = uuid.MustParse(categoryId.String)
			category.Name = categoryName.String
			category.Slug = categorySlug.String
			if categoryDesc.Valid {
				category.Description = &categoryDesc.String
			}
			product.ProductCategory = &category
		}

		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

// GetAllProductsWithPagination implements ProductRepository
func (r *productRepository) GetAllProductsWithPagination(ctx context.Context, offset, limit int) ([]model.Product, int, error) {
	// Get total count
	var totalCount int
	countQuery := `SELECT COUNT(*) FROM products`
	err := r.db.QueryRowContext(ctx, countQuery).Scan(&totalCount)
	if err != nil {
		if ctx.Err() != nil {
			return nil, 0, ctx.Err()
		}
		return nil, 0, err
	}

	// Get paginated results
	query := `
		SELECT 
			p.id, p.product_category_id, p.name, p.description, p.image_url, p.is_active, p.created_at, p.updated_at,
			pc.id, pc.name, pc.slug, pc.description
		FROM products p
		LEFT JOIN product_categories pc ON p.product_category_id = pc.id
		ORDER BY p.created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		if ctx.Err() != nil {
			return nil, 0, ctx.Err()
		}
		return nil, 0, err
	}
	defer rows.Close()

	var products []model.Product
	for rows.Next() {
		select {
		case <-ctx.Done():
			return nil, 0, ctx.Err()
		default:
		}

		var product model.Product
		var category model.ProductCategory
		var categoryId sql.NullString
		var categoryName sql.NullString
		var categorySlug sql.NullString
		var categoryDesc sql.NullString

		err := rows.Scan(
			&product.Id, &product.ProductCategoryId, &product.Name, &product.Description,
			&product.ImageUrl, &product.IsActive, &product.CreatedAt, &product.UpdatedAt,
			&categoryId, &categoryName, &categorySlug, &categoryDesc,
		)
		if err != nil {
			return nil, 0, err
		}

		if categoryId.Valid {
			category.Id = uuid.MustParse(categoryId.String)
			category.Name = categoryName.String
			category.Slug = categorySlug.String
			if categoryDesc.Valid {
				category.Description = &categoryDesc.String
			}
			product.ProductCategory = &category
		}

		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return products, totalCount, nil
}

// GetProductById implements ProductRepository
func (r *productRepository) GetProductById(ctx context.Context, id uuid.UUID) (model.Product, error) {
	query := `
		SELECT 
			p.id, p.product_category_id, p.name, p.description, p.image_url, p.is_active, p.created_at, p.updated_at,
			pc.id, pc.name, pc.slug, pc.description
		FROM products p
		LEFT JOIN product_categories pc ON p.product_category_id = pc.id
		WHERE p.id = $1
	`

	var product model.Product
	var category model.ProductCategory
	var categoryId sql.NullString
	var categoryName sql.NullString
	var categorySlug sql.NullString
	var categoryDesc sql.NullString

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&product.Id, &product.ProductCategoryId, &product.Name, &product.Description,
		&product.ImageUrl, &product.IsActive, &product.CreatedAt, &product.UpdatedAt,
		&categoryId, &categoryName, &categorySlug, &categoryDesc,
	)

	if err != nil {
		if ctx.Err() != nil {
			return model.Product{}, ctx.Err()
		}
		return model.Product{}, err
	}

	if categoryId.Valid {
		category.Id = uuid.MustParse(categoryId.String)
		category.Name = categoryName.String
		category.Slug = categorySlug.String
		if categoryDesc.Valid {
			category.Description = &categoryDesc.String
		}
		product.ProductCategory = &category
	}

	return product, nil
}

// GetProductsByCategory implements ProductRepository
func (r *productRepository) GetProductsByCategory(ctx context.Context, categoryId uuid.UUID) ([]model.Product, error) {
	query := `
		SELECT 
			p.id, p.product_category_id, p.name, p.description, p.image_url, p.is_active, p.created_at, p.updated_at,
			pc.id, pc.name, pc.slug, pc.description
		FROM products p
		LEFT JOIN product_categories pc ON p.product_category_id = pc.id
		WHERE p.product_category_id = $1
		ORDER BY p.created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, categoryId)
	if err != nil {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		return nil, err
	}
	defer rows.Close()

	var products []model.Product
	for rows.Next() {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		var product model.Product
		var category model.ProductCategory
		var categoryIdStr sql.NullString
		var categoryName sql.NullString
		var categorySlug sql.NullString
		var categoryDesc sql.NullString

		err := rows.Scan(
			&product.Id, &product.ProductCategoryId, &product.Name, &product.Description,
			&product.ImageUrl, &product.IsActive, &product.CreatedAt, &product.UpdatedAt,
			&categoryIdStr, &categoryName, &categorySlug, &categoryDesc,
		)
		if err != nil {
			return nil, err
		}

		if categoryIdStr.Valid {
			category.Id = uuid.MustParse(categoryIdStr.String)
			category.Name = categoryName.String
			category.Slug = categorySlug.String
			if categoryDesc.Valid {
				category.Description = &categoryDesc.String
			}
			product.ProductCategory = &category
		}

		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

// UpdateProduct implements ProductRepository
func (r *productRepository) UpdateProduct(ctx context.Context, payload model.Product) (model.Product, error) {
	var product model.Product
	query := `UPDATE products SET product_category_id = $1, name = $2, description = $3, image_url = $4, is_active = $5, updated_at = $6 
			  WHERE id = $7 
			  RETURNING id, product_category_id, name, description, image_url, is_active, created_at, updated_at`

	err := r.db.QueryRowContext(ctx, query, payload.ProductCategoryId, payload.Name, payload.Description,
		payload.ImageUrl, payload.IsActive, time.Now(), payload.Id).
		Scan(&product.Id, &product.ProductCategoryId, &product.Name, &product.Description,
			&product.ImageUrl, &product.IsActive, &product.CreatedAt, &product.UpdatedAt)

	if err != nil {
		if ctx.Err() != nil {
			return model.Product{}, ctx.Err()
		}
		return model.Product{}, err
	}

	return product, nil
}

// DeleteProduct implements ProductRepository
func (r *productRepository) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM products WHERE id = $1`, id)
	if err != nil {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		return err
	}
	return nil
} // Cr

func (r *productRepository) CreateProductAffiliateLink(ctx context.Context, payload model.ProductAffiliateLink) (model.ProductAffiliateLink, error) {
	newId := uuid.Must(uuid.NewV7())
	var link model.ProductAffiliateLink

	query := `INSERT INTO product_affiliate_links (id, product_id, platform_name, url, created_at, updated_at) 
			  VALUES ($1, $2, $3, $4, $5, $6) 
			  RETURNING id, product_id, platform_name, url, created_at, updated_at`

	err := r.db.QueryRowContext(ctx, query, newId, payload.ProductId, payload.PlatformName, payload.Url, time.Now(), time.Now()).
		Scan(&link.Id, &link.ProductId, &link.PlatformName, &link.Url, &link.CreatedAt, &link.UpdatedAt)

	if err != nil {
		if ctx.Err() != nil {
			return model.ProductAffiliateLink{}, ctx.Err()
		}
		return model.ProductAffiliateLink{}, err
	}

	return link, nil
}

// GetAffiliateLinksbyProductId implements ProductRepository
func (r *productRepository) GetAffiliateLinksbyProductId(ctx context.Context, productId uuid.UUID) ([]model.ProductAffiliateLink, error) {
	query := `SELECT id, product_id, platform_name, url, created_at, updated_at 
			  FROM product_affiliate_links 
			  WHERE product_id = $1 
			  ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, productId)
	if err != nil {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		return nil, err
	}
	defer rows.Close()

	var links []model.ProductAffiliateLink
	for rows.Next() {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		var link model.ProductAffiliateLink
		err := rows.Scan(&link.Id, &link.ProductId, &link.PlatformName, &link.Url, &link.CreatedAt, &link.UpdatedAt)
		if err != nil {
			return nil, err
		}
		links = append(links, link)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return links, nil
}

// UpdateProductAffiliateLink implements ProductRepository
func (r *productRepository) UpdateProductAffiliateLink(ctx context.Context, payload model.ProductAffiliateLink) (model.ProductAffiliateLink, error) {
	var link model.ProductAffiliateLink
	query := `UPDATE product_affiliate_links SET platform_name = $1, url = $2, updated_at = $3 
			  WHERE id = $4 
			  RETURNING id, product_id, platform_name, url, created_at, updated_at`

	err := r.db.QueryRowContext(ctx, query, payload.PlatformName, payload.Url, time.Now(), payload.Id).
		Scan(&link.Id, &link.ProductId, &link.PlatformName, &link.Url, &link.CreatedAt, &link.UpdatedAt)

	if err != nil {
		if ctx.Err() != nil {
			return model.ProductAffiliateLink{}, ctx.Err()
		}
		return model.ProductAffiliateLink{}, err
	}

	return link, nil
}

// DeleteProductAffiliateLink implements ProductRepository
func (r *productRepository) DeleteProductAffiliateLink(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM product_affiliate_links WHERE id = $1`, id)
	if err != nil {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		return err
	}
	return nil
}

// AddProductToArticle implements ProductRepository
func (r *productRepository) AddProductToArticle(ctx context.Context, articleId, productId uuid.UUID) error {
	query := `INSERT INTO article_product (article_id, product_id, created_at, updated_at) 
			  VALUES ($1, $2, $3, $4) 
			  ON CONFLICT (article_id, product_id) DO NOTHING`

	_, err := r.db.ExecContext(ctx, query, articleId, productId, time.Now(), time.Now())
	if err != nil {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		return err
	}
	return nil
}

// RemoveProductFromArticle implements ProductRepository
func (r *productRepository) RemoveProductFromArticle(ctx context.Context, articleId, productId uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM article_product WHERE article_id = $1 AND product_id = $2`, articleId, productId)
	if err != nil {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		return err
	}
	return nil
}

// GetProductsByArticleId implements ProductRepository
func (r *productRepository) GetProductsByArticleId(ctx context.Context, articleId uuid.UUID) ([]model.Product, error) {
	query := `
		SELECT 
			p.id, p.product_category_id, p.name, p.description, p.image_url, p.is_active, p.created_at, p.updated_at,
			pc.id, pc.name, pc.slug, pc.description
		FROM products p
		LEFT JOIN product_categories pc ON p.product_category_id = pc.id
		INNER JOIN article_product ap ON p.id = ap.product_id
		WHERE ap.article_id = $1
		ORDER BY ap.created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, articleId)
	if err != nil {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		return nil, err
	}
	defer rows.Close()

	var products []model.Product
	for rows.Next() {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		var product model.Product
		var category model.ProductCategory
		var categoryId sql.NullString
		var categoryName sql.NullString
		var categorySlug sql.NullString
		var categoryDesc sql.NullString

		err := rows.Scan(
			&product.Id, &product.ProductCategoryId, &product.Name, &product.Description,
			&product.ImageUrl, &product.IsActive, &product.CreatedAt, &product.UpdatedAt,
			&categoryId, &categoryName, &categorySlug, &categoryDesc,
		)
		if err != nil {
			return nil, err
		}

		if categoryId.Valid {
			category.Id = uuid.MustParse(categoryId.String)
			category.Name = categoryName.String
			category.Slug = categorySlug.String
			if categoryDesc.Valid {
				category.Description = &categoryDesc.String
			}
			product.ProductCategory = &category
		}

		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

// GetArticlesByProductId implements ProductRepository
func (r *productRepository) GetArticlesByProductId(ctx context.Context, productId uuid.UUID) ([]model.Article, error) {
	query := `
		SELECT 
			a.id, a.title, a.slug, a.content, a.user_id, a.category_id, a.views, a.status, a.created_at, a.updated_at,
			u.id, u.name, u.email, u.role,
			c.id, c.name
		FROM articles a
		JOIN users u ON a.user_id = u.id
		JOIN categories c ON a.category_id = c.id
		INNER JOIN article_product ap ON a.id = ap.article_id
		WHERE ap.product_id = $1
		ORDER BY ap.created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, productId)
	if err != nil {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		return nil, err
	}
	defer rows.Close()

	var articles []model.Article
	for rows.Next() {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		var article model.Article
		var user model.User
		var category model.Category

		err := rows.Scan(
			&article.Id, &article.Title, &article.Slug, &article.Content,
			&article.UserId, &article.CategoryId, &article.Views, &article.Status,
			&article.CreatedAt, &article.UpdatedAt,
			&user.Id, &user.Name, &user.Email, &user.Role,
			&category.Id, &category.Name,
		)
		if err != nil {
			return nil, err
		}

		article.User = &user
		article.Category = &category
		articles = append(articles, article)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return articles, nil
}
