import { useState } from 'react'
import { Link, useLocation, useNavigate } from 'react-router-dom'

export default function NavBar() {
  const [menuOpen, setMenuOpen] = useState(false)
  const navigate = useNavigate()
  const location = useLocation()

  const logout = () => {
    localStorage.removeItem('token')
    navigate('/auth/login')
  }

  return (
    <nav className="navbar">
      <div className="navLeft">
        <h2 className="title">Haul</h2>

        <div className="navButtons">
          <Link
            to="/dashboard"
            className={`navButton ${location.pathname === '/dashboard' ? 'active' : ''}`}
          >
            Receipts
          </Link>

          <Link
            to="/recipes"
            className={`navButton ${location.pathname === '/recipes' ? 'active' : ''}`}
          >
            Recipes
          </Link>
        </div>
      </div>

      <div className="userMenu">
        <button
          className="userIconBtn"
          onClick={() => setMenuOpen(prev => !prev)}
          aria-label="User menu"
        >
          <i className="ti ti-user" aria-hidden="true" />
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