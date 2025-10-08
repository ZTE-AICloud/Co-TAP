DEFAULT_MAX_TOKENS = 8192
DEFAULT_TEMPERATURE = 0


# 为了保持向后兼容性，保留原有的LLMConfig类
# 新的配置管理器在 mek.config.manager 中


class LLMConfig:
    """
    Configuration class for the Language Learning Model (LLM).

    This class encapsulates the necessary parameters to interact with an LLM API,
    such as OpenAI's GPT models. It stores the API key, model name, and base URL
    for making requests to the LLM service.
    """

    def __init__(
            self,
            api_key: str,
            provider: str,
            model: str,
            base_url: str,
            temperature: float = DEFAULT_TEMPERATURE,
            max_tokens: int = DEFAULT_MAX_TOKENS,
    ):
        """
        Initialize the LLMConfig with the provided parameters.

        Args:
                api_key (str): The authentication key for accessing the LLM API.
                model (str, optional): The specific LLM model to use for generating responses.
                base_url (str, optional): The base URL of the LLM API service.
        """
        self.base_url = base_url
        self.api_key = api_key
        self.provider = provider
        self.model = model
        self.temperature = temperature
        self.max_tokens = max_tokens

    def __repr__(self):
        return f'[LLMConfig] provider: {self.provider}, model name: {self.model}, url: {self.base_url}, temperature: {self.temperature}, max_tokens: {self.max_tokens}.'
