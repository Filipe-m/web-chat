import React, {
  createContext,
  useState,
  useContext,
  ReactNode,
  useEffect,
} from "react";
import AsyncStorage from "@react-native-async-storage/async-storage";

type User = {
  id: number | null;
};

type UserContextType = {
  user: User;
  setUser: (user: User) => void;
  clearUser: () => void;
  extractUserInfo: (token: string) => void;
};

const UserContext = createContext<UserContextType | undefined>(undefined);

export function useUser() {
  const context = useContext(UserContext);
  if (!context) {
    throw new Error("useUser must be used within a UserProvider");
  }
  return context;
}

export const UserProvider = ({ children }: { children: ReactNode }) => {
  const [user, setUser] = useState<User>({ id: null });

  useEffect(() => {
    const loadUser = async () => {
      const storedUser = await AsyncStorage.getItem("user");
      if (storedUser) {
        setUser(JSON.parse(storedUser));
      }
    };

    loadUser();
  }, []);

  const extractUserInfo = (token: string) => {
    const parseJwt = (token: string) => {
      try {
        return JSON.parse(atob(token.split(".")[1]));
      } catch (error) {
        console.log(error);
        return null;
      }
    };

    const decoded = parseJwt(token);
    const userId = decoded?.id || null;
    console.log(userId);

    setUser({ id: userId });
    AsyncStorage.setItem("user", JSON.stringify({ id: userId }));
  };

  const clearUser = () => {
    setUser({ id: null });
    AsyncStorage.removeItem("user");
  };

  return (
    <UserContext.Provider value={{ user, setUser, clearUser, extractUserInfo }}>
      {children}
    </UserContext.Provider>
  );
};
