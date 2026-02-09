import { useAuth } from "./features/auth";
import AuthPage from "./pages/AuthPage.tsx";
import Dashboard from "./pages/Dashboard";

function AppRoutes() {
  const { user } = useAuth();

  if (!user) {
    return <AuthPage />;
  }

  return <Dashboard />;
}

export default AppRoutes;
