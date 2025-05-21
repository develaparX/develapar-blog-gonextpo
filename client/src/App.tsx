import { createBrowserRouter, RouterProvider } from "react-router-dom";
import Homepage from "./pages/Homepage";
import ArticleDetailPage from "./pages/ArticleDetailPage";
import LoginPage from "./pages/LoginPage";
import NotFoundPage from "./pages/NotFoundPage";
import ArticleList from "./pages/articles/index";
import MainLayout from "./components/layout/MainLayout";
import { articlesLoader } from "./routes/articles";

const router = createBrowserRouter(
  [
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
          loader: articlesLoader,
        },
        {
          path: "article/:slug",
          element: <ArticleDetailPage />,
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
  ],
  {
    future: {
      v7_partialHydration: true,
    } as any,
  }
);

function App() {
  return <RouterProvider router={router} />;
}

export default App;
