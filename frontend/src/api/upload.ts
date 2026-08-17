import request from '@/utils/request'

export async function uploadImage(file: File): Promise<string> {
  const form = new FormData()
  form.append('file', file)
  const res = await request.post<{ code: number; message: string; data: { url: string } }>('/api/v1/upload/image', form, {
    headers: { 'Content-Type': 'multipart/form-data' },
  })
  return res.data.data.url
}
