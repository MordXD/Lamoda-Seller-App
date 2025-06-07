import { createContext, useContext, useState, useEffect } from 'react';
import type { ReactNode } from 'react';
import { v4 as uuidv4 } from 'uuid';
import { login as apiLogin } from '../api/auth';
import type { Account, AuthContextType } from '../types/auth';

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
  const [isInitialized, setIsInitialized] = useState(false);

  // Загрузка данных из localStorage при инициализации
  useEffect(() => {
    console.log('Инициализация AuthContext...');
    
    const savedAccounts = localStorage.getItem('accounts');
    const activeAccountId = localStorage.getItem('active_account_id');
    
    console.log('Сохраненные аккаунты:', savedAccounts);
    console.log('ID активного аккаунта:', activeAccountId);
    
    if (savedAccounts) {
      try {
        const parsedAccounts: Account[] = JSON.parse(savedAccounts);
        setAccounts(parsedAccounts);
        console.log('Загружены аккаунты:', parsedAccounts);
        
        if (activeAccountId) {
          const active = parsedAccounts.find(acc => acc.id === activeAccountId);
          if (active) {
            setActiveAccount(active);
            localStorage.setItem('active_token', active.token);
            console.log('Восстановлен активный аккаунт:', active);
          } else {
            console.log('Активный аккаунт не найден среди сохраненных');
          }
        }
      } catch (error) {
        console.error('Ошибка при парсинге сохраненных аккаунтов:', error);
        localStorage.removeItem('accounts');
        localStorage.removeItem('active_account_id');
        localStorage.removeItem('active_token');
      }
    } else {
      console.log('Сохраненные аккаунты не найдены');
    }
    
    setIsInitialized(true);
    console.log('Инициализация AuthContext завершена');
  }, []);

  const login = async (usernameOrEmail: string, password: string): Promise<void> => {
    setIsLoading(true);
    try {
      console.log('Отправляем запрос на авторизацию:', { 
        email: usernameOrEmail,
        password_length: password.length,
        password_first_char: password.charAt(0),
        password_last_char: password.charAt(password.length - 1)
      });
      
      const loginData = { email: usernameOrEmail, password };
      console.log('Данные для отправки:', loginData);
      
      const response = await apiLogin(loginData);
      
      console.log('Ответ сервера:', response);
      
      const { token } = response;
      
      if (!token) {
        console.error('Токен отсутствует в ответе сервера:', response);
        throw new Error('Токен не получен от сервера');
      }
      
      console.log('Получен токен:', token);
      
      // Создаем новый аккаунт с минимальными данными (бэкенд возвращает только токен)
      const newAccount: Account = {
        id: uuidv4(), // Генерируем временный ID
        shopName: usernameOrEmail, // Используем email как имя магазина
        token: token,
        user: {
          id: uuidv4(), // Генерируем временный ID пользователя
          email: usernameOrEmail,
          shopName: usernameOrEmail,
        },
      };
      
      console.log('Создан новый аккаунт:', newAccount);
      
      // Проверяем, есть ли уже такой аккаунт (по email)
      const existingAccountIndex = accounts.findIndex(
        acc => acc.user?.email === usernameOrEmail
      );
      
      let updatedAccounts: Account[];
      if (existingAccountIndex !== -1) {
        // Обновляем существующий аккаунт
        updatedAccounts = [...accounts];
        updatedAccounts[existingAccountIndex] = newAccount;
        console.log('Обновляем существующий аккаунт');
      } else {
        // Добавляем новый аккаунт
        updatedAccounts = [...accounts, newAccount];
        console.log('Добавляем новый аккаунт');
      }
      
      // Сохраняем в localStorage
      localStorage.setItem('accounts', JSON.stringify(updatedAccounts));
      localStorage.setItem('active_account_id', newAccount.id);
      localStorage.setItem('active_token', newAccount.token);
      
      console.log('Данные сохранены в localStorage');
      
      // Обновляем состояние
      setAccounts(updatedAccounts);
      setActiveAccount(newAccount);
      
      console.log('Авторизация успешна');
      
    } catch (error) {
      console.error('Ошибка авторизации:', error);
      if (error && typeof error === 'object' && 'response' in error) {
        const axiosError = error as any;
        console.error('Детали ошибки:', axiosError.response?.data);
        console.error('Статус ошибки:', axiosError.response?.status);
        console.error('Заголовки ответа:', axiosError.response?.headers);
        console.error('Конфигурация запроса:', axiosError.config);
      }
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

  const logoutAll = () => {
    // Удаляем все данные аккаунтов
    localStorage.removeItem('accounts');
    localStorage.removeItem('active_account_id');
    localStorage.removeItem('active_token');
    setAccounts([]);
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
    isInitialized,
    login,
    logout,
    logoutAll,
    switchAccount,
    addAccount,
  };

  return (
    <AuthContext.Provider value={value}>
      {children}
    </AuthContext.Provider>
  );
}; 