import { useState } from "react";
import LoginForm from "../features/auth/components/LoginForm";
import RegisterForm from "../features/auth/components/RegisterForm";
import GoogleLoginButton from "../features/auth/components/GoogleLoginButton";

function AuthPage() {
  const [isLogin, setIsLogin] = useState(true);

  const toggleForm = () => {
    setIsLogin(!isLogin);
  };

  return (
    <div className="min-h-screen from-blue-50 to-indigo-100 flex items-center justify-center p-4">
      <div className=" max-w-md">
        {/* Logo/Header */}
        <div className="text-center mb-8">
          <h1 className="text-4xl font-bold text-gray-900 mb-2">
            📚 Readers App
          </h1>
          <p className="text-gray-600">
            {isLogin
              ? "Inicia sesión para continuar"
              : "Crea tu cuenta gratuita"}
          </p>
        </div>

        {/* Google Login Button (siempre visible) */}
        <div className="mb-6">
          <GoogleLoginButton />
        </div>

        {/* Divisor */}
        <div className="relative mb-6">
          <div className="absolute inset-0 flex items-center">
            <div className="w-full border-t border-gray-300"></div>
          </div>
          <div className="relative flex justify-center text-sm">
            <span className="px-2 bg-gradient-to-br from-blue-50 to-indigo-100 text-gray-500">
              O continúa con email
            </span>
          </div>
        </div>

        {/* Formulario dinámico (Login o Register) */}
        {isLogin ? (
          <LoginForm onToggleForm={toggleForm} />
        ) : (
          <RegisterForm onToggleForm={toggleForm} />
        )}
      </div>
    </div>
  );
}

export default AuthPage;
