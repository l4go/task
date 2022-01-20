# type Mission
複雑な依存関係をもつgoroutine群を管理するためのモジュールです。  
階層構造の依存関係をもつ処理群を、一連のミッションとして見立てて管理することを目指したので、このような名前になってます。

goroutineの関係が複雑になると、go標準のパッケージ(sync.WaitGroupやcontext.Context)で実装すると、select文まみれになるので、その対策として作られました。  
階層的に管理できるので、グループ単位部分的なキャンセルも、全体を一斉にキャンセルすることも可能になっています。

go標準のsync.WaitGroupおよび、context.Contextとの大きな違いは以下のとおり。

* 追加のコード無しに、複数のタスクグループを連携させることが出来る。(sync.WaitGroupとの違い)
* 安全に動的なgoroutineの追加に対応できる。(sync.WaitGroupとの違い)
* 複数の箇所からキャンセル出来る。(context.Contextとの違い)
* context.ContextのDeadline()やValue()という余計な機能はない。
    * タイムアウト処理は`github.com/l4go/timer`を使う前提。

## import
```go
import "github.com/l4go/task"
```
vendoringして使うことを推奨します。

## 機能の概略
Missionが、提供する機能は、以下のものとその組み合わせです。

1. タスクの終了をDone()メソッドで、親のタスクに通知する。
1. 子供のタスクの終了を、親のDone()で待ってから終了する。
1. Cancel()メソッドで、子孫のタスクにキャンセルを通知する。
1. Abort()メソッドで、先祖も含めた全部のタスクにキャンセルを通知する。
1. 子供のタスク終了の通知を、Done()、Activate()、及び、Recv()メソッドの実行まで、遅延させる。
1. キャンセル通知をRecvCancel()からのchannelで受け取って、処理の中断を実装できるようにする。
1. 子供のタスク終了をRecvDone()もしくはRecv()からのchannelで受け取り、selectで他の処理と一緒に終了待ち出来る。

sync.WaitGroupでは、生成タスク数が未定の場合、複数タスク生成完了前に生成済みタスクがすべて終了すると、間違った完了通知が起こる問題があります。  
この問題への対策が、子供のタスク終了の通知を遅延する機能です。

## 利用サンプル
孫のworkerがある２重のworker構造へ、goroutineにキャンセル処理を実装したコードサンプルです。  
比較的シンプルに、複雑なgorutineの完了管理が実装できてることが分かると思います。

[example](../examples/ex_mission/ex_mission.go)

## メソッド概略

### func NewMission() \*Mission
トップのタスク管理の\*Missionを生成します。
\*Mission生成直後は、小の\*Missionの終了通知は遅延状態になっています。

New()メソッドで、子供の\*Missionを生成できます。

### func (p \*Mission) New() \*Mission
小タスク管理用のMissionを生成します。

### func (p \*Mission) NewCancel() Canceller
サブタスク管理用の\*Cancelを生成します。

生成された\*Cancelは、\*Mission側のCancel処理に連動する状態になっています。
\*Cancel側のCancel処理は、Mission側へ連動させていないので、サブタスクの影響を受けません。

### func (p \*Mission) Link(Canceller)
\*MissionのCancel処理をCancellerへ連動させます。

### func (m \*Mission) Done()
子のMission終了通知を有効にし、子のMissionの終了を待ってから、自分のMissionを終了させます。

### func (m \*Mission) NowaitDone()
子のMission終了通知を有効にし、自分のMissionを終了させます。子のMissionの終了は待ちません。

子のMissionの終了を待てない特殊な場合に使いますが、非推奨です。

### func (m \*Mission) Cancel()
キャンセル情報を、自分及び子孫のMissionに伝播させます。親および祖先には伝播しません。複数箇所から実行されても問題なく動作します。

### func (m \*Mission) Abort()
キャンセル情報を、トップのMissionも含めてすべてのMissionに伝播します。

トップのMissionでCancel()メソッドを実行した時と同じ動作です。

### func (m \*Mission) IsCanceled() bool
キャンセルの有無を確認します。

処理を開始する前に実施の判断をしたい時に利用します。

### func (m \*Mission) RecvCancel() <-chan struct{}
キャンセルを通知するchannelを取得します。

キャンセルされるとこのchannelがcloseされます。

### func (m \*Mission) Recv() <-chan struct{}
子のMission終了通知を有効にしてから、子のMissionがすべて終了した状態を示すchannelを取得します。
子のMissionをすべて生成する(New()メゾッドを呼び終わる)前にRecv()メソッドを呼ぶと誤動作します。  
別のgoroutingなどで事前に終了待ちの処理を実行したい場合は、RecvDone()メソッドを利用して、後から子のMissionの終了通知を有効にください。

通常は、Done()メソッドもしくは、RecvDone()メソッドの利用が推奨されます。
他のタスク管理と連携させる処理のみに利用するメゾッドで、通常は、Down()メソッドを使えば解決します。  
main()関数などで、終了用signalと並行して待つ状況などが想定された適切な用途です。

### func (m \*Mission) RecvDone() <-chan struct{}
子のMission終了通知するchannelを取得します。

子のMissionを生成する前に、別のgoroutineから通知用のチャンネルを利用する場合に使います。  
Activate()メソッド、Done()メソッド、Recv()メソッドのどれかが実行されると、通知は開始されます。

### func (m \*Mission) Activate()
子のMissionがすべて終了した状態の通知を開始します。

Recv()メソッドやDone()メソッドを実行する前に、通知を開始したい場合に利用します。

### func (m \*Mission) AsContext() \*context.Context
Missionをgo標準のcontext\.Contextインタフェース互換形式に変換します。

go標準の関数と連携するために使います。

### func (m \*Mission) Parson() \*Mission
親のMissionを取得します。親のMissionが無い場合は、nilが返ります。

親のMissionをキャンセルするために使います。  
全部のMissionをキャンセルするのには、Abort()メソッドを利用してください。
