package main

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
	"github.com/gorilla/mux"
)

type Book struct{
	ID					string		`json:"id"`
	Title				string		`json:"title"`
	Author			string		`json:"author"`
	PublishDate	string		`json:"publish_date"`
	ISBN				string		`json:"isbn"`
}

type BookCheckout struct{
	BookId				string	`json:"book_id"`
	User					string	`json:"user"`
	CheckoutDate	string	`json:"checkout_date"`
	IsGenesis			bool		`json:"is_genesis"`
}

type Block struct{
	Position	int
	Data			BookCheckout
	PrevHash	string
	TimeStamp	string	
	Hash			string
}

type Blockchain struct{
	Block[] *Block
}

 var blockchain *Blockchain

func (b *Block)generateHash(){
	bytes,_ := json.Marshal(b.Data)
	data := string(b.Position) + b.TimeStamp + string(bytes) + b.PrevHash

	hash := sha256.New()
	hash.Write([]byte(data))
	b.Hash = hex.EncodeToString(hash.Sum(nil))
}

func CreateBlock(prevBlock *Block, data BookCheckout)(*Block){
	block := &Block{}
	block.Position = prevBlock.Position + 1
	block.Data = data
	block.PrevHash = prevBlock.Hash
	block.TimeStamp = time.Now().String()
	block.generateHash()
	return block
}

func (bc *Blockchain)AddBlock(data BookCheckout){
	prevBlock := bc.Block[bc.Block(length)-1]
	Block := CreateBlock(prevBlock,data)

	if validBlock(prevBlock,Block){
		bc.Block = append(bc.Block, Block)
	}
}

func validBlock(prevBlock *Block,block *Block)(bool){

	if prevBlock.Hash != block.Hash{
		return false
	}

	if !block.ValidateBlock(block.Hash){
		return false
	}

	if prevBlock.Position + 1 != block.Position{
		return false
	}

	return true
}

func (b *Block)ValidateBlock(Hash string)(bool){
	b.generateHash()
	if b.Hash != Hash {
		return false
	}

	return true
}

func newBook(w http.ResponseWriter, r *http.Request){
	var book Book

	if err := json.NewDecoder(r.Body).Decode(&book);err!=nil{
		w.WriteHeader(http.StatusInternalServerError)
		log.Print(err)
		w.Write([]byte("could not create new book"))
	}

	h := md5.New()
	io.WriteString(h,book.ISBN+book.PublishDate)
	book.ID = fmt.Sprintf("%x", h.Sum(nil))

	res,err := json.MarshalIndent(book,""," ")
	if err !=nil{
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("failed to load the payload: %v",err)
		w.Write([]byte("could not save the book"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func writeBlock(w http.ResponseWriter, r *http.Request){
	var CheckoutItem BookCheckout
	if err := json.NewDecoder(r.Body).Decode(&CheckoutItem);err!=nil{
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("could  not write to block:%v",err)
		w.Write([]byte("could not write block"))
	}
	blockchain.AddBlock(CheckoutItem)
}

func getBlockchain(w http.ResponseWriter, r *http.Request){

	jbytes,err := json.MarshalIndent(blockchain.Block,""," ")
	if err !=nil{
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err)
		return 
	}
	io.WriteString(w,string(jbytes))
}

func GenesisBlock() (*Block){
	return CreateBlock(&Block{},&BookCheckout{IsGenesis:true})
}

func NewBlockChain()(*Blockchain){
	return &Blockchain{[]*Block{GenesisBlock()}}
}

func main(){

	blockchain = NewBlockChain()

	r := mux.NewRouter()
	r.HandleFunc("/",getBlockchain).Methods("GET")
	r.HandleFunc("/",writeBlock).Methods("POST")
	r.HandleFunc("/new", newBook).Methods("POST")
	
	go func(){
		for _,block := range blockchain.Block{
			fmt.Printf("Prev. Hash: %x\n",block.PrevHash)
			bytes,_ := json.MarshalIndent(block.Data,""," ")
			fmt.Printf("Data:%v\n",string(bytes))
			fmt.Printf("Hash:%v\n",block.Hash)
		}
	}()

	log.Println("Listening on port 3000")

	log.Fatal(http.ListenAndServe(":3000", r))	
}