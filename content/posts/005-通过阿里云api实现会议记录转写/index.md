---
title: "通过阿里云api实现会议记录转写"
date: 2025-11-04T10:14:16+08:00
lastmod: 2025-11-04T10:34:38+08:00
draft: true
categories:
tags:
---
## 引入库

```python
from dashscope.audio.asr import Transcription
import dashscope
dashscope.api_key = api_key
```
文档： https://help.aliyun.com/zh/model-studio/paraformer-recorded-speech-recognition-python-sdk

## 发起转换任务

```python
task_response = Transcription.async_call(
    model='paraformer-v2',
    file_urls=[file_url],
    language_hints=['zh', 'en'],
    diarization_enabled=True,
)

if task_response.status_code != HTTPStatus.OK:
    raise Exception(f"任务提交失败: {task_response.status_code}")

task_id = task_response.output.task_id
transcribe_response = Transcription.wait(task=task_id)
```

- 需要一个可以直接访问url，来访问要转码的音频。
- diarization_enabled 区分发言人。需要把音频预处理为单声道(可以用ffmpeg)。
- 任务是异步的，需要等待。

## 获取转化结果

```python
if transcribe_response.status_code == HTTPStatus.OK:
    results = transcribe_response.output.results
    if results and len(results) > 0:
        for result in results:
            if isinstance(result, dict) and result.get('subtask_status') == 'SUCCEEDED':
                transcription_url = result.get('transcription_url')
                if transcription_url:
                    self.transcription_queue.put(('status', "正在获取转录结果..."))
                    try:
                        import requests
                        response = requests.get(transcription_url)
                        if response.status_code == 200:
                            transcription_data = response.json()
                        else:
                            self.transcription_queue.put(('error', f"获取转录结果失败: HTTP {response.status_code}"))
                    except Exception as e:
                        self.transcription_queue.put(('error', f"获取转录结果时出错: {str(e)}"))
            else:
                self.transcription_queue.put(('error', f"转录任务失败，状态: {result.get('subtask_status', 'UNKNOWN')}"))
                break
    else:
        self.transcription_queue.put(('error', "转录结果为空"))
else:
    error_msg = f"语音转文字失败: {transcribe_response.status_code}"
    if hasattr(transcribe_response, 'output') and transcribe_response.output:
        if hasattr(transcribe_response.output, 'message'):
            error_msg += f" - {transcribe_response.output.message}"
    raise Exception(error_msg)
    
except Exception as e:
self.transcription_queue.put(('error', str(e)))
```

- 转化结果是以json文件提供的

### 对json文件格式化

```python
"""将转录JSON数据转换为对话格式"""
    formatted_lines = []
    
    # 检查转录数据格式
    if not isinstance(transcription_data, dict):
        return None
    
    # 尝试多种可能的数据结构
    sentences = None
    
    # 可能的结构1: transcripts[0].sentences
    if 'transcripts' in transcription_data and isinstance(transcription_data['transcripts'], list):
        if len(transcription_data['transcripts']) > 0:
            first_transcript = transcription_data['transcripts'][0]
            if isinstance(first_transcript, dict) and 'sentences' in first_transcript:
                sentences = first_transcript['sentences']
    
    # 可能的结构2: 直接在根级别有sentences
    if sentences is None and 'sentences' in transcription_data:
        sentences = transcription_data['sentences']
    
    # 可能的结构3: results数组中的sentences
    if sentences is None and 'results' in transcription_data:
        for result in transcription_data['results']:
            if isinstance(result, dict) and 'sentences' in result:
                sentences = result['sentences']
                break
    
    if sentences is None:
        print("未找到sentences数据，转为简单文本模式")
        # 如果没有sentence结构，尝试提取文本
        transcription_text_parts = []
        if 'transcripts' in transcription_data:
            transcripts = transcription_data['transcripts']
            if isinstance(transcripts, list) and len(transcripts) > 0:
                for transcript in transcripts:
                    if isinstance(transcript, dict) and 'text' in transcript:
                        transcription_text_parts.append(transcript['text'])
                if transcription_text_parts:
                    transcription_text = " ".join(transcription_text_parts)
                    return transcription_text.strip()
        return None
    
    # 按时间排序
    if isinstance(sentences, list) and len(sentences) > 0:
        sentences.sort(key=lambda x: x.get('begin_time', 0))
        
        # 生成对话格式的输出
        last_speaker_id = None
        for sentence in sentences:
            begin_time = format_time(sentence.get('begin_time', 0))
            speaker_id = format_speaker_id(sentence.get('speaker_id', '未知'))
            if last_speaker_id != speaker_id:
                last_speaker_id = speaker_id
                formatted_lines.append('')
            text = sentence.get('text', '')
            
            line = f"{begin_time} {speaker_id} : {text}"
            formatted_lines.append(line)
        
        return '\n'.join(formatted_lines)
    
    return None
```
- 转为对阅读更友好的对话模式