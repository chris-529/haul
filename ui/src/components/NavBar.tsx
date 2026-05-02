import { useState } from 'react'
import { useNavigate } from 'react-router-dom'

export default function NavBar() {
  const [menuOpen, setMenuOpen] = useState(false)
  const navigate = useNavigate()

  const logout = () => {
    localStorage.removeItem('token')
    navigate('/auth/login')
  }

  return (
    <nav className="navbar">
      <h2 className="title">Haul</h2>

      <div className="userMenu">
        <button
          className="userIconBtn"
          onClick={() => setMenuOpen(prev => !prev)}
          aria-label="User menu"
        >
          👤
        </button>

        {menuOpen && (
          <div className="userDropdown">
            <button onClick={logout}>Logout</button>
          </div>
        )}
      </div>
    </nav>
  )
}