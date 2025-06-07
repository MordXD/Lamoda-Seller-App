import { createContext, useContext, useState, useEffect } from 'react';
import type { ReactNode } from 'react';
import { v4 as uuidv4 } from 'uuid';
import apiClient from '../api/axios';
import type { Account, AuthContextType, LoginResponse } from '../types/auth';

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};

interface AuthProviderProps {
  children: ReactNode;
}

export const AuthProvider = ({ children }: AuthProviderProps) => {
  const [accounts, setAccounts] = useState<Account[]>([]);
  const [activeAccount, setActiveAccount] = useState<Account | null>(null);
  const [isLoading, setIsLoading] = useState(false);

  // Загрузка данных из localStorage при инициализации
  useEffect(() => {
    const savedAccounts = localStorage.getItem('accounts');
    const activeAccountId = localStorage.getItem('active_account_id');
    
    if (savedAccounts) {
      const parsedAccounts: Account[] = JSON.parse(savedAccounts);
      setAccounts(parsedAccounts);
      
      if (activeAccountId) {
        const active = parsedAccounts.find(acc => acc.id === activeAccountId);
        if (active) {
          setActiveAccount(active);
          localStorage.setItem('active_token', active.token);
        }
      }
    }
  }, []);

  const login = async (email: string, password: string): Promise<void> => {
    setIsLoading(true);
    try {
      const response = await apiClient.post<LoginResponse>('/api/auth/login', {
        email,
        password,
      });
      
      const { accessToken, user } = response.data;
      
      // Создаем новый аккаунт
      const newAccount: Account = {
        id: uuidv4(),
        shopName: user.shopName,
        token: accessToken,
        user,
      };
      
      // Проверяем, есть ли уже такой аккаунт
      const existingAccountIndex = accounts.findIndex(
        acc => acc.user?.email === user.email
      );
      
      let updatedAccounts: Account[];
      if (existingAccountIndex !== -1) {
        // Обновляем существующий аккаунт
        updatedAccounts = [...accounts];
        updatedAccounts[existingAccountIndex] = newAccount;
      } else {
        // Добавляем новый аккаунт
        updatedAccounts = [...accounts, newAccount];
      }
      
      // Сохраняем в localStorage
      localStorage.setItem('accounts', JSON.stringify(updatedAccounts));
      localStorage.setItem('active_account_id', newAccount.id);
      localStorage.setItem('active_token', newAccount.token);
      
      // Обновляем состояние
      setAccounts(updatedAccounts);
      setActiveAccount(newAccount);
      
    } catch (error) {
      console.error('Login error:', error);
      throw error;
    } finally {
      setIsLoading(false);
    }
  };

  const logout = () => {
    // Удаляем активный аккаунт и токен
    localStorage.removeItem('active_account_id');
    localStorage.removeItem('active_token');
    setActiveAccount(null);
  };

  const switchAccount = (accountId: string) => {
    const account = accounts.find(acc => acc.id === accountId);
    if (account) {
      localStorage.setItem('active_account_id', account.id);
      localStorage.setItem('active_token', account.token);
      // Перезагружаем страницу для обновления всех данных
      window.location.reload();
    }
  };

  const addAccount = (account: Account) => {
    const updatedAccounts = [...accounts, account];
    setAccounts(updatedAccounts);
    localStorage.setItem('accounts', JSON.stringify(updatedAccounts));
  };

  const value: AuthContextType = {
    accounts,
    activeAccount,
    isLoading,
    login,
    logout,
    switchAccount,
    addAccount,
  };

  return (
    <AuthContext.Provider value={value}>
      {children}
    </AuthContext.Provider>
  );
}; 