# type Finish
任意のgoroutineから、任意の1つ以上のgoroutineの完了を行うためのモジュールです。
Missionの軽量版で、完了通知に特化した機能のみを有します。  
Doneを持たせたgoroutine全てに対して、どこかのgoroutineから一度に完了を通知することができます。  

context.Contextで同様のことを行う場合、Doneの活用されかたによってpanicする問題があります。  
Doneは上記問題を解決しつつ、goroutineの完了通知処理だけに特化したモジュールです。  

go標準のcontext.Contextとの違いは以下の通り
* 複数の箇所から完了通知が可能
* 完了通知以外の機能を有さない

## import
```go
import "github.com/l4go/task"
```
vendoringして使うことを推奨します。

## 機能の概略
Finishが、提供する機能です。

1. Done()メソッドで、全てのDone待ちに完了を通知します。
1. 完了の通知をRecvDone()で受け取って、処理の完了を実装できるようにしています。
	* なお、受け取った側の処理の完了を通知元のgoroutineで待ちたい場合は、より高機能な[task.Mission](./Mission.md) を用いる必要があります

sync.WaitGroupでは、生成タスク数が未定の場合、複数タスク生成完了前に生成済みタスクがすべて完了すると、間違った完了通知が起こる問題があります。  
この問題への対策が、子供のタスク完了の通知を遅延する機能です。

## 利用サンプル
複数のworkerがある構造のgoroutineに完了処理を実装したコードサンプルです。  

[example](../examples/ex_finish/ex_finish.go)

## メソッド概略

### func NewFinish() \*Finish

Finishを生成します。

### func (f \*Finish) Done()

完了を通知します。  

### func (f \*Finish) RecvDone() <-chan struct{}

完了を通知するchannelを取得します。

### func (f \*Finish) AsContext() \*context.Context

Finishをgo標準のcontext\.Contextインタフェース互換形式に変換します。
