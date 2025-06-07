import { useState, useEffect, useCallback } from 'react';
import { getProfile, updateProfile, changePassword, getBalance, addBalance, withdrawBalance } from '../api/profile';
import { getLinkedAccounts, switchAccount as apiSwitchAccount } from '../api/account';
import { useAuth } from '../context/AuthContext';
import type { User, UpdateProfileRequest, ChangePasswordRequest, AddBalanceRequest, WithdrawBalanceRequest } from '../api/profile';
import type { LinkedAccount } from '../api/account';

export const useProfile = () => {
  const { logout } = useAuth();
  const [user, setUser] = useState<User | null>(null);
  const [balance, setBalance] = useState<number>(0);
  const [linkedAccounts, setLinkedAccounts] = useState<LinkedAccount[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const loadProfile = useCallback(async () => {
    try {
      setIsLoading(true);
      setError(null);
      
      const userData = await getProfile();
      setUser(userData);
      setBalance(userData.balance_kopecks);
    } catch (err) {
      console.error('Ошибка загрузки профиля:', err);
      setError('Не удалось загрузить данные профиля');
    } finally {
      setIsLoading(false);
    }
  }, []);

  const loadBalance = useCallback(async () => {
    try {
      const balanceData = await getBalance();
      setBalance(balanceData.balance_kopecks);
    } catch (err) {
      console.error('Ошибка загрузки баланса:', err);
    }
  }, []);

  const loadLinkedAccounts = useCallback(async () => {
    try {
      const accounts = await getLinkedAccounts();
      setLinkedAccounts(accounts);
    } catch (err) {
      console.error('Ошибка загрузки связанных аккаунтов:', err);
    }
  }, []);

  const updateUserProfile = useCallback(async (data: UpdateProfileRequest) => {
    try {
      setError(null);
      await updateProfile(data);
      await loadProfile(); // Перезагружаем профиль после обновления
      return true;
    } catch (err) {
      console.error('Ошибка обновления профиля:', err);
      setError('Не удалось обновить профиль');
      return false;
    }
  }, [loadProfile]);

  const changeUserPassword = useCallback(async (data: ChangePasswordRequest) => {
    try {
      setError(null);
      await changePassword(data);
      return true;
    } catch (err) {
      console.error('Ошибка смены пароля:', err);
      setError('Не удалось изменить пароль');
      return false;
    }
  }, []);

  const addUserBalance = useCallback(async (data: AddBalanceRequest) => {
    try {
      setError(null);
      await addBalance(data);
      await loadBalance(); // Перезагружаем баланс после пополнения
      return true;
    } catch (err) {
      console.error('Ошибка пополнения баланса:', err);
      setError('Не удалось пополнить баланс');
      return false;
    }
  }, [loadBalance]);

  const withdrawUserBalance = useCallback(async (data: WithdrawBalanceRequest) => {
    try {
      setError(null);
      await withdrawBalance(data);
      await loadBalance(); // Перезагружаем баланс после снятия
      return true;
    } catch (err) {
      console.error('Ошибка снятия средств:', err);
      setError('Не удалось снять средства');
      return false;
    }
  }, [loadBalance]);

  const switchUserAccount = useCallback(async (targetUserId: string) => {
    try {
      setError(null);
      const response = await apiSwitchAccount({ target_user_id: targetUserId });
      
      // Обновляем токен в localStorage
      localStorage.setItem('active_token', response.token);
      
      // Перезагружаем страницу для обновления всех данных с новым токеном
      window.location.reload();
      
      return true;
    } catch (err) {
      console.error('Ошибка переключения аккаунта:', err);
      setError('Не удалось переключить аккаунт');
      return false;
    }
  }, []);

  const logoutUser = useCallback(() => {
    logout();
    // Перенаправление на страницу входа будет обработано в компоненте
  }, [logout]);

  useEffect(() => {
    loadProfile();
    loadLinkedAccounts();
  }, [loadProfile, loadLinkedAccounts]);

  return {
    user,
    balance,
    linkedAccounts,
    isLoading,
    error,
    loadProfile,
    loadBalance,
    loadLinkedAccounts,
    updateProfile: updateUserProfile,
    changePassword: changeUserPassword,
    addBalance: addUserBalance,
    withdrawBalance: withdrawUserBalance,
    switchAccount: switchUserAccount,
    logout: logoutUser,
  };
}; 