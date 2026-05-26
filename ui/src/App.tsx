import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import './App.css'

import Login from './pages/login'
import Register from './pages/register'
import Receipts from './pages/receipts'
import ProtectedRoute from './components/ProtectedRoute'
import Recipes from './pages/recipes'
import '@tabler/icons-webfont/dist/tabler-icons.min.css'

export default function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Navigate to="/dashboard" replace />} />

        <Route path="/auth/login" element={<Login />} />
        <Route path="/auth/register" element={<Register />} />

        <Route
          path="/dashboard"
          element={
            <ProtectedRoute>
              <Receipts />
            </ProtectedRoute>
          }
        />

        <Route
          path="/recipes"
          element={
            <ProtectedRoute>
              <Recipes />
            </ProtectedRoute>
          }
        />
      </Routes>
    </BrowserRouter>
  )
}