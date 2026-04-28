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

type ReceiptCardProps = {
  receipt: Receipt
}

export default function ReceiptCard({ receipt }: ReceiptCardProps) {
  return (
    <div className="receiptUploadCard">
      <div className="receiptHeader">
        <h3>{receipt.store}</h3>
        <span>{receipt.status}</span>
      </div>

      <ul className="receiptItems">
        {receipt.items.map((item, index) => (
          <li key={item.id || index} className="receiptItem">
            <span>{item.name}</span>
            <span>
              {item.quantity} — ${item.price}
            </span>
          </li>
        ))}
      </ul>
    </div>
  )
}