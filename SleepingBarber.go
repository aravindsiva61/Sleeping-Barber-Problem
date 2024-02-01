package main

import (
    "fmt"
    "math/rand"
    "sync"
    "time"
)

const (
    numberOfBarbers  = 2
    numberOfChairs   = 5
    shopWorkingHours = 10 * time.Second // Simulating shop working hours
)

//Structure to represent barbershop
type BarberShop struct {
    waitingRoom chan int
    barberReady chan bool
    mutex       sync.Mutex
    isClosed    bool
    wg          sync.WaitGroup
}

//Function to instantiate a new barber shop
func NewBarberShop() *BarberShop {
	return &BarberShop{
        waitingRoom: make(chan int, numberOfChairs),
        barberReady: make(chan bool, numberOfBarbers),
        mutex:       sync.Mutex{},
        isClosed:    false,
        wg:          sync.WaitGroup{},
    }
}

//Function to open the shop
func (shop *BarberShop) openShop() {
    for i := 0; i < numberOfBarbers; i++ {
        shop.wg.Add(1)
        go shop.barber(i)
    }

    time.Sleep(shopWorkingHours)
    shop.closeShop()
}

//Function to show barber cutting the hair and going home after working hours
func (shop *BarberShop) barber(id int) {
	defer shop.wg.Done()
    for {
        select {
        case customer, ok := <-shop.waitingRoom:
            if !ok {
                fmt.Printf("Barber %d is going home\n", id)
                return
            }
            fmt.Printf("Barber %d is cutting hair of customer %d\n", id, customer)
            time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
            fmt.Printf("Barber %d finished with customer %d\n", id, customer)
            shop.barberReady <- true
        }
    }
}

//Function to show customer arriving and waiting to get haircut
func (shop *BarberShop) customer(id int) {
    shop.mutex.Lock()
    defer shop.mutex.Unlock()
    if shop.isClosed {
        fmt.Printf("Customer %d found the shop closed and is leaving\n", id)
        return
    }

    select {
    case shop.waitingRoom <- id:
        fmt.Printf("Customer %d is waiting in the waiting room\n", id)
    default:
        fmt.Printf("Customer %d found no empty chairs and is leaving\n", id)
    }
}

//Function to close the shop
func (shop *BarberShop) closeShop() {
	shop.mutex.Lock()
    shop.isClosed = true
    remainingCustomers := len(shop.waitingRoom)
    shop.mutex.Unlock()

    for i := 0; i < remainingCustomers; i++ {
        <-shop.barberReady
    }

    shop.wg.Wait()
    close(shop.waitingRoom)
    fmt.Println("The barber shop is now closed")
}

//Main function to orchestate the sleeping barber code flow
func main() {
    rand.Seed(time.Now().UnixNano())
    shop := NewBarberShop()
    go shop.openShop()

    customerID := 1
    for {
        time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
        go shop.customer(customerID)
        customerID++

        shop.mutex.Lock()
        if shop.isClosed {
            shop.mutex.Unlock()
            break
        }
        shop.mutex.Unlock()
    }

    time.Sleep(1 * time.Second)
}
