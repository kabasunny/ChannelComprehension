package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func worker(id int, errCh chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()

	randGen := rand.New(rand.NewSource(time.Now().UnixNano())) // ローカル乱数生成器を初期化

	startTime := time.Now()
	fmt.Printf("[%s] Worker %d started\n", startTime.Format("2006-01-02 15:04:05.000"), id)

	time.Sleep(time.Second) // 1秒スリープ

	// 25% の確率でエラーを発生
	if randGen.Intn(4) == 0 {
		errMsg := fmt.Sprintf("[%s] Worker %d encountered an error", time.Now().Format("2006-01-02 15:04:05.000"), id)
		// エラーチャネルに送信（タイムアウト付き）
		select {
		case errCh <- errMsg:
			fmt.Println(errMsg)
		case <-time.After(time.Millisecond * 100): // 100ms待っても送信できなかったら
			fmt.Printf("[%s] Error channel is full, dropping error from worker %d\n", time.Now().Format("2006-01-02 15:04:05.000"), id)
			return // エラーを破棄してゴルーチンを終了
		}

		time.Sleep(3 * time.Second) // 3秒スリープ (エラーからの復旧をシミュレート)
		fmt.Printf("[%s] Worker %d recovered from error\n", time.Now().Format("2006-01-02 15:04:05.000"), id)
		return
	}

	fmt.Printf("[%s] Worker %d finished successfully\n", time.Now().Format("2006-01-02 15:04:05.000"), id)
}
