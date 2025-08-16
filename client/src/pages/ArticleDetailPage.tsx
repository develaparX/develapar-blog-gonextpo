import { useArticleBySlug } from "@/hooks/useApi";
import { useParams } from "react-router-dom";

const ArticleDetailPage = () => {
  const { slug } = useParams();
  const { data: article, isLoading } = useArticleBySlug(slug);

  return (
    <div className="max-w-3xl mx-auto px-4 py-12">
      <article>
        <h1 className="text-4xl font-bold mb-4 leading-tight">
          {article?.title}
        </h1>

        <div className="text-gray-500 text-sm mb-8">
          By <span className="font-semibold">{article?.user.name}</span> ·
          {article?.created_at} · 5 min read
        </div>
        <div>category : {article?.category.name}</div>

        <img
          src="https://source.unsplash.com/random/900x400?article"
          alt="Article Banner"
          className="w-full h-72 object-cover rounded-lg mb-8"
        />

        <div className="prose prose-lg max-w-none">
          <p>{article?.content}</p>
        </div>

        <div>{}</div>
      </article>
    </div>
  );
};

export default ArticleDetailPage;
