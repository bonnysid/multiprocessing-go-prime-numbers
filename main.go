package main

import (
	"fmt"
	"math"
	"sync"
	"time"
)

var limit = 10000
var numThreads = 100
var n = limit / 2
var start, end = n, 2 * n

// Функция для проверки числа на простоту
func isPrime(n int) bool {
	if n <= 1 {
		return false
	}
	sqrt := int(math.Sqrt(float64(n)))
	for i := 2; i <= sqrt; i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}

// Функция для последовательного поиска простых чисел
func sequentialCheckPrimes(start, end int, basePrimes []int) []int {
	var result []int
	for i := start; i <= end; i++ {
		isPrime := true
		sqrt := int(math.Sqrt(float64(i)))
		for _, prime := range basePrimes {
			if prime > sqrt {
				break
			}
			if i%prime == 0 {
				isPrime = false
				break
			}
		}
		if isPrime {
			result = append(result, i)
		}
	}
	return result
}

// Функция для параллельного поиска простых чисел
func parallelPrimeSearch(limit int) []int {
	var primes []int
	var wg sync.WaitGroup
	ch := make(chan int)

	for i := 0; i <= limit; i++ {
		wg.Add(1)
		go func(num int) {
			defer wg.Done()
			if isPrime(num) {
				ch <- num
			}
		}(i)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	for prime := range ch {
		primes = append(primes, prime)
	}

	return primes
}

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

// Функция для параллельной проверки простых чисел в интервале от n до 2n
func parallelCheckPrimesTask0(n int, primes []int) []int {
	var wg sync.WaitGroup
	var mutex sync.Mutex
	var result []int

	for i := n; i <= 2*n; i++ {
		wg.Add(1)
		go func(num int) {
			defer wg.Done()
			isPrime := true
			sqrt := int(math.Sqrt(float64(num)))
			mutex.Lock()
			for _, prime := range primes {
				if prime > sqrt {
					break
				}
				if num%prime == 0 {
					isPrime = false
					break
				}
			}
			mutex.Unlock()
			if isPrime {
				mutex.Lock()
				result = append(result, num)
				mutex.Unlock()
			}
		}(i)
	}

	wg.Wait()
	return result
}

// Функция для параллельной проверки простых чисел с использованием базовых простых чисел
func parallelCheckPrimesWithPrimePerThread(start, end int, basePrimes []int, wg *sync.WaitGroup, mutex *sync.Mutex, result *[]int, currentIndex *int) {
	defer wg.Done()

	for {
		mutex.Lock()
		if *currentIndex >= len(basePrimes) {
			mutex.Unlock()
			break
		}
		prime := basePrimes[*currentIndex]
		*currentIndex++
		mutex.Unlock()

		for i := start; i <= end; i++ {
			isPrime := true
			sqrt := int(math.Sqrt(float64(i)))
			for _, basePrime := range basePrimes {
				if basePrime > sqrt {
					break
				}
				if i%basePrime == 0 && basePrime != prime {
					isPrime = false
					break
				}
			}
			if isPrime && i%prime == 0 {
				mutex.Lock()
				*result = append(*result, i)
				mutex.Unlock()
			}
		}
	}
}

// Функция для проверки простоты числа с использованием базовых простых чисел
func checkPrimeWithBasePrimes(start, end int, basePrimes []int, wg *sync.WaitGroup, mutex *sync.Mutex, result *[]int) {
	defer wg.Done()

	for i := start; i <= end; i++ {
		isPrime := true
		sqrt := int(math.Sqrt(float64(i)))
		for _, prime := range basePrimes {
			if prime > sqrt {
				break
			}
			if i%prime == 0 {
				isPrime = false
				break
			}
		}
		if isPrime {
			mutex.Lock()
			*result = append(*result, i)
			mutex.Unlock()
		}
	}
}

// Функция для параллельной проверки простых чисел в интервале от start до end
func parallelCheckPrimesWithSubsetOfPrimes(start, end int, basePrimesSubset []int, wg *sync.WaitGroup, mutex *sync.Mutex, result *[]int) {
	defer wg.Done()

	for i := start; i <= end; i++ {
		isPrime := true
		sqrt := int(math.Sqrt(float64(i)))
		for _, prime := range basePrimesSubset {
			if prime > sqrt {
				break
			}
			if i%prime == 0 {
				isPrime = false
				break
			}
		}
		if isPrime {
			mutex.Lock()
			*result = append(*result, i)
			mutex.Unlock()
		}
	}
}

// Функция для параллельной проверки простых чисел в интервале от start до end
func parallelCheckPrimes(start, end int, basePrimes []int, wg *sync.WaitGroup, mutex *sync.Mutex, result *[]int) {
	defer wg.Done()

	for i := start; i <= end; i++ {
		isPrime := true
		sqrt := int(math.Sqrt(float64(i)))
		for _, prime := range basePrimes {
			if prime > sqrt {
				break
			}
			if i%prime == 0 {
				isPrime = false
				break
			}
		}
		if isPrime {
			mutex.Lock()
			*result = append(*result, i)
			mutex.Unlock()
		}
	}
}

func algorithm0() time.Duration {
	fmt.Println("Модифицированный последовательный алгоритм поиска")
	startSequential := time.Now()
	basePrimes := sieveOfEratosthenes(limit)

	_ = parallelCheckPrimesTask0(start, basePrimes)
	return time.Since(startSequential)
}

func algorithm1() time.Duration {
	fmt.Println("Параллельный алгоритм №1: декомпозиция по данным")
	startSequential := time.Now()
	basePrimes := sieveOfEratosthenes(limit)

	var wg sync.WaitGroup
	var mutex sync.Mutex
	var result []int

	rangeSize := (end - start + 1) / numThreads
	wg.Add(numThreads)

	for i := 0; i < numThreads; i++ {
		subStart := start + i*rangeSize
		subEnd := subStart + rangeSize - 1
		if i == numThreads-1 {
			subEnd = end
		}
		go parallelCheckPrimes(subStart, subEnd, basePrimes, &wg, &mutex, &result)
	}

	wg.Wait()
	return time.Since(startSequential)
}

func algorithm2() time.Duration {
	fmt.Println("Параллельный алгоритм №2: декомпозиция набора простых чисел")
	startSequential := time.Now()
	basePrimes := sieveOfEratosthenes(limit)

	var wg sync.WaitGroup
	var mutex sync.Mutex
	var result []int

	wg.Add(numThreads)

	// Разделение базовых простых чисел на части для каждого потока
	basePrimesPerThread := len(basePrimes) / numThreads
	for i := 0; i < numThreads; i++ {
		subStart := i * basePrimesPerThread
		subEnd := subStart + basePrimesPerThread - 1
		if i == numThreads-1 {
			subEnd = len(basePrimes) - 1
		}
		basePrimesSubset := basePrimes[subStart : subEnd+1]

		go parallelCheckPrimesWithSubsetOfPrimes(start, end, basePrimesSubset, &wg, &mutex, &result)
	}

	wg.Wait()
	return time.Since(startSequential)
}

func algorithm3() time.Duration {
	fmt.Println("Параллельный алгоритм №3: применение пула потоков")
	startSequential := time.Now()
	basePrimes := sieveOfEratosthenes(limit)

	var mutex sync.Mutex
	var result []int

	var wg sync.WaitGroup
	wg.Add(numThreads)

	rangeSize := (end - start + 1) / numThreads
	for i := 0; i < numThreads; i++ {
		subStart := start + i*rangeSize
		subEnd := subStart + rangeSize - 1
		if i == numThreads-1 {
			subEnd = end
		}
		go parallelCheckPrimes(subStart, subEnd, basePrimes, &wg, &mutex, &result)
	}

	wg.Wait()
	return time.Since(startSequential)
}

func algorithm4() time.Duration {
	fmt.Println("Параллельный алгоритм №4: последовательный перебор простых чисел")
	startSequential := time.Now()
	basePrimes := sieveOfEratosthenes(limit)

	var mutex sync.Mutex
	var result []int
	var currentIndex int

	var wg sync.WaitGroup
	wg.Add(numThreads)

	for i := 0; i < numThreads; i++ {
		go parallelCheckPrimesWithPrimePerThread(start, end, basePrimes, &wg, &mutex, &result, &currentIndex)
	}

	wg.Wait()
	return time.Since(startSequential)
}

func measureSequential() time.Duration {
	startTime := time.Now()
	basePrimes := sieveOfEratosthenes(limit)
	_ = sequentialCheckPrimes(start, end, basePrimes) // Замените аргументы на ваши тестовые данные
	duration := time.Since(startTime)
	return duration
}

func main() {
	sequentialDuration := measureSequential()
	parallelDurationAlgorithm0 := algorithm0()
	parallelDurationAlgorithm1 := algorithm1()
	parallelDurationAlgorithm2 := algorithm2()
	parallelDurationAlgorithm3 := algorithm3()
	parallelDurationAlgorithm4 := algorithm4()

	// Расчет ускорения и эффективности
	speedup0 := float64(sequentialDuration) / float64(parallelDurationAlgorithm0)
	speedup1 := float64(sequentialDuration) / float64(parallelDurationAlgorithm1)
	speedup2 := float64(sequentialDuration) / float64(parallelDurationAlgorithm2)
	speedup3 := float64(sequentialDuration) / float64(parallelDurationAlgorithm3)
	speedup4 := float64(sequentialDuration) / float64(parallelDurationAlgorithm4)

	efficiency0 := speedup0 / float64(numThreads)
	efficiency1 := speedup1 / float64(numThreads)
	efficiency2 := speedup2 / float64(numThreads)
	efficiency3 := speedup3 / float64(numThreads)
	efficiency4 := speedup4 / float64(numThreads)

	fmt.Printf("Время выполнения последовательного кода: %v\n", sequentialDuration)
	fmt.Println()

	fmt.Printf("Время выполнения модифицированного алгоритма: %v\n", parallelDurationAlgorithm0)
	fmt.Printf("Ускорение модифицированного алгоритма: %.2f\n", speedup0)
	fmt.Printf("Эффективность модифицированного алгоритма: %.2f\n", efficiency0)
	fmt.Println()

	fmt.Printf("Время выполнения параллельного алгоритма №1 декомпозиция по данным: %v\n", parallelDurationAlgorithm1)
	fmt.Printf("Ускорение модифицированного алгоритма: %.2f\n", speedup1)
	fmt.Printf("Эффективность модифицированного алгоритма: %.2f\n", efficiency1)
	fmt.Println()

	fmt.Printf("Время выполнения параллельного алгоритма №2 декомпозиция набора простых чисел: %v\n", parallelDurationAlgorithm2)
	fmt.Printf("Ускорение модифицированного алгоритма: %.2f\n", speedup2)
	fmt.Printf("Эффективность модифицированного алгоритма: %.2f\n", efficiency2)
	fmt.Println()

	fmt.Printf("Время выполнения параллельного алгоритма №3 применение пула потоков: %v\n", parallelDurationAlgorithm3)
	fmt.Printf("Ускорение модифицированного алгоритма: %.2f\n", speedup3)
	fmt.Printf("Эффективность модифицированного алгоритма: %.2f\n", efficiency3)
	fmt.Println()

	fmt.Printf("Время выполнения параллельного алгоритма №4 последовательный перебор простых чисел: %v\n", parallelDurationAlgorithm4)
	fmt.Printf("Ускорение модифицированного алгоритма: %.2f\n", speedup4)
	fmt.Printf("Эффективность модифицированного алгоритма: %.2f\n", efficiency4)
	fmt.Println()
}
