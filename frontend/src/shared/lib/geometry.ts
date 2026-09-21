/** `points` attribute for an SVG <polygon> shaped like an n-pointed star. */
export function starPoints(cx: number, cy: number, outer: number, inner: number, n: number): string {
  return Array.from({ length: n * 2 }, (_, i) => {
    const angle = (i * Math.PI) / n - Math.PI / 2
    const r = i % 2 === 0 ? outer : inner
    return `${cx + r * Math.cos(angle)},${cy + r * Math.sin(angle)}`
  }).join(' ')
}
