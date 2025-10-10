from pydantic import BaseModel, Field
from typing import Optional, List, Dict, Any
# Shared properties

class HAIChatResponse(BaseModel):
    """聊天响应模型"""
    status: str = Field(description="响应状态")
    message: str = Field(description="响应消息")
    data: Optional[List[Dict[str, Any]]] = Field(default=None, description="聊天数据")

class HAIMessage(BaseModel):
    """聊天消息模型"""
    id: str = Field(description="消息ID")
    role: str = Field(description="角色")
    content: str = Field(description="内容")

class HAIAgentRequest(BaseModel):
    """聊天请求模型"""
    threadId: str = Field(description="对话ID")
    runId: str = Field(description="运行ID")
    messages: List[HAIMessage] = Field(default=[], description="消息")
    context: List[HAIMessage] = Field(default=[], description="上下文")
    tools: List[Dict[str, Any]] = Field(default=[], description="工具")
    forwardedProps: Dict[str, Any] = Field(default={}, description="转发属性")
