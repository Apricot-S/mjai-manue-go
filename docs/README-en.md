\* This is the translation of [the original document (Japanese)](https://github.com/gimite/mjai-manue/blob/master/README.md). Additionally, the broken links have been fixed.

## Overview

Mahjong AI for [Mjai Mahjong AI match server](https://gimite.net/pukiwiki/index.php?Mjai%20%E9%BA%BB%E9%9B%80AI%E5%AF%BE%E6%88%A6%E3%82%B5%E3%83%BC%E3%83%90)

## How It Works

[Sample game record of a self-match](https://gimite.net/mjai/samples/manue011.tonnan/2013-11-26-143619.mjson.html)

First, calculate the following scores for each possible discard. These scores can be seen in the debug output of the game record above.

* horaProb / Hora probability / Win rate
  * The probability of winning in this round if that tile is discarded.
  * Calculated using Monte Carlo simulations. If there are N self-draws until the end of the round, draw N tiles randomly and check if a winning hand can be formed with 13 tiles in the hand plus N tiles. This process is repeated 1000 times.
  * In practice, to speed up the process, "the necessary tiles to win from the current hand" are pre-calculated, and it is checked if these necessary tiles are included in the N tiles drawn randomly.
* avgHoraPt / Average hora points / Average win points
  * The average winning points when winning.
  * Calculated using Monte Carlo simulations simultaneously with horaProb. It is the average points of the winning hands formed with 13 tiles in the hand plus N tiles.
* unsafeProb / Unsafe probability / Deal-in rate
  * The probability of dealing into another player's hand with that discard.
  * Currently, only considers dealing into a player who has declared Riichi.
  * Estimated using decision tree learning. Features include "Honors", "Suji", etc. Training data is from Tenhou's game records. [Analysis of Mahjong dangerous tile using statistics](https://gimite.net/pukiwiki/index.php?%E7%B5%B1%E8%A8%88%E3%81%AB%E3%82%88%E3%82%8B%E9%BA%BB%E9%9B%80%E5%8D%B1%E9%99%BA%E7%89%8C%E5%88%86%E6%9E%90) for more information.
* avgHojuPt / Average hoju points / Average deal-in points
  * The average points paid if dealing into another player's hand.
  * Currently a fixed value of 6265 points derived from self-match logs. Not included in the game record debug output.

以上の数値から、この局で自分が得る点数の期待値(expPt)を求めることができます。

* expPt = (1 - unsafeProb) * horaProb * avgHoraPt - unsafeProb * avgHojuPt

このexpPtが最大となる打牌を採用します。

「鳴くか、鳴かないか」「リーチか、ダマか」も同様の方法で判断します。

## License

"New BSD Licence"\[sic]

## Author

[Hiroshi Ichikawa](https://gimite.net/pukiwiki/index.php?%E9%80%A3%E7%B5%A1%E5%85%88)
