import React from "react";
import { motion } from "framer-motion";
import { Link } from "react-router-dom";

const NotFoundPage = () => {
  return (
    <div className="relative min-h-screen bg-black overflow-hidden flex flex-col items-center justify-center text-white px-4">
      {/* Glowing Light Effect */}
      <div className="absolute inset-0 z-0">
        <div className="absolute top-0 left-1/2 transform -translate-x-1/2 w-[500px] h-[500px] bg-gradient-to-br from-purple-500 via-pink-500 to-green-400 opacity-30 blur-3xl rounded-full animate-pulse" />
      </div>

      {/* Content */}
      <motion.div
        className="z-10 text-center"
        initial={{ opacity: 0, scale: 0.95 }}
        animate={{ opacity: 1, scale: 1 }}
        transition={{ duration: 0.8 }}
      >
        <h1 className="text-[120px] font-extrabold leading-none text-white mb-6 drop-shadow-[0_0_20px_rgba(255,255,255,0.2)]">
          404
        </h1>
        <p className="text-lg text-gray-300 mb-8">
          The page you are looking for doesn’t exist or has been moved.
          <br />
          Please go back to the homepage.
        </p>
        <Link
          to="/"
          className="inline-block px-6 py-3 text-sm font-semibold rounded-full bg-white text-black hover:bg-gray-100 transition"
        >
          Go back home
        </Link>
      </motion.div>

      {/* Footer Info (Optional) */}
      <div className="absolute bottom-4 left-4 text-xs text-gray-500">
        © 2025
      </div>
      <div className="absolute bottom-4 right-4 text-xs text-gray-500">
        Made in React & Framer Motion
      </div>
    </div>
  );
};

export default NotFoundPage;
