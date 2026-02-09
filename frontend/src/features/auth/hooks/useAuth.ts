import { useContext } from "react";
import { AuthContext } from "../context/AuthContext";
import type { AuthContextType } from "../types";

export const useAuth = (): AuthContextType => {
  const context = useContext(AuthContext);

  // Si alguien intenta usar useAuth() fuera del Provider, lanzamos error
  // Esto ayuda a detectar bugs temprano en desarrollo
  if (context === undefined) {
    throw new Error("useAuth debe usarse dentro de un AuthProvider");
  }

  return context;
};
