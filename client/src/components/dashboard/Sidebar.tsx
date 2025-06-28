import { useState } from "react";
import {
  Home,
  Users,
  Settings,
  BarChart3,
  FileText,
  Menu,
  X,
  ChevronDown,
  ChevronRight,
  ChevronLeft,
  type LucideIcon,
} from "lucide-react";

interface SubmenuItem {
  title: string;
  path: string;
}

interface MenuItem {
  id: string;
  title: string;
  icon: LucideIcon;
  path: string;
  submenu?: SubmenuItem[];
}

interface SidebarProps {
  isOpen: boolean;
  setIsOpen: (isOpen: boolean) => void;
  isMinimized?: boolean;
  setIsMinimized?: (isMinimized: boolean) => void;
}

const Sidebar: React.FC<SidebarProps> = ({
  isOpen,
  setIsOpen,
  isMinimized = false,
  setIsMinimized,
}) => {
  const [activeMenu, setActiveMenu] = useState<string>("dashboard");
  const [expandedMenus, setExpandedMenus] = useState<Record<string, boolean>>(
    {}
  );

  const menuItems: MenuItem[] = [
    {
      id: "dashboard",
      title: "Dashboard",
      icon: Home,
      path: "/admin/dashboard",
    },
    {
      id: "users",
      title: "Users",
      icon: Users,
      path: "/admin/users",
      submenu: [
        { title: "All Users", path: "/admin/users" },
        { title: "Add User", path: "/admin/users/add" },
        { title: "User Roles", path: "/admin/users/roles" },
      ],
    },
    {
      id: "reports",
      title: "Reports",
      icon: BarChart3,
      path: "/admin/reports",
      submenu: [
        { title: "Analytics", path: "/admin/reports/analytics" },
        { title: "Sales Report", path: "/admin/reports/sales" },
        { title: "User Activity", path: "/admin/reports/activity" },
      ],
    },
    {
      id: "documents",
      title: "Documents",
      icon: FileText,
      path: "/admin/documents",
    },
    {
      id: "settings",
      title: "Settings",
      icon: Settings,
      path: "/admin/settings",
    },
  ];

  const handleToggleMinimize = () => {
    if (setIsMinimized) {
      setIsMinimized(!isMinimized);
      // Close expanded menus when minimizing
      if (!isMinimized) {
        setExpandedMenus({});
      }
    }
  };

  const toggleSubmenu = (menuId: string): void => {
    setExpandedMenus((prev) => ({
      ...prev,
      [menuId]: !(prev[menuId] ?? false),
    }));
  };

  const handleMenuClick = (menuId: string, hasSubmenu: boolean): void => {
    setActiveMenu(menuId);
    if (hasSubmenu && !isMinimized) {
      toggleSubmenu(menuId);
    }
  };

  return (
    <>
      {/* Mobile overlay */}
      {isOpen && (
        <div
          className="fixed inset-0 bg-black bg-opacity-50 z-40 lg:hidden"
          onClick={() => setIsOpen(false)}
        />
      )}

      {/* Sidebar */}
      <div
        className={`
        fixed inset-y-0 left-0 z-50 bg-white shadow-sm transform transition-all duration-300 ease-in-out
        lg:translate-x-0 lg:static lg:inset-0
        ${isOpen ? "translate-x-0" : "-translate-x-full"}
        ${isMinimized ? "w-16" : "w-64"}
      `}
      >
        {/* Header */}
        <div className="flex items-center justify-between h-16 px-4 border-b border-gray-200">
          {!isMinimized && (
            <h1 className="text-xl font-bold text-gray-800">develapar.</h1>
          )}

          <div className="flex items-center space-x-2">
            {/* Toggle minimize button - only show on desktop */}
            {setIsMinimized && (
              <button
                onClick={handleToggleMinimize}
                className="hidden lg:block p-2 rounded-md text-gray-400 hover:text-gray-600 hover:bg-gray-100"
                title={isMinimized ? "Expand sidebar" : "Minimize sidebar"}
              >
                <ChevronLeft
                  size={20}
                  className={`transform transition-transform ${
                    isMinimized ? "rotate-180" : ""
                  }`}
                />
              </button>
            )}

            {/* Mobile close button */}
            <button
              onClick={() => setIsOpen(false)}
              className="lg:hidden p-2 rounded-md text-gray-400 hover:text-gray-600 hover:bg-gray-100"
            >
              <X size={20} />
            </button>
          </div>
        </div>

        {/* Navigation */}
        <nav className="flex-1 px-4 py-4 space-y-1">
          {menuItems.map((item) => {
            const Icon = item.icon;
            const hasSubmenu = Boolean(item.submenu && item.submenu.length > 0);
            const isExpanded = Boolean(expandedMenus[item.id]);

            return (
              <div key={item.id}>
                {/* Main menu item */}
                <button
                  onClick={() => handleMenuClick(item.id, hasSubmenu)}
                  className={`
                    w-full flex items-center justify-between px-3 py-2 text-sm font-medium rounded-lg transition-colors
                    ${
                      activeMenu === item.id
                        ? "bg-blue-100 text-blue-700"
                        : "text-gray-700 hover:bg-gray-100 hover:text-gray-900"
                    }
                    ${isMinimized ? "justify-center" : ""}
                  `}
                  title={isMinimized ? item.title : ""}
                >
                  <div
                    className={`flex items-center ${
                      isMinimized ? "justify-center" : ""
                    }`}
                  >
                    <Icon size={18} className={isMinimized ? "" : "mr-3"} />
                    {!isMinimized && <span>{item.title}</span>}
                  </div>
                  {hasSubmenu && !isMinimized && (
                    <div className="ml-2">
                      {isExpanded ? (
                        <ChevronDown size={16} />
                      ) : (
                        <ChevronRight size={16} />
                      )}
                    </div>
                  )}
                </button>

                {/* Submenu */}
                {hasSubmenu && isExpanded && item.submenu && !isMinimized && (
                  <div className="ml-6 mt-1 space-y-1">
                    {item.submenu.map((subItem: SubmenuItem, index: number) => (
                      <a
                        key={index}
                        href={subItem.path}
                        className="block px-3 py-2 text-sm text-gray-600 rounded-lg hover:bg-gray-100 hover:text-gray-900 transition-colors"
                      >
                        {subItem.title}
                      </a>
                    ))}
                  </div>
                )}
              </div>
            );
          })}
        </nav>

        {/* Footer */}
        <div className="p-4 border-t border-gray-200">
          <div
            className={`flex items-center ${
              isMinimized ? "justify-center" : ""
            }`}
          >
            <div className="w-8 h-8 bg-blue-500 rounded-full flex items-center justify-center">
              <span className="text-white text-sm font-medium">A</span>
            </div>
            {!isMinimized && (
              <div className="ml-3 flex-1 min-w-0">
                <p className="text-sm font-medium text-gray-700 truncate">
                  Admin User
                </p>
                <p className="text-xs text-gray-500 truncate">
                  admin@example.com
                </p>
              </div>
            )}
          </div>
        </div>
      </div>
    </>
  );
};

export default Sidebar;
