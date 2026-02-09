import { AuthProvider } from "./features/auth";
import AppRoutes from "./AppRoutes.tsx";

function App() {
  return (
    <AuthProvider>
      <AppRoutes />
    </AuthProvider>
  );
}

export default App;
