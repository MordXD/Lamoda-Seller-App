import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';

export default function LoginPage() {
  const [login, setLogin] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const { login: authLogin, isLoading } = useAuth();
  const navigate = useNavigate();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');

    if (!login || !password) {
      setError('Пожалуйста, заполните все поля');
      return;
    }

    try {
      await authLogin(login, password);
      navigate('/dashboard');
    } catch (error) {
      console.error('Login failed:', error);
      setError('Неверный логин или пароль');
    }
  };

  return (
    <div className="min-h-screen bg-gray-50 flex items-center justify-center px-4 font-montserrat-regular">
      <div className="max-w-md w-full bg-white rounded-2xl shadow-sm p-8">
        {/* Логотип */}
        <div className="text-center">
          <div className="flex items-center justify-center mb-2">
            <img
              className="h-[40px] w-[220px]"
              src="/icons/lamoda-icon.svg"
              alt="Войти"
            />
          </div>
        </div>

        {/* Форма входа */}
        <form onSubmit={handleSubmit} className="space-y-4">
          {/* Поле логина */}
          <div className="relative">
            <label className="block text-xs text-gray-500 mb-2 uppercase tracking-wide font-montserrat-regular">
              телефон / email
            </label>
            <input
              type="text"
              value={login}
              onChange={(e) => setLogin(e.target.value)}
              className="w-full px-0 py-3 text-lg bg-transparent border-0 border-b-2 border-gray-200 focus:border-gray-400 focus:outline-none focus:ring-0 placeholder-gray-400 font-montserrat-regular"
              placeholder="Логин"
              disabled={isLoading}
              required
            />
          </div>

          {/* Поле пароля */}
          <div className="relative">
            <label className="block text-xs text-gray-500 mb-2 uppercase tracking-wide font-montserrat-regular">
              Пароль
            </label>
            <input
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              className="w-full px-0 py-3 text-lg bg-transparent border-0 border-b-2 border-gray-200 focus:border-gray-400 focus:outline-none focus:ring-0 placeholder-gray-400 font-montserrat-regular"
              placeholder="Пароль"
              disabled={isLoading}
              required
              minLength={8}
            />
          </div>

          {/* Сообщение об ошибке */}
          {error && (
            <div className="text-red-600 text-sm text-center font-montserrat-regular">{error}</div>
          )}

          {/* Кнопка входа */}
          <div className="pt-4">
            <button
              type="submit"
              disabled={isLoading}
              className="w-full bg-black text-white py-4 px-6 text-lg font-montserrat-semibold rounded-lg hover:bg-gray-800 focus:outline-none focus:ring-2 focus:ring-gray-500 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
            >
              {isLoading ? 'Входим...' : 'Войти'}
            </button>
          </div>

          {/* Ссылки */}
          <div className="flex justify-between items-center pt-6">
            <button
              type="button"
              className="text-gray-600 hover:text-gray-800 text-base transition-colors font-montserrat-regular"
              onClick={() => {
                // Логика восстановления пароля
                console.log('Forgot password clicked');
              }}
            >
              Забыли пароль
            </button>
            <button
              type="button"
              className="text-gray-600 hover:text-gray-800 text-base underline transition-colors font-montserrat-regular"
              onClick={() => {
                // Переход на регистрацию
                navigate('/register');
              }}
            >
              Регистрация
            </button>
          </div>
        </form>

        {/* Политика конфиденциальности */}
        <div className="text-center mt-12">
          <button
            type="button"
            className="text-gray-500 hover:text-gray-700 text-sm transition-colors font-montserrat-regular"
            onClick={() => {
              // Переход на политику конфиденциальности
              window.open('/privacy-policy', '_blank');
            }}
          >
            Политика конфиденциальности
          </button>
        </div>
      </div>
    </div>
  );
} 