import { createBrowserRouter, RouterProvider } from "react-router-dom";
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
import CommentsPage from "./pages/CommentsPage";
import LikesPage from "./pages/LikesPage";
import TagsPage from "./pages/TagsPage";
import UsersPage from "./pages/UsersPage";

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
        path: "articledetail",
        element: <ArticleDetailPage />,
      },
      {
        path: "bookmarks",
        element: <BookmarksPage />,
      },
      {
        path: "comments",
        element: <CommentsPage />,
      },
      {
        path: "likes",
        element: <LikesPage />,
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
        element: <TagsPage />,
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
  return <RouterProvider router={router} />;
}

export default App;
