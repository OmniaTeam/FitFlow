import os
import json
from typing import Dict, List
import requests
from app.core.database import SessionLocal
from app.core.models import User

YANDEX_GPT_API_KEY = os.getenv("YANDEX_GPT_API_KEY")
YANDEX_GPT_API_URL = "https://llm.api.cloud.yandex.net/foundationModels/v1/completion"

def get_user_data(user_id: int) -> Dict:
    db = SessionLocal()
    try:
        user = db.query(User).filter(User.user_id == user_id).first()
        if not user:
            return {}
        return {
            "goals": user.goals,
            "preferences": user.preferences,
            "past_workouts": user.past_workouts,
            "past_diet": user.past_diet
        }
    finally:
        db.close()

def generate_workout_plan(user_id: int) -> Dict:
    user_data = get_user_data(user_id)
    if not user_data:
        return {}
    
    prompt = f"""
    Based on the following user data, generate a detailed workout plan:
    Goals: {user_data['goals']}
    Preferences: {user_data['preferences']}
    Past Workouts: {user_data['past_workouts']}
    
    Generate a weekly workout plan with specific exercises, sets, reps, and rest periods.
    """
    
    response = requests.post(
        YANDEX_GPT_API_URL,
        headers={
            "Authorization": f"Api-Key {YANDEX_GPT_API_KEY}",
            "Content-Type": "application/json"
        },
        json={
            "modelUri": "gpt://b1g.../yandexgpt",
            "completionOptions": {
                "temperature": 0.7,
                "maxTokens": 2000
            },
            "messages": [
                {
                    "role": "system",
                    "text": "You are a professional fitness trainer creating personalized workout plans."
                },
                {
                    "role": "user",
                    "text": prompt
                }
            ]
        }
    )
    
    if response.status_code == 200:
        return response.json()["result"]["alternatives"][0]["message"]["text"]
    return {}

def generate_diet_plan(user_id: int) -> Dict:
    user_data = get_user_data(user_id)
    if not user_data:
        return {}
    
    prompt = f"""
    Based on the following user data, generate a detailed diet plan:
    Goals: {user_data['goals']}
    Preferences: {user_data['preferences']}
    Past Diet: {user_data['past_diet']}
    
    Generate a weekly diet plan with specific meals, portions, and nutritional information.
    """
    
    response = requests.post(
        YANDEX_GPT_API_URL,
        headers={
            "Authorization": f"Api-Key {YANDEX_GPT_API_KEY}",
            "Content-Type": "application/json"
        },
        json={
            "modelUri": "gpt://b1g.../yandexgpt",
            "completionOptions": {
                "temperature": 0.7,
                "maxTokens": 2000
            },
            "messages": [
                {
                    "role": "system",
                    "text": "You are a professional nutritionist creating personalized diet plans."
                },
                {
                    "role": "user",
                    "text": prompt
                }
            ]
        }
    )
    
    if response.status_code == 200:
        return response.json()["result"]["alternatives"][0]["message"]["text"]
    return {} 