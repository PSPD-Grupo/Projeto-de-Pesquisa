from fastapi import FastAPI, WebSocket, WebSocketDisconnect
from fastapi.middleware.cors import CORSMiddleware
from fastapi.staticfiles import StaticFiles
from pydantic import BaseModel
import grpc.aio
import asyncio
import time
import os

from app.config import settings
from app.grpc_clients.servidor_a import FeedClient
from app.grpc_clients.servidor_b import ReacaoClient
from app.grpc_generated import chat_pb2, chat_pb2_grpc

app = FastAPI(title="Stub FastAPI gRPC", version="0.1.0")

# Enable CORS for frontend requests
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)


class DesabafoRequest(BaseModel):
    texto: str


class ReacaoRequest(BaseModel):
    desabafo_id: str


@app.get("/health")
def health() -> dict[str, str]:
    return {"status": "ok"}


@app.post("/desabafos")
def postar_desabafo_param(texto: str):
    return FeedClient().postar_desabafo(texto)


@app.post("/publicar")
def postar_desabafo(req: DesabafoRequest):
    return FeedClient().postar_desabafo(req.texto)


@app.get("/feed")
def listar_feed(quant: int = 100):
    return FeedClient().listar_feed(quant)


@app.get("/desabafos/{desabafo_id}/reacoes")
def buscar_reacoes(desabafo_id: str):
    return ReacaoClient().buscar_reacoes(desabafo_id)


@app.post("/reagir")
def reagir(req: ReacaoRequest):
    return ReacaoClient().reagir(req.desabafo_id)


@app.websocket("/sala/{room_id}")
async def websocket_chat(websocket: WebSocket, room_id: str, nickname: str):
    await websocket.accept()
    
    try:
        async with grpc.aio.insecure_channel(settings.servidor_a_host) as channel:
            stub = chat_pb2_grpc.ChatServiceStub(channel)
            
            # Open bidirectional stream
            stream = stub.Chat()
            
            # Send initial message to register the user in the room
            initial_msg = chat_pb2.ChatMessage(
                message_id="",
                text="entrou na sala",
                room_id=room_id,
                sender_nickname=nickname,
                send_at=int(time.time() * 1000)
            )
            await stream.write(initial_msg)
            
            # Task to read from gRPC and send to WebSocket
            async def grpc_to_ws():
                try:
                    async for grpc_msg in stream:
                        await websocket.send_json({
                            "message_id": grpc_msg.message_id,
                            "text": grpc_msg.text,
                            "room_id": grpc_msg.room_id,
                            "sender_nickname": grpc_msg.sender_nickname,
                            "send_at": grpc_msg.send_at
                        })
                except asyncio.CancelledError:
                    pass
                except Exception as e:
                    print(f"Error in grpc_to_ws: {e}")
            
            # Task to read from WebSocket and send to gRPC
            async def ws_to_grpc():
                try:
                    while True:
                        data = await websocket.receive_json()
                        text = data.get("text", "")
                        msg = chat_pb2.ChatMessage(
                            message_id="",
                            text=text,
                            room_id=room_id,
                            sender_nickname=nickname,
                            send_at=int(time.time() * 1000)
                        )
                        await stream.write(msg)
                except WebSocketDisconnect:
                    pass
                except asyncio.CancelledError:
                    pass
                except Exception as e:
                    print(f"Error in ws_to_grpc: {e}")
            
            grpc_task = asyncio.create_task(grpc_to_ws())
            ws_task = asyncio.create_task(ws_to_grpc())
            
            done, pending = await asyncio.wait(
                [grpc_task, ws_task],
                return_when=asyncio.FIRST_COMPLETED
            )
            
            for task in pending:
                task.cancel()
                
            await stream.done_writing()
            
    except Exception as e:
        print(f"Connection to gRPC ChatServer failed: {e}")
        try:
            await websocket.close()
        except Exception:
            pass


# Mount static files at root (must be mounted last so API routes take precedence)
static_path = os.path.join(os.path.dirname(os.path.dirname(__file__)), "static")
if os.path.exists(static_path):
    app.mount("/", StaticFiles(directory=static_path, html=True), name="static")
