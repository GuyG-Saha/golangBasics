package main

import (
	"Blockchain/blockchain"
	"fmt"
	"time"
)

func main() {
	bc := blockchain.NewBlockchain()
	// Add some transactions
	bc.AddTransaction(blockchain.Transaction{
		Sender:   "Alice",
		Receiver: "Bob",
		Amount:   50.0,
	})
	bc.AddTransaction(blockchain.Transaction{
		Sender:   "Bob",
		Receiver: "Charlie",
		Amount:   34.0,
	})
	// Let's see how long mining takes
	fmt.Println("Mining new block...")
	start := time.Now()
	// Mine a new block
	block := bc.MineBlock(6)
	duration := time.Since(start)

	// Print results
	fmt.Printf("\nNew Block Mined in %v!\n", duration)
	fmt.Printf("Index: %d\n", block.Index)
	fmt.Printf("Timestamp: %v\n", block.Timestamp)
	fmt.Printf("Proof: %d\n", block.Proof)
	fmt.Printf("Previous Hash: %s\n", blockchain.FormatHash(block.PrevHash))
	fmt.Printf("Number of Transactions: %d\n", len(block.Transactions))

	// Print all transactions in the block
	fmt.Println("\nTransactions in this block:")
	for i, tx := range block.Transactions {
		fmt.Printf("%d. %s -> %s: %.2f\n",
			i+1, tx.Sender, tx.Receiver, tx.Amount)
	}
	// New Block
	bc.AddTransaction(blockchain.Transaction{
		Sender:   "Delta",
		Receiver: "Eilon",
		Amount:   76.1,
	})
	bc.AddTransaction(blockchain.Transaction{
		Sender:   "Fitzgerald",
		Receiver: "Galil",
		Amount:   184.7,
	})
	fmt.Println("Mining new block...")
	start2 := time.Now()
	block2 := bc.MineBlock(6)
	duration = time.Since(start2)
	// Print results
	fmt.Printf("\nNew Block Mined in %v!\n", duration)
	fmt.Printf("Index: %d\n", block2.Index)
	fmt.Printf("Timestamp: %v\n", block2.Timestamp)
	fmt.Printf("Proof: %d\n", block2.Proof)
	fmt.Printf("Previous Hash: %s\n", blockchain.FormatHash(block2.PrevHash))
	fmt.Printf("Number of Transactions: %d\n", len(block2.Transactions))

	// Print all transactions in the block
	fmt.Println("\nTransactions in this block:")
	for i, tx := range block2.Transactions {
		fmt.Printf("%d. %s -> %s: %.2f\n",
			i+1, tx.Sender, tx.Receiver, tx.Amount)
	}
}
