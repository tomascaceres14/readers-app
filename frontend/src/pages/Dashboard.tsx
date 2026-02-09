// src/pages/Dashboard.tsx

/**
 * DASHBOARD - PÁGINA PRINCIPAL (PROTEGIDA)
 *
 * Esta página solo se muestra a usuarios autenticados.
 * Aquí irá tu aplicación principal.
 *
 * Por ahora es un placeholder simple que muestra:
 * - Información del usuario
 * - Botón de logout
 */

import { useAuth } from "../features/auth";

function Dashboard() {
  const { user, logout } = useAuth();

  const handleLogout = async () => {
    try {
      await logout();
    } catch (error) {
      console.error("Error al cerrar sesión:", error);
    }
  };

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <header className="bg-white shadow">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4 flex justify-between items-center">
          <h1 className="text-2xl font-bold text-gray-900">Readers App</h1>
          <button
            onClick={handleLogout}
            className="px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700 transition-colors"
          >
            Cerrar Sesión
          </button>
        </div>
      </header>

      {/* Main Content */}
      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="bg-white rounded-lg shadow p-6">
          <h2 className="text-xl font-semibold mb-4">¡Bienvenido! 👋</h2>

          <div className="space-y-2 text-gray-600">
            <p>
              <span className="font-medium">Email:</span> {user?.email}
            </p>
            <p>
              <span className="font-medium">Nombre:</span>{" "}
              {user?.displayName || "No definido"}
            </p>
            <p>
              <span className="font-medium">UID:</span> {user?.uid}
            </p>
            <p>
              <span className="font-medium">Email verificado:</span>{" "}
              {user?.emailVerified ? "✅ Sí" : "❌ No"}
            </p>
          </div>

          <div className="mt-6 p-4 bg-blue-50 rounded-lg">
            <p className="text-sm text-blue-800">
              🎉 <strong>¡Autenticación funcionando!</strong> Aquí construirás
              tu aplicación principal.
            </p>
          </div>
        </div>
      </main>
    </div>
  );
}

export default Dashboard;
