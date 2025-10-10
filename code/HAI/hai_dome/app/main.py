import argparse
from pathlib import Path

from fastapi import FastAPI
from fastapi.routing import APIRoute

from app.api.main import api_router

def custom_generate_unique_id(route: APIRoute) -> str:
    return f"{route.tags[0]}-{route.name}"


# if settings.SENTRY_DSN and settings.ENVIRONMENT != "local":
#     sentry_sdk.init(dsn=str(settings.SENTRY_DSN), enable_tracing=True)

app = FastAPI(
    title="HAI AI DOME",
    openapi_url="/api/v1/hai/openapi.json",
    generate_unique_id_function=custom_generate_unique_id,
)
from fastapi.staticfiles import StaticFiles

# 确保 静态文件 目录存在
Path("ai-ui").mkdir(parents=True, exist_ok=True)

# 挂载静态文件目录
app.mount("/aiui", StaticFiles(directory="ai-ui", html=True), name="static")

app.include_router(api_router, prefix="/api/v1/hai")

if __name__ == "__main__":
    import uvicorn

    # 解析命令行参数
    parser = argparse.ArgumentParser(description='HAI AI DOME')
    parser.add_argument('--dev', action='store_true', help='使用内部模型进行开发模式')
    parser.add_argument('--port', type=int, default=8000, help='服务监听端口 (默认: 8000)')
    args = parser.parse_args()

    uvicorn.run(app, host="10.137.216.137", port=args.port, reload_excludes="workspace/*")
