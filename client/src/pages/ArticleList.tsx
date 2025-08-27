import { usePaginatedArticles, useCategories } from "@/hooks/useApi";
import { useState } from "react";
import type {
  Article,
  APIResponse,
  PaginatedArticlesResponse,
} from "@/services/api";

const ITEMS_PER_PAGE = 10;

const ArticleList = () => {
  const [page, setPage] = useState(1);

  const {
    data: articlesResponse,
    isLoading: articlesLoading,
    error: articlesError,
    isPlaceholderData,
  } = usePaginatedArticles(page, ITEMS_PER_PAGE);

  const { data: categories, isLoading: categoriesLoading } = useCategories();

  console.log("Fetched articles response:", articlesResponse);

  // Extract articles and pagination from response
  const typedResponse =
    articlesResponse as APIResponse<PaginatedArticlesResponse>;
  const articles: Article[] = typedResponse?.data?.articles || [];
  const pagination = typedResponse?.pagination;

  const goToPrev = () => {
    if (pagination?.has_prev) {
      setPage((prev) => Math.max(prev - 1, 1));
    }
  };

  const goToNext = () => {
    if (pagination?.has_next) {
      setPage((prev) => prev + 1);
    }
  };

  if (articlesLoading && !isPlaceholderData) {
    return (
      <div className="max-w-7xl mx-auto px-4 py-12">
        <div className="flex items-center justify-center min-h-[400px]">
          <div className="text-center">
            <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto mb-4"></div>
            <p className="text-gray-600">Loading articles...</p>
          </div>
        </div>
      </div>
    );
  }

  if (articlesError) {
    return (
      <div className="max-w-7xl mx-auto px-4 py-12">
        <div className="flex items-center justify-center min-h-[400px]">
          <div className="text-center">
            <p className="text-red-600 mb-4">Error loading articles</p>
            <p className="text-gray-600 text-sm">{articlesError.message}</p>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="max-w-7xl mx-auto px-4 py-12">
      <div className="grid grid-cols-1 md:grid-cols-4 gap-10">
        {/* Article List */}
        <div className="md:col-span-3 space-y-8">
          {articles.length === 0 ? (
            <div className="text-center py-12">
              <p className="text-gray-500 text-lg">No articles found</p>
              <p className="text-gray-400 text-sm mt-2">
                Check back later for new content
              </p>
            </div>
          ) : (
            articles.map((article) => (
              <div
                key={article.id}
                className="flex flex-col md:flex-row gap-6 border-b pb-6 hover:bg-gray-50 rounded-md p-4 transition-all duration-300"
              >
                <div className="w-full md:w-64 h-48 bg-gray-200 rounded-lg shadow flex items-center justify-center">
                  <span className="text-gray-400 text-sm">No Image</span>
                </div>
                <div className="flex-1">
                  <h2 className="text-2xl font-bold mb-2 hover:text-blue-600 transition cursor-pointer">
                    {article.title}
                  </h2>
                  <p className="text-gray-600 text-sm mb-3 line-clamp-3">
                    {article.content.slice(0, 150)}...
                  </p>
                  <div className="flex items-center justify-between">
                    <div className="text-gray-500 text-sm">
                      By{" "}
                      <span className="font-medium">{article.user.name}</span>
                    </div>
                    <div className="text-gray-400 text-xs">
                      {new Date(article.updated_at).toLocaleDateString(
                        "en-US",
                        {
                          year: "numeric",
                          month: "short",
                          day: "numeric",
                        }
                      )}
                    </div>
                  </div>
                  {article.category && (
                    <div className="mt-2">
                      <span className="inline-block bg-blue-100 text-blue-800 text-xs px-2 py-1 rounded">
                        {article.category.name}
                      </span>
                    </div>
                  )}
                </div>
              </div>
            ))
          )}

          {/* Pagination */}
          {pagination && (
            <div className="flex items-center justify-between mt-10">
              <div className="flex items-center gap-4">
                <button
                  onClick={goToPrev}
                  disabled={!pagination.has_prev || isPlaceholderData}
                  className={`px-4 py-2 rounded-md border ${
                    !pagination.has_prev || isPlaceholderData
                      ? "bg-gray-200 text-gray-400 cursor-not-allowed"
                      : "bg-white hover:bg-gray-100 text-gray-700"
                  }`}
                >
                  Previous
                </button>
                <button
                  onClick={goToNext}
                  disabled={!pagination.has_next || isPlaceholderData}
                  className={`px-4 py-2 rounded-md border ${
                    !pagination.has_next || isPlaceholderData
                      ? "bg-gray-200 text-gray-400 cursor-not-allowed"
                      : "bg-white hover:bg-gray-100 text-gray-700"
                  }`}
                >
                  Next
                </button>
              </div>

              <div className="text-sm text-gray-600">
                <span>
                  Page {pagination.page} of {pagination.total_pages}
                </span>
                <span className="ml-4">Total: {pagination.total} articles</span>
              </div>
            </div>
          )}

          {/* Loading indicator for page changes */}
          {isPlaceholderData && (
            <div className="flex items-center justify-center py-4">
              <div className="animate-spin rounded-full h-6 w-6 border-b-2 border-blue-600 mr-2"></div>
              <span className="text-gray-600">Loading...</span>
            </div>
          )}
        </div>

        {/* Sidebar */}
        <aside className="md:col-span-1 space-y-4">
          <h3 className="text-lg font-semibold border-b pb-2">Categories</h3>
          {categoriesLoading ? (
            <div className="animate-pulse space-y-2">
              {[...Array(5)].map((_, i) => (
                <div key={i} className="h-4 bg-gray-200 rounded"></div>
              ))}
            </div>
          ) : (
            <ul className="space-y-2">
              {categories?.map((category) => (
                <li
                  key={category.id}
                  className="text-sm text-blue-600 hover:underline cursor-pointer"
                >
                  {category.name}
                </li>
              ))}
            </ul>
          )}
        </aside>
      </div>
    </div>
  );
};

export default ArticleList;
