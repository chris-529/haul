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
    <div>
      <h3>{receipt.store}</h3>
      <p>{receipt.status}</p>

      <ul>
        {receipt.items.map((item, index) => (
          <li key={item.id || index}>
            {item.name}, {item.quantity} — ${item.price}
          </li>
        ))}
      </ul>
    </div>
  )
}