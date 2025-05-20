import { Routes, Route } from "react-router-dom";
import Homepage from "./pages/Homepage";
import ArticleDetailPage from "./pages/ArticleDetailPage";
import LoginPage from "./pages/LoginPage";
import NotFoundPage from "./pages/NotFoundPage";
import ArticleList from "./pages/articles/index";

function App() {
  return (
    <Routes>
      <Route path="/" element={<Homepage />} />
      <Route path="/articles" element={<ArticleList />} />

      <Route path="/articles/:slug" element={<ArticleDetailPage />} />
      <Route path="/login" element={<LoginPage />} />
      <Route path="*" element={<NotFoundPage />} />
    </Routes>
  );
}

export default App;
