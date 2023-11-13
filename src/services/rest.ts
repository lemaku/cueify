import type {
  InspectBody,
  InspectResult,
  SummarizeBody,
  SummarizeResult,
  ValidationBody,
  ValidationResult
} from '@/types/api'

export async function validate(path: string[], current: any): Promise<ValidationResult> {
  const body: ValidationBody = {
    path: path,
    value: current
  }

  const response = await fetch(`${window.location.origin}/api/validate`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(body)
  })

  const json = await response.json()

  if (response.ok) {
    return Promise.resolve(json)
  } else {
    return Promise.reject('err')
  }
}

export async function inspect(path: string[], current: any): Promise<InspectResult> {
  const body: InspectBody = {
    path: path,
    value: current
  }

  const response = await fetch(`${window.location.origin}/api/inspect`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(body)
  })

  const json = await response.json()

  if (response.ok) {
    return Promise.resolve(json)
  } else {
    return Promise.reject('err')
  }
}

export async function summarize(current: any): Promise<SummarizeResult> {
  const body: SummarizeBody = {
    value: current
  }

  const response = await fetch(`${window.location.origin}/api/summarize`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(body)
  })

  const json = await response.json()

  if (response.ok) {
    return Promise.resolve(json)
  } else {
    return Promise.reject('err')
  }
}
