import { Outlet, Link } from "react-router-dom";
import logo from "../../assets/logo.png";
import { Bell, User } from "lucide-react";

const MainLayout = () => {
  const loggedIn = false;

  return (
    <div className="flex flex-col min-h-screen">
      {/* Navbar */}
      <nav className="b w-full py-5">
        <div className="flex justify-evenly items-center ">
          <div className="flex gap-8 items-center">
            <div className="max-w-60">
              <img src={logo} width={"140"}></img>
            </div>

            <div className="flex gap-8 border-l-2">
              <div className="pl-4">Stories</div>
              <div>Idea</div>
              <div>Gaming</div>
              <div>Tech</div>
            </div>
          </div>
          <div className="flex items-center gap-5">
            {loggedIn ? (
              <>
                <button className="bg-gray-200 hover:bg-gray-700 hover:text-white rounded-lg px-5">
                  Write here
                </button>

                <div>
                  <Bell size={20} />
                </div>
                <div>
                  <User size={20} />
                </div>
              </>
            ) : (
              <div className="border hover:bg-black hover:text-white  rounded-lg py-1 px-10">
                Sign Up!
              </div>
            )}
          </div>
        </div>
      </nav>

      {/* Main content */}
      <main className="flex-1">
        <Outlet />
      </main>

      {/* Footer */}
      <footer className="bg-black text-white py-2 mt-20">
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
