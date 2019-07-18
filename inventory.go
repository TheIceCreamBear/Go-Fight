package main

import "fmt"

const inventorySize = 10

type Inventory struct {
	armorSlot *Item
	itemSlots [inventorySize]*Item
}

func NewInventory() *Inventory {
	inv := new(Inventory)
	return inv
}

func (inv *Inventory) addItem(item *Item) bool {
	if inv.isFull() {
		return false
	}
	index := inv.findFirstEmpty()
	if index == -1 {
		fmt.Println("THIS SHOULDNT BE POSSIBLE")
		fmt.Println("addItem(*Item)bool: index == -1")
		return false
	}
	inv.itemSlots[index] = item
	return true
}

// DO NOT CALL THIS METHOD ON A FULL INVENTORY
func (inv *Inventory) findFirstEmpty() int {
	for i := 0; i < inventorySize; i++ {
		if inv.itemSlots[i] == nil {
			return i
		}
	}
	return -1
}

func (inv *Inventory) isFull() bool {
	for i := 0; i < inventorySize; i++ {
		if inv.itemSlots[i] == nil {
			return false
		}
	}
	return true
}

func (inv *Inventory) slotsUsed() int {
	count := 0
	for i := 0; i < inventorySize; i++ {
		if inv.itemSlots[i] != nil {
			count++
		}
	}
	return count
}

func (inv *Inventory) printFullInventory() {
	fmt.Println("\nPrinting Inventory:")
	inv.printItemInventory()
	if inv.armorSlot == nil {
		fmt.Printf("ArmorSlot: Empty\n")
	} else {
		fmt.Printf("ArmorSlot: DefenseBoost=%-9.6f", inv.armorSlot.effect)
	}
}

func (inv *Inventory) printItemAt(index int) {
	if index >= 0 && index < inventorySize {
		current := inv.itemSlots[index]
		fmt.Printf("ItemSlot%1d: Type=%-7s Effect=%7.3f\n", index, getStringFromItemType(current.iType), current.effect)
	}
}

func (inv *Inventory) printItemInventory() {
	fmt.Printf("Slots used:%2d/%-2d\n", inv.slotsUsed(), inventorySize)
	for i := 0; i < inventorySize; i++ {
		current := inv.itemSlots[i]
		if current == nil {
			fmt.Printf("ItemSlot%1d: Empty\n", i)
		} else {
			fmt.Printf("ItemSlot%1d: Type=%-7s Effect=%7.3f\n", i, getStringFromItemType(current.iType), current.effect)
		}
	}
}
