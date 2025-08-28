import { Outlet } from "react-router-dom";
import Sidebar from "../dashboard/Sidebar";
import { useState } from "react";
import Navbar from "../dashboard/Navbar";

const AdminLayout = () => {
  const [isSidebarOpen, setIsSidebarOpen] = useState<boolean>(false);
  const [isSidebarMinimized, setIsSidebarMinimized] = useState<boolean>(false);

  const toggleSidebar = (): void => {
    setIsSidebarOpen(!isSidebarOpen);
  };

  return (
    <div className="h-screen overflow-hidden  flex">
      {/* Sidebar */}
      <Sidebar
        isOpen={isSidebarOpen}
        setIsOpen={setIsSidebarOpen}
        isMinimized={isSidebarMinimized}
        setIsMinimized={setIsSidebarMinimized}
      />

      {/* Content */}

      {/* Main content from routes */}
      <div className="flex-1 flex flex-col overflow-hidden">
        {/* Navbar */}
        <Navbar onMenuClick={toggleSidebar} />

        {/* Main content from routes */}
        <main className="flex-1 overflow-y-auto p-6">
          <Outlet />
        </main>
      </div>
    </div>
  );
};

export default AdminLayout;
