import React, { useState } from "react";

type Post = {
  id: number;
  title: string;
  excerpt: string;
  author: string;
  date: string;
  image: string;
  featured: boolean;
};

const allPosts: Post[] = [
  // === Featured (5 total) ===
  {
    id: 1,
    title: "The Future of Clean UI",
    excerpt: "Why minimal design will define the next decade of interfaces.",
    author: "Aria Blake",
    date: "June 10, 2025",
    image:
      "https://images.unsplash.com/photo-1743701168271-15d33fac46e8?q=80&w=1374&auto=format&fit=crop&ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D",
    featured: true,
  },
  {
    id: 2,
    title: "Dark Mode Isn’t Just a Trend",
    excerpt: "It’s accessibility, battery, and mood — all in one.",
    author: "Derek Lin",
    date: "June 11, 2025",
    image:
      "https://images.unsplash.com/photo-1743701168271-15d33fac46e8?q=80&w=1374&auto=format&fit=crop&ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D",
    featured: true,
  },
  {
    id: 3,
    title: "Typography That Talks",
    excerpt: "Fonts speak louder than words. Use them wisely.",
    author: "Rina Ko",
    date: "June 9, 2025",
    image:
      "https://plus.unsplash.com/premium_photo-1749668819614-85f8f8b1396a?q=80&w=1470&auto=format&fit=crop&ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D",
    featured: true,
  },
  {
    id: 4,
    title: "Design Systems: Friend or Foe?",
    excerpt: "Are they empowering creativity or limiting it?",
    author: "Leo Grant",
    date: "June 8, 2025",
    image:
      "https://plus.unsplash.com/premium_photo-1749668819614-85f8f8b1396a?q=80&w=1470&auto=format&fit=crop&ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D",
    featured: true,
  },
  {
    id: 5,
    title: "Motion UI Done Right",
    excerpt: "Animation should guide, not distract.",
    author: "Sophie Moore",
    date: "June 7, 2025",
    image:
      "https://plus.unsplash.com/premium_photo-1749668819614-85f8f8b1396a?q=80&w=1470&auto=format&fit=crop&ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D",
    featured: true,
  },

  // === Regular Posts ===
  ...Array.from({ length: 20 }).map((_, i) => ({
    id: i + 6,
    title: `Regular Post #${i + 1}`,
    excerpt: "Lorem ipsum dolor sit amet, consectetur adipiscing elit.",
    author: "Author " + (i + 1),
    date: `June ${6 - (i % 5)}, 2025`,
    image:
      "https://plus.unsplash.com/premium_photo-1749668819614-85f8f8b1396a?q=80&w=1470&auto=format&fit=crop&ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D",
    featured: false,
  })),
];

const POSTS_PER_PAGE = 6;

const Homepage = () => {
  const featured = allPosts.filter((p) => p.featured);
  const regular = allPosts.filter((p) => !p.featured);

  const [page, setPage] = useState(1);
  const totalPages = Math.ceil(regular.length / POSTS_PER_PAGE);

  const currentRegularPosts = regular.slice(
    (page - 1) * POSTS_PER_PAGE,
    page * POSTS_PER_PAGE
  );

  return (
    <main className="max-w-7xl mx-auto px-4 py-12 space-y-16">
      {/* Featured Section */}
      <section>
        <div className="grid lg:grid-cols-3 gap-6">
          {/* Large post */}
          <div className="lg:col-span-2 relative h-[400px] rounded-2xl overflow-hidden group shadow-lg">
            <img
              src={featured[0].image}
              alt={featured[0].title}
              className="w-full h-full object-cover group-hover:scale-105 transition duration-700"
            />
            <div className="absolute inset-0 bg-gradient-to-t from-black/70 via-black/30 to-transparent p-8 flex flex-col justify-end">
              <h2 className="text-3xl font-bold text-white mb-2">
                {featured[0].title}
              </h2>
              <p className="text-white/90">{featured[0].excerpt}</p>
              <span className="text-sm text-white/70 mt-2">
                By {featured[0].author} · {featured[0].date}
              </span>
            </div>
          </div>

          {/* Smaller posts */}
          <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
            {featured.slice(1).map((post) => (
              <div
                key={post.id}
                className="group relative h-48 rounded-xl overflow-hidden shadow-md"
              >
                <img
                  src={post.image}
                  alt={post.title}
                  className="w-full h-full object-cover group-hover:scale-105 transition duration-500"
                />
                <div className="absolute inset-0 bg-gradient-to-t from-black/70 to-transparent p-4 flex flex-col justify-end">
                  <h3 className="text-lg font-semibold text-white">
                    {post.title}
                  </h3>
                </div>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* Regular Posts */}
      <section>
        <h3 className="text-2xl font-bold mb-6">Latest Posts</h3>
        <div className="grid md:grid-cols-3 gap-8">
          {currentRegularPosts.map((post) => (
            <div
              key={post.id}
              className="rounded-xl overflow-hidden shadow hover:shadow-lg transition"
            >
              <img
                src={post.image}
                alt={post.title}
                className="h-48 w-full object-cover"
              />
              <div className="p-4 bg-white">
                <h4 className="text-lg font-semibold mb-1 hover:underline cursor-pointer">
                  {post.title}
                </h4>
                <p className="text-gray-600 text-sm mb-2">{post.excerpt}</p>
                <span className="text-sm text-gray-400">
                  By {post.author} · {post.date}
                </span>
              </div>
            </div>
          ))}
        </div>

        {/* Pagination */}
        <div className="flex justify-center items-center gap-3 mt-10">
          {Array.from({ length: totalPages }).map((_, i) => (
            <button
              key={i}
              onClick={() => setPage(i + 1)}
              className={`px-4 py-2 rounded-lg text-sm ${
                page === i + 1
                  ? "bg-black text-white"
                  : "bg-gray-100 hover:bg-gray-200"
              }`}
            >
              {i + 1}
            </button>
          ))}
        </div>
      </section>
    </main>
  );
};

export default Homepage;
