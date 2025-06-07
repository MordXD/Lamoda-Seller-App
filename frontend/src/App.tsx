import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { AuthProvider } from './context/AuthContext';
import LoginPage from './pages/LoginPage';
import ProductsPage from './pages/ProductsPage';
import ProtectedRoute from './components/ProtectedRoute';

function App() {
  return (
    <AuthProvider>
      <Router>
        <Routes>
          <Route path="/login" element={<LoginPage />} />
          <Route element={<ProtectedRoute />}>
            <Route path="/products" element={<ProductsPage />} />
          </Route>
          <Route path="*" element={<Navigate to="/products" />} />
        </Routes>
      </Router>
    </AuthProvider>
  );
}

export default App;
