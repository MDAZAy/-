# Usage Guide

## Installation In TradingView

1. Open TradingView.
2. Go to `Pine Editor`.
3. Open the file `dazay_meta_ai_v6_elite.pine`.
4. Copy all code.
5. Paste it into `Pine Editor`.
6. Click `Save`.
7. Click `Add to chart`.

## Core Input Groups

### `META-AI Core`

- `Base Learning Rate`: controls how fast the model adapts.
- `Weight Memory`: higher values make adaptation smoother and slower.
- `Reward Lookback`: how many bars later the model judges signal quality.
- `Buy Threshold` / `Sell Threshold`: main trigger thresholds.
- `Cooldown Bars`: minimum bars between valid entries.

### `META-AI Pro Layer`

- `Use Higher Timeframe Bias`: adds directional confirmation from a higher timeframe.
- `Higher Timeframe`: timeframe used for top-down bias.
- `HTF Bias Weight`: strength of higher-timeframe influence.
- `Use Adaptive Signal Smoothing`: stabilizes final signal.
- `Signal Persistence Bars`: requires short-term continuation before accepting entries.
- `Long/Short Edge Bias`: lets the model adjust thresholds based on side-specific performance.

### `Quant Layer`

- `Use Regime Filtering`: blocks weak entries outside strong conditions.
- `Min Trend Strength (ADX)`: minimum trend strength.
- `Noise Sensitivity`: stronger filtering in noisy markets.
- `Dynamic Threshold Multiplier`: widens thresholds in difficult conditions.
- `Min Signal Quality`: minimum quality score required before entry.

### `META-AI Health / State`

- `Entropy Limit`: defines disagreement tolerance between internal engines.
- `Freeze On Instability`: blocks entries when the model is unstable.
- `Freeze Entropy Trigger`: hard-stop threshold for instability.

### `Risk Management`

- `Risk Mode`: ATR-based or percent-based stop logic.
- `Risk Size (SL)`: stop distance input.
- `Risk : Reward Ratio`: take-profit multiple.
- `High Volatility Stop Buffer`: widens stop under high volatility.
- `Use Confidence Position Scaling`: outputs a dynamic `size_scale`.

## What The Dashboard Shows

- confidence
- side
- entry / SL / TP / BE
- market regime
- ADX
- AI state
- entropy
- rolling accuracy
- quality score
- position size scale
- higher timeframe bias
- long/short side edge statistics

## Alert Payload

The alert JSON contains:

- `type`
- `price`
- `sl`
- `tp`
- `be`
- `conf`
- `quality`
- `size_scale`
- `htf_bias`
- `long_edge`
- `short_edge`
- `state`
- `regime`

## Recommended Workflow

1. Start with default settings.
2. Pick the chart timeframe you actually trade.
3. Set the higher timeframe for directional confirmation.
4. Validate alerts on replay or paper trading.
5. Then optimize thresholds, cooldown, and risk settings per market.

## Caution

- Lower thresholds create more signals but usually more noise.
- Strong smoothing reduces noise but can delay entries.
- A high `HTF Bias Weight` can filter too aggressively on reversals.
- Always test per asset class.

