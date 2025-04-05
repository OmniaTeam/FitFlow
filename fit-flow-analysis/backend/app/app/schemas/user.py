from typing import List, Dict, Optional
from pydantic import BaseModel

class Preferences(BaseModel):
    diet: str
    exercise_type: str
    additional_preferences: Optional[Dict[str, str]] = None

class PastWorkout(BaseModel):
    date: str
    exercises: List[Dict[str, str]]
    duration: int
    calories_burned: float

class PastDiet(BaseModel):
    date: str
    meals: List[Dict[str, str]]
    calories: float
    nutrients: Dict[str, float]

class UserBase(BaseModel):
    user_id: int
    goals: List[str]
    preferences: Preferences
    past_workouts: List[PastWorkout]
    past_diet: List[PastDiet]

class UserCreate(UserBase):
    pass

class User(UserBase):
    class Config:
        from_attributes = True 