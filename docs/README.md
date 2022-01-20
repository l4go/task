# golib/task ライブラリ

goroutineの処理(task)の管理を支援するモジュール群です。

* [task.Mission](Mission.md)
    * 親子関係のあるgoroutine群の終了管理を、簡素かつ柔軟に行うためのモジュールです。
* [task.Cancel](Cancel.md)
    * 1つのgoroutineから、複数のgoroutineのキャンセルを行うためのモジュールです。
* [task.Pool](Pool.md)
    * goroutine pool(thread poolのgoroutine版)を実装するモジュールです。
  
# type task.Canceller interface
task.Cancelおよびtask.Missionのキャンセル機能のみを取り出した互換コード用のinterface(task.Cancelおよびtask.Missionから変換可能)

# func task.IsCanceled(cc task.Canceller) bool
キャンセルの有無を確認します。task.Cancelおよびtask.Missionの両方に利用できます。  
処理を開始する前に実施の判断をしたい時に利用します。
