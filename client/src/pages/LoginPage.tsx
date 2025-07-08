import { useEffect, useState } from "react";

const images = [
  "https://4kwallpapers.com/images/walls/thumbs_3t/22156.jpg",
  "https://4kwallpapers.com/images/walls/thumbs_3t/19468.jpg",
  "https://4kwallpapers.com/images/walls/thumbs_3t/19461.jpg",
  "https://4kwallpapers.com/images/walls/thumbs_3t/19165.jpg",
];

const LoginPage = () => {
  const [activeTab, setActiveTab] = useState("login");
  const [currentImage, setCurrentImage] = useState(0);

  useEffect(() => {
    const interval = setInterval(() => {
      setCurrentImage((prev) => (prev + 1) % images.length);
    }, 1000); // ganti setiap 5 detik
    return () => clearInterval(interval);
  }, []);

  return (
    <div className="flex items-center justify-center min-h-screen min-w-full bg-gray-50">
      <div className="fixed z-50 top-[8%] text-4xl md:text-6xl font-bold text-gray-800 typing-animation">
        Mari mencurahkan ide bersama!
      </div>
      <div className="flex justify-evenly items-center w-screen gap-7 ">
        {/* KIRI */}
        <div className="  bg-white rounded-2xl shadow-lg p-8 w-1/2 max-w-md ">
          {/* Tab */}
          <div className="flex justify-evenly mb-6 ">
            <button
              id="login"
              className={`px-4 py-2 font-semibold ${
                activeTab === "login"
                  ? "border-b-2 border-blue-500 text-blue-600"
                  : "text-gray-500"
              }`}
              onClick={() => setActiveTab("login")}
            >
              Login
            </button>
            <button
              id="register"
              className={`ml-4 px-4 py-2 font-semibold ${
                activeTab === "signup"
                  ? "border-b-2 border-blue-500 text-blue-600"
                  : "text-gray-500"
              }`}
              onClick={() => setActiveTab("signup")}
            >
              Sign Up
            </button>
          </div>

          {/* Form */}
          {activeTab === "login" ? (
            <form className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700">
                  Email
                </label>
                <input
                  type="email"
                  className="mt-1 w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-400"
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700">
                  Password
                </label>
                <input
                  type="password"
                  className="mt-1 w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-400"
                />
              </div>
              <button
                type="submit"
                className="w-full bg-blue-600 text-white py-2 rounded-md hover:bg-blue-700"
              >
                Login
              </button>
            </form>
          ) : (
            <form className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700">
                  Name
                </label>
                <input
                  type="text"
                  className="mt-1 w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-400"
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700">
                  Email
                </label>
                <input
                  type="email"
                  className="mt-1 w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-400"
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700">
                  Password
                </label>
                <input
                  type="password"
                  className="mt-1 w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-400"
                />
              </div>
              <button
                type="submit"
                className="w-full bg-blue-600 text-white py-2 rounded-md hover:bg-blue-700"
              >
                Sign Up
              </button>
            </form>
          )}
        </div>

        {/* KANAN */}
        <div className="w-1/2 h-full relative">
          <img
            src={images[currentImage]}
            alt="Slide"
            className="w-full h-full object-cover rounded-lg"
          />
        </div>
      </div>
    </div>
  );
};

export default LoginPage;
