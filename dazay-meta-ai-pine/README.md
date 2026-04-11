# Dazay META-AI Pine

Private repository for the Pine Script indicator `Dazay Pro: META-AI v6 [Elite Quant]`.

## What It Is

This indicator combines:

- adaptive EOT-based signal weighting
- WaveTrend confirmation
- market regime and noise filtering
- AI health state management
- higher-timeframe directional bias
- adaptive signal smoothing
- dynamic risk and alert payloads

The goal is to reduce weak entries in noisy conditions and improve signal quality for manual trading or bot-based alert execution.

## Main Improvements Over The Original

- more stable learning logic for internal weights
- `HEALTHY / ADAPTING / DEGRADED / FROZEN` AI state machine
- higher-timeframe bias confirmation
- long/short edge tracking
- stronger anti-whipsaw filtering
- richer JSON alerts for automation
- cleaner dashboard and risk context

## Repository Structure

- `dazay_meta_ai_v6_elite.pine` - indicator source
- `docs/USAGE.md` - quick usage guide
- `LICENSE` - MPL 2.0 license text

## Quick Start

1. Open TradingView.
2. Open `Pine Editor`.
3. Copy the full contents of `dazay_meta_ai_v6_elite.pine`.
4. Paste into the editor and click `Add to chart`.
5. Tune parameters for your market and timeframe.

## Best Use Cases

- trend-following intraday workflows
- higher-timeframe directional filtering
- alert-based execution
- discretionary confirmation for manual entries

## Important Notes

- This is an indicator, not a guaranteed profitable system.
- Parameters should be tuned per symbol and timeframe.
- Use paper trading or backtesting before live deployment.

