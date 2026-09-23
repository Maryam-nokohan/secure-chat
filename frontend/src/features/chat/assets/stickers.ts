export interface Sticker {
  id: string
  emoji: string
  label: string
  gradient: [string, string]
}

/**
 * A self-contained sticker pack: each "sticker" renders as a large emoji on a
 * gradient card (see StickerPicker/StickerBubble), so nothing needs to be
 * fetched from a third party. Sent as a message body of the form
 * [[sticker:<id>]], the same convention files use ([[file:...]]).
 */
export const STICKERS: Sticker[] = [
  { id: 'wave', emoji: '👋', label: 'Wave', gradient: ['#8b6fff', '#c9b8ff'] },
  { id: 'heart', emoji: '❤️', label: 'Heart', gradient: ['#f43f5e', '#fda4af'] },
  { id: 'laugh', emoji: '😂', label: 'Laughing', gradient: ['#e8920a', '#f5c060'] },
  { id: 'fire', emoji: '🔥', label: 'Fire', gradient: ['#f43f5e', '#e8920a'] },
  { id: 'thumbsup', emoji: '👍', label: 'Thumbs up', gradient: ['#10b981', '#6ee7b7'] },
  { id: 'clap', emoji: '👏', label: 'Clap', gradient: ['#c8962b', '#e8c46a'] },
  { id: 'party', emoji: '🎉', label: 'Party', gradient: ['#8b6fff', '#f43f5e'] },
  { id: 'think', emoji: '🤔', label: 'Thinking', gradient: ['#6b6880', '#a09cb0'] },
  { id: 'cry', emoji: '😭', label: 'Crying', gradient: ['#4a7fc1', '#8b6fff'] },
  { id: 'cool', emoji: '😎', label: 'Cool', gradient: ['#10b981', '#c8962b'] },
  { id: 'star', emoji: '🌟', label: 'Star', gradient: ['#e8920a', '#8b6fff'] },
  { id: 'moon', emoji: '🌙', label: 'Moon', gradient: ['#4a7fc1', '#0a0a0f'] },
  { id: 'coffee', emoji: '☕', label: 'Coffee', gradient: ['#c8962b', '#6b6880'] },
  { id: 'rocket', emoji: '🚀', label: 'Rocket', gradient: ['#f43f5e', '#8b6fff'] },
  { id: 'cat', emoji: '🐱', label: 'Cat', gradient: ['#e8920a', '#fda4af'] },
  { id: 'ok', emoji: '👌', label: 'OK', gradient: ['#10b981', '#e8c46a'] },
]

export const FILE_STICKER_MARKER = /^\[\[sticker:([a-z]+)\]\]$/
