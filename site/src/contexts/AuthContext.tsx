
import { createContext, useContext, useState, ReactNode, useEffect } from 'react';
import { UserData } from '@/types/user';
import { loginUser, registerUser } from '@/lib/userApi';
import { toast } from "@/components/ui/sonner";

interface AuthContextType {
  currentUser: UserData | null;
  isAuthenticated: boolean;
  login: (login: string, password: string) => Promise<boolean>;
  register: (name: string, login: string, password: string) => Promise<boolean>;
  logout: () => void;
  updateUserData: (userData: UserData) => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [currentUser, setCurrentUser] = useState<UserData | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetch('http://localhost:8002/auth/me', {
      credentials: 'include',
    })
      .then(res => {
        if (!res.ok) throw new Error('not authenticated');
        return res.json();
      })
      .then((user: UserData) => setCurrentUser(user))
      .catch(() => setCurrentUser(null))
      .finally(() => setLoading(false));
  }, []);

  const login = async (login: string, password: string) => {
    try {
      const userData = await loginUser({ login, password });
      setCurrentUser(userData);
      toast.success("Logged in successfully");
      return true;
    } catch (error) {
      console.error("Login error:", error);
      toast.error("Invalid login or password");
      return false;
    }
  };

  const register = async (name: string, login: string, password: string) => {
    try {
      await registerUser({ login, password, name });
      toast.success("Registration successful");
      return true;
    } catch (error) {
      console.error("Registration error:", error);
      toast.error("Registration failed");
      return false;
    }
  };

  const logout = () => {
    setCurrentUser(null);
    toast.success("Logged out successfully");
  };

  const updateUserData = (userData: UserData) => {
    setCurrentUser(userData);
  };
  
  if (currentUser && loading) {
    return <div>Loading…</div>;
  }

  return (
    <AuthContext.Provider
      value={{
        currentUser,
        isAuthenticated: !!currentUser,
        login,
        register,
        logout,
        updateUserData
      }}
    >
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
}
