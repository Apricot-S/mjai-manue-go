**This is a translation of [the original document (Japanese)](https://github.com/gimite/mjai-manue/blob/master/doc/terminology.txt). It has also been changed from plain text to a description list.**

<dl>
  <dt>probability distribution (probDist, dist)</dt>
  <dd>確率分布。</dt>

  <dt>probability (prob)</dt>
  <dd>確率。</dt>

  <dt>hora</dt>
  <dd>和了。</dt>

  <dt>score</dt>
  <dd>あるプレーヤのある時点での点数。</dt>

  <dt>points</dt>
  <dd>和了の点数。</dt>

  <dt>fu</dt>
  <dd>30符とかの符。</dt>

  <dt>fan</dt>
  <dd>飜。</dt>

  <dt>score change</dt>
  <dd>あるプレーヤのある局におけるscoreの変動。</dt>

  <dt>score changes</dt>
  <dd>scoreChanges[player.id]がplayerのscore changeとなるような4要素のベクトル(配列)。 e.g., [8000, -8000, 0, 0]</dt>

  <dt>player ID</dt>
  <dd>0～3のプレーヤID。</dt>

  <dt>hora factors</dt>
  <dd>horaPoints * horaFactors[player.id] = scoreChanges[player.id] となるような4要素のベクトル(配列)。 <br>
  ロンなら[1, -1, 0, 0]、子のツモなら[1, -1/2, -1/4, -1/4]など。</dd>

  <dt>furo</dt>
  <dd>副露。なき。</dt>

  <dt>pai ID (pid)</dt>
  <dd>牌の種類を表す0～33の整数。</dt>

  <dt>action</dt>
  <dd>自摸とか打牌とかチーとか。</dt>

  <dt>metric</dt>
  <dd>あるアクション(2mを打牌、など)の結果についての様々な統計値/推定値。</dt>

  <dt>count vector</dt>
  <dd>牌のmulti setを表すデータ構造の1つ。countVector[pai.id]がpaiの個数となるような配列。</dt>

  <dt>bit vectors</dt>
  <dd>牌のmulti setを表すデータ構造の1つ。bitVectors[i][pai.id] = (count(pai) > i)となるようなBitVectorの配列。</dt>

  <dt>rank</dt>
  <dd>順位。1～4の整数。</dt>

  <dt>statistics (stats)</dt>
  <dd>あらかじめ牌譜から収集された統計情報。</dt>
</dl>
