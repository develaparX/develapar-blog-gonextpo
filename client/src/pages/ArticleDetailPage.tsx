import { useArticleBySlug, useProductsByArticleId } from "@/hooks/useApi";
import { useParams, Link } from "react-router-dom";
import {
  ArrowLeft,
  Calendar,
  Eye,
  User,
  Tag,
  ExternalLink,
  ShoppingCart,
} from "lucide-react";

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
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto mb-4"></div>
          <p className="text-gray-600">Loading article...</p>
        </div>
      </div>
    );
  }

  if (articleError) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center">
          <div className="text-red-500 text-6xl mb-4">⚠️</div>
          <h2 className="text-2xl font-bold text-gray-800 mb-2">
            Article Not Found
          </h2>
          <p className="text-gray-600 mb-6">{articleError.message}</p>
          <Link
            to="/articles"
            className="inline-flex items-center px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
          >
            <ArrowLeft className="w-4 h-4 mr-2" />
            Back to Articles
          </Link>
        </div>
      </div>
    );
  }

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString("en-US", {
      year: "numeric",
      month: "long",
      day: "numeric",
    });
  };

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header Navigation */}
      <div className="bg-white shadow-sm border-b">
        <div className="max-w-4xl mx-auto px-4 py-4">
          <Link
            to="/articles"
            className="inline-flex items-center text-gray-600 hover:text-blue-600 transition-colors"
          >
            <ArrowLeft className="w-4 h-4 mr-2" />
            Back to Articles
          </Link>
        </div>
      </div>

      <div className="max-w-4xl mx-auto px-4 py-8">
        {/* Article Header */}
        <header className="mb-8">
          <div className="mb-4">
            <div className="inline-flex items-center rounded-md border px-2.5 py-0.5 text-xs font-semibold transition-colors focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2 border-transparent bg-primary text-primary-foreground hover:bg-primary/80">
              <Tag className="w-3 h-3 mr-1" />
              {article?.category?.name}
            </div>
          </div>

          <h1 className="scroll-m-20 text-4xl font-extrabold tracking-tight lg:text-5xl mb-6">
            {article?.title}
          </h1>

          <div className="flex flex-wrap items-center gap-6 text-sm text-muted-foreground">
            <div className="flex items-center">
              <User className="w-4 h-4 mr-2" />
              <span className="font-medium">{article?.user?.name}</span>
            </div>
            <div className="flex items-center">
              <Calendar className="w-4 h-4 mr-2" />
              <span>{formatDate(article?.updated_at || "")}</span>
            </div>
            <div className="flex items-center">
              <Eye className="w-4 h-4 mr-2" />
              <span>{article?.views} views</span>
            </div>
          </div>
        </header>

        {/* Article Content */}
        <article className="mb-12">
          <div
            className="prose prose-lg max-w-none prose-headings:scroll-m-20 prose-headings:font-semibold prose-headings:tracking-tight prose-h1:text-4xl prose-h2:text-3xl prose-h3:text-2xl prose-h4:text-xl prose-p:leading-7 prose-blockquote:border-l-2 prose-blockquote:pl-6 prose-blockquote:italic prose-code:relative prose-code:rounded prose-code:bg-muted prose-code:px-[0.3rem] prose-code:py-[0.2rem] prose-code:font-mono prose-code:text-sm prose-pre:overflow-x-auto"
            dangerouslySetInnerHTML={{ __html: article?.content || "" }}
          />
        </article>

        {/* Related Products Section */}
        {products && products.length > 0 && (
          <section className="mt-12">
            <div className="mb-6">
              <h2 className="scroll-m-20 text-3xl font-semibold tracking-tight flex items-center">
                <ShoppingCart className="w-6 h-6 mr-3 text-primary" />
                Related Products
              </h2>
              <p className="text-muted-foreground mt-2">
                Products mentioned in this article
              </p>
            </div>

            <div className="grid gap-6 md:grid-cols-2">
              {products.map((product) => (
                <div
                  key={product.id}
                  className="rounded-lg border bg-card text-card-foreground shadow-sm hover:shadow-md transition-shadow"
                >
                  <div className="p-6">
                    <div className="flex items-start justify-between mb-4">
                      <div className="flex-1">
                        <h3 className="scroll-m-20 text-xl font-semibold tracking-tight mb-2">
                          {product.name}
                        </h3>
                        <p className="text-sm text-muted-foreground mb-3">
                          {product.description}
                        </p>
                        <div className="inline-flex items-center rounded-md border px-2.5 py-0.5 text-xs font-semibold transition-colors focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2 border-transparent bg-secondary text-secondary-foreground">
                          {product.product_category.name}
                        </div>
                      </div>
                      {product.image_url && (
                        <img
                          src={product.image_url}
                          alt={product.name}
                          className="w-16 h-16 object-cover rounded-md ml-4"
                        />
                      )}
                    </div>

                    {/* Affiliate Links */}
                    {product.affiliate_links &&
                      product.affiliate_links.length > 0 && (
                        <div className="space-y-3 pt-4 border-t">
                          <h4 className="text-sm font-medium">Available at:</h4>
                          <div className="space-y-2">
                            {product.affiliate_links.map((link) => (
                              <a
                                key={link.id}
                                href={link.url}
                                target="_blank"
                                rel="noopener noreferrer"
                                className="flex items-center justify-between p-3 rounded-md border bg-muted/50 hover:bg-muted transition-colors group"
                              >
                                <div className="flex items-center">
                                  <span className="font-medium">
                                    {link.platform_name}
                                  </span>
                                  {link.price && (
                                    <span className="ml-2 text-sm font-semibold text-green-600">
                                      {link.currency} {link.price}
                                    </span>
                                  )}
                                </div>
                                <ExternalLink className="w-4 h-4 text-muted-foreground group-hover:text-primary transition-colors" />
                              </a>
                            ))}
                          </div>
                        </div>
                      )}
                  </div>
                </div>
              ))}
            </div>
          </section>
        )}

        {productsLoading && (
          <section className="mt-12">
            <div className="rounded-lg border bg-card text-card-foreground shadow-sm p-6">
              <div className="flex items-center justify-center">
                <div className="animate-spin rounded-full h-6 w-6 border-b-2 border-primary mr-3"></div>
                <span className="text-muted-foreground">
                  Loading related products...
                </span>
              </div>
            </div>
          </section>
        )}

        {productsError && (
          <section className="mt-12">
            <div className="rounded-lg border bg-card text-card-foreground shadow-sm p-6">
              <div className="text-center text-muted-foreground">
                <p>Unable to load related products</p>
              </div>
            </div>
          </section>
        )}
      </div>
    </div>
  );
};

export default ArticleDetailPage;
