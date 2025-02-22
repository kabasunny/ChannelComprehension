package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func worker(id int, errCh chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()

	startTime := time.Now()
	fmt.Printf("[%s] Worker %d started\n", startTime.Format("2006-01-02 15:04:05.000"), id)

	time.Sleep(time.Second) // 1秒スリープ

	// 10% の確率でエラーを発生
	if rand.Intn(10) == 0 {
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

func main() {
	rand.Seed(time.Now().UnixNano()) // 乱数シードを初期化

	errCh := make(chan string, 5) // バッファサイズ5のエラーチャネル
	var wg sync.WaitGroup

	// タイムアウト用のコンテキストを作成
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for i := 0; i < 200; i++ {
		select {
		case <-timeoutCtx.Done():
			fmt.Println("Timeout reached, exiting...")
			return // タイムアウトしたらループを抜ける
		default:
			// 通常の処理を続行
		}

		numWorkers := rand.Intn(3) + 1 // 1〜3 のランダムな数のゴルーチンを起動
		for j := 0; j < numWorkers; j++ {
			wg.Add(1)
			go worker(i*10+j, errCh, &wg) // ゴルーチンIDが重複しないように調整
		}

		// エラーチャネルを監視
		select {
		case errMsg := <-errCh:
			fmt.Println("Main goroutine received error:", errMsg)
		default: // チャネルが空の場合は何もしない
		}

		time.Sleep(time.Millisecond * 100) // メインループを少しスローダウン
	}

	wg.Wait() // 全てのゴルーチンの終了を待つ
	fmt.Println("All workers finished")

	// 残っているエラーメッセージを処理
	for {
		select {
		case errMsg := <-errCh:
			fmt.Println("Main goroutine received remaining error:", errMsg)
		default:
			fmt.Println("No more errors in the channel")
			return
		}
	}
}
