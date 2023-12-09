package circuit_breaker

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

//type Circuit func(context.Context) (string, error)
//
//func Breaker(circuit Circuit, failureThreshold uint) Circuit {
//	var consecutiveFailures int = 0
//	var lastAttempt = time.Now()
//	var m sync.RWMutex
//
//	return func(ctx context.Context) (string, error) {
//		m.RLock() // Установить "блокировку чтения"
//		d := consecutiveFailures - int(failureThreshold)
//
//		if d >= 0 {
//			shouldRetryAt := lastAttempt.Add(time.Second * 2 << d)
//			if !time.Now().After(shouldRetryAt) {
//				m.RUnlock()
//				return "", errors.New("service unreachable")
//			}
//		}
//
//		m.RUnlock() // Освободить блокировку чтения
//
//		response, err := circuit(ctx) // Послать запрос, как обычно
//
//		m.Lock() // Заблокировать общие ресурсы
//		defer m.Unlock()
//
//		lastAttempt = time.Now() // Зафиксировать время попытки
//		if err != nil {          // Если Circuit вернула ошибку,
//			consecutiveFailures++ // увеличить счетчик ошибок
//			return response, err  // и вернуть ошибку
//		}
//
//		consecutiveFailures = 0 // Сбросить счетчик ошибок
//		return response, nil
//	}
//}

// Mocked API function
func MockedAPI(ctx context.Context) (string, error) {
	if rand.Float32() < 0.5 { // 50% вероятность ошибки
		return "", errors.New("mocked API error")
	}
	return "success", nil
}

type Circuit func(context.Context) (string, error)

// Implementation of the Circuit Breaker pattern
func Breaker(circuit Circuit, failureThreshold uint) Circuit {
	var consecutiveFailures int = 0
	var lastAttempt = time.Now()
	var m sync.RWMutex

	return func(ctx context.Context) (string, error) {
		m.RLock()
		d := consecutiveFailures - int(failureThreshold)

		if d >= 0 {
			shouldRetryAt := lastAttempt.Add(time.Second * 2 << d) // экспоненциальный рост
			if !time.Now().After(shouldRetryAt) {
				m.RUnlock()
				return "", errors.New("service unreachable")
			}
		}

		m.RUnlock()

		response, err := circuit(ctx)

		m.Lock()
		defer m.Unlock()

		lastAttempt = time.Now()
		if err != nil {
			consecutiveFailures++
			return response, err
		}

		consecutiveFailures = 0
		return response, nil
	}
}

func TestCircuitBreaker() {
	rand.Seed(time.Now().UnixNano()) // Инициализация генератора случайных чисел

	breaker := Breaker(MockedAPI, 2) // Настройка Circuit Breaker с порогом в 2 ошибки

	for i := 0; i < 10; i++ {
		result, err := breaker(context.Background())
		if err != nil {
			fmt.Printf("Attempt %d: Failed - %s\n", i+1, err)
		} else {
			fmt.Printf("Attempt %d: Success - %s\n", i+1, result)
		}
		time.Sleep(500 * time.Millisecond)
	}
}
