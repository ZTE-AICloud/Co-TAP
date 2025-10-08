import json
from abc import ABC, abstractmethod
from typing import Any, Optional, List, Dict

import httpx
from openai import RateLimitError
from pydantic import BaseModel
from tenacity import retry, stop_after_attempt, wait_random_exponential, retry_if_exception

from mek.llm_client.config import LLMConfig
from mek.llm_client.models import Message
from mek.utils.log_util import logger


def is_server_or_retry_error(exception: BaseException) -> bool:
    if isinstance(exception, RateLimitError | json.decoder.JSONDecodeError):
        return True
    if isinstance(exception, httpx.HTTPStatusError) and 500 <= exception.response.status_code < 600:
        return True
    return False


class LLMClient(ABC):
    def __init__(self, config: LLMConfig | None):
        if not config:
            config = LLMConfig()

        self.config = config
        self.provider = config.provider
        self.model = config.model
        self.temperature = config.temperature
        self.max_tokens = config.max_tokens

    def _clean_input(self, content: str) -> str:
        """Clean content string of invalid unicode and control characters.

        Args:
            content: Raw content string to be cleaned

        Returns:
            Cleaned string safe for LLM processing
        """
        # Clean any invalid Unicode
        cleaned = content.encode('utf-8', errors='ignore').decode('utf-8')

        # Remove zero-width characters and other invisible unicode
        zero_width = '\u200b\u200c\u200d\ufeff\u2060'
        for char in zero_width:
            cleaned = cleaned.replace(char, '')

        # Remove control characters except newlines, returns, and tabs
        cleaned = ''.join(char for char in cleaned if ord(char) >= 32 or char in '\n\r\t')

        return cleaned

    @retry(
        stop=stop_after_attempt(4),
        wait=wait_random_exponential(multiplier=10, min=5, max=120),
        retry=retry_if_exception(is_server_or_retry_error),
        after=lambda retry_state: logger.warning(
            f'Retrying {retry_state.fn.__name__ if retry_state.fn else "function"} after {retry_state.attempt_number} attempts...'
        )
        if retry_state.attempt_number > 1
        else None,
        reraise=True,
    )
    def _generate_response_with_retry(
            self,
            messages: list[Message],
            response_model: type[BaseModel] | None = None,
            stream: bool = False,
            tools: Optional[List[Dict]] = None,
            tool_choice: str = "auto"
    ) -> dict[str, Any]:
        try:
            return self._generate_response(messages, response_model, stream, tools, tool_choice)
        except (httpx.HTTPStatusError, RateLimitError) as e:
            raise e

    @abstractmethod
    def _generate_response(
            self,
            messages: list[Message],
            response_model: type[BaseModel] | None = None,
            stream: bool = False,
            tools: Optional[List[Dict]] = None,
            tool_choice: str = "auto",
            log_probs: bool = False,
            top_log_probs: Optional[int] = -1
    ) -> dict[str, Any]:
        raise NotImplementedError

    @abstractmethod
    def _parse_response(
            self,
            response,
            messages: list[Message],
            response_model: type[BaseModel] | None = None,
            tools: Optional[List[Dict]] = None
    ) -> dict[str, Any]:
        raise NotImplementedError

    def __repr__(self):
        return repr(self.config)

    def generate_response(
            self,
            messages: list[Message],
            response_model: type[BaseModel] | None = None,
            stream: bool = False,
            tools: Optional[List[Dict]] = None,
            tool_choice: str = "auto"
    ) -> dict[str, Any]:
        if response_model is not None:
            serialized_model = json.dumps(response_model.model_json_schema())
            messages[
                -1
            ].content += (
                f'\n\nRespond with a JSON object in the following format:\n\n{serialized_model}'
            )

        for message in messages:
            message.content = self._clean_input(message.content)

        return self._generate_response_with_retry(messages, response_model, stream, tools, tool_choice)


class EmbedderClient(ABC):
    @abstractmethod
    def create(self, input_data: str | list[str]) -> list[float]:
        pass

    def create_batch(self, input_data_list: list[str]) -> list[list[float]]:
        raise NotImplementedError()
