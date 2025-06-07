export interface User {
  id: string;
  email: string;
  shopName: string;
}

export interface Account {
  id: string;
  shopName: string;
  token: string;
  user?: User;
}

// Ответ сервера (фактический)
export interface ServerLoginResponse {
  token: string;
}

// Ожидаемый ответ (для будущего использования)
export interface LoginResponse {
  accessToken: string;
  user: User;
}

export interface LoginData {
  username_or_email: string;
  password: string;
}

export interface AuthContextType {
  accounts: Account[];
  activeAccount: Account | null;
  isLoading: boolean;
  isInitialized: boolean;
  login: (usernameOrEmail: string, password: string) => Promise<void>;
  logout: () => void;
  logoutAll: () => void;
  switchAccount: (accountId: string) => void;
  addAccount: (account: Account) => void;
} 