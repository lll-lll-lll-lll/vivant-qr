![](./sample.png)

# 何これ

付箋に写っている数列を特定のパターンでしか中身の値が取得できな cli

# 微妙なところ

ocr の精度に依存している...

# 作り始めたきっかけ

vivant が面白くて、あるワンシーンで付箋に書いてある数列をカメラで映すとあるサーバの場所がわかる的なシーンがあり
それに近いものを作りたくなったから

# use

```sh
vivantqr --apikey=OO --file ./output.png --read true
vivantqr --apikey=OO --file ./output.png --write true
```
