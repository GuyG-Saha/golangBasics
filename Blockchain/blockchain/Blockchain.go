package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"runtime"
	"strings"
	"time"
)

const (
	MIN_DIFFICULTY_FOR_USING_CPU_CORES = 6
)

type Transaction struct {
	Sender   string
	Receiver string
	Amount   float64
}

type Block struct {
	Index        int
	Timestamp    time.Time
	Proof        int
	PrevHash     string
	Transactions []Transaction
}

type Blockchain struct {
	Chain        []Block
	Transactions []Transaction // pending transactions pool
	Nodes        map[string]bool
	Reward       float64
}

func (bc *Blockchain) MineBlock(difficulty int) Block {
	lastBlock := bc.GetLastBlock()
	previousProof := lastBlock.Proof
	var newProof int
	if difficulty >= MIN_DIFFICULTY_FOR_USING_CPU_CORES {
		newProof = ProofOfWorkConcurrent(previousProof, difficulty)
	} else {
		newProof = ProofOfWork(previousProof, difficulty)
	}
	previousHash := Hash(lastBlock)
	return bc.CreateBlock(newProof, previousHash)
}

func (bc *Blockchain) CreateBlock(proof int, previousHash string) Block {
	block := Block{
		Index:        len(bc.Chain) + 1,
		Timestamp:    time.Now(),
		Proof:        proof,
		Transactions: bc.Transactions,
		PrevHash:     previousHash,
	}
	bc.Transactions = make([]Transaction, 0) // Reset pending transactions
	bc.Chain = append(bc.Chain, block)
	return block
}

func NewBlockchain() *Blockchain {
	bc := &Blockchain{
		Chain:        make([]Block, 0),
		Transactions: make([]Transaction, 0),
		Nodes:        make(map[string]bool),
		Reward:       1.6,
	}
	bc.CreateBlock(1, "0") // Genesis block
	return bc
}

func (bc *Blockchain) GetLastBlock() Block {
	return bc.Chain[len(bc.Chain)-1]
}

func ProofOfWork(previousProof int, difficulty int) int {
	target := strings.Repeat("0", difficulty)
	newProof := 0
	for {
		guess := fmt.Sprintf("%d%d", previousProof, newProof)
		guessHash := sha256.Sum256([]byte(guess))
		if hex.EncodeToString(guessHash[:])[:difficulty] == target {
			return newProof
		}
		newProof++
	}
}

func ProofOfWorkConcurrent(previousProof int, difficulty int) int {
	target := strings.Repeat("0", difficulty)
	numWorkers := runtime.NumCPU() / 2 // Get half the number of CPU cores
	resultChan := make(chan int)       // Channel for successful proof
	// Start workers with different ranges to search
	for i := 0; i < numWorkers; i++ {
		go func(workerID int) {
			// Each worker starts at different offset and increments by numWorkers
			for proof := workerID; ; proof += numWorkers {
				guess := fmt.Sprintf("%d%d", previousProof, proof)
				guessHash := sha256.Sum256([]byte(guess))
				if hex.EncodeToString(guessHash[:])[:difficulty] == target {
					resultChan <- proof
					return
				}
			}
		}(i)
	}
	// Wait for first worker to find valid proof
	return <-resultChan
}

func (bc *Blockchain) AddTransaction(newTransaction Transaction) int {
	bc.Transactions = append(bc.Transactions, newTransaction)
	return bc.Chain[len(bc.Chain)-1].Index + 1
}

func Hash(block Block) string {
	h := sha256.New()
	h.Write([]byte(fmt.Sprintf("%d", block.Index)))
	h.Write([]byte(block.Timestamp.String()))
	h.Write([]byte(fmt.Sprintf("%d", block.Proof)))
	h.Write([]byte(block.PrevHash))
	for _, tx := range block.Transactions {
		h.Write([]byte(tx.Sender))
		h.Write([]byte(tx.Receiver))
		h.Write([]byte(fmt.Sprintf("%f", tx.Amount)))
	}
	return hex.EncodeToString(h.Sum(nil))
}
func FormatHash(hash string) string {
	return fmt.Sprintf("%064s", hash)
}
