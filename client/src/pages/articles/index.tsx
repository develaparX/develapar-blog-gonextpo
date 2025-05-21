import { useArticleStore } from "../../stores/articleStore";

const ArticleListPage = () => {
  const articles = useArticleStore((state) => state.articles);
  // const notification = useArticleStore((state) => state.notification);

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
