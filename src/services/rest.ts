import type { InspectBody, InspectResult, ValidationBody, ValidationResult } from '@/types/api'

export async function validate(path: string[], current: any): Promise<ValidationResult> {

  const body: ValidationBody = {
    path: path,
    value: current
  }

  const response = await fetch(`api/validate`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(body)
  })

  const json = await response.json()

  if (response.ok) {
    return Promise.resolve(json);
  } else {
    return Promise.reject("err")
  }

}


export async function inspect(path: string[], current: any): Promise<InspectResult> {

  const body: InspectBody = {
    path: path,
    value: current
  }

  console.log('current', JSON.stringify(body));

  const response = await fetch(`api/inspect`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(body)
  })

  const json = await response.json()

  console.log('body', json);

  if (response.ok) {
    return Promise.resolve(json);
  } else {
    return Promise.reject("err")
  }
}

