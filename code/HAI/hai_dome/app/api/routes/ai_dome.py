from pathlib import Path
from typing import Any
from fastapi import APIRouter, Request
from fastapi.responses import StreamingResponse
import asyncio
from app.utils.logger import logger
import json

from app.models import (
    HAIChatResponse,
    HAIAgentRequest
)

router = APIRouter(tags=["ai dome"], prefix="/ai_dome")

current_dir = Path(__file__).parent
storage_path = current_dir.resolve()
dataPath = storage_path / "ai_dome_data"


class DisconnectWatcher:
    def __init__(self, request: Request):
        self.request = request
        self._disconnected = False

    async def watch(self):
        try:
            # 持续接收空数据检测连接状态
            while not await self.request.is_disconnected():
                await asyncio.sleep(0.1)  # 每0.5秒检测一次
        except asyncio.CancelledError:
            # 被取消时设置断开标志
            self._disconnected = True

    @property
    def disconnected(self) -> bool:
        return self._disconnected

async def translate_sse(watcher, filepath):
    watch_task = asyncio.create_task(watcher.watch())

    try:
        # 根据agent类型传递不同的参数
        data = getChatMockData(filepath)
        # ChatAgent需要dialogueId, data, multiMedia参数
        for contentInfo in data:
            if watcher.disconnected:
                logger.info("客户端已断开连接，停止SSE流")
                break
            data = {k: v for k, v in contentInfo.items() if v is not None}
            sse_data = f"data: {json.dumps(data, ensure_ascii=False)}\n\n"
            yield sse_data
            await asyncio.sleep(0.3)

    except asyncio.CancelledError:
        logger.info("SSE流被取消，客户端可能已断开连接")
        raise
    except Exception as e:
        logger.error(f"SSE流处理过程中发生错误: {e}")
        raise
    finally:
        if not watch_task.done():
            watch_task.cancel()
        await asyncio.sleep(0.2)


@router.post("/chat")
def ai_app_chat(body: HAIAgentRequest, request: Request) -> HAIChatResponse:
    filepath = dataPath / "chat-answer.json"
    watcher = DisconnectWatcher(request)

    return StreamingResponse(
        translate_sse(watcher, filepath),
        media_type="text/event-stream"
    )

@router.post("/task")
def ai_app_task(body: HAIAgentRequest, request: Request) -> HAIChatResponse:
    filepath = dataPath / "task-answer.json"
    watcher = DisconnectWatcher(request)

    return StreamingResponse(
        translate_sse(watcher, filepath),
        media_type="text/event-stream"
    )


def getChatMockData(filepath) -> Any:
    """
    读取指定路径的JSON文件并返回解析后的数据

    Args:
        filepath: JSON文件的路径

    Returns:
        解析后的JSON数据
    """
    try:
        logger.info(f"开始读取文件: {filepath}")
        with open(filepath, 'r', encoding='utf-8') as file:
            data = json.load(file)
            logger.info(f"成功读取文件，数据长度: {len(data) if isinstance(data, list) else 'not a list'}")
            for item in data:
                logger.info(f"item: {item}")
                if (item.get("delta") is not None) and (item.get("type") != 'STATE_DELTA'):
                    item["delta"] = json.dumps(item["delta"], ensure_ascii=False)
            return data
    except FileNotFoundError:
        logger.error(f"文件未找到: {filepath}")
        return None
    except json.JSONDecodeError as e:
        logger.error(f"JSON解析错误: {e}")
        return None
    except Exception as e:
        logger.error(f"读取文件时发生错误: {e}")
        return None