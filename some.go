package main

import (
	"crypto/rand"
	"math/big"
	"log"
)

func sixDigits() int64 {
	max := big.NewInt(9999)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		log.Fatal(err)
	}
	return n.Int64()
}

func main() {
	//client := redis.NewClient(&redis.Options{
	//	Addr:     "localhost:6379",
	//	Password: "", // no password set
	//	DB:       0,  // use default DB
	//})
	//
	//pong, err := client.Ping().Result()
	//fmt.Println(pong, err)
	//// Output: PONG <nil>
	//err = client.Set("ololo", "trololo", 0).Err()
	//if err != nil {
	//	panic(err)
	//}
	//
	//val, err := client.Get("ololo").Result()
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println("key", val)
	//fmt.Println(fmt.Sprintf("%04d",rand.Intn(1000)))
	for i := 0; i < 100; i++ {
		s := sixDigits()
		log.Printf("%04d\n", s)
	}
}
