import json
import time
import base64
import requests
from PIL import Image 
import io

class FusionBrainAPI:
  def __init__(self, url, api_key, secret_key):
    self.URL = url
    self.AUTH_HEADERS = {
      'X-Key': f'Key {api_key}',
      'X-Secret': f'Secret {secret_key}',
    }

  def get_pipeline(self):
    response = requests.get(self.URL + 'key/api/v1/pipelines', headers=self.AUTH_HEADERS)
    data = response.json()
    
    return data[0]['id']
  
  def get_pipelines(self):
    response = requests.get(self.URL + 'key/api/v1/pipelines', headers=self.AUTH_HEADERS)
    data = response.json()
    
    return data

  def generate(self, prompt, pipeline_id, numImages=1, width=1024, height=1024):    
    params = {
      "type": "GENERATE",
      "numImages": numImages,
      "width": width,
      "height": height,
      "generateParams": {}
    }
    params['generateParams']['query'] = prompt

    data = {
      'pipeline_id': (None, pipeline_id),
      'params': (None, json.dumps(params), 'application/json')
    }
    response = requests.post(self.URL + 'key/api/v1/pipeline/run', headers=self.AUTH_HEADERS, files=data)
    data = response.json()
    return data['uuid']

  def check_generation(self, request_id, delay=5):
    response = requests.get(self.URL + 'key/api/v1/pipeline/status/' + request_id, headers=self.AUTH_HEADERS)
    data = response.json()
    
    print(f'status is {data["status"]}')

    if data['status'] == 'DONE':
      for file in data['result']['files']:
        api.base64_to_image(file)
      
      print('Status: DONE')
      print(f"Количество изображений: {len(data['result']['files'])}")
      return
    else:
      time.sleep(delay)
      self.check_generation(request_id)

  def base64_to_image(self, base64_string, output_path="output.png"):
      """
      Преобразует строку base64 в изображение и сохраняет его.

      Args:
        base64_string: Строка в формате base64, содержащая данные изображения.
        output_path: Путь, по которому будет сохранено изображение.
      """
      try:
        # Декодируем строку base64 в байты
        image_data = base64.b64decode(base64_string)

        # Открываем байтовый поток
        image_stream = io.BytesIO(image_data)

        # Открываем изображение с помощью Pillow
        image = Image.open(image_stream)

        # Сохраняем изображение
        image.save(output_path)

        print(f"Изображение успешно сохранено в {output_path}")

      except Exception as e:
        print(f"Ошибка при преобразовании: {e}")

if __name__ == '__main__':
  api = FusionBrainAPI('https://api-key.fusionbrain.ai/', 'D3B80F41D98430D2E371D1F0833D8A71', 'EB06F16AB38B65A0FD3E7F9A9BC53E2B')
  pipelines = api.get_pipelines()
  print(f'pipelines: {pipelines}')
  
  pipeline_id = api.get_pipeline()
  print(f'pipeline_id: {pipeline_id}')

  promt = input("Введите промт для генерации изображения: ")

  uuid = api.generate(promt, pipeline_id)
  print(f'uuid: {uuid}')
  
  api.check_generation(uuid)
