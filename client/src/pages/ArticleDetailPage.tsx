import { useArticleBySlug, useProductsByArticleId } from "@/hooks/useApi";
import { useParams } from "react-router-dom";

const ArticleDetailPage = () => {
  const { slug } = useParams();
  const {
    data: article,
    isLoading: articleLoading,
    error: articleError,
  } = useArticleBySlug(slug);

  // Convert article.id to string for UUID compatibility
  const articleIdString = article?.id?.toString();
  const {
    data: products,
    isLoading: productsLoading,
    error: productsError,
  } = useProductsByArticleId(articleIdString);

  if (articleLoading) {
    return <div>Loading article...</div>;
  }

  if (articleError) {
    return <div>Error loading article: {articleError.message}</div>;
  }

  return (
    <div className="max-w-3xl mx-auto px-4 py-12">
      <article>
        {/* ARTICLE DATA SECTION */}
        <h1>ARTICLE DATA:</h1>
        <div>
          <h2>Article ID: {article?.id}</h2>
          <h2>Title: {article?.title}</h2>
          <h2>Slug: {article?.slug}</h2>
          <h2>Content: {article?.content}</h2>
          <h2>Views: {article?.views}</h2>
          <h2>Created At: {article?.created_at}</h2>
          <h2>Updated At: {article?.updated_at}</h2>

          <h3>User Info:</h3>
          <div>User ID: {article?.user.id}</div>
          <div>User Name: {article?.user.name}</div>
          <div>User Email: {article?.user.email}</div>

          <h3>Category Info:</h3>
          <div>Category ID: {article?.category.id}</div>
          <div>Category Name: {article?.category.name}</div>
        </div>

        <hr style={{ margin: "40px 0", border: "2px solid #000" }} />

        {/* PRODUCTS DATA SECTION */}
        <h1>PRODUCTS DATA:</h1>
        {productsLoading && <div>Loading products...</div>}
        {productsError && (
          <div>Error loading products: {productsError.message}</div>
        )}

        {products && products.length > 0 ? (
          <div>
            <h2>Found {products.length} products for this article:</h2>
            {products.map((product, index) => (
              <div
                key={product.id}
                style={{
                  border: "1px solid #ccc",
                  margin: "20px 0",
                  padding: "20px",
                }}
              >
                <h3>Product #{index + 1}:</h3>
                <div>Product ID: {product.id}</div>
                <div>Name: {product.name}</div>
                <div>Description: {product.description}</div>
                <div>Image URL: {product.image_url}</div>
                <div>Is Active: {product.is_active ? "Yes" : "No"}</div>
                <div>Created At: {product.created_at}</div>
                <div>Updated At: {product.updated_at}</div>
                <div>Product Category ID: {product.product_category_id}</div>

                <h4>Product Category Info:</h4>
                <div>Category ID: {product.product_category.id}</div>
                <div>Category Name: {product.product_category.name}</div>
                <div>Category Slug: {product.product_category.slug}</div>
                <div>
                  Category Description: {product.product_category.description}
                </div>
                <div>
                  Category Created At: {product.product_category.created_at}
                </div>
                <div>
                  Category Updated At: {product.product_category.updated_at}
                </div>

                <h4>Affiliate Links:</h4>
                {product.affiliate_links &&
                product.affiliate_links.length > 0 ? (
                  product.affiliate_links.map((link, linkIndex) => (
                    <div
                      key={link.id}
                      style={{ marginLeft: "20px", marginBottom: "10px" }}
                    >
                      <div>Link #{linkIndex + 1}:</div>
                      <div>Link ID: {link.id}</div>
                      <div>Platform: {link.platform_name}</div>
                      <div>
                        URL:{" "}
                        <a
                          href={link.url}
                          target="_blank"
                          rel="noopener noreferrer"
                        >
                          {link.url}
                        </a>
                      </div>
                      <div>
                        Price:{" "}
                        {link.price
                          ? `${link.currency} ${link.price}`
                          : "Not specified"}
                      </div>
                    </div>
                  ))
                ) : (
                  <div>No affiliate links available</div>
                )}
              </div>
            ))}
          </div>
        ) : (
          !productsLoading && <div>No products found for this article</div>
        )}

        <hr style={{ margin: "40px 0", border: "2px solid #000" }} />

        {/* RAW JSON DATA FOR DEBUGGING */}
        <h1>RAW JSON DATA (for debugging):</h1>
        <h2>Article JSON:</h2>
        <pre
          style={{ background: "#f5f5f5", padding: "20px", overflow: "auto" }}
        >
          {JSON.stringify(article, null, 2)}
        </pre>

        <h2>Products JSON:</h2>
        <pre
          style={{ background: "#f5f5f5", padding: "20px", overflow: "auto" }}
        >
          {JSON.stringify(products, null, 2)}
        </pre>
      </article>
    </div>
  );
};

export default ArticleDetailPage;
