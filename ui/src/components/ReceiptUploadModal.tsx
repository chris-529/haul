import { useState } from 'react'
import type { Receipt } from '../types'

type ReceiptUploadModalProps = {
  onClose: () => void
  onUploadSuccess: (receipt: Receipt) => void
}

export default function ReceiptUploadModal({
  onClose,
  onUploadSuccess,
}: ReceiptUploadModalProps) {
  const [file, setFile] = useState<File | null>(null)

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
    onUploadSuccess(data)
    onClose()
  }

  return (
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

        <button className="cancelBtn" onClick={onClose}>
          Cancel
        </button>
      </div>
    </div>
  )
}