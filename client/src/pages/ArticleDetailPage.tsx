import { useParams } from "react-router-dom";

const ArticleDetailPage = () => {
  const { slug } = useParams();

  return <div>ArticleDetailPage : {slug}</div>;
};

export default ArticleDetailPage;
