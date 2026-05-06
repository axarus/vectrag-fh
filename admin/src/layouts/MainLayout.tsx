import { Link, Outlet, useLocation } from 'react-router-dom';

export default function MainLayout() {
  const location = useLocation();

  const isActive = (path: string) => location.pathname === path;

  const menuItems = [
    { path: '/', label: 'Dashboard', icon: '📊' },
    { path: '/models', label: 'Models', icon: '📦' },
  ];

  return (
    <div className="min-h-screen bg-base-200">
      <div className="drawer lg:drawer-open">
        <input id="drawer-toggle" type="checkbox" className="drawer-toggle" />
        
        <div className="drawer-content flex flex-col">
          <div className="navbar bg-base-100 shadow-lg lg:hidden">
            <div className="flex-none">
              <label htmlFor="drawer-toggle" className="btn btn-square btn-ghost">
                <svg
                  className="w-6 h-6"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M4 6h16M4 12h16M4 18h16"
                  />
                </svg>
              </label>
            </div>
            <div className="flex-1">
              <a className="btn btn-ghost text-xl">VectraG</a>
            </div>
          </div>

          <div className="navbar bg-base-100 shadow-lg hidden lg:flex">
            <div className="flex-1">
              <a className="btn btn-ghost text-xl">VectraG</a>
            </div>
            <div className="flex-none">
              <div className="badge badge-primary">Development Mode</div>
            </div>
          </div>

          <main className="flex-1 p-6">
            <Outlet />
          </main>
        </div>

        <div className="drawer-side">
          <label htmlFor="drawer-toggle" className="drawer-overlay"></label>
          <aside className="w-64 min-h-full bg-base-100">
            <div className="p-4">
              <h2 className="text-2xl font-bold mb-6">VectraG</h2>
              <ul className="menu menu-vertical w-full">
                {menuItems.map((item) => (
                  <li key={item.path}>
                    <Link
                      to={item.path}
                      className={isActive(item.path) ? 'active' : ''}
                    >
                      <span className="text-xl mr-2">{item.icon}</span>
                      {item.label}
                    </Link>
                  </li>
                ))}
              </ul>
            </div>
          </aside>
        </div>
      </div>
    </div>
  );
}

