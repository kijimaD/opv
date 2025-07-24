org-mode pomodoro viewer

## 概要

Emacsでorg-pomodoroを使って計測している。LAN内の別の端末から表示するために、HTMLでページを表示する。

## メモ

- 凝る必要はない
- kindle paper white
- 変動するデータだけをREST API(JSON形式)でやりとりし、更新する。30秒ごと

表示項目。

- タスクタイトル
- 残り時間(mm:ss)
- 残り時間のビジュアル表示

## 参考

(org-pomodoro-remaining-seconds) : 残り時間(秒)
org-pomodoro-length : 1セッションあたりの長さ(分)
(org-pomodoro-active-p) : pomodoroがアクティブかどうか
(cl-case org-pomodoro-state ...) : 現在のpomodoroステート
kd/pmd-today-point : 今日達成した数
org-clock-heading : 実施中のタイトル

Emacsとの連携は https://kijimad.github.io/.emacs.d/#org19f0828 を参考にする

## TODO

- [ ] Emacs pomodoroから残り時間だけ取得して(10秒ごと)、秒数カウントはクライアント側でやらせるのがいいのかもしれない。1秒間隔でリクエストを送るのは微妙
- タスクが終了したときにはクリア演出を表示したい
- [ ] サーバ側のイベントは頻度が低いのでリアルタイム通知したい(開始、終了)
- タイトルのPomodoro Timerはいらない。タスク名と秒数を大きく表示する
- 完了タスク数を、数字とタスク数分のトマトをビジュアル表示する。2個完了していれば2個のトマト
- 開始中は、背景を赤くしたい
- 5分タイマーのときも表示できるか確認する
