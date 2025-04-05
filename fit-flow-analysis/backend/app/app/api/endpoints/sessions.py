from fastapi import APIRouter, Depends, HTTPException
from sqlalchemy.orm import Session
from typing import List
from app.core.database import get_db
from app.schemas.session import SessionCreate, Session
from app.core.models import Session as SessionModel
from app.utils.yandex_gpt import generate_workout_plan, generate_diet_plan
from datetime import datetime

router = APIRouter()

@router.post("/sessions", response_model=Session)
def create_session(session: SessionCreate, db: Session = Depends(get_db)):
    db_session = SessionModel(user_id=session.user_id)
    db.add(db_session)
    db.commit()
    db.refresh(db_session)
    return db_session

@router.get("/sessions/{session_id}", response_model=Session)
def get_session(session_id: int, db: Session = Depends(get_db)):
    db_session = db.query(SessionModel).filter(SessionModel.session_id == session_id).first()
    if db_session is None:
        raise HTTPException(status_code=404, detail="Session not found")
    return db_session

@router.put("/sessions/{session_id}", response_model=Session)
def update_session(session_id: int, db: Session = Depends(get_db)):
    db_session = db.query(SessionModel).filter(SessionModel.session_id == session_id).first()
    if db_session is None:
        raise HTTPException(status_code=404, detail="Session not found")
    
    # Generate new plans using YandexGPT
    workout_plan = generate_workout_plan(db_session.user_id)
    diet_plan = generate_diet_plan(db_session.user_id)
    
    db_session.workout_plan = workout_plan
    db_session.diet_plan = diet_plan
    db_session.updated_at = datetime.now()
    
    db.commit()
    db.refresh(db_session)
    return db_session

@router.delete("/sessions/{session_id}")
def delete_session(session_id: int, db: Session = Depends(get_db)):
    db_session = db.query(SessionModel).filter(SessionModel.session_id == session_id).first()
    if db_session is None:
        raise HTTPException(status_code=404, detail="Session not found")
    db.delete(db_session)
    db.commit()
    return {"message": "Session deleted successfully"} 