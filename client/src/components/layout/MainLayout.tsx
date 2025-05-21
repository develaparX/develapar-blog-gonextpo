import { Outlet } from "react-router-dom";

const MainLayout = () => {
  return (
    <div>
      <header>Navbar atau header</header>
      <main>
        <Outlet />
      </main>
      <footer>Footer</footer>
    </div>
  );
};

export default MainLayout;
