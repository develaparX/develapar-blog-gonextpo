import { Outlet, Link, useNavigate } from "react-router-dom";
import logo from "../../assets/logo.png";
import { Bell, User } from "lucide-react";
import { Button } from "../ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "../ui/dropdown-menu";

const MainLayout = () => {
  const navigate = useNavigate();
  const loggedIn = true;

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
              <div className="pl-4 cursor-pointer ">Stories</div>
              <div>Idea</div>
              <div>Gaming</div>
              <div>Tech</div>
            </div>
          </div>
          <div className="flex items-center gap-5">
            {loggedIn ? (
              <>
                <Button
                  variant={"secondary"}
                  onClick={() => navigate("/create")}
                >
                  Create Post
                </Button>

                <DropdownMenu>
                  <DropdownMenuTrigger>
                    <Bell size={20} />
                  </DropdownMenuTrigger>
                  <DropdownMenuContent align="end">
                    <DropdownMenuItem
                      onClick={() => navigate("/notifications")}
                    >
                      Notifications
                    </DropdownMenuItem>
                    <DropdownMenuItem onClick={() => navigate("/messages")}>
                      Messages
                    </DropdownMenuItem>
                  </DropdownMenuContent>
                </DropdownMenu>

                <DropdownMenu>
                  <DropdownMenuTrigger>
                    <User size={20} />
                  </DropdownMenuTrigger>
                  <DropdownMenuContent align="start">
                    <DropdownMenuItem onClick={() => navigate("/profile")}>
                      Profile
                    </DropdownMenuItem>
                    <DropdownMenuItem onClick={() => navigate("/settings")}>
                      Settings
                    </DropdownMenuItem>
                    <DropdownMenuItem onClick={() => navigate("/logout")}>
                      Logout
                    </DropdownMenuItem>
                  </DropdownMenuContent>
                </DropdownMenu>
              </>
            ) : (
              <Button variant={"outline"} onClick={() => navigate("/login")}>
                Sign Up!
              </Button>
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
          <p className="text-md font-semibold">
            develapar © {new Date().getFullYear()}
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
