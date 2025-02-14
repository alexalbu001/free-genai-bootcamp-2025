from flask import request, jsonify, g
from flask_cors import cross_origin
import json

def load(app):
  # Endpoint: GET /words with pagination (50 words per page)
  @app.route('/api/words', methods=['GET'])
  @cross_origin()
  def get_words():
    try:
      cursor = app.db.cursor()
      page = request.args.get('page', 1, type=int)
      per_page = request.args.get('per_page', 10, type=int)
      sort_by = request.args.get('sort_by', 'kanji')
      order = request.args.get('order', 'asc').upper()
      
      offset = (page - 1) * per_page
      
      cursor.execute(f'''
        SELECT 
            w.*,
            COALESCE(wr.correct_count, 0) as correct_count,
            COALESCE(wr.wrong_count, 0) as wrong_count
        FROM words w
        LEFT JOIN word_reviews wr ON w.id = wr.word_id
        ORDER BY {sort_by} {order}
        LIMIT ? OFFSET ?
      ''', (per_page, offset))
      
      words = cursor.fetchall()
      
      # Get total count
      cursor.execute('SELECT COUNT(*) as count FROM words')
      total = cursor.fetchone()['count']
      
      return jsonify({
        "words": [{
          "id": word['id'],
          "kanji": word['kanji'],
          "romaji": word['romaji'],
          "english": word['english'],
          "correct_count": word['correct_count'],
          "wrong_count": word['wrong_count']
        } for word in words],
        "total_pages": (total + per_page - 1) // per_page,
        "current_page": page,
        "total_words": total
      })
      
    except Exception as e:
      return jsonify({"error": str(e)}), 500
    finally:
      app.db.close()

  # Endpoint: GET /api/words/:id to get a single word with its details
  @app.route('/api/words/<int:word_id>', methods=['GET'])
  @cross_origin()
  def get_word(word_id):
    try:
      cursor = app.db.cursor()
      
      # Query to fetch the word and its details
      cursor.execute('''
        SELECT w.id, w.kanji, w.romaji, w.english,
               COALESCE(r.correct_count, 0) AS correct_count,
               COALESCE(r.wrong_count, 0) AS wrong_count,
               GROUP_CONCAT(DISTINCT g.id || '::' || g.name) as groups
        FROM words w
        LEFT JOIN word_reviews r ON w.id = r.word_id
        LEFT JOIN word_groups wg ON w.id = wg.word_id
        LEFT JOIN groups g ON wg.group_id = g.id
        WHERE w.id = ?
        GROUP BY w.id
      ''', (word_id,))
      
      word = cursor.fetchone()
      
      if not word:
        return jsonify({"error": "Word not found"}), 404
      
      # Parse the groups string into a list of group objects
      groups = []
      if word["groups"]:
        for group_str in word["groups"].split(','):
          group_id, group_name = group_str.split('::')
          groups.append({
            "id": int(group_id),
            "name": group_name
          })
      
      return jsonify({
        "word": {
          "id": word["id"],
          "kanji": word["kanji"],
          "romaji": word["romaji"],
          "english": word["english"],
          "correct_count": word["correct_count"],
          "wrong_count": word["wrong_count"],
          "groups": groups
        }
      })
      
    except Exception as e:
      return jsonify({"error": str(e)}), 500