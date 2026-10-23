import { useCallback, useEffect, useRef, useState } from 'react'
import { summarizeLevels } from '../lib/voiceWaveform'

const CANDIDATE_MIME_TYPES = [
  'audio/webm;codecs=opus',
  'audio/ogg;codecs=opus',
  'audio/webm',
  'audio/mp4', 
  'audio/ogg',
]
const MAX_RECORDING_MS = 5 * 60 * 1000
const MIN_RECORDING_MS = 500
const SAMPLE_INTERVAL_MS = 80
const LIVE_BARS = 36

function pickMimeType(): string | null {
  if (typeof MediaRecorder === 'undefined') return null
  for (const type of CANDIDATE_MIME_TYPES) {
    if (MediaRecorder.isTypeSupported(type)) return type
  }
  return null
}

export interface VoiceRecording {
  blob: Blob
  mimeType: string
  durationMs: number
  waveform: Uint8Array
}

export interface VoicePreview extends VoiceRecording {
  url: string
}

export type VoiceRecorderState = 'idle' | 'recording' | 'preview' | 'error'

export interface VoiceRecorder {
  state: VoiceRecorderState
  elapsedMs: number
  levels: number[]
  error: string | null
  preview: VoicePreview | null
  start: () => Promise<void>
  stop: (toPreview: boolean) => Promise<VoiceRecording | null>
  discard: () => void
}

export function useVoiceRecorder(): VoiceRecorder {
  const [state, setState] = useState<VoiceRecorderState>('idle')
  const [elapsedMs, setElapsedMs] = useState(0)
  const [levels, setLevels] = useState<number[]>([])
  const [error, setError] = useState<string | null>(null)
  const [preview, setPreview] = useState<VoicePreview | null>(null)

  const recorderRef = useRef<MediaRecorder | null>(null)
  const streamRef = useRef<MediaStream | null>(null)
  const audioCtxRef = useRef<AudioContext | null>(null)
  const analyserRef = useRef<AnalyserNode | null>(null)
  const chunksRef = useRef<Blob[]>([])
  const samplesRef = useRef<number[]>([])
  const startedAtRef = useRef(0)
  const mimeRef = useRef('')
  const tickRef = useRef<number | null>(null)
  const previewUrlRef = useRef<string | null>(null)
  const stopRef = useRef<(toPreview: boolean) => Promise<VoiceRecording | null>>(() => Promise.resolve(null))

  const teardown = useCallback(() => {
    if (tickRef.current !== null) {
      window.clearInterval(tickRef.current)
      tickRef.current = null
    }
    streamRef.current?.getTracks().forEach((t) => t.stop())
    streamRef.current = null
    void audioCtxRef.current?.close().catch(() => undefined)
    audioCtxRef.current = null
    analyserRef.current = null
    recorderRef.current = null
  }, [])

  const clearPreview = useCallback(() => {
    if (previewUrlRef.current) {
      URL.revokeObjectURL(previewUrlRef.current)
      previewUrlRef.current = null
    }
    setPreview(null)
  }, [])

  const start = useCallback(async () => {
    setError(null)
    clearPreview()
    const mimeType = pickMimeType()
    if (!mimeType) {
      setError('Voice recording is not supported in this browser.')
      setState('error')
      return
    }
    try {
      const stream = await navigator.mediaDevices.getUserMedia({
        audio: { echoCancellation: true, noiseSuppression: true },
      })
      streamRef.current = stream
      chunksRef.current = []
      samplesRef.current = []
      mimeRef.current = mimeType

      const recorder = new MediaRecorder(stream, { mimeType, audioBitsPerSecond: 32000 })
      recorder.ondataavailable = (e) => {
        if (e.data.size > 0) chunksRef.current.push(e.data)
      }
      recorderRef.current = recorder

      const ctx = new AudioContext()
      const analyser = ctx.createAnalyser()
      analyser.fftSize = 512
      ctx.createMediaStreamSource(stream).connect(analyser)
      audioCtxRef.current = ctx
      analyserRef.current = analyser
      const buf = new Uint8Array(analyser.fftSize)

      recorder.start()
      startedAtRef.current = Date.now()
      setElapsedMs(0)
      setLevels([])
      tickRef.current = window.setInterval(() => {
        const elapsed = Date.now() - startedAtRef.current
        analyser.getByteTimeDomainData(buf)
        let sum = 0
        for (let i = 0; i < buf.length; i++) {
          const v = (buf[i] - 128) / 128
          sum += v * v
        }
        const level = Math.min(1, Math.sqrt(sum / buf.length) * 3)
        samplesRef.current.push(level)
        setLevels((prev) => [...prev.slice(-(LIVE_BARS - 1)), level])
        setElapsedMs(elapsed)
        if (elapsed >= MAX_RECORDING_MS) void stopRef.current(true)
      }, SAMPLE_INTERVAL_MS)
      setState('recording')
    } catch {
      setError('Microphone access was denied.')
      setState('error')
      teardown()
    }
  }, [clearPreview, teardown])

  const stop = useCallback(
    (toPreview: boolean): Promise<VoiceRecording | null> => {
      const recorder = recorderRef.current
      if (!recorder || recorder.state === 'inactive') return Promise.resolve(null)
      const mimeType = recorder.mimeType || mimeRef.current
      const durationMs = Date.now() - startedAtRef.current
      const waveform = summarizeLevels(samplesRef.current)
      return new Promise((resolve) => {
        recorder.onstop = () => {
          const blob = new Blob(chunksRef.current, { type: mimeType })
          teardown()
          if (blob.size === 0 || durationMs < MIN_RECORDING_MS) {
            setState('idle')
            resolve(null)
            return
          }
          const rec: VoiceRecording = { blob, mimeType, durationMs, waveform }
          if (toPreview) {
            const url = URL.createObjectURL(blob)
            previewUrlRef.current = url
            setPreview({ ...rec, url })
            setState('preview')
          } else {
            setState('idle')
          }
          resolve(rec)
        }
        recorder.stop()
      })
    },
    [teardown],
  )
  stopRef.current = stop

  const discard = useCallback(() => {
    const recorder = recorderRef.current
    if (recorder && recorder.state !== 'inactive') {
      recorder.onstop = null
      recorder.stop()
    }
    teardown()
    clearPreview()
    setState('idle')
    setElapsedMs(0)
    setLevels([])
  }, [teardown, clearPreview])

  useEffect(
    () => () => {
      const recorder = recorderRef.current
      if (recorder && recorder.state !== 'inactive') {
        recorder.onstop = null
        recorder.stop()
      }
      teardown()
      if (previewUrlRef.current) URL.revokeObjectURL(previewUrlRef.current)
    },
    [teardown],
  )

  return { state, elapsedMs, levels, error, preview, start, stop, discard }
}