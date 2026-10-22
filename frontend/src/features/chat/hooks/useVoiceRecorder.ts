import { useCallback, useRef, useState } from 'react'

const CANDIDATE_MIME_TYPES = ['audio/webm;codecs=opus', 'audio/ogg;codecs=opus', 'audio/webm', 'audio/ogg']

function pickMimeType(): string | null {
  for (const type of CANDIDATE_MIME_TYPES) {
    if (MediaRecorder.isTypeSupported(type)) return type
  }
  return null
}

export interface VoiceRecording {
  blob: Blob
  mimeType: string
  durationMs: number
}

export type VoiceRecorderState = 'idle' | 'recording' | 'error'

export interface VoiceRecorder {
  state: VoiceRecorderState
  elapsedMs: number
  error: string | null
  start: () => Promise<void>
  stop: () => Promise<VoiceRecording | null>
  cancel: () => void
}

export function useVoiceRecorder(): VoiceRecorder {
  const [state, setState] = useState<VoiceRecorderState>('idle')
  const [elapsedMs, setElapsedMs] = useState(0)
  const [error, setError] = useState<string | null>(null)

  const recorderRef = useRef<MediaRecorder | null>(null)
  const streamRef = useRef<MediaStream | null>(null)
  const chunksRef = useRef<Blob[]>([])
  const startedAtRef = useRef(0)
  const tickRef = useRef<number | null>(null)

  const releaseStream = useCallback(() => {
    if (tickRef.current !== null) {
      window.clearInterval(tickRef.current)
      tickRef.current = null
    }
    streamRef.current?.getTracks().forEach((t) => t.stop())
    streamRef.current = null
    recorderRef.current = null
  }, [])

  const start = useCallback(async () => {
    setError(null)
    const mimeType = pickMimeType()
    if (!mimeType) {
      setError('Voice recording is not supported in this browser.')
      setState('error')
      return
    }
    try {
      const stream = await navigator.mediaDevices.getUserMedia({ audio: true })
      streamRef.current = stream
      chunksRef.current = []
      const recorder = new MediaRecorder(stream, { mimeType })
      recorder.ondataavailable = (e) => {
        if (e.data.size > 0) chunksRef.current.push(e.data)
      }
      recorderRef.current = recorder
      recorder.start()
      startedAtRef.current = Date.now()
      setElapsedMs(0)
      tickRef.current = window.setInterval(() => {
        setElapsedMs(Date.now() - startedAtRef.current)
      }, 200)
      setState('recording')
    } catch {
      setError('Microphone access was denied.')
      setState('error')
      releaseStream()
    }
  }, [releaseStream])

  const stop = useCallback((): Promise<VoiceRecording | null> => {
    const recorder = recorderRef.current
    if (!recorder || state !== 'recording') return Promise.resolve(null)
    const mimeType = recorder.mimeType
    const durationMs = Date.now() - startedAtRef.current
    return new Promise((resolve) => {
      recorder.onstop = () => {
        const blob = new Blob(chunksRef.current, { type: mimeType })
        releaseStream()
        setState('idle')
        resolve(blob.size > 0 ? { blob, mimeType, durationMs } : null)
      }
      recorder.stop()
    })
  }, [releaseStream, state])

  const cancel = useCallback(() => {
    const recorder = recorderRef.current
    if (recorder && recorder.state !== 'inactive') {
      recorder.onstop = null
      recorder.stop()
    }
    releaseStream()
    setState('idle')
    setElapsedMs(0)
  }, [releaseStream])

  return { state, elapsedMs, error, start, stop, cancel }
}
