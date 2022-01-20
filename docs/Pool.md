# type Pool
goroutine pool(thread poolのgoroutine版)のためのモジュールです。  
Missionと連携して、より複雑な処理の管理を可能にします。

## import
```go
import "github.com/l4go/task"
```
vendoringして使うことを推奨します。

## 利用サンプル
Poolを使って実行(Do())と、弱実行(WeakDo())を行うサンプルコードです。  
弱実行のコードがシンプルなのがわかるかと思います。

[example](../examples/ex_pool/ex_pool.go)

## メソッド概略

### type PoolFunc func(*Mission, ...interface{})
goroutine poolで実行される関数の型です。
Poolがキャンセルされると、渡したMissionをキャンセルした上で呼び出します。

### type PoolWeakFunc func(...interface{})
goroutine poolで弱実行される関数の型です。  
弱実行は、キャンセルされると実行自体が呼び出されないで捨てます。

キャンセル処理も、リソース開放も必要ない処理をシンプルに実装するための関数形式です。

### func NewPool(m *Mission, cnt int) *Pool
指定した数のgoroutine poolを作成します。

### func (p *Pool) Close()
すべての処理がMission経由でキャンセルされます。

### func (p *Pool) Do(f PoolFunc, m *Mission, args ...interface{})
PoolFuncをgoroutine poolで実行します。渡したMissionを使えば個別にキャンセルすることも可能です。
また、渡したMissionはPoolからもキャンセルされることがあります。

### func (p *Pool) WeakDo(f PoolWeakFunc, args ...interface{})
PoolWeakFuncをgoroutine poolで弱実行します。
