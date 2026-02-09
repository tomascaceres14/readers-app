// src/features/auth/components/RegisterForm.tsx

/**
 * REGISTER FORM - FORMULARIO DE REGISTRO
 *
 * Similar a LoginForm pero con campos adicionales:
 * - Email
 * - Password
 * - Confirm Password (debe coincidir)
 * - Display Name (opcional)
 *
 * Validaciones:
 * - Email válido
 * - Password mínimo 6 caracteres
 * - Passwords deben coincidir
 */

import { useForm } from "@tanstack/react-form";
import { useState } from "react";
import { useAuth } from "../hooks/useAuth";

interface RegisterFormValues {
  email: string;
  password: string;
  confirmPassword: string;
  displayName: string;
}

interface RegisterFormProps {
  onToggleForm: () => void;
}

function RegisterForm({ onToggleForm }: RegisterFormProps) {
  const { register } = useAuth();
  const [error, setError] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(false);

  const form = useForm({
    defaultValues: {
      email: "",
      password: "",
      confirmPassword: "",
      displayName: "",
    } as RegisterFormValues,
    onSubmit: async ({ value }) => {
      setError(null);
      setIsLoading(true);

      try {
        await register({
          email: value.email,
          password: value.password,
          displayName: value.displayName || undefined,
        });
      } catch (err: any) {
        const errorMessage = getFirebaseErrorMessage(err.code);
        setError(errorMessage);
      } finally {
        setIsLoading(false);
      }
    },
  });

  return (
    <div className="w-full max-w-md">
      <div className="bg-white rounded-lg shadow-lg p-8">
        <h2 className="text-2xl font-bold text-gray-900 mb-6 text-center">
          Crear Cuenta
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
            name="displayName"
            children={(field) => (
              <div>
                <label
                  htmlFor="displayName"
                  className="block text-sm font-medium text-gray-700 mb-1"
                >
                  Nombre <span className="text-gray-400">(opcional)</span>
                </label>
                <input
                  id="displayName"
                  type="text"
                  value={field.state.value}
                  onBlur={field.handleBlur}
                  onChange={(e) => field.handleChange(e.target.value)}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  placeholder="Tu nombre"
                />
              </div>
            )}
          />

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

          <form.Field
            name="confirmPassword"
            validators={{
              onChangeListenTo: ["password"],
              onChange: ({ value, fieldApi }) => {
                const password = fieldApi.form.getFieldValue("password");
                if (!value) return "Confirma tu contraseña";
                if (value !== password) {
                  return "Las contraseñas no coinciden";
                }
                return undefined;
              },
            }}
            children={(field) => (
              <div>
                <label
                  htmlFor="confirmPassword"
                  className="block text-sm font-medium text-gray-700 mb-1"
                >
                  Confirmar Contraseña
                </label>
                <input
                  id="confirmPassword"
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
            {isLoading ? "Creando cuenta..." : "Crear Cuenta"}
          </button>
        </form>

        <div className="mt-4 text-center">
          <button
            onClick={onToggleForm}
            className="text-sm text-blue-600 hover:text-blue-700"
          >
            ¿Ya tienes cuenta?{" "}
            <span className="font-medium">Inicia sesión</span>
          </button>
        </div>
      </div>
    </div>
  );
}

function getFirebaseErrorMessage(errorCode: string): string {
  switch (errorCode) {
    case "auth/email-already-in-use":
      return "Este email ya está registrado";
    case "auth/invalid-email":
      return "Email inválido";
    case "auth/operation-not-allowed":
      return "Registro no permitido. Contacta soporte";
    case "auth/weak-password":
      return "Contraseña muy débil. Usa al menos 6 caracteres";
    case "auth/network-request-failed":
      return "Error de conexión. Verifica tu internet";
    default:
      return "Error al crear la cuenta. Intenta nuevamente";
  }
}

export default RegisterForm;
