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
  <dd>Fu like 30 Fu.</dd>

  <dt>fan</dt>
  <dd>Han.</dd>

  <dt>score change</dt>
  <dd>Change in a player's score in a certain round.</dd>

  <dt>score changes</dt>
  <dd>A 4-element vector (array) where scoreChanges[player.id] is the player's score change. e.g., [8000, -8000, 0, 0]</dd>

  <dt>player ID</dt>
  <dd>Player ID between 0 and 3.</dd>

  <dt>hora factors</dt>
  <dd>
    A 4-element vector (array) where horaPoints * horaFactors[player.id] = scoreChanges[player.id].<br>
    [1, -1, 0, 0] for Ron, [1, -1/2, -1/4, -1/4] for non-dealer's Tsumo, etc.
  </dd>

  <dt>furo</dt>
  <dd>Meld. Call.</dd>

  <dt>pai ID (pid)</dt>
  <dd>An integer between 0 and 33 representing the type of tile.</dd>

  <dt>action</dt>
  <dd>Self-draw, discard, Chii, etc.</dd>

  <dt>metric</dt>
  <dd>Various statistics/estimates about the outcome of a certain action (e.g., discarding 2m).</dd>

  <dt>count vector</dt>
  <dd>A data structure representing a multi set of tiles. An array where countVector[pai.id] is the number of pai.</dd>

  <dt>bit vectors</dt>
  <dd>A data structure representing a multi set of tiles. An array of BitVector where bitVectors[i][pai.id] = (count(pai) > i).</dd>

  <dt>rank</dt>
  <dd>Rank. An integer between 1 and 4.</dd>

  <dt>statistics (stats)</dt>
  <dd>Statistics collected in advance from game records.</dd>
</dl>
