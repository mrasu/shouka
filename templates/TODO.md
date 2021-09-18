
* cliのinteractive部分を作る
* "sk-" のプレフィックスを "sk-yyyymmdd"にする (ECSのTaskの定義が昔の番号を保存しているせいで、同じ名前を使うとappspec.ymlの番号も変える必要が有るから)
* CircleCIを使ったり、Pipelineを使うようになるなどでツールやフェーズによって作るファイルが分岐するので、それに耐えられるようにする
* migrationのフェーズも作る
