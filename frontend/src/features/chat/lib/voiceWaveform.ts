export const WAVE_BARS = 48

export function encodeWaveform(bytes: Uint8Array): string {
  let bin = ''
  for (let i = 0; i < bytes.length; i++) bin += String.fromCharCode(bytes[i])
  return btoa(bin)
}

export function fallbackBars(seed: string, n = WAVE_BARS): number[] {
  let h = 2166136261
  for (let i = 0; i < seed.length; i++) {
    h ^= seed.charCodeAt(i)
    h = Math.imul(h, 16777619)
  }
  if (h === 0) h = 1
  const out: number[] = []
  for (let i = 0; i < n; i++) {
    h ^= h << 13
    h ^= h >>> 17
    h ^= h << 5
    out.push(0.25 + (((h >>> 0) % 1000) / 1000) * 0.6)
  }
  return out
}

export function decodeWaveform(b64: string, seed = ''): number[] {
  try {
    const bin = atob(b64)
    const out: number[] = []
    for (let i = 0; i < bin.length; i++) out.push(bin.charCodeAt(i) / 255)
    return out.length ? out : fallbackBars(seed)
  } catch {
    return fallbackBars(seed)
  }
}

export function summarizeLevels(samples: number[], bars = WAVE_BARS): Uint8Array {
  const out = new Uint8Array(bars)
  if (samples.length === 0) return out
  for (let i = 0; i < bars; i++) {
    const from = Math.floor((i * samples.length) / bars)
    const to = Math.max(from + 1, Math.floor(((i + 1) * samples.length) / bars))
    let peak = 0
    for (let j = from; j < to && j < samples.length; j++) peak = Math.max(peak, samples[j])
    out[i] = Math.round(peak * 255)
  }
  let max = 0
  for (const v of out) max = Math.max(max, v)
  if (max > 0 && max < 255) {
    const k = 255 / max
    for (let i = 0; i < bars; i++) out[i] = Math.round(out[i] * k)
  }
  for (let i = 0; i < bars; i++) out[i] = Math.max(out[i], 20) 
  return out
}