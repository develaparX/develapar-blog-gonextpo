import { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Hash, Search, ArrowLeft } from "lucide-react";
import { useTags } from "@/hooks/useApi";
import type { Tag } from "@/services/api";

const AllTagsPage = () => {
  const navigate = useNavigate();
  const [filteredTags, setFilteredTags] = useState<Tag[]>([]);
  const [searchQuery, setSearchQuery] = useState("");

  // Use TanStack Query hook
  const { data: tags = [], isLoading: loading, error } = useTags();

  useEffect(() => {
    // Filter tags based on search query
    if (searchQuery.trim()) {
      const filtered = tags.filter((tag) =>
        tag.name.toLowerCase().includes(searchQuery.toLowerCase())
      );
      setFilteredTags(filtered);
    } else {
      setFilteredTags(tags);
    }
  }, [tags, searchQuery]);

  const handleTagClick = (tagId: number, tagName: string) => {
    navigate(`/tag/${tagId}?name=${encodeURIComponent(tagName)}`);
  };

  if (loading) {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="max-w-4xl mx-auto">
          <div className="animate-pulse space-y-6">
            <div className="h-8 bg-gray-200 rounded w-1/3"></div>
            <div className="h-10 bg-gray-200 rounded"></div>
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
              {[1, 2, 3, 4, 5, 6].map((i) => (
                <div key={i} className="h-24 bg-gray-200 rounded"></div>
              ))}
            </div>
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

          <div className="flex items-center gap-2 mb-4">
            <Hash size={28} className="text-muted-foreground" />
            <h1 className="text-3xl font-bold">All Tags</h1>
          </div>

          <p className="text-muted-foreground mb-6">
            Browse articles by topic. Click on any tag to see related articles.
          </p>

          {/* Search Tags */}
          <div className="relative mb-6">
            <Input
              type="text"
              placeholder="Search tags..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="pl-10"
            />
            <Search
              size={16}
              className="absolute left-3 top-1/2 transform -translate-y-1/2 text-muted-foreground"
            />
          </div>

          <p className="text-sm text-muted-foreground">
            Showing {filteredTags.length} of {tags.length} tags
          </p>
        </div>

        {/* Tags Grid */}
        {filteredTags.length > 0 ? (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {filteredTags.map((tag) => (
              <Card
                key={tag.id}
                className="hover:shadow-lg transition-all cursor-pointer hover:scale-105"
                onClick={() => handleTagClick(tag.id, tag.name)}
              >
                <CardContent className="p-4">
                  <div className="flex items-center gap-2">
                    <Hash size={16} className="text-muted-foreground" />
                    <span className="font-medium">{tag.name}</span>
                  </div>
                  {tag.description && (
                    <p className="text-sm text-muted-foreground mt-2">
                      {tag.description}
                    </p>
                  )}
                </CardContent>
              </Card>
            ))}
          </div>
        ) : (
          <Card>
            <CardContent className="pt-6 text-center">
              <Search
                size={48}
                className="mx-auto mb-4 text-muted-foreground"
              />
              <h3 className="text-lg font-semibold mb-2">No tags found</h3>
              <p className="text-muted-foreground mb-4">
                {searchQuery
                  ? `No tags match "${searchQuery}". Try a different search term.`
                  : "No tags are available at the moment."}
              </p>
              {searchQuery && (
                <Button variant="outline" onClick={() => setSearchQuery("")}>
                  Clear Search
                </Button>
              )}
            </CardContent>
          </Card>
        )}
      </div>
    </div>
  );
};

export default AllTagsPage;
