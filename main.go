package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func main() {
	randGen := rand.New(rand.NewSource(time.Now().UnixNano())) // ローカル乱数生成器を初期化
	errCh := make(chan string, 2)                              // バッファサイズ2のエラーチャネル
	var wg sync.WaitGroup

	// タイムアウト用のコンテキストを作成 20秒
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	for i := 1; i < 10; i++ {
		select {
		case <-timeoutCtx.Done(): // <-chan struct{}型のチャネルを返す
			fmt.Println("Timeout reached, exiting...")
			return // タイムアウトしたらループを抜ける
		default:
			fmt.Println("i...", i)
			// 通常の処理を続行
		}

		numWorkers := randGen.Intn(5) + 4 // 5〜10 のランダムな数のゴルーチンを起動
		for j := 1; j < numWorkers; j++ {
			wg.Add(1)
			go worker(i*10+j, errCh, &wg) // ゴルーチンIDが重複しないように調整
		}

		// エラーチャネルを監視
		select {
		case errMsg := <-errCh:
			fmt.Println("Main goroutine received error:", errMsg)
		default: // チャネルが空の場合は何もしない
		}

		time.Sleep(time.Millisecond * 1000) // メインループを少しスローダウン 1000ミリ秒
	}

	// 全てのワーカーが終了するまで待機し、その後エラーチャネルをクローズ
	go func() {
		fmt.Println("Waiting for all workers to finish...")
		wg.Wait()
		fmt.Println("All workers finished.")

		// エラーの集計
		errorCount := 0
		for errMsg := range errCh {
			fmt.Println("Processing remaining error:", errMsg)
			errorCount++
		}

		fmt.Printf("Total errors encountered: %d\n", errorCount)
		fmt.Println("Closing error channel.")
		close(errCh)
	}()

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
