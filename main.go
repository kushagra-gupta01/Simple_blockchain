package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
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

func (bc *Blockchain)AddBlock(data BookCheckout){

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
	json.NewDecoder(r.Body).Decode(&CheckoutItem);err!=nil{
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("could  not write to block:%v",err)
		w.Write([]byte("could not write block"))
	}
	Blockchain.AddBlock(CheckoutItem)
}

func getBlockchain(w http.ResponseWriter, r *http.Request){

}

func main(){
	r := mux.NewRouter()
	r.HandleFunc("/",getBlockchain).Methods("GET")
	r.HandleFunc("/",writeBlock).Methods("POST")
	r.HandleFunc("/new", newBook).Methods("POST")
	
	log.Println("Listening on port 3000")

	log.Fatal(http.ListenAndServe(":3000", r))	
}