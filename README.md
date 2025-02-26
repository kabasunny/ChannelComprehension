# ChannelComprehension


## 仕様

### メインループ
- 200回のループを実行

### ゴルーチン
- 各ループで、ランダムな数 (1〜3) のゴルーチンを起動

### ゴルーチンの処理
1. 1秒間スリープ
2. 10% の確率でエラーを発生
   - **エラーが発生した場合**:
     - エラーメッセージをエラーチャネル (errCh) に送信
     - 3秒間スリープ (エラーからの復旧をシミュレート)
     - 復旧メッセージを標準出力に出力
   - **エラーが発生しなかった場合**:
     - 正常終了メッセージを標準出力に出力

### エラーチャネル
- バッファサイズは5
- メインゴルーチンは、エラーチャネルを監視し、エラーメッセージを受信したら標準出力に出力

### タイムアウト
- 処理全体で30秒経過したら、強制的に終了

### 出力
- 各ゴルーチンの開始、正常終了、エラー発生、エラー復旧のメッセージを、タイムスタンプ付きで標準出力に出力
- エラーチャネルがフルになった場合は、その旨を標準出力に出力

## コードのポイント

### エラーチャネルのバッファ
- `errCh := make(chan string, 5)` で、バッファサイズを5に設定

### ゴルーチン ID
- 各ゴルーチンにユニークな ID (i*10 + j) を割り当てているため、どのゴルーチンでエラーが発生したかを区別可能

### エラーチャネルへの送信 (タイムアウト付き)
- `worker` 関数内で、`select` 文を使ってエラーチャネルへの送信にタイムアウト (100ms) を設けているため、チャネルがフルの場合にゴルーチンが無期限にブロックされるのを防ぎ、エラーを破棄して処理を続行可能

### main関数でのタイムアウト
- 全体処理が30秒を超えたら強制終了

### メインループのスローダウン
- `time.Sleep(time.Millisecond * 100)` をメインループに追加しているため、ゴルーチンの処理が追いつかないほど高速にループが回るのを防ぎ、出力を見やすくしている

### エラーチャネルの監視
- メインゴルーチンは、各ループで `select` 文を使ってエラーチャネルを監視しているため、エラーメッセージを受信したら、それを出力可能
- `default` ケースがあるため、チャネルが空の場合でもブロックされない

### 処理終了後の後処理
- `wg.Wait()` で全てのゴルーチンの処理が終わった後、チャネルに残存しているエラーがないか確認して出力

## 実行方法

1. 上記のコードを `main.go` ファイルとして保存
2. ターミナルで `go run main.go` を実行

