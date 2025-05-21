import type { LoaderFunctionArgs } from "react-router-dom";
import { getAllArticles } from "../services/articleService";
import { useArticleStore } from "../stores/articleStore";

export async function articlesLoader(_: LoaderFunctionArgs) {
  try {
    const { articles, message } = await getAllArticles();
    useArticleStore.getState().setArticles(articles);
    useArticleStore.getState().setNotification(message);
    return null;
  } catch (error) {
    throw new Response("Gagal mengambil data artikel", { status: 500 });
  }
}
