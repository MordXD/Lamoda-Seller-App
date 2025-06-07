import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';

export default function LoginPage() {
  const navigate = useNavigate();
  const { login, isLoading, activeAccount, isInitialized } = useAuth();
  
  const [usernameOrEmail, setUsernameOrEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');

  // Перенаправляем на дашборд, если пользователь уже авторизован
  useEffect(() => {
    if (isInitialized && activeAccount) {
      navigate('/dashboard');
    }
  }, [isInitialized, activeAccount, navigate]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');

    if (!usernameOrEmail.trim() || !password.trim()) {
      setError('Пожалуйста, заполните все поля');
      return;
    }

    try {
      await login(usernameOrEmail.trim(), password);
      navigate('/dashboard');
    } catch (error: any) {
      console.error('Login failed:', error);
      setError(error.response?.data?.message || error.message || 'Неверные данные для входа');
    }
  };

  // Показываем загрузку пока инициализируется контекст
  if (!isInitialized) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-gray-900"></div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50 flex items-center justify-center py-12 px-4 font-montserrat-regular">
      <div className="max-w-md w-full bg-white rounded-2xl shadow-sm p-8">
        <div className="text-center">
          <img
            className="h-[40px] w-[220px] mx-auto mb-6"
            src="/icons/lamoda-icon.svg"
            alt="Lamoda"
          />
          <h2 className="text-3xl font-bold text-gray-900 mb-2">
            Вход в систему
          </h2>
          <p className="text-gray-600 mb-8">
            Войдите в свой аккаунт продавца
          </p>
        </div>

        <form onSubmit={handleSubmit} className="space-y-6">
          {/* Username or Email */}
          <div>
            <label htmlFor="usernameOrEmail" className="block text-sm text-gray-700 mb-2 font-montserrat-medium">
              Email или имя пользователя <span className="text-red-500">*</span>
            </label>
            <input
              id="usernameOrEmail"
              type="text"
              required
              value={usernameOrEmail}
              onChange={(e) => setUsernameOrEmail(e.target.value)}
              className="w-full px-4 py-3 text-lg bg-white border-2 border-gray-200 rounded-lg focus:border-gray-400 focus:outline-none focus:ring-0 placeholder-gray-400 font-montserrat-regular"
              placeholder="Введите email или имя пользователя"
              disabled={isLoading}
            />
          </div>

          {/* Password */}
          <div>
            <label htmlFor="password" className="block text-sm text-gray-700 mb-2 font-montserrat-medium">
              Пароль <span className="text-red-500">*</span>
            </label>
            <input
              id="password"
              type="password"
              required
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              className="w-full px-4 py-3 text-lg bg-white border-2 border-gray-200 rounded-lg focus:border-gray-400 focus:outline-none focus:ring-0 placeholder-gray-400 font-montserrat-regular"
              placeholder="Введите пароль"
              disabled={isLoading}
            />
          </div>

          {/* Error message */}
          {error && (
            <div className="p-4 bg-red-50 border border-red-200 rounded-lg">
              <p className="text-red-600 text-sm">{error}</p>
            </div>
          )}

          {/* Submit button */}
          <button
            type="submit"
            disabled={isLoading}
            className="w-full bg-black text-white py-4 px-6 text-lg font-montserrat-semibold rounded-lg hover:bg-gray-800 focus:outline-none focus:ring-2 focus:ring-gray-500 focus:ring-offset-2 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {isLoading ? 'Вход...' : 'Войти'}
          </button>

          <div className="text-center pt-4">
            <p className="text-sm text-gray-600">
              Нет аккаунта?{' '}
              <button
                type="button"
                onClick={() => navigate('/register')}
                className="text-black font-medium hover:underline"
              >
                Зарегистрироваться
              </button>
            </p>
          </div>
        </form>
      </div>
    </div>
  );
} 