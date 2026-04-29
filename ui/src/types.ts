export type Item = {
  id?: string
  receipt_id?: string
  name: string
  price: number
  quantity: number
  unit: string
}

export type Receipt = {
  id?: string
  user_id?: string
  store: string
  status: string
  created_at?: string
  items: Item[]
}