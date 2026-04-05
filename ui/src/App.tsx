import { useState } from 'react'

export default function App() {
  const [file, setFile] = useState<File | null>(null)

  const upload = async () => {
    if (!file) return
    const body = new FormData()
    body.append('receipt_image', file)
    const res = await fetch('/api/receipts', { method: 'POST', body })
    console.log(await res.json())
  }

  const s = {
    app: { background: '#121212', color: '#ccc', minHeight: '100vh', padding: '2rem', fontFamily: 'sans-serif' },
    input: { background: '#222', border: '1px solid #333', color: '#ccc', padding: '6px', borderRadius: '4px' },
    btn: { background: '#333', color: '#fff', border: 'none', padding: '8px 16px', borderRadius: '4px', cursor: 'pointer', marginLeft: '8px' }
  }

  return (
    <div style={s.app}>
      <h2 style={{ color: '#fff' }}>Haul</h2>
      <input 
        type="file" 
        style={s.input} 
        onChange={e => setFile(e.target.files?.[0] || null)} 
      />
      <button 
        onClick={upload} 
        style={s.btn}
      >
        Upload
      </button>
      {file && <p style={{ color: '#666', fontSize: '14px' }}>{file.name}</p>}
    </div>
  )
}