from enum import Enum
from typing import Optional, Any

from pydantic import BaseModel


class RoleType(str, Enum):
    """Role enum class definition."""
    User = "user"
    System = "system"
    AI = "assistant"
    Tool = "tool"

    def __repr__(self):
        """Get a string representation."""
        return f'"{self.value}"'


class Message(BaseModel):
    role: Optional[RoleType] = None
    content: Optional[str] = None
    reasoning_content: Optional[str] = None
    name: Optional[str] = None
    tool_call_id: Optional[str] = None
    tool_calls: Optional[dict[str, Any]] = None
