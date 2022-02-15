# type Cancel
任意のgoroutineから、任意の1つ以上のgoroutineのキャンセルを行うためのモジュールです。  
Missionの軽量版で、キャンセルに特化した機能のみを有します。  
Cancelを持たせたgoroutine全てに対して、どこかのgoroutineに一度にキャンセルすることができます。

context.Contextで同様のことを行う場合、Cancelの活用されかたによってpanicする問題があります。  
Cancelは上記問題を解決しつつ、goroutineキャンセル処理だけに特化したモジュールです。  

go標準のcontext.Contextとの違いは以下の通り
* 複数の箇所からキャンセルが可能
* キャンセル以外の機能を有さない

## import
```go
import "github.com/l4go/task"
```
vendoringして使うことを推奨します。

## 機能の概略
Cancelが、提供する機能です。

1. Cancel()メソッドで、全てのCancelにキャンセルを通知します。
1. キャンセル通知をRecvCancel()で受け取って、しょrいの中断を実装できるようにしています。

sync.WaitGroupでは、生成タスク数が未定の場合、複数タスク生成完了前に生成済みタスクがすべて終了すると、間違った完了通知が起こる問題があります。  
この問題への対策が、子供のタスク終了の通知を遅延する機能です。

## 利用サンプル
workerと、その子worker(`sub_worker`)がある構造へ、goroutineにキャンセル処理を実装したコードサンプルです。  

[example](../examples/ex_cancel/ex_cancel.go)

## メソッド概略

### func NewCancel() \*Cancel

Cancelを生成します。

### func (c \*Cancel) Cancel()

キャンセル情報を設定します。

### func (c \*Cancel) RecvCancel() <-chan struct{}

キャンセルを通知するchannelを取得します。

### func (c \*Cancel) AsContext() \*context.Context

Missionをgo標準のcontext\.Contextインタフェース互換形式に変換します。
