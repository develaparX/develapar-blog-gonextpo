import { useState, useEffect } from "react";
import { useSearchParams, useNavigate } from "react-router-dom";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Search, Calendar, User, Eye } from "lucide-react";
import { useManualSearch } from "@/hooks/useApi";
import type { Article } from "@/services/api";

const SearchPage = () => {
  const [searchParams, setSearchParams] = useSearchParams();
  const navigate = useNavigate();
  const [searchQuery, setSearchQuery] = useState(searchParams.get("q") || "");

  // Use TanStack Query hook for manual search
  const {
    data: articles,
    isLoading: loading,
    error,
    search,
  } = useManualSearch();

  useEffect(() => {
    const query = searchParams.get("q");
    if (query) {
      setSearchQuery(query);
      search(query);
    }
  }, [searchParams, search]);

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    if (searchQuery.trim()) {
      setSearchParams({ q: searchQuery.trim() });
    }
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString("en-US", {
      year: "numeric",
      month: "long",
      day: "numeric",
    });
  };

  const truncateContent = (content: string, maxLength: number = 200) => {
    if (content.length <= maxLength) return content;
    return content.substring(0, maxLength) + "...";
  };

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="max-w-4xl mx-auto">
        {/* Search Header */}
        <div className="mb-8">
          <h1 className="text-3xl font-bold mb-4">Search Articles</h1>

          {/* Search Form */}
          <form onSubmit={handleSearch} className="relative mb-4">
            <Input
              type="text"
              placeholder="Search articles by title or content..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="w-full pr-12 h-12 text-lg"
            />
            <Button
              type="submit"
              size="sm"
              className="absolute right-2 top-1/2 transform -translate-y-1/2"
              disabled={loading}
            >
              <Search size={16} />
            </Button>
          </form>

          {/* Search Results Info */}
          {searchParams.get("q") && (
            <p className="text-muted-foreground">
              {loading
                ? "Searching..."
                : `Found ${articles.length} result${
                    articles.length !== 1 ? "s" : ""
                  } for "${searchParams.get("q")}"`}
            </p>
          )}
        </div>

        {/* Error State */}
        {error && (
          <Card className="mb-6 border-red-200 bg-red-50">
            <CardContent className="pt-6">
              <p className="text-red-600">{error}</p>
            </CardContent>
          </Card>
        )}

        {/* Loading State */}
        {loading && (
          <div className="space-y-4">
            {[1, 2, 3].map((i) => (
              <Card key={i} className="animate-pulse">
                <CardHeader>
                  <div className="h-6 bg-gray-200 rounded w-3/4"></div>
                </CardHeader>
                <CardContent>
                  <div className="space-y-2">
                    <div className="h-4 bg-gray-200 rounded"></div>
                    <div className="h-4 bg-gray-200 rounded w-5/6"></div>
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>
        )}

        {/* Search Results */}
        {!loading && articles.length > 0 && (
          <div className="space-y-6">
            {articles.map((article) => (
              <Card
                key={article.id}
                className="hover:shadow-lg transition-shadow"
              >
                <CardHeader>
                  <div className="flex items-start justify-between">
                    <div className="flex-1">
                      <CardTitle
                        className="text-xl mb-2 hover:text-blue-600 cursor-pointer"
                        onClick={() => navigate(`/article/${article.slug}`)}
                      >
                        {article.title}
                      </CardTitle>
                      <div className="flex items-center gap-4 text-sm text-muted-foreground">
                        <div className="flex items-center gap-1">
                          <User size={14} />
                          <span>{article.user.name}</span>
                        </div>
                        <div className="flex items-center gap-1">
                          <Calendar size={14} />
                          <span>{formatDate(article.created_at)}</span>
                        </div>
                        <div className="flex items-center gap-1">
                          <Eye size={14} />
                          <span>{article.views} views</span>
                        </div>
                      </div>
                    </div>
                    <Badge variant="secondary">{article.category.name}</Badge>
                  </div>
                </CardHeader>
                <CardContent>
                  <p className="text-muted-foreground leading-relaxed">
                    {truncateContent(article.content)}
                  </p>
                  <Button
                    variant="link"
                    className="mt-2 p-0 h-auto"
                    onClick={() => navigate(`/article/${article.slug}`)}
                  >
                    Read more →
                  </Button>
                </CardContent>
              </Card>
            ))}
          </div>
        )}

        {/* No Results */}
        {!loading &&
          searchParams.get("q") &&
          articles.length === 0 &&
          !error && (
            <Card>
              <CardContent className="pt-6 text-center">
                <Search
                  size={48}
                  className="mx-auto mb-4 text-muted-foreground"
                />
                <h3 className="text-lg font-semibold mb-2">
                  No articles found
                </h3>
                <p className="text-muted-foreground mb-4">
                  Try adjusting your search terms or browse our categories.
                </p>
                <Button variant="outline" onClick={() => navigate("/articles")}>
                  Browse All Articles
                </Button>
              </CardContent>
            </Card>
          )}
      </div>
    </div>
  );
};

export default SearchPage;
