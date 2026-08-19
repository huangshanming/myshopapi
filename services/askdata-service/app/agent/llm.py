from langchain.chat_models import init_chat_model

from app.conf.app_config import app_config

llm = init_chat_model(
    model=app_config.llm.model_name,
    model_provider="openai",
    api_key=app_config.llm.api_key,
    base_url=app_config.llm.base_url,
    # 召回扩展、SQL 生成更看重稳定性，所以这里关闭随机发散
    temperature=0,
)

if __name__ == "__main__":
    # 单独运行本文件时，用一个最小请求验证模型名、密钥和 base_url 是否配置正确
    print(llm.invoke("你好").content)