from typing import Dict, List, Optional
from pydantic import BaseModel
from datetime import datetime

class WorkoutPlan(BaseModel):
    Monday: List[Dict[str, str]]
    Tuesday: List[Dict[str, str]]
    Wednesday: List[Dict[str, str]]
    Thursday: List[Dict[str, str]]
    Friday: List[Dict[str, str]]
    Saturday: List[Dict[str, str]]
    Sunday: List[Dict[str, str]]

class MealPlan(BaseModel):
    breakfast: List[Dict[str, str]]
    lunch: List[Dict[str, str]]
    dinner: List[Dict[str, str]]

class DietPlan(BaseModel):
    Monday: MealPlan
    Tuesday: MealPlan
    Wednesday: MealPlan
    Thursday: MealPlan
    Friday: MealPlan
    Saturday: MealPlan
    Sunday: MealPlan

class SessionBase(BaseModel):
    user_id: int
    workout_plan: WorkoutPlan
    diet_plan: DietPlan
    created_at: datetime
    updated_at: datetime

class SessionCreate(BaseModel):
    user_id: int

class Session(SessionBase):
    session_id: int

    class Config:
        from_attributes = True 