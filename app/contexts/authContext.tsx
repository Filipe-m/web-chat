import {
  useContext,
  createContext,
  type PropsWithChildren,
  useEffect,
} from "react";
import { useStorageState } from "@/hooks/useStorageState";
import { useUser } from "./userContext";
import AsyncStorage from "@react-native-async-storage/async-storage";

const AuthContext = createContext<{
  signIn: (email: string, password: string) => void;
  signOut: () => void;
  session?: string | null;
  isLoading: boolean;
}>({
  signIn: (email, password) => null,
  signOut: () => null,
  session: null,
  isLoading: false,
});

export function useSession() {
  const value = useContext(AuthContext);
  if (process.env.NODE_ENV !== "production") {
    if (!value) {
      throw new Error("useSession must be wrapped in a <SessionProvider />");
    }
  }

  return value;
}

export function SessionProvider({ children }: PropsWithChildren) {
  const [[isLoading, session], setSession] = useStorageState("session");
  const { extractUserInfo } = useUser();

  useEffect(() => {
    if (session) {
      extractUserInfo(session);
    }
  }, []);

  async function signIn(email: string, password: string) {
    const payload = {
      email: email,
      password: password,
    };

    const url: string | undefined = process.env.EXPO_PUBLIC_API_URL;

    try {
      const response = await fetch(url + "login", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(payload),
      });

      if (response.status == 404) {
        alert("Email ou senha inválidos");
      }

      const json = await response.json();

      setSession(json.token);
      extractUserInfo(json.token);
    } catch (error) {
      console.error(error);
      alert("Credenciais inválidas!");
    }
  }

  const signOut = () => {
    setSession(null);
  };

  return (
    <AuthContext.Provider
      value={{
        signIn,
        signOut,
        session,
        isLoading,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
}
