import type { Receipt } from '../types'

type ReceiptCardProps = {
  receipt: Receipt
  onDeleteReceipt: (receiptId: string) => void
}

export default function ReceiptCard({ receipt, onDeleteReceipt }: ReceiptCardProps) {
  return (
    <div className="receiptCard">
      <div className="receiptHeader">
        <h3>{receipt.store}</h3>
        <span className="receiptStatus">{receipt.status}</span>

        <span className="receiptDate">
          {receipt.created_at
            ? new Date(receipt.created_at).toLocaleDateString()
            : 'No date'}
        </span>

        {receipt.id && (
          <button
            className="receiptDeleteBtn"
            onClick={() => onDeleteReceipt(receipt.id!)}
          >
            Delete
          </button>
        )}
      </div>

      <div className="receiptItemsHeader">
        <span>Item</span>
        <span>Qty</span>
        <span>Price</span>
        <span>Unit</span>
      </div>

      <ul className="receiptItems">
        {receipt.items.map((item, index) => (
          <li key={item.id || index} className="receiptItem">
            <span>{item.name}</span>
            <span>{item.quantity}</span>
            <span>${item.price.toFixed(2)}</span>
            <span>{item.unit}</span>
          </li>
        ))}
      </ul>
    </div>
  )
}