# 仕組み

- 画像処理で数列の紙を読み込む。
- どの順番で数列を読み込むか指定
- 順番で読みこみ復号処理

```sh
export LIBRARY_PATH="/opt/homebrew/lib"
export CPATH="/opt/homebrew/include"
```

opencv. load video -> send frame -> python. recognize numeric image -> grpc -> golang
