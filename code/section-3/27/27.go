package _27

import "time"

// demonstrates a race condition

var balance = 100

func main() {
	go spend(30)
	go spend(40)
}

func spend(amount int) {
	b := balance
	time.Sleep(time.Second)
	b -= amount
	balance = b
}
