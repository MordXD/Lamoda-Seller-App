import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { register } from '../api/auth';
import type { RegisterData } from '../api/auth';

export default function RegisterPage() {
  const navigate = useNavigate();
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState('');
  const [generatedPassword, setGeneratedPassword] = useState('');

  // Основные поля - только имя и email
  const [name, setName] = useState('');
  const [email, setEmail] = useState('');

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');

    if (!name.trim() || !email.trim()) {
      setError('Пожалуйста, заполните все обязательные поля');
      return;
    }

    setIsLoading(true);

    try {
      const registerData: RegisterData = {
        name: name.trim(),
        email: email.trim(),
      };

      const response = await register(registerData);
      setGeneratedPassword(response.temporary_password);
    } catch (error: any) {
      console.error('Registration failed:', error);
      setError(error.response?.data?.message || error.message || 'Что-то пошло не так');
    } finally {
      setIsLoading(false);
    }
  };

  const handleCopyToClipboard = () => {
    navigator.clipboard.writeText(generatedPassword);
    alert('Пароль скопирован в буфер обмена!');
  };

  if (generatedPassword) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center px-4 font-montserrat-regular">
        <div className="max-w-md w-full bg-white rounded-2xl shadow-sm p-8 text-center">
          <h2 className="text-2xl font-bold mb-4">Регистрация успешна!</h2>
          <p className="mb-4">Ваш временный пароль для входа:</p>
          <div className="relative p-4 bg-gray-100 rounded-lg mb-4">
            <p className="text-lg font-mono break-all">{generatedPassword}</p>
            <button
              onClick={handleCopyToClipboard}
              className="absolute top-2 right-2 text-gray-500 hover:text-gray-800"
              title="Скопировать пароль"
            >
              <svg xmlns="http://www.w3.org/2000/svg" className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
              </svg>
            </button>
          </div>
          <p className="text-sm text-gray-500 mb-6">
            Сохраните его в надежном месте. Используйте этот пароль и ваш email для входа в систему.
          </p>
          <button
            onClick={() => navigate('/login')}
            className="w-full bg-black text-white py-4 px-6 text-lg font-montserrat-semibold rounded-lg hover:bg-gray-800 focus:outline-none focus:ring-2 focus:ring-gray-500 focus:ring-offset-2 transition-colors"
          >
            Перейти на страницу входа
          </button>
        </div>
      </div>
    );
  }

  const formInputClass = "w-full px-4 py-3 text-lg bg-white border-2 border-gray-200 rounded-lg focus:border-gray-400 focus:outline-none focus:ring-0 placeholder-gray-400 font-montserrat-regular";
  const formLabelClass = "block text-sm text-gray-700 mb-2 font-montserrat-medium";

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
            Регистрация продавца
          </h2>
          <p className="text-gray-600 mb-8">
            Создайте аккаунт для работы с маркетплейсом Lamoda
          </p>
        </div>

        <form onSubmit={handleSubmit} className="space-y-6">
          {/* Ваше имя */}
          <div>
            <label htmlFor="name" className={formLabelClass}>
              Ваше имя <span className="text-red-500">*</span>
            </label>
            <input 
              id="name" 
              type="text" 
              required 
              value={name} 
              onChange={(e) => setName(e.target.value)} 
              className={formInputClass} 
              placeholder="Введите ваше имя" 
              disabled={isLoading} 
            />
          </div>

          {/* E-mail */}
          <div>
            <label htmlFor="email" className={formLabelClass}>
              E-mail <span className="text-red-500">*</span>
            </label>
            <input 
              id="email" 
              type="email" 
              required 
              value={email} 
              onChange={(e) => setEmail(e.target.value)} 
              className={formInputClass} 
              placeholder="Введите ваш email" 
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
            {isLoading ? 'Регистрация...' : 'Зарегистрироваться'}
          </button>

          <div className="text-center pt-4">
            <p className="text-sm text-gray-600">
              Уже есть аккаунт?{' '}
              <button
                type="button"
                onClick={() => navigate('/login')}
                className="text-black font-medium hover:underline"
              >
                Войти
              </button>
            </p>
          </div>
        </form>
      </div>
    </div>
  );
} 