import { useEffect, useState } from 'react'
import '../App.css'
import NavBar from '../components/NavBar'
import ReceiptList from '../components/ReceiptList'
import type { Receipt, Item } from '../types'

export default function Recipes() {
  const [receipts, setReceipts] = useState<Receipt[]>([])
  const [selectedReceipt, setSelectedReceipt] = useState<Receipt | null>(null)
  const [recipeItems, setRecipeItems] = useState<Item[]>([])
  const [isDraggingOver, setIsDraggingOver] = useState(false)

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

  const handleDragStart = (
    e: React.DragEvent<HTMLButtonElement>,
    item: Item
  ) => {
    e.dataTransfer.setData('application/json', JSON.stringify(item))
    e.dataTransfer.effectAllowed = 'copy'
  }

  const handleDrop = (e: React.DragEvent<HTMLDivElement>) => {
    e.preventDefault()
    setIsDraggingOver(false)

    const rawItem = e.dataTransfer.getData('application/json')
    if (!rawItem) return

    const item = JSON.parse(rawItem) as Item
    setRecipeItems(prev => [...prev, item])
  }

  return (
    <div className="app">
      <NavBar />

      <div className="recipeBuilderLayout">
        <div className="receiptListPanel">
          <div className="receiptListHeader">
            <h2>Receipts</h2>
          </div>

          <ReceiptList
            receipts={receipts}
            selectedReceiptId={selectedReceipt?.id}
            onSelectReceipt={setSelectedReceipt}
          />
        </div>

        <div className="recipeItemsPanel">
          <div className="receiptListHeader">
            <h2>Items</h2>
          </div>

          {!selectedReceipt && (
            <p className="emptyText">Select a receipt to view items.</p>
          )}

          {selectedReceipt && (
            <div className="draggableItems">
              {selectedReceipt.items.map(item => (
                <button
                  key={item.id}
                  className="ingredientButton"
                  draggable
                  onDragStart={e => handleDragStart(e, item)}
                >
                  {item.name}
                </button>
              ))}
            </div>
          )}
        </div>

        <div
          className={`recipeCreatePanel ${isDraggingOver ? 'dragOver' : ''}`}
          onDragOver={e => {
            e.preventDefault()
            setIsDraggingOver(true)
          }}
          onDragLeave={() => setIsDraggingOver(false)}
          onDrop={handleDrop}
        >
          <div className="receiptListHeader">
            <h2>Create Recipe</h2>
          </div>

          {recipeItems.length === 0 ? (
            <p className="emptyText">Drag ingredients here.</p>
          ) : (
            <div className="recipeIngredientList">
              {recipeItems.map((item, index) => (
                <div key={`${item.id}-${index}`} className="recipeIngredient">
                  {item.name}
                </div>
              ))}
            </div>
          )}
        </div>
      </div>
    </div>
  )
}