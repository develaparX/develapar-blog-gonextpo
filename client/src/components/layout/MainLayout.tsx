import { Outlet, Link } from "react-router-dom";
// import logo from "../../assets/logo.png";
// import { Bell, User } from "lucide-react";
// import { Button } from "../ui/button";
// import {
//   DropdownMenu,
//   DropdownMenuContent,
//   DropdownMenuItem,
//   DropdownMenuTrigger,
// } from "../ui/dropdown-menu";
import MainNavbar from "../MainNavbar";

const MainLayout = () => {
  return (
    <div className="flex flex-col min-h-screen">
      {/* Navbar */}
      <MainNavbar />

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
