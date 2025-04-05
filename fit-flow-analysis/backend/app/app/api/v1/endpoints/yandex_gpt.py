from fastapi import APIRouter, HTTPException, Depends
from pydantic import BaseModel
from typing import List, Dict, Any, Optional
from app.core.config import settings
from yandex_cloud_ml_sdk import YCloudML

router = APIRouter()


class Message(BaseModel):
    role: str
    text: str


class YandexGPTRequest(BaseModel):
    messages: List[Message]
    temperature: Optional[float] = 0.3
    model_version: Optional[str] = "rc"


class YandexGPTResponse(BaseModel):
    result: List[Dict[str, Any]]


@router.post("/", response_model=YandexGPTResponse)
async def generate_text(request: YandexGPTRequest):
    """
    Generate text using Yandex GPT.
    """
    try:
        # Initialize Yandex ML SDK
        sdk = YCloudML(
            folder_id=settings.YANDEX_FOLDER_ID, 
            auth=settings.YANDEX_API_KEY
        )

        # Configure and run the model
        model = sdk.models.completions("yandexgpt", model_version=request.model_version)
        model = model.configure(temperature=request.temperature)
        
        # Convert request messages to the format expected by the API
        formatted_messages = [
            {"role": msg.role, "text": msg.text} for msg in request.messages
        ]
        
        # Run the model with the formatted messages
        result = model.run(formatted_messages)
        
        # Return the result
        return {"result": [alt for alt in result]}
    
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"Error calling Yandex GPT API: {str(e)}") 