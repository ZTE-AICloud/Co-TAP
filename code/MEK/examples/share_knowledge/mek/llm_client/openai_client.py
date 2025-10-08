import json
from typing import Optional, List, Dict, Any

from openai import RateLimitError, OpenAI
from pydantic import BaseModel

from mek.llm_client.base import LLMClient
from mek.llm_client.config import LLMConfig
from mek.llm_client.models import Message, RoleType
from mek.utils.data_util import json_loads, generate_template_from_model
from mek.utils.log_util import logger


class OpenAIClient(LLMClient):
    """
    OpenAIClient is a client class for interacting with large language models.
    """

    def __init__(
            self,
            config: LLMConfig | None = None,
            client: Any = None,
            max_retries: int = 2
    ):
        """
        Initialize the OpenAIClient with the provided configuration, and client.

        Args:
            config (LLMConfig | None): The configuration for the LLM client, including API key, model, base URL, temperature, and max tokens.
            client (Any | None): An optional async client instance to use. If not provided, a new OpenAIClient is created.

        """

        super().__init__(config)
        self.max_retires = max_retries
        if not client:
            assert config.api_key is not None, "API_KEY must be provided ！"
            self.client = OpenAI(
                api_key=config.api_key,
                base_url=config.base_url
            )
        else:
            self.client = client

    def _parse_response(
            self,
            response,
            messages: list[Message],
            response_model: type[BaseModel] | None = None,
            tools: Optional[List[Dict]] = None
    ) -> dict[str, Any]:
        try:
            content = response.choices[0].message.content
            processed_response = {
                "content": json_loads(content) if response_model else content,
                "json_format": bool(response_model)
            }
            messages.append(
                Message(
                    role=RoleType.AI,
                    content=content
                )
            )
            return processed_response
        except Exception as e:
            logger.error(f'Error in parse model response ({self.model}): {e}')
            raise

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
        print(f"[llm client]: use {self.model}")
        openai_messages: list = []
        for m in messages:
            m.content = self._clean_input(m.content)
            if m.role == RoleType.User or m.role == RoleType.AI:
                openai_messages.append({'role': m.role, 'content': m.content})
            elif m.role == RoleType.System:
                if response_model:
                    template = generate_template_from_model(response_model)
                    format_str = json.dumps(template, indent=2, ensure_ascii=False)
                    openai_messages.append({'role': m.role, 'content': m.content + f'\n\nRespond with a JSON object in the following format:\n\n{format_str}'})
                else:
                    openai_messages.append({'role': m.role, 'content': m.content})

        parameters: dict[str, Any] = {
            "model": self.model,
            "messages": openai_messages,
            "temperature": self.temperature,
            "max_tokens": self.max_tokens,
            "stream": stream,
        }
        if response_model:
            parameters["response_format"] = {"type": "json_object"}
        if tools:
            parameters["tools"] = tools
            parameters["tool_choice"] = tool_choice
        if log_probs:
            parameters["logprobs"] = True
            if top_log_probs >= 0:
                parameters["top_logprobs"] = top_log_probs
        try:
            response = self.client.chat.completions.create(**parameters)

            return self._parse_response(response, messages, response_model, tools)
        except RateLimitError as e:
            raise RateLimitError from e
        except Exception as e:
            logger.error(f'Error in generating LLM response: {e}')
            raise
