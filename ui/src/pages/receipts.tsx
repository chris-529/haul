import { useState } from 'react'
import '../App.css'
import NavBar from '../components/NavBar'
import ReceiptCard from '../components/ReceiptCard'

type Item = {
  id?: string
  name: string
  price: number
  quantity: number
}

type Receipt = {
  store: string
  status: string
  items: Item[]
}

export default function Receipts() {
  const [file, setFile] = useState<File | null>(null)
  const [receipt, setReceipt] = useState<Receipt | null>(null)

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

    const data = await res.json()
    setReceipt(data)
  }

  return (
    <div className="app">
      <NavBar />

      <div className="receiptForm">
        <p className="receiptFormTitle">Upload a receipt</p>

        <div className="receiptControls">
            <label className="fileButton">
                Choose receipt
                <input
                type="file"
                onChange={e => setFile(e.target.files?.[0] || null)}
                />
            </label>

            <button onClick={upload} className="btn">
                Upload
            </button>
            </div>

            <span className="fileName">
            {file ? file.name : 'No file chosen'}
            </span>

        </div>


      {receipt && <ReceiptCard receipt={receipt} />}
    </div>
  )
}