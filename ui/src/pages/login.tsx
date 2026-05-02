import { useState, type SubmitEvent } from 'react'
import { Link, useNavigate } from 'react-router-dom'

export default function Login() {
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [message, setMessage] = useState('')

  const navigate = useNavigate()

  const login = async (e: SubmitEvent<HTMLFormElement>) => {
    e.preventDefault()
    setMessage('')

    const res = await fetch('/login', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email, password }),
    })

    if (!res.ok) {
      setMessage('Invalid email or password')
      return
    }

    const data = await res.json()
    localStorage.setItem('token', data.token)
    navigate('/dashboard')
  }

  return (
    <div className="authPage">
      <h1 className="authBrand">Haul</h1>
      <div className="authBox">
        <h2>Login</h2>

        <form className="authForm" onSubmit={login}>
          <input
            type="email"
            placeholder="email"
            value={email}
            onChange={e => setEmail(e.target.value)}
          />

          <input
            type="password"
            placeholder="password"
            value={password}
            onChange={e => setPassword(e.target.value)}
          />

          <button className="btn" type="submit">Login</button>
        </form>

        {message && <p className="authMessage">{message}</p>}

        <p className="authSwitch">
          Don&apos;t have an account? <Link to="/auth/register">Register</Link>
        </p>
      </div>
    </div>
  )
}