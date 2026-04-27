import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import './App.css'

import Login from './pages/login'
import Register from './pages/register'
import Receipts from './pages/receipts'

export default function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Navigate to="/dashboard" />} />
        <Route path="/auth/login" element={<Login />} />
        <Route path="/auth/register" element={<Register />} />
        <Route path="/dashboard" element={<Receipts />} />
      </Routes>
    </BrowserRouter>
  )
}