// src/features/auth/context/AuthContext.tsx

/**
 * AUTHCONTEXT - EL CEREBRO DE LA AUTENTICACIÓN
 *
 * Este archivo es el MÁS IMPORTANTE del sistema de auth.
 *
 * ¿Qué hace?
 * 1. Mantiene el estado global del usuario autenticado
 * 2. Provee funciones de login, registro, logout
 * 3. Escucha cambios en Firebase (si el usuario cierra sesión en otra pestaña)
 * 4. Proporciona el token JWT para tu backend
 *
 * PATRÓN: Context + Provider
 * - Context: Canal de comunicación global
 * - Provider: Componente que envuelve tu app y provee los datos
 * - Hook (useAuth): Forma fácil de consumir el context
 *
 * Flujo:
 * 1. Usuario hace login → llamamos a Firebase
 * 2. Firebase devuelve User → actualizamos estado
 * 3. Cualquier componente usa useAuth() → ve el usuario actualizado
 */

import React, {
  createContext,
  useEffect,
  useState,
  type ReactNode,
} from "react";
import {
  signInWithEmailAndPassword,
  createUserWithEmailAndPassword,
  signOut,
  onAuthStateChanged,
  GoogleAuthProvider,
  signInWithPopup,
  updateProfile,
} from "firebase/auth";
import { auth } from "../../../lib/firebase";

// Tipos se importan con 'type' (se eliminan en build)
import type {
  User,
  AuthContextType,
  LoginCredentials,
  RegisterCredentials,
} from "../types";

// Funciones/valores se importan normalmente (se mantienen en build)
import { mapFirebaseUser } from "../types";

// 1. Creamos el Context (el "canal de comunicación")
// undefined porque inicialmente no tiene valor hasta que el Provider lo provee
const AuthContext = createContext<AuthContextType | undefined>(undefined);

// 2. Props del Provider (solo necesita children)
interface AuthProviderProps {
  children: ReactNode;
}

// 3. AuthProvider - Componente que envuelve tu app
export const AuthProvider = ({ children }: AuthProviderProps) => {
  // Estados locales del Provider
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true); // true mientras verifica sesión
  const [token, setToken] = useState<string | null>(null);

  // EFECTO: Escucha cambios de autenticación en Firebase
  // onAuthStateChanged se ejecuta:
  // - Cuando la app carga (verifica si hay sesión activa)
  // - Cuando el usuario hace login/logout
  // - Cuando cambia en otra pestaña
  useEffect(() => {
    const unsubscribe = onAuthStateChanged(auth, async (firebaseUser) => {
      if (firebaseUser) {
        // Usuario autenticado
        const mappedUser = mapFirebaseUser(firebaseUser);
        setUser(mappedUser);

        // Obtenemos el token JWT
        const idToken = await firebaseUser.getIdToken();
        setToken(idToken);
      } else {
        // No hay usuario autenticado
        setUser(null);
        setToken(null);
      }
      setLoading(false); // Terminamos de cargar
    });

    // Cleanup: cuando el componente se desmonta, dejamos de escuchar
    return unsubscribe;
  }, []);

  // FUNCIÓN: Login con email y password
  const login = async (credentials: LoginCredentials): Promise<User> => {
    const { email, password } = credentials;
    const userCredential = await signInWithEmailAndPassword(
      auth,
      email,
      password,
    );
    const mappedUser = mapFirebaseUser(userCredential.user);
    return mappedUser;
  };

  // FUNCIÓN: Registro con email y password
  const register = async (credentials: RegisterCredentials): Promise<User> => {
    const { email, password, displayName } = credentials;

    // Creamos el usuario en Firebase
    const userCredential = await createUserWithEmailAndPassword(
      auth,
      email,
      password,
    );

    // Si se proporcionó displayName, actualizamos el perfil
    if (displayName && userCredential.user) {
      await updateProfile(userCredential.user, { displayName });
    }

    const mappedUser = mapFirebaseUser(userCredential.user);
    return mappedUser;
  };

  // FUNCIÓN: Login con Google
  const loginWithGoogle = async (): Promise<User> => {
    const provider = new GoogleAuthProvider();
    // signInWithPopup abre una ventana popup de Google
    const userCredential = await signInWithPopup(auth, provider);
    const mappedUser = mapFirebaseUser(userCredential.user);
    return mappedUser;
  };

  // FUNCIÓN: Logout
  const logout = async (): Promise<void> => {
    await signOut(auth);
    // onAuthStateChanged detectará el cambio y actualizará el estado
  };

  // FUNCIÓN: Obtener token actualizado
  // Úsala antes de cada request a tu backend
  const getToken = async (): Promise<string | null> => {
    if (!auth.currentUser) return null;

    // getIdToken(true) fuerza refresh si expiró
    const idToken = await auth.currentUser.getIdToken(true);
    setToken(idToken);
    return idToken;
  };

  // Valor que se provee a todos los componentes hijos
  const value: AuthContextType = {
    user,
    loading,
    token,
    login,
    register,
    loginWithGoogle,
    logout,
    getToken,
  };

  // Mientras carga, puedes mostrar un spinner
  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="text-xl">Cargando...</div>
      </div>
    );
  }

  // Provider envuelve a children y les da acceso al value
  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};

// Exportamos también el Context para que el hook pueda usarlo
export { AuthContext };
