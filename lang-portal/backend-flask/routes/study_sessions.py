from flask import request, jsonify, g
from flask_cors import cross_origin
from datetime import datetime
import math

def load(app):
  # todo /study_sessions POST

  @app.route('/api/study-sessions', methods=['GET'])
  @cross_origin()
  def get_study_sessions():
    try:
      cursor = app.db.cursor()
      
      # Get pagination parameters
      page = request.args.get('page', 1, type=int)
      per_page = request.args.get('per_page', 10, type=int)
      offset = (page - 1) * per_page

      # Get total count
      cursor.execute('''
        SELECT COUNT(*) as count 
        FROM study_sessions ss
        JOIN groups g ON g.id = ss.group_id
        JOIN study_activities sa ON sa.id = ss.study_activity_id
      ''')
      total_count = cursor.fetchone()['count']

      # Get paginated sessions
      cursor.execute('''
        SELECT 
          ss.id,
          ss.group_id,
          g.name as group_name,
          sa.id as activity_id,
          sa.name as activity_name,
          ss.created_at,
          COUNT(wri.id) as review_items_count
        FROM study_sessions ss
        JOIN groups g ON g.id = ss.group_id
        JOIN study_activities sa ON sa.id = ss.study_activity_id
        LEFT JOIN word_review_items wri ON wri.study_session_id = ss.id
        GROUP BY ss.id
        ORDER BY ss.created_at DESC
        LIMIT ? OFFSET ?
      ''', (per_page, offset))
      sessions = cursor.fetchall()

      return jsonify({
        'items': [{
          'id': session['id'],
          'group_id': session['group_id'],
          'group_name': session['group_name'],
          'activity_id': session['activity_id'],
          'activity_name': session['activity_name'],
          'start_time': session['created_at'],
          'end_time': session['created_at'],  # For now, just use the same time since we don't track end time
          'review_items_count': session['review_items_count']
        } for session in sessions],
        'total': total_count,
        'page': page,
        'per_page': per_page,
        'total_pages': math.ceil(total_count / per_page)
      })
    except Exception as e:
      return jsonify({"error": str(e)}), 500

  @app.route('/api/study-sessions/<id>', methods=['GET'])
  @cross_origin()
  def get_study_session(id):
    try:
      cursor = app.db.cursor()
      
      # Get session details
      cursor.execute('''
        SELECT 
          ss.id,
          ss.group_id,
          g.name as group_name,
          sa.id as activity_id,
          sa.name as activity_name,
          ss.created_at,
          COUNT(wri.id) as review_items_count
        FROM study_sessions ss
        JOIN groups g ON g.id = ss.group_id
        JOIN study_activities sa ON sa.id = ss.study_activity_id
        LEFT JOIN word_review_items wri ON wri.study_session_id = ss.id
        WHERE ss.id = ?
        GROUP BY ss.id
      ''', (id,))
      
      session = cursor.fetchone()
      if not session:
        return jsonify({"error": "Study session not found"}), 404

      # Get pagination parameters
      page = request.args.get('page', 1, type=int)
      per_page = request.args.get('per_page', 10, type=int)
      offset = (page - 1) * per_page

      # Get the words reviewed in this session with their review status
      cursor.execute('''
        SELECT 
          w.*,
          COALESCE(SUM(CASE WHEN wri.correct = 1 THEN 1 ELSE 0 END), 0) as session_correct_count,
          COALESCE(SUM(CASE WHEN wri.correct = 0 THEN 1 ELSE 0 END), 0) as session_wrong_count
        FROM words w
        JOIN word_review_items wri ON wri.word_id = w.id
        WHERE wri.study_session_id = ?
        GROUP BY w.id
        ORDER BY w.kanji
        LIMIT ? OFFSET ?
      ''', (id, per_page, offset))
      
      words = cursor.fetchall()

      # Get total count of words
      cursor.execute('''
        SELECT COUNT(DISTINCT w.id) as count
        FROM words w
        JOIN word_review_items wri ON wri.word_id = w.id
        WHERE wri.study_session_id = ?
      ''', (id,))
      
      total_count = cursor.fetchone()['count']

      return jsonify({
        'session': {
          'id': session['id'],
          'group_id': session['group_id'],
          'group_name': session['group_name'],
          'activity_id': session['activity_id'],
          'activity_name': session['activity_name'],
          'start_time': session['created_at'],
          'end_time': session['created_at'],  # For now, just use the same time
          'review_items_count': session['review_items_count']
        },
        'words': [{
          'id': word['id'],
          'kanji': word['kanji'],
          'romaji': word['romaji'],
          'english': word['english'],
          'correct_count': word['session_correct_count'],
          'wrong_count': word['session_wrong_count']
        } for word in words],
        'total': total_count,
        'page': page,
        'per_page': per_page,
        'total_pages': math.ceil(total_count / per_page)
      })
    except Exception as e:
      return jsonify({"error": str(e)}), 500

  @app.route('/api/study-sessions/<int:id>/review', methods=['POST'])
  @cross_origin()
  def create_session_review(id):
    try:
      data = request.get_json()
      word_id = data.get('word_id')
      correct = data.get('correct')

      if word_id is None or correct is None:
        return jsonify({"error": "word_id and correct are required"}), 400

      cursor = app.db.cursor()

      # Verify the study session exists
      cursor.execute('SELECT id FROM study_sessions WHERE id = ?', (id,))
      if not cursor.fetchone():
        return jsonify({"error": "Study session not found"}), 404

      # Create the review item
      cursor.execute('''
        INSERT INTO word_review_items (word_id, study_session_id, correct, created_at)
        VALUES (?, ?, ?, datetime('now'))
      ''', (word_id, id, correct))

      # Update the word_reviews table
      cursor.execute('''
        INSERT INTO word_reviews (word_id, correct_count, wrong_count)
        VALUES (?, ?, ?)
        ON CONFLICT(word_id) DO UPDATE SET
        correct_count = correct_count + ?,
        wrong_count = wrong_count + ?,
        last_reviewed = datetime('now')
      ''', (
        word_id,
        1 if correct else 0,
        0 if correct else 1,
        1 if correct else 0,
        0 if correct else 1
      ))

      app.db.commit()

      return jsonify({
        "success": True,
        "word_id": word_id,
        "study_session_id": id,
        "correct": correct,
        "created_at": datetime.now().isoformat()
      })

    except Exception as e:
      return jsonify({"error": str(e)}), 500

  @app.route('/api/study-sessions/<int:id>/raw', methods=['GET'])
  @cross_origin()
  def get_session_raw_data(id):
    try:
      cursor = app.db.cursor()

      # Get session details with group and activity info
      cursor.execute('''
        SELECT 
          ss.id,
          ss.group_id,
          g.name as group_name,
          sa.id as activity_id,
          sa.name as activity_name,
          sa.url as activity_url,
          ss.created_at
        FROM study_sessions ss
        JOIN groups g ON g.id = ss.group_id
        JOIN study_activities sa ON sa.id = ss.study_activity_id
        WHERE ss.id = ?
      ''', (id,))

      session = cursor.fetchone()
      if not session:
        return jsonify({"error": "Study session not found"}), 404

      # Get all words in the group with their review status for this session
      cursor.execute('''
        SELECT 
          w.*,
          COALESCE(wri.correct, NULL) as review_status
        FROM words w
        JOIN word_groups wg ON w.id = wg.word_id
        LEFT JOIN word_review_items wri ON 
          w.id = wri.word_id AND 
          wri.study_session_id = ?
        WHERE wg.group_id = ?
      ''', (id, session['group_id']))

      words = cursor.fetchall()

      return jsonify({
        "session": {
          "id": session['id'],
          "group_id": session['group_id'],
          "group_name": session['group_name'],
          "activity_id": session['activity_id'],
          "activity_name": session['activity_name'],
          "activity_url": session['activity_url'],
          "created_at": session['created_at']
        },
        "words": [{
          "id": word['id'],
          "kanji": word['kanji'],
          "romaji": word['romaji'],
          "english": word['english'],
          "parts": word['parts'],
          "review_status": word['review_status']
        } for word in words]
      })

    except Exception as e:
      return jsonify({"error": str(e)}), 500

  @app.route('/api/study-sessions/reset', methods=['POST'])
  @cross_origin()
  def reset_study_sessions():
    try:
      cursor = app.db.cursor()
      
      # First delete all word review items since they have foreign key constraints
      cursor.execute('DELETE FROM word_review_items')
      
      # Then delete all study sessions
      cursor.execute('DELETE FROM study_sessions')
      
      app.db.commit()
      
      return jsonify({"message": "Study history cleared successfully"}), 200
    except Exception as e:
      return jsonify({"error": str(e)}), 500

  @app.route('/api/study-sessions', methods=['POST'])
  @cross_origin()
  def create_study_session():
    try:
      # Get request data
      data = request.get_json()
      group_id = data.get('group_id')
      activity_id = data.get('activity_id')
      
      # Validate required fields
      if not group_id or not activity_id:
        return jsonify({"error": "group_id and activity_id are required"}), 400
        
      cursor = app.db.cursor()
      
      # Insert new study session
      cursor.execute('''
        INSERT INTO study_sessions (group_id, study_activity_id, created_at)
        VALUES (?, ?, datetime('now'))
      ''', (group_id, activity_id))
      
      session_id = cursor.lastrowid
      app.db.commit()
      
      # Return the created session
      cursor.execute('''
        SELECT 
          ss.id,
          ss.group_id,
          g.name as group_name,
          sa.id as activity_id,
          sa.name as activity_name,
          ss.created_at
        FROM study_sessions ss
        JOIN groups g ON g.id = ss.group_id
        JOIN study_activities sa ON sa.id = ss.study_activity_id
        WHERE ss.id = ?
      ''', (session_id,))
      
      session = cursor.fetchone()
      
      return jsonify({
        'id': session['id'],
        'group_id': session['group_id'],
        'group_name': session['group_name'],
        'activity_id': session['activity_id'],
        'activity_name': session['activity_name'],
        'start_time': session['created_at'],
        'end_time': session['created_at']  # For now, same as start_time
      }), 201
      
    except Exception as e:
      return jsonify({"error": str(e)}), 500