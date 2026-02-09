// src/features/auth/components/LoginForm.tsx

/**
 * LOGIN FORM - FORMULARIO DE INICIO DE SESIÓN
 *
 * Usa TanStack Form para manejo de formularios con:
 * - Validación en tiempo real
 * - Type safety completo
 * - Mensajes de error personalizados
 *
 * Campos:
 * - Email: Requerido, formato email válido
 * - Password: Requerido, mínimo 6 caracteres
 *
 * Flujo:
 * 1. Usuario llena el formulario
 * 2. Validación en tiempo real mientras escribe
 * 3. Al enviar → llama a login() del AuthContext
 * 4. Si éxito → AuthContext actualiza estado → AppRoutes redirige a Dashboard
 * 5. Si error → muestra mensaje de error
 */

import { useForm } from "@tanstack/react-form";
import { useState } from "react";
import { useAuth } from "../hooks/useAuth";
import type { LoginCredentials } from "../types";

interface LoginFormProps {
  onToggleForm: () => void;
}

function LoginForm({ onToggleForm }: LoginFormProps) {
  const { login } = useAuth();
  const [error, setError] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(false);

  const form = useForm({
    defaultValues: {
      email: "",
      password: "",
    } as LoginCredentials,
    onSubmit: async ({ value }) => {
      setError(null);
      setIsLoading(true);

      try {
        await login(value);
      } catch (err: any) {
        const errorMessage = getFirebaseErrorMessage(err.code);
        setError(errorMessage);
      } finally {
        setIsLoading(false);
      }
    },
  });

  return (
    <div className="bg-red-500 w-full max-w-md">
      <div className="bg-red-500 rounded-lg shadow-lg p-8 mx-auto">
        <h2 className="text-2xl font-bold text-gray-900 mb-6 text-center">
          Iniciar Sesión
        </h2>

        {error && (
          <div className="mb-4 p-3 bg-red-50 border border-red-200 rounded-lg">
            <p className="text-sm text-red-600">{error}</p>
          </div>
        )}

        <form
          onSubmit={(e) => {
            e.preventDefault();
            e.stopPropagation();
            form.handleSubmit();
          }}
          className="space-y-4"
        >
          <form.Field
            name="email"
            validators={{
              onChange: ({ value }) => {
                if (!value) return "El email es requerido";
                if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(value)) {
                  return "Ingresa un email válido";
                }
                return undefined;
              },
            }}
            children={(field) => (
              <div>
                <label
                  htmlFor="email"
                  className="block text-sm font-medium text-gray-700 mb-1"
                >
                  Email
                </label>
                <input
                  id="email"
                  type="email"
                  value={field.state.value}
                  onBlur={field.handleBlur}
                  onChange={(e) => field.handleChange(e.target.value)}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  placeholder="tu@email.com"
                />
                {field.state.meta.errors.length > 0 && (
                  <p className="mt-1 text-sm text-red-600">
                    {field.state.meta.errors[0]}
                  </p>
                )}
              </div>
            )}
          />

          <form.Field
            name="password"
            validators={{
              onChange: ({ value }) => {
                if (!value) return "La contraseña es requerida";
                if (value.length < 6) {
                  return "La contraseña debe tener al menos 6 caracteres";
                }
                return undefined;
              },
            }}
            children={(field) => (
              <div>
                <label
                  htmlFor="password"
                  className="block text-sm font-medium text-gray-700 mb-1"
                >
                  Contraseña
                </label>
                <input
                  id="password"
                  type="password"
                  value={field.state.value}
                  onBlur={field.handleBlur}
                  onChange={(e) => field.handleChange(e.target.value)}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  placeholder="••••••••"
                />
                {field.state.meta.errors.length > 0 && (
                  <p className="mt-1 text-sm text-red-600">
                    {field.state.meta.errors[0]}
                  </p>
                )}
              </div>
            )}
          />

          <button
            type="submit"
            disabled={isLoading}
            className="w-full py-2 px-4 bg-blue-600 hover:bg-blue-700 disabled:bg-blue-300 text-white font-medium rounded-lg transition-colors"
          >
            {isLoading ? "Iniciando sesión..." : "Iniciar Sesión"}
          </button>
        </form>

        <div className="mt-4 text-center">
          <button
            onClick={onToggleForm}
            className="text-sm text-blue-600 hover:text-blue-700"
          >
            ¿No tienes cuenta? <span className="font-medium">Regístrate</span>
          </button>
        </div>
      </div>
    </div>
  );
}

function getFirebaseErrorMessage(errorCode: string): string {
  switch (errorCode) {
    case "auth/user-not-found":
      return "No existe una cuenta con este email";
    case "auth/wrong-password":
      return "Contraseña incorrecta";
    case "auth/invalid-email":
      return "Email inválido";
    case "auth/user-disabled":
      return "Esta cuenta ha sido deshabilitada";
    case "auth/too-many-requests":
      return "Demasiados intentos. Intenta más tarde";
    case "auth/network-request-failed":
      return "Error de conexión. Verifica tu internet";
    default:
      return "Error al iniciar sesión. Intenta nuevamente";
  }
}

export default LoginForm;
