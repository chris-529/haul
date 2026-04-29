import type { Receipt } from '../types'

type ReceiptListProps = {
  receipts: Receipt[]
  selectedReceiptId?: string
  onSelectReceipt: (receipt: Receipt) => void
}

export default function ReceiptList({
  receipts,
  selectedReceiptId,
  onSelectReceipt,
}: ReceiptListProps) {
  if (receipts.length === 0) {
    return (
      <div className="receiptList">
        <p className="emptyText">No receipts yet</p>
      </div>
    )
  }

  return (
    <div className="receiptList">
      {receipts.map((receipt, index) => (
        <button
          key={receipt.id || index}
          className={
            receipt.id === selectedReceiptId
              ? 'receiptListItem selected'
              : 'receiptListItem'
          }
          onClick={() => onSelectReceipt(receipt)}
        >
          <span>{receipt.store}</span>
          <span>{receipt.status}</span>
        </button>
      ))}
    </div>
  )
}