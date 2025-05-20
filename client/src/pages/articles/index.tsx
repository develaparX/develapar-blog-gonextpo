import { useEffect } from "react";
import { useArticleStore } from "../../stores/articleStore";

const ArticleListPage = () => {
  const { articles, notification, fetchArticles } = useArticleStore();

  useEffect(() => {
    const fetchData = async () => {
      try {
        await fetchArticles();
      } catch (error) {
        console.error("Error di komponen fetchArticles:", error);
      }
    };
    fetchData();
  }, [fetchArticles]); // tambahin fetchArticles di dependency supaya ESLint senang

  useEffect(() => {
    if (notification) {
      // Contoh sederhana: tampilkan notifikasi di console,
      // kamu bisa ganti ini dengan toast notification kalau mau
      console.log("Notifikasi:", notification);
    }
  }, [notification]);

  return (
    <div className="p-4">
      <h1 className="text-2xl font-bold mb-4">Daftar Artikel</h1>
      <div className="grid gap-4">
        {articles.map((article) => (
          <div key={article.id}>
            <a
              href={`/article/${article.slug}`}
              className="text-blue-500 hover:underline"
            >
              {article.title}
            </a>
            <p>{article.content}</p>
          </div>
        ))}
      </div>
    </div>
  );
};

export default ArticleListPage;
