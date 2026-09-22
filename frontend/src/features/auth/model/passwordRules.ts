/**
 * Mirrors pkg.ValidatePassword on the server (8+ chars, upper, lower, digit,
 * punctuation/symbol). The server stays the authority; this only gives feedback earlier.
 */
export interface PasswordCheck {
  label: string
  ok: boolean
}

export function checkPassword(password: string): PasswordCheck[] {
  return [
    { label: '8+ characters', ok: password.length >= 8 },
    { label: 'an uppercase letter', ok: /\p{Lu}/u.test(password) },
    { label: 'a lowercase letter', ok: /\p{Ll}/u.test(password) },
    { label: 'a number', ok: /\p{Nd}/u.test(password) },
    { label: 'a special character', ok: /[\p{P}\p{S}]/u.test(password) },
  ]
}

/** 0 = empty, 1 = weak, 2 = fair, 3 = meets every server rule. */
export function passwordStrength(password: string): 0 | 1 | 2 | 3 {
  if (password.length === 0) return 0
  const checks = checkPassword(password)
  const passed = checks.filter((c) => c.ok).length
  if (passed === checks.length) return 3
  if (checks[0].ok && passed >= 3) return 2
  return 1
}
