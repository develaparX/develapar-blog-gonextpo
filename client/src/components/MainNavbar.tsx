import { Bell, User, Search } from "lucide-react";
import logo from "@/assets/logo.png";
import {
  NavigationMenu,
  NavigationMenuContent,
  NavigationMenuItem,
  NavigationMenuLink,
  NavigationMenuList,
  NavigationMenuTrigger,
  navigationMenuTriggerStyle,
} from "./ui/navigation-menu";
import { useNavigate } from "react-router-dom";
import { Button } from "./ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "./ui/dropdown-menu";
import { Input } from "./ui/input";
import { useState } from "react";
import { useCategories, useTags } from "@/hooks/useApi";
const MainNavbar = () => {
  const navigate = useNavigate();
  const loggedIn = true;

  // State for search
  const [searchQuery, setSearchQuery] = useState("");

  // Use centralized API hooks
  const { data: categories, loading: categoriesLoading } = useCategories();
  const { data: tags, loading: tagsLoading } = useTags();

  // Ensure we have arrays to work with
  const categoriesList = categories || [];
  const tagsList = tags || [];

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    if (searchQuery.trim()) {
      navigate(`/search?q=${encodeURIComponent(searchQuery.trim())}`);
    }
  };

  const handleCategoryClick = (categoryId: number, categoryName: string) => {
    navigate(`/category/${encodeURIComponent(categoryName)}`);
  };

  const handleTagClick = (tagId: number, tagName: string) => {
    navigate(`/tag/${tagId}?name=${encodeURIComponent(tagName)}`);
  };

  return (
    <div className="w-full mt-5 z-100">
      <NavigationMenu
        viewport={false}
        className="w-full max-w-none justify-between items-center px-5"
      >
        <div className="max-w-60">
          <img src={logo} width={"140"}></img>
        </div>
        <NavigationMenuList>
          <NavigationMenuItem>
            <NavigationMenuLink
              asChild
              className={navigationMenuTriggerStyle()}
              onClick={() => navigate("/")}
            >
              <div>Home</div>
            </NavigationMenuLink>
          </NavigationMenuItem>

          {/* Dynamic Categories Menu */}
          <NavigationMenuItem>
            <NavigationMenuTrigger>Categories</NavigationMenuTrigger>
            <NavigationMenuContent>
              <ul className="grid w-[400px] gap-2 md:w-[500px] md:grid-cols-2 lg:w-[600px] p-4">
                {categoriesList.length > 0 ? (
                  categoriesList.map((category) => (
                    <CategoryItem
                      key={category.id}
                      title={category.name}
                      description={
                        category.description ||
                        `Browse articles in ${category.name}`
                      }
                      onClick={() =>
                        handleCategoryClick(category.id, category.name)
                      }
                    />
                  ))
                ) : (
                  <li className="text-muted-foreground">
                    Loading categories...
                  </li>
                )}
              </ul>
            </NavigationMenuContent>
          </NavigationMenuItem>

          {/* Tags Menu */}
          <NavigationMenuItem>
            <NavigationMenuTrigger>Tags</NavigationMenuTrigger>
            <NavigationMenuContent>
              <div className="w-[400px] p-4">
                <div className="grid grid-cols-3 gap-2">
                  {tagsList.length > 0 ? (
                    tagsList.slice(0, 12).map((tag) => (
                      <Button
                        key={tag.id}
                        variant="ghost"
                        size="sm"
                        className="justify-start h-8 text-xs"
                        onClick={() => handleTagClick(tag.id, tag.name)}
                      >
                        #{tag.name}
                      </Button>
                    ))
                  ) : (
                    <span className="text-muted-foreground col-span-3">
                      Loading tags...
                    </span>
                  )}
                </div>
                {tagsList.length > 12 && (
                  <Button
                    variant="link"
                    size="sm"
                    className="mt-2 p-0 h-auto"
                    onClick={() => navigate("/tags")}
                  >
                    View all tags →
                  </Button>
                )}
              </div>
            </NavigationMenuContent>
          </NavigationMenuItem>

          <NavigationMenuItem>
            <NavigationMenuLink
              asChild
              className={navigationMenuTriggerStyle()}
              onClick={() => navigate("/articles")}
            >
              <div>All Articles</div>
            </NavigationMenuLink>
          </NavigationMenuItem>
        </NavigationMenuList>

        {/* Search Bar */}
        <div className="flex items-center gap-3">
          <form onSubmit={handleSearch} className="relative">
            <Input
              type="text"
              placeholder="Search articles..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="w-64 pr-10"
            />
            <Button
              type="submit"
              size="sm"
              variant="ghost"
              className="absolute right-1 top-1/2 transform -translate-y-1/2 h-8 w-8 p-0"
            >
              <Search size={16} />
            </Button>
          </form>
        </div>

        <div className="flex items-center gap-5">
          {loggedIn ? (
            <>
              <Button variant={"secondary"} onClick={() => navigate("/create")}>
                Create Post
              </Button>

              <DropdownMenu>
                <DropdownMenuTrigger>
                  <Bell size={20} />
                </DropdownMenuTrigger>
                <DropdownMenuContent align="end" className="mt-3">
                  <DropdownMenuItem onClick={() => navigate("/notifications")}>
                    Notifications
                  </DropdownMenuItem>
                  <DropdownMenuItem onClick={() => navigate("/messages")}>
                    Messages
                  </DropdownMenuItem>
                </DropdownMenuContent>
              </DropdownMenu>

              <DropdownMenu>
                <DropdownMenuTrigger>
                  <User size={20} />
                </DropdownMenuTrigger>
                <DropdownMenuContent align="end" className="mt-3">
                  <DropdownMenuItem onClick={() => navigate("/profile")}>
                    Profile
                  </DropdownMenuItem>
                  <DropdownMenuItem onClick={() => navigate("/settings")}>
                    Settings
                  </DropdownMenuItem>
                  <DropdownMenuItem onClick={() => navigate("/logout")}>
                    Logout
                  </DropdownMenuItem>
                </DropdownMenuContent>
              </DropdownMenu>
            </>
          ) : (
            <Button variant={"outline"} onClick={() => navigate("/login")}>
              Sign Up!
            </Button>
          )}
        </div>
      </NavigationMenu>
    </div>
  );
};

export default MainNavbar;

function CategoryItem({
  title,
  description,
  onClick,
  ...props
}: React.ComponentPropsWithoutRef<"li"> & {
  title: string;
  description: string;
  onClick: () => void;
}) {
  return (
    <li {...props}>
      <NavigationMenuLink asChild>
        <div
          className="cursor-pointer hover:bg-accent hover:text-accent-foreground p-3 rounded-md transition-colors"
          onClick={onClick}
        >
          <div className="text-sm leading-none font-medium">{title}</div>
          <p className="text-muted-foreground line-clamp-2 text-sm leading-snug mt-1">
            {description}
          </p>
        </div>
      </NavigationMenuLink>
    </li>
  );
}
