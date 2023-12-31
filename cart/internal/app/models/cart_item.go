package models

type CartItem struct {
	User  int64
	Sku   uint32
	Count uint64
}

type CartItemWithInfo struct {
	Name  string
	User  int64
	Sku   uint32
	Price uint32
	Count uint64
}

type CartItemInfo struct {
	Name  string
	Sku   uint32
	Price uint32
}
