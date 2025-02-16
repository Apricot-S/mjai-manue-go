\* This is a translation of [the original document (Japanese)](https://github.com/gimite/mjai-manue/blob/master/doc/terminology.txt). It has also been changed from plain text to a description list.

<dl>
  <dt>probability distribution (probDist, dist)</dt>
  <dd>Probability distribution.</dd>

  <dt>probability (prob)</dt>
  <dd>Probability.</dd>

  <dt>hora</dt>
  <dd>Win.</dd>

  <dt>score</dt>
  <dd>The score of a player at a certain point in time.</dd>

  <dt>points</dt>
  <dd>The points of the win.</dd>

  <dt>fu</dt>
  <dd>Fu like 30 fu.</dd>

  <dt>fan</dt>
  <dd>Han.</dd>

  <dt>score change</dt>
  <dd>Change in a player's score in a certain round.</dd>

  <dt>score changes</dt>
  <dd>A 4-element vector (array) where scoreChanges[player.id] is the player's score change. e.g., [8000, -8000, 0, 0]</dd>

  <dt>player ID</dt>
  <dd>Player ID between 0 and 3.</dd>

  <dt>hora factors</dt>
  <dd>horaPoints * horaFactors[player.id] = scoreChanges[player.id] となるような4要素のベクトル(配列)。 <br>
  ロンなら[1, -1, 0, 0]、子のツモなら[1, -1/2, -1/4, -1/4]など。</dd>

  <dt>furo</dt>
  <dd>Meld. Call.</dd>

  <dt>pai ID (pid)</dt>
  <dd>An integer between 0 and 33 representing the type of tile.</dd>

  <dt>action</dt>
  <dd>Self-draw, discard, chi, etc.</dd>

  <dt>metric</dt>
  <dd>あるアクション(2mを打牌、など)の結果についての様々な統計値/推定値。</dd>

  <dt>count vector</dt>
  <dd>牌のmulti setを表すデータ構造の1つ。countVector[pai.id]がpaiの個数となるような配列。</dd>

  <dt>bit vectors</dt>
  <dd>牌のmulti setを表すデータ構造の1つ。bitVectors[i][pai.id] = (count(pai) > i)となるようなBitVectorの配列。</dd>

  <dt>rank</dt>
  <dd>Rank. An integer between 1 and 4.</dd>

  <dt>statistics (stats)</dt>
  <dd>Statistics collected in advance from game records.</dd>
</dl>
