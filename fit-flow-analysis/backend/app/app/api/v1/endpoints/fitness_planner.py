from fastapi import APIRouter, HTTPException, Depends
from pydantic import BaseModel, Field, RootModel
from typing import List, Dict, Any, Optional, Union, ForwardRef
from app.core.config import settings
from yandex_cloud_ml_sdk import YCloudML
import json

router = APIRouter()


class ExercisePreference(BaseModel):
    exercise_type: str = Field(..., description="Type of exercise: strength, cardio, etc.")


class DietPreference(BaseModel):
    diet_type: str = Field(..., description="Type of diet: vegetarian, keto, etc.")


class UserPreferences(BaseModel):
    diet: DietPreference
    exercise: ExercisePreference


class Workout(BaseModel):
    date: str
    exercises: List[Dict[str, Any]]


class Exercise(BaseModel):
    exercise: str
    sets: int
    reps: Optional[int] = None
    weight: Optional[float] = None
    duration: Optional[str] = None


class Meal(BaseModel):
    name: str
    weight: float
    calories: float
    proteins: float
    fats: float
    carbs: float
    ingredients: Optional[List[Dict[str, Any]]] = None


class DailyMeals(BaseModel):
    breakfast: List[Meal]
    lunch: List[Meal]
    dinner: List[Meal]
    snacks: Optional[List[Meal]] = None


class FitnessPlanRequest(BaseModel):
    user_id: int
    goals: List[str]
    preferences: UserPreferences
    past_workouts: Optional[List[Workout]] = None
    past_diet: Optional[Dict[str, DailyMeals]] = None


class WorkoutPlan(RootModel):
    # Updated to use RootModel instead of __root__ field
    root: Dict[str, List[Exercise]]


class DietPlan(RootModel):
    # Updated to use RootModel instead of __root__ field
    root: Dict[str, DailyMeals]


class FitnessPlanResponse(BaseModel):
    workout_plan: WorkoutPlan
    diet_plan: DietPlan


@router.post("/generate", response_model=FitnessPlanResponse)
async def generate_fitness_plan(request: FitnessPlanRequest):
    """
    Generate personalized workout and diet plans using Yandex GPT.
    """
    try:
        # Initialize Yandex ML SDK
        sdk = YCloudML(
            folder_id=settings.YANDEX_FOLDER_ID, 
            auth=settings.YANDEX_API_KEY
        )

        # Configure the model
        model = sdk.models.completions("yandexgpt", model_version="rc")
        model = model.configure(temperature=0.3)
        
        # Create a system prompt that explains the task
        system_prompt = """
        Ты - эксперт по фитнесу и питанию. Твоя задача - создать персонализированный план тренировок и питания,
        учитывая цели пользователя, его предпочтения, прошлые тренировки и питание.
        
        План тренировок должен включать упражнения, подходы, повторения и веса.
        План питания должен включать приемы пищи на каждый день с указанием КБЖУ и веса порций.
        
        Ответ должен быть в формате JSON со следующей структурой:
        {
          "workout_plan": {
            "Monday": [
              {
                "exercise": "Push-ups",
                "sets": 4,
                "reps": 12,
                "weight": null
              }
            ],
            "Tuesday": [
              {
                "exercise": "Running",
                "sets": 1,
                "reps": null,
                "weight": null,
                "duration": "40 minutes"
              }
            ]
          },
          "diet_plan": {
            "Monday": {
              "breakfast": [
                {
                  "name": "oatmeal",
                  "weight": 180,
                  "calories": 216,
                  "proteins": 7.2,
                  "fats": 3.6,
                  "carbs": 33.6
                }
              ],
              "lunch": [...],
              "dinner": [...]
            },
            "Tuesday": {
              "breakfast": [...],
              "lunch": [...],
              "dinner": [...]
            }
          }
        }
        """
        
        # Prepare user data as a request to GPT
        user_data_json = request.model_dump_json(indent=2)
        user_prompt = f"Создай план тренировок и питания на основе следующих данных о пользователе: {user_data_json}"
        
        # Run the model with the formatted messages
        result = model.run([
            {"role": "system", "text": system_prompt},
            {"role": "user", "text": user_prompt}
        ])
        
        # Process the response
        # Get the first alternative from the result
        response_text = next(result)
        
        # Extract JSON from the response
        # This handles cases where the model might return some text before or after the JSON
        try:
            # Try to find JSON in the response
            start_idx = response_text.find('{')
            end_idx = response_text.rfind('}') + 1
            
            if start_idx >= 0 and end_idx > start_idx:
                json_str = response_text[start_idx:end_idx]
                response_data = json.loads(json_str)
            else:
                # If no JSON delimiters found, try parsing the whole text
                response_data = json.loads(response_text)
                
            # Validate the response structure
            if not isinstance(response_data, dict) or 'workout_plan' not in response_data or 'diet_plan' not in response_data:
                raise ValueError("Response missing required workout_plan or diet_plan fields")
                
            # Create a validated response object with RootModel
            return FitnessPlanResponse(
                workout_plan=WorkoutPlan(root=response_data['workout_plan']),
                diet_plan=DietPlan(root=response_data['diet_plan'])
            )
            
        except json.JSONDecodeError as json_err:
            raise HTTPException(
                status_code=500, 
                detail=f"Failed to parse JSON from GPT response: {str(json_err)}"
            )
        except ValueError as val_err:
            raise HTTPException(
                status_code=500,
                detail=f"Invalid response structure: {str(val_err)}"
            )
    
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"Error generating fitness plan: {str(e)}") 