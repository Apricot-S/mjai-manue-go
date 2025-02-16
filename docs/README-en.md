\* This is the translation of [the original document (Japanese)](https://github.com/gimite/mjai-manue/blob/master/README.md). Additionally, the broken links have been fixed.

## Overview

[Mjai Mahjong AI competition server](https://gimite.net/pukiwiki/index.php?Mjai%20%E9%BA%BB%E9%9B%80AI%E5%AF%BE%E6%88%A6%E3%82%B5%E3%83%BC%E3%83%90) 用の麻雀AIです。

## How It Works

[Sample game record of a self-match](https://gimite.net/mjai/samples/manue011.tonnan/2013-11-26-143619.mjson.html)

まず、それぞれの打牌をした場合について、以下の数値を算出します。これらのスコアは、上の牌譜のデバッグ出力で確認できます。

* horaProb / Hora probability / 和了率
  * その打牌をした場合に、この局で自分が和了できる確率。
  * モンテカルロで求める。終局までにNツモあるとすると、ランダムにN枚引いて、手牌13枚+N枚で和了を作れるかどうかをチェック。これを1000回繰り返す。
  * 実際には高速化のために「今の手牌から和了するための必要牌」をあらかじめ求めておき、ランダムに引いたN枚に必要牌が含まれるかをチェックしている。
* avgHoraPt / Average hora points / 平均和了点
  * 自分が和了した場合の平均和了点。
  * horaProbと同時にモンテカルロで求める。手牌13枚+N枚で作れた和了の点数の平均。
* unsafeProb / Unsafe probability / 放銃率
  * その打牌で誰かに放銃する確率。
  * 今のところ、リーチしている人への放銃だけを考慮。
  * 決定木学習を使って推定。特徴量は「字牌」「スジ」など。学習データは天鳳の牌譜。[Analysis of Mahjong dangerous tile using statistics](https://gimite.net/pukiwiki/index.php?%E7%B5%B1%E8%A8%88%E3%81%AB%E3%82%88%E3%82%8B%E9%BA%BB%E9%9B%80%E5%8D%B1%E9%99%BA%E7%89%8C%E5%88%86%E6%9E%90)参照。
* avgHojuPt / Average hoju points / 平均放銃点
  * 放銃した場合に払う額の平均。
  * 今のところは自己対戦のログから求めた固定値6265点。牌譜のデバッグ出力にはない。

以上の数値から、この局で自分が得る点数の期待値(expPt)を求めることができます。

* expPt = (1 - unsafeProb) * horaProb * avgHoraPt - unsafeProb * avgHojuPt

このexpPtが最大となる打牌を採用します。

「鳴くか、鳴かないか」「リーチか、ダマか」も同様の方法で判断します。

## License

"New BSD Licence"\[sic]

## Author

[Hiroshi Ichikawa](https://gimite.net/pukiwiki/index.php?%E9%80%A3%E7%B5%A1%E5%85%88)
