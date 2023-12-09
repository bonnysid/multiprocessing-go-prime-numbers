package main

import (
	"math"
	"sync"
	"time"
)
import "fmt"

var limit = 1000000
var numThreads = 10

// Функция для выполнения решета Эратосфена в заданном диапазоне
func sieveOfEratosthenes(limit int) []int {
	prime := make([]bool, limit+1)

	for i := 2; i <= limit; i++ {
		prime[i] = true
	}

	for p := 2; p*p <= limit; p++ {
		if prime[p] == true {
			for i := p * p; i <= limit; i += p {
				prime[i] = false
			}
		}
	}

	var primes []int
	for p := 2; p <= limit; p++ {
		if prime[p] == true {
			primes = append(primes, p)
		}
	}
	return primes
}

func algorithmEratosthenes() {
	var startTime = time.Now()
	sieveOfEratosthenes(limit)
	fmt.Printf("Время выполнения алгортима Эратосфена: %.2f\n", time.Since(startTime).Seconds())
}

// ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

func parallelSieve(start, end int, primes []int) []int {
	var primesRes []int
	for num := start; num <= end; num++ {
		isComposite := false
		for _, prime := range primes {
			if num%prime == 0 {
				isComposite = true
				break
			}
		}
		if !isComposite {
			primesRes = append(primesRes, num)
		}
	}

	return primesRes
}

// Параллельный алгоритм №1: декомпозиция по данным
func algorithmParallel1() {
	var startTime = time.Now()
	var eratosthenesLimit = int(math.Ceil(float64(limit) * 0.1))
	var primes = sieveOfEratosthenes(eratosthenesLimit)
	var start = eratosthenesLimit + 1
	var end = limit
	var intervalSize = (end - start + 1) / numThreads
	var wg sync.WaitGroup

	for i := 0; i < numThreads; i++ {
		wg.Add(1)
		go func(start, end int, primes []int) {
			parallelSieve(start, end, primes)
			wg.Done()
		}(start+i*intervalSize, start+(i+1)*intervalSize-1, primes)
	}

	wg.Wait()

	fmt.Printf("Время выполнения Параллельный алгоритм №1: декомпозиция по данным: %.2f\n", time.Since(startTime).Seconds())
}

// ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

// Параллельный алгоритм №2: декомпозиция набора простых чисел
func algorithmParallel2() {
	var startTime = time.Now()
	var eratosthenesLimit = int(math.Ceil(float64(limit) * 0.1))
	var primes = sieveOfEratosthenes(eratosthenesLimit)
	var start = eratosthenesLimit + 1
	var end = limit
	primeSets := make([][]int, numThreads)
	for i, prime := range primes {
		primeSets[i%numThreads] = append(primeSets[i%numThreads], prime)
	}
	var wg sync.WaitGroup

	for i := 0; i < numThreads; i++ {
		wg.Add(1)
		go func(start, end int, primes []int) {
			parallelSieve(start, end, primes)
			wg.Done()
		}(start, end, primeSets[i])
	}

	wg.Wait()

	fmt.Printf("Время выполнения Параллельный алгоритм №2: декомпозиция набора простых чисел: %.2f\n", time.Since(startTime).Seconds())
}

// ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

func worker(num int, primes []int, results chan<- int) {
	for _, prime := range primes {
		if num%prime == 0 && num != prime {
			results <- num
			break
		}
	}
}

// Параллельный алгоритм №3: применение пула потоков
func algorithmParallel3() {
	var startTime = time.Now()
	var eratosthenesLimit = int(math.Ceil(float64(limit) * 0.1))
	var primes = sieveOfEratosthenes(eratosthenesLimit)
	var start = eratosthenesLimit + 1
	var end = limit
	var wg sync.WaitGroup

	var results = make(chan int)

	for i := start; i <= end; i++ {
		wg.Add(1)
		go func(num int) {
			defer wg.Done()
			worker(num, primes, results)
		}(i)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	var compositeNumbers = []int{}
	for res := range results {
		compositeNumbers = append(compositeNumbers, res)
	}

	fmt.Printf("Время выполнения Параллельный алгоритм №3: применение пула потоков: %.2f\n", time.Since(startTime).Seconds())
}

// ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

func worker2(prime int, start, end int, results chan<- int) {
	for num := start; num <= end; num++ {
		if num%prime == 0 && num != prime {
			results <- num
		}
	}
}

// Параллельный алгоритм №4: последовательный перебор простых чисел
func algorithmParallel4() {
	var startTime = time.Now()
	var eratosthenesLimit = int(math.Ceil(float64(limit) * 0.1))
	var primes = sieveOfEratosthenes(eratosthenesLimit)
	var start = eratosthenesLimit + 1
	var end = limit
	var wg sync.WaitGroup

	var results = make(chan int)

	for _, prime := range primes {
		wg.Add(1)
		go func(p int) {
			defer wg.Done()
			worker2(p, start, end, results)
		}(prime)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	compositeNumbers := make(map[int]bool)
	for res := range results {
		compositeNumbers[res] = true
	}

	fmt.Printf("Время выполнения Параллельный алгоритм №4: последовательный перебор простых чисел: %.2f\n", time.Since(startTime).Seconds())
}

// ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

func main() {
	var startTime = time.Now()
	var wg sync.WaitGroup

	fmt.Printf("Граница от 0 до %d\n", limit)
	fmt.Printf("Колличетво ядер %d\n", numThreads)

	wg.Add(5)
	go func() {
		defer wg.Done()
		algorithmEratosthenes()
	}()
	go func() {
		defer wg.Done()
		algorithmParallel1()
	}()
	go func() {
		defer wg.Done()
		algorithmParallel2()
	}()
	go func() {
		defer wg.Done()
		algorithmParallel3()
	}()
	go func() {
		defer wg.Done()
		algorithmParallel4()
	}()
	wg.Wait()

	fmt.Printf("Полное время: %.2f\n", time.Since(startTime).Seconds())
}
