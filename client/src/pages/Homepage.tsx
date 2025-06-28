const featuredBig = {
  title: "Discover the Power of Edge Computing",
  excerpt: "A new wave in tech infrastructure is changing the way data moves.",
  image:
    "https://plus.unsplash.com/premium_photo-1749668819614-85f8f8b1396a?q=80&w=1470&auto=format&fit=crop&ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D",
};

const featuredSmalls = [
  {
    title: "Figma vs Adobe XD: Still a Debate?",
    image:
      "https://plus.unsplash.com/premium_photo-1749668819614-85f8f8b1396a?q=80&w=1470&auto=format&fit=crop&ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D",
  },
  {
    title: "Digital Nomad Life in Bali",
    image:
      "https://plus.unsplash.com/premium_photo-1749668819614-85f8f8b1396a?q=80&w=1470&auto=format&fit=crop&ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D",
  },
  {
    title: "React 19 Preview: What's New?",
    image:
      "https://plus.unsplash.com/premium_photo-1749668819614-85f8f8b1396a?q=80&w=1470&auto=format&fit=crop&ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D",
  },
  {
    title: "Rebranding Strategies in 2025",
    image:
      "https://plus.unsplash.com/premium_photo-1749668819614-85f8f8b1396a?q=80&w=1470&auto=format&fit=crop&ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D",
  },
];

const regularPosts = [
  {
    title: "Async Rust: The New Hype or Just Noise?",
    excerpt: "Why async in Rust is so different and what to look out for.",
    image: "https://source.unsplash.com/random/300x200?rust",
  },
  {
    title: "How to Make Your Portfolio Stand Out",
    excerpt: "Design tricks and storytelling that attract hiring managers.",
    image: "https://source.unsplash.com/random/300x200?portfolio",
  },
  {
    title: "Learning DevOps in 2025: What Matters Most?",
    excerpt: "A guide to mastering CI/CD pipelines and infrastructure.",
    image: "https://source.unsplash.com/random/300x200?devops",
  },
];

const categories = [
  "Tech",
  "Design",
  "Productivity",
  "Remote Work",
  "Business",
  "Marketing",
];

const Homepage = () => {
  return (
    <div className="max-w-7xl mx-auto px-4 py-5 space-y-16">
      {/* FEATURED SECTION */}
      <section className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Big Featured */}
        <div className="relative h-[400px] lg:col-span-2 rounded-xl overflow-hidden shadow-md group">
          <img
            src={featuredBig.image}
            alt="Featured"
            className="w-full h-full object-cover group-hover:scale-105 transition duration-500"
          />
          <div className="absolute inset-0  bg-opacity-50 p-8 flex flex-col justify-end text-white">
            <h2 className="text-3xl font-bold">{featuredBig.title}</h2>
            <p className="text-sm text-white/80 mt-2">{featuredBig.excerpt}</p>
          </div>
        </div>

        {/* Small Featured Posts */}
        <div className="space-y-4">
          {featuredSmalls.map((post, index) => (
            <div
              key={index}
              className="flex items-center gap-4 bg-white shadow hover:shadow-md rounded-lg p-3"
            >
              <img
                src={post.image}
                alt="small"
                className="w-16 h-16 rounded object-cover"
              />
              <h4 className="text-sm font-medium">{post.title}</h4>
            </div>
          ))}
        </div>
      </section>

      {/* REGULAR POSTS + CATEGORIES */}
      <section className="grid grid-cols-1 lg:grid-cols-3 gap-10">
        {/* Regular Posts */}
        <div className="lg:col-span-2 space-y-6">
          {regularPosts.map((post, index) => (
            <div
              key={index}
              className="flex gap-4 bg-white border border-gray-200 rounded-md p-4 shadow  hover:shadow-md"
            >
              <img
                src={post.image}
                alt="post"
                className="w-28 h-20 object-cover rounded"
              />
              <div>
                <h3 className="text-lg font-semibold">{post.title}</h3>
                <p className="text-sm text-gray-600">{post.excerpt}</p>
              </div>
            </div>
          ))}

          {/* Pagination */}
          <div className="flex justify-center pt-6">
            <div className="flex items-center space-x-2">
              <button className="px-4 py-2 bg-gray-200 hover:bg-gray-300 rounded">
                Prev
              </button>
              <button className="px-4 py-2 bg-blue-600 text-white rounded">
                1
              </button>
              <button className="px-4 py-2 bg-gray-200 hover:bg-gray-300 rounded">
                2
              </button>
              <button className="px-4 py-2 bg-gray-200 hover:bg-gray-300 rounded">
                Next
              </button>
            </div>
          </div>
        </div>

        {/* Categories Section */}
        <aside className="space-y-4">
          <h4 className="text-xl font-bold mb-4">Categories</h4>
          <ul className="space-y-2 ml-9">
            {categories.map((cat, i) => (
              <li key={i}>
                <button className="w-[55%] text-center  bg-gray-100 hover:bg-gray-200 px-4 py-2 rounded-md text-sm">
                  {cat}
                </button>
              </li>
            ))}
          </ul>
        </aside>
      </section>
    </div>
  );
};

export default Homepage;
