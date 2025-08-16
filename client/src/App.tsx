import { createBrowserRouter, RouterProvider } from "react-router-dom";
import { QueryClientProvider } from "@tanstack/react-query";
import { ReactQueryDevtools } from "@tanstack/react-query-devtools";
import { queryClient } from "./lib/queryClient";
import Homepage from "./pages/Homepage";
import ArticleDetailPage from "./pages/ArticleDetailPage";
import LoginPage from "./pages/LoginPage";
import NotFoundPage from "./pages/NotFoundPage";
import MainLayout from "./components/layout/MainLayout";
import ArticleList from "./pages/ArticleList";
import Dashboard from "./pages/dashboard/Dashboard";
import AdminLayout from "./components/layout/AdminLayout";
import AddArticle from "./pages/AddArticle";
import ArticleTagsPage from "./pages/ArticleTagsPage";
import BookmarksPage from "./pages/BookmarksPage";
import CategoriesPage from "./pages/CategoriesPage";
import TagsPageAdmin from "./pages/TagPage";
import UsersPage from "./pages/UsersPage";
// New search and filter pages
import SearchPage from "./pages/SearchPage";
import CategoryPage from "./pages/CategoryPage";
import TagPage from "./pages/TagPage";
import AllTagsPagePublic from "./pages/AllTagsPage";

const router = createBrowserRouter([
  {
    path: "/",
    element: <MainLayout />,
    children: [
      {
        index: true,
        element: <Homepage />,
      },
      {
        path: "articles",
        element: <ArticleList />,
      },
      {
        path: "article/:slug",
        element: <ArticleDetailPage />,
      },
      {
        path: "bookmarks",
        element: <BookmarksPage />,
      },
      {
        path: "search",
        element: <SearchPage />,
      },
      {
        path: "category/:categoryName",
        element: <CategoryPage />,
      },
      {
        path: "tag/:tagId",
        element: <TagPage />,
      },
      {
        path: "tags",
        element: <AllTagsPagePublic />,
      },
    ],
  },
  {
    path: "/dashboard",
    element: <AdminLayout />,
    children: [
      {
        index: true,
        element: <Dashboard />,
      },
      {
        path: "add-article",
        element: <AddArticle />,
      },
      {
        path: "article-tags",
        element: <ArticleTagsPage />,
      },
      {
        path: "categories",
        element: <CategoriesPage />,
      },
      {
        path: "tags",
        element: <TagsPageAdmin />,
      },
      {
        path: "users",
        element: <UsersPage />,
      },
    ],
  },
  {
    path: "/login",
    element: <LoginPage />,
  },
  {
    path: "*",
    element: <NotFoundPage />,
  },
]);

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <RouterProvider router={router} />
      {/* React Query Devtools - only shows in development */}
      <ReactQueryDevtools initialIsOpen={false} />
    </QueryClientProvider>
  );
}

export default App;
