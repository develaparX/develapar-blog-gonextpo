import { useParams } from "react-router-dom";

type Props = {};

const ArticleDetailPage = (props: Props) => {
  const { slug } = useParams();

  return <div>ArticleDetailPage : {slug}</div>;
};

export default ArticleDetailPage;
