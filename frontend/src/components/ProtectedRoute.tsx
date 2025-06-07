import { Navigate, Outlet } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';

export default function ProtectedRoute() {
  const { activeAccount, isInitialized } = useAuth();

  // Показываем загрузку, пока не завершена инициализация
  if (!isInitialized) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-gray-900 mx-auto"></div>
          <p className="mt-4 text-gray-600">Загрузка...</p>
        </div>
      </div>
    );
  }

  // Если активного аккаунта нет после инициализации, перенаправляем на страницу входа
  if (!activeAccount) {
    return <Navigate to="/login" replace />;
  }

  // Если аккаунт есть, рендерим дочерние компоненты
  return <Outlet />;
} 