import { useEffect, useState } from 'react'
import '../App.css'
import NavBar from '../components/NavBar'
import ReceiptList from '../components/ReceiptList'
import ReceiptCard from '../components/ReceiptCard'
import ReceiptUploadModal from '../components/ReceiptUploadModal'
import type { Receipt } from '../types'

export default function Receipts() {
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

  const handleUploadSuccess = (newReceipt: Receipt) => {
    setReceipt(newReceipt)
    setReceipts(prev => [newReceipt, ...prev])
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

          <ReceiptList
            receipts={receipts}
            selectedReceiptId={receipt?.id}
            onSelectReceipt={setReceipt}
          />
        </div>

        <div className="receiptDetailPanel">
          {receipt && <ReceiptCard receipt={receipt} />}
        </div>
      </div>

      {showUpload && (
        <ReceiptUploadModal
          onClose={() => setShowUpload(false)}
          onUploadSuccess={handleUploadSuccess}
        />
      )}
    </div>
  )
}