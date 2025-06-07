import { Navigate, Outlet } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';

export default function ProtectedRoute() {
  const { activeAccount } = useAuth();

  // Если активного аккаунта нет, перенаправляем на страницу входа
  if (!activeAccount) {
    return <Navigate to="/login" replace />;
  }

  // Если аккаунт есть, рендерим дочерние компоненты
  return <Outlet />;
} 