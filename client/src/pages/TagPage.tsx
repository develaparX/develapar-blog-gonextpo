import { useParams, useSearchParams, useNavigate } from "react-router-dom";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Calendar, User, Eye, ArrowLeft, Hash } from "lucide-react";
import { useArticlesByTag } from "@/hooks/useApi";
import type { Article } from "@/services/api";

const TagPage = () => {
  const { tagId } = useParams<{ tagId: string }>();
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();
  const tagName = searchParams.get("name") || "";

  // Use TanStack Query hook
  const {
    data: articles = [],
    isLoading: loading,
    error,
  } = useArticlesByTag(tagId ? parseInt(tagId) : undefined);

  const articleList = articles || [];

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

  if (loading) {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="max-w-4xl mx-auto">
          <div className="animate-pulse space-y-6">
            <div className="h-8 bg-gray-200 rounded w-1/3"></div>
            <div className="h-4 bg-gray-200 rounded w-2/3"></div>
            {[1, 2, 3].map((i) => (
              <Card key={i}>
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
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="max-w-4xl mx-auto">
          <Card className="border-red-200 bg-red-50">
            <CardContent className="pt-6 text-center">
              <p className="text-red-600 mb-4">{error}</p>
              <Button variant="outline" onClick={() => navigate("/")}>
                Go Home
              </Button>
            </CardContent>
          </Card>
        </div>
      </div>
    );
  }

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="max-w-4xl mx-auto">
        {/* Header */}
        <div className="mb-8">
          <Button
            variant="ghost"
            className="mb-4"
            onClick={() => window.history.back()}
          >
            <ArrowLeft size={16} className="mr-2" />
            Back
          </Button>

          <div className="flex items-center gap-2 mb-2">
            <Hash size={24} className="text-muted-foreground" />
            <h1 className="text-3xl font-bold">{tagName || `Tag ${tagId}`}</h1>
          </div>

          <p className="text-muted-foreground">
            {articleList.length} article{articleList.length !== 1 ? "s" : ""}{" "}
            tagged with this topic
          </p>
        </div>

        {/* Articles */}
        {articleList.length > 0 ? (
          <div className="space-y-6">
            {articleList.map((article) => (
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
        ) : (
          <Card>
            <CardContent className="pt-6 text-center">
              <Hash size={48} className="mx-auto mb-4 text-muted-foreground" />
              <h3 className="text-lg font-semibold mb-2">No articles found</h3>
              <p className="text-muted-foreground mb-4">
                There are no articles with this tag yet.
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

export default TagPage;
