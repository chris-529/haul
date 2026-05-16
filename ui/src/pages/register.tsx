import { useState, type SubmitEvent } from 'react'
import { Link, useNavigate } from 'react-router-dom'

export default function Register() {
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [message, setMessage] = useState('')
  const [inviteCode, setInviteCode] = useState('')

  const navigate = useNavigate()

  const register = async (e: SubmitEvent<HTMLFormElement>) => {
    e.preventDefault()
    setMessage('')

    const res = await fetch('/register', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        email,
        password,
        invite_code: inviteCode,
      }),
    })

    if (!res.ok) {
      setMessage('Could not create account')
      return
    }

    navigate('/auth/login')
  }

  return (
    <div className="authPage">
      <h1 className="authBrand">Haul</h1>
      <div className="authBox">
        <h2>Register</h2>

        <form className="authForm" onSubmit={register}>
          <input
            type="email"
            placeholder="Email"
            value={email}
            onChange={e => setEmail(e.target.value)}
          />

          <input
            type="password"
            placeholder="Password"
            value={password}
            onChange={e => setPassword(e.target.value)}
          />

          <input
            type="text"
            placeholder="Invite Code"
            value={inviteCode}
            onChange={e => setInviteCode(e.target.value)}
          />

          <button className="btn" type="submit">Register</button>
        </form>

        {message && <p className="authMessage">{message}</p>}

        <p className="authSwitch">
          Already have an account? <Link to="/auth/login">Login</Link>
        </p>
      </div>
    </div>
  )
}