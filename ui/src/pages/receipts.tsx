import { useEffect, useState } from 'react'
import '../App.css'
import NavBar from '../components/NavBar'
import ReceiptCard from '../components/ReceiptCard'
import type { Receipt } from '../types'

export default function Receipts() {
  const [file, setFile] = useState<File | null>(null)
  const [receipt, setReceipt] = useState<Receipt | null>(null)
  const [receipts, setReceipts] = useState<Receipt[]>([])
  const [showUpload, setShowUpload] = useState(false)

  useEffect(() => {
    const loadReceipts = async () => {
      const token = localStorage.getItem('token')

      const res = await fetch('/receipts/', {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      })

      if (!res.ok) return

      const data = await res.json()
      setReceipts(data)
    }

    loadReceipts()
  }, [])

  const upload = async () => {
    if (!file) return

    const body = new FormData()
    body.append('receipt_image', file)

    const token = localStorage.getItem('token')

    const res = await fetch('/receipts/', {
      method: 'POST',
      headers: {
        Authorization: `Bearer ${token}`,
      },
      body,
    })

    if (!res.ok) return

    const data = await res.json()

    setReceipt(data)
    setReceipts(prev => [data, ...prev])
    setFile(null)
    setShowUpload(false)
  }

  return (
    <div className="app">
      <NavBar />

      <div className="receiptsLayout">
        <div className="receiptListPanel">
          <div className="receiptListHeader">
            <h2>Receipts</h2>

            <button className="btn" onClick={() => setShowUpload(true)}>
              New
            </button>
          </div>

          <div className="receiptList">
            {receipts.length === 0 ? (
              <p className="emptyText">No receipts yet</p>
            ) : (
              receipts.map((r, index) => (
                <button
                  key={r.id || index}
                  className="receiptListItem"
                  onClick={() => setReceipt(r)}
                >
                  <span>{r.store}</span>
                  <span>{r.status}</span>
                </button>
              ))
            )}
          </div>
        </div>

        <div className="receiptDetailPanel">
          {receipt && <ReceiptCard receipt={receipt} />}
        </div>
      </div>

      {showUpload && (
        <div className="modalOverlay">
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

            <button className="cancelBtn" onClick={() => setShowUpload(false)}>
              Cancel
            </button>
          </div>
        </div>
      )}
    </div>
  )
}