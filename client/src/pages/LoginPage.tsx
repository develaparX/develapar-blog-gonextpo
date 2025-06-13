import React, { useState } from "react";
import { motion } from "framer-motion";
import { Eye, EyeOff } from "lucide-react";

const LoginPage = () => {
  const [showPassword, setShowPassword] = useState(false);

  return (
    <div className="relative min-h-screen bg-black text-white flex items-center justify-center overflow-hidden px-4">
      {/* Background Glow */}
      <div className="absolute inset-0 z-0">
        <div className="absolute top-1/3 left-1/2 -translate-x-1/2 w-[600px] h-[600px] bg-gradient-to-br from-purple-600 via-indigo-500 to-teal-400 opacity-30 blur-3xl rounded-full animate-pulse" />
      </div>

      {/* Login Card */}
      <motion.div
        className="z-10 w-full max-w-md bg-white/5 backdrop-blur-md rounded-2xl p-8 shadow-xl"
        initial={{ y: 30, opacity: 0 }}
        animate={{ y: 0, opacity: 1 }}
        transition={{ duration: 0.8, ease: "easeOut" }}
      >
        <h2 className="text-3xl font-bold text-center mb-8">Welcome Back 👋</h2>

        <form className="space-y-6">
          {/* Email Field */}
          <div className="relative">
            <input
              type="email"
              placeholder="Email"
              className="w-full px-4 py-3 bg-transparent border border-gray-700 rounded-lg text-white focus:outline-none focus:border-purple-500 transition"
              required
            />
          </div>

          {/* Password Field */}
          <div className="relative">
            <input
              type={showPassword ? "text" : "password"}
              placeholder="Password"
              className="w-full px-4 py-3 bg-transparent border border-gray-700 rounded-lg text-white focus:outline-none focus:border-purple-500 transition"
              required
            />
            <div
              onClick={() => setShowPassword(!showPassword)}
              className="absolute right-4 top-1/2 -translate-y-1/2 text-gray-400 hover:text-white cursor-pointer"
            >
              {showPassword ? <EyeOff size={20} /> : <Eye size={20} />}
            </div>
          </div>

          {/* Button */}
          <motion.button
            whileHover={{ scale: 1.05 }}
            whileTap={{ scale: 0.98 }}
            type="submit"
            className="w-full py-3 bg-purple-600 hover:bg-purple-700 transition text-white font-semibold rounded-lg"
          >
            Login
          </motion.button>
        </form>

        <div className="mt-6 text-center text-sm text-gray-400">
          Don't have an account?{" "}
          <span className="text-purple-400 hover:underline cursor-pointer">
            Sign up
          </span>
        </div>
      </motion.div>
    </div>
  );
};

export default LoginPage;
