import { Outlet, Link } from "react-router-dom";

const MainLayout = () => {
  return (
    <div className="flex flex-col min-h-screen">
      {/* Floating Navbar */}
      <header className="fixed top-6 left-1/2 -translate-x-1/2 z-50">
        <nav className="backdrop-blur-md bg-white/80 dark:bg-black/10 shadow-lg px-6 py-2 rounded-full flex gap-6 items-center text-sm font-medium border border-gray-200 dark:border-gray-300">
          <Link to="/" className="hover:text-blue-600 transition">
            Home
          </Link>
          <Link to="/about" className="hover:text-blue-600 transition">
            About
          </Link>
          <Link to="/contact" className="hover:text-blue-600 transition">
            Contact
          </Link>
          <Link to="/blog" className="hover:text-blue-600 transition">
            Blog
          </Link>
        </nav>
      </header>

      {/* Main content */}
      <main className="flex-1 pt-24">
        <Outlet />
      </main>

      {/* Footer */}
      <footer className="bg-black text-white py-10 mt-20">
        <div className="max-w-6xl mx-auto px-4 text-center space-y-4">
          <p className="text-lg font-semibold">
            YourBlog © {new Date().getFullYear()}
          </p>
          <p className="text-sm text-gray-400">
            Made with 💙 by Your Name. All rights reserved.
          </p>
          <div className="flex justify-center gap-4 text-sm text-gray-400">
            <Link to="/privacy" className="hover:text-white transition">
              Privacy
            </Link>
            <Link to="/terms" className="hover:text-white transition">
              Terms
            </Link>
          </div>
        </div>
      </footer>
    </div>
  );
};

export default MainLayout;
