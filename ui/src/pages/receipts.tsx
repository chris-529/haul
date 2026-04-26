import { useState } from 'react'
import '../App.css'

export default function App() {
  const [file, setFile] = useState<File | null>(null)

  const upload = async () => {
    if (!file) return
    const body = new FormData()
    body.append('receipt_image', file)

    const token = localStorage.getItem('token')
    const res = await fetch('/receipts', {
      method: 'POST',
      headers: {
        Authorization: `Bearer ${token}`,
      },
      body,
    })
    console.log(await res.json())
  }

  return (
    <div className="app">
      <h2 className="title">Haul</h2>

      <input
        type="file"
        className="input"
        onChange={e => setFile(e.target.files?.[0] || null)}
      />

      <button onClick={upload} className="btn">
        Upload
      </button>

      {file && <p className="filename">{file.name}</p>}
    </div>
  )
}