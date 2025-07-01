import { useState } from "react";

type Article = {
  id: number;
  title: string;
  description: string;
  image: string;
  author: string;
  date: string;
};

const allArticles: Article[] = [
  // --- Simulasi banyak artikel ---
  {
    id: 1,
    title: "Understanding React Server Components",
    description:
      "A deep dive into RSC and how they can optimize your frontend performance.",
    image: "https://source.unsplash.com/random/600x400?tech",
    author: "John Doe",
    date: "June 10, 2025",
  },
  {
    id: 2,
    title: "Tailwind Tips for Better UI Design",
    description:
      "Quick tips to make your UI more consistent using Tailwind CSS.",
    image: "https://source.unsplash.com/random/600x400?design",
    author: "Jane Smith",
    date: "June 8, 2025",
  },
  {
    id: 3,
    title: "Why TypeScript Makes React Development Safer",
    description:
      "Explore the benefits of using TypeScript in large-scale React apps.",
    image: "https://source.unsplash.com/random/600x400?typescript",
    author: "Sarah Lin",
    date: "June 5, 2025",
  },
  {
    id: 4,
    title: "How to Build a Blog with Next.js",
    description:
      "Step-by-step building your own blog using Markdown and Next.js",
    image: "https://source.unsplash.com/random/600x400?nextjs",
    author: "Alex Kim",
    date: "June 2, 2025",
  },
  {
    id: 5,
    title: "From JavaScript to Go: Learning Backend the Smart Way",
    description: "How a frontend developer can pick up backend using Golang.",
    image: "https://source.unsplash.com/random/600x400?golang",
    author: "Emma Stone",
    date: "May 28, 2025",
  },
  {
    id: 6,
    title: "How AI Tools Are Changing Software Development",
    description:
      "Exploring the impact of AI on coding workflows and productivity.",
    image: "https://source.unsplash.com/random/600x400?ai",
    author: "Liam Brown",
    date: "May 20, 2025",
  },
  // tambahkan lebih banyak jika perlu
];

const categories = [
  "Tech",
  "Design",
  "TypeScript",
  "React",
  "CSS",
  "Startup",
  "Tools",
];

const ITEMS_PER_PAGE = 3;

const ArticleList = () => {
  const [page, setPage] = useState(1);

  const totalPages = Math.ceil(allArticles.length / ITEMS_PER_PAGE);
  const paginatedArticles = allArticles.slice(
    (page - 1) * ITEMS_PER_PAGE,
    page * ITEMS_PER_PAGE
  );

  const goToPrev = () => setPage((prev) => Math.max(prev - 1, 1));
  const goToNext = () => setPage((prev) => Math.min(prev + 1, totalPages));

  return (
    <div className="max-w-7xl mx-auto px-4 py-12">
      <div className="grid grid-cols-1 md:grid-cols-4 gap-10">
        {/* Article List */}
        <div className="md:col-span-3 space-y-8">
          {paginatedArticles.map((article) => (
            <div
              key={article.id}
              className="flex flex-col md:flex-row gap-6 border-b pb-6 hover:bg-gray-50 rounded-md p-4 transition-all duration-300"
            >
              <img
                src={article.image}
                alt={article.title}
                className="w-full md:w-64 h-48 object-cover rounded-lg shadow"
              />
              <div>
                <h2 className="text-2xl font-bold mb-2 hover:text-blue-600 transition">
                  {article.title}
                </h2>
                <p className="text-gray-600 text-sm mb-3">
                  {article.description}
                </p>
                <div className="text-gray-500 text-sm">
                  By <span className="font-medium">{article.author}</span> on{" "}
                  {article.date}
                </div>
              </div>
            </div>
          ))}

          {/* Pagination */}
          <div className="flex items-center justify-center gap-4 mt-10">
            <button
              onClick={goToPrev}
              disabled={page === 1}
              className={`px-4 py-2 rounded-md border ${
                page === 1
                  ? "bg-gray-200 text-gray-400 cursor-not-allowed"
                  : "bg-white hover:bg-gray-100 text-gray-700"
              }`}
            >
              Previous
            </button>
            <span className="text-sm text-gray-600">
              Page {page} of {totalPages}
            </span>
            <button
              onClick={goToNext}
              disabled={page === totalPages}
              className={`px-4 py-2 rounded-md border ${
                page === totalPages
                  ? "bg-gray-200 text-gray-400 cursor-not-allowed"
                  : "bg-white hover:bg-gray-100 text-gray-700"
              }`}
            >
              Next
            </button>
          </div>
        </div>

        {/* Sidebar */}
        <aside className="md:col-span-1 space-y-4">
          <h3 className="text-lg font-semibold border-b pb-2">Categories</h3>
          <ul className="space-y-2">
            {categories.map((cat, idx) => (
              <li
                key={idx}
                className="text-sm text-blue-600 hover:underline cursor-pointer"
              >
                {cat}
              </li>
            ))}
          </ul>
        </aside>
      </div>
    </div>
  );
};

export default ArticleList;
