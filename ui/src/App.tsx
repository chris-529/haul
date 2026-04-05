import { useState } from 'react'

export default function App() {
  const [file, setFile] = useState<File | null>(null)

  const upload = async () => {
    if (!file) return
    const body = new FormData()
    body.append('receipt_image', file)

    const res = await fetch('/api/receipts', { method: 'POST', body })
    console.log('Response:', await res.json())
  }

  return (
    <div style={{ padding: '2rem', fontFamily: 'sans-serif' }}>
      <h1>Haul</h1>
      <input type="file" onChange={e => setFile(e.target.files?.[0] || null)} />
      <button onClick={upload} style={{ marginLeft: '10px' }}>Upload</button>
    </div>
  )
}