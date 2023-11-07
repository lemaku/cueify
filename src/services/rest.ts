import type { ValidationBody, ValidationResult } from '@/types/api'

export async function validate(path: string[], tmp: any): Promise<ValidationResult> {

  const body: ValidationBody = {
    path: path,
    value: tmp
  }

  const response = await fetch(`api/validate`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(body)
  })

  const json = await response.json()

  console.log(json);

  if (response.ok) {
    return Promise.resolve(json as ValidationResult);
  } else {
    return Promise.reject("err")
  }

}
