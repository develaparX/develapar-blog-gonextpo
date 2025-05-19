import { Routes, Route } from "react-router";
import Homepage from "./pages/Homepage";

function App() {
  return (
    <Routes>
      <Route path="/" element={<Homepage />} />
      {/* <Route path="/articles/:slug" element={<ArticlePage />} />
      <Route path="/login" element={<LoginPage />} />
      <Route path="*" element={<NotFoundPage />} /> */}
    </Routes>
  );
}

export default App;
