from flask import jsonify
from flask_cors import cross_origin
from datetime import datetime, timedelta

def load(app):
    @app.route('/api/dashboard/recent-session', methods=['GET'])
    @cross_origin()
    def get_recent_session():
        try:
            cursor = app.db.cursor()
            
            # Get the most recent study session with activity name and results
            cursor.execute('''
                SELECT 
                    ss.id,
                    ss.group_id,
                    sa.name as activity_name,
                    ss.created_at,
                    COUNT(CASE WHEN wri.correct = 1 THEN 1 END) as correct_count,
                    COUNT(CASE WHEN wri.correct = 0 THEN 1 END) as wrong_count
                FROM study_sessions ss
                JOIN study_activities sa ON ss.study_activity_id = sa.id
                LEFT JOIN word_review_items wri ON ss.id = wri.study_session_id
                GROUP BY ss.id
                ORDER BY ss.created_at DESC
                LIMIT 1
            ''')
            
            session = cursor.fetchone()
            
            if not session:
                return jsonify(None)
            
            return jsonify({
                "id": session["id"],
                "group_id": session["group_id"],
                "activity_name": session["activity_name"],
                "created_at": session["created_at"],
                "correct_count": session["correct_count"],
                "wrong_count": session["wrong_count"]
            })
            
        except Exception as e:
            return jsonify({"error": str(e)}), 500

    @app.route('/api/dashboard/stats', methods=['GET'])
    @cross_origin()
    def get_stats():
        try:
            cursor = app.db.cursor()
            
            # Get total words
            cursor.execute('SELECT COUNT(*) as count FROM words')
            total_words = cursor.fetchone()['count']
            
            # Get studied words count
            cursor.execute('''
                SELECT COUNT(DISTINCT word_id) as count 
                FROM word_review_items
            ''')
            studied_words = cursor.fetchone()['count']
            
            # Get success rate
            cursor.execute('''
                SELECT 
                    ROUND(AVG(CASE WHEN correct = 1 THEN 100 ELSE 0 END), 2) as rate
                FROM word_review_items
            ''')
            success_rate = cursor.fetchone()['rate'] or 0
            
            # Get total sessions
            cursor.execute('SELECT COUNT(*) as count FROM study_sessions')
            total_sessions = cursor.fetchone()['count']
            
            # Get active groups
            cursor.execute('''
                SELECT COUNT(DISTINCT group_id) as count 
                FROM study_sessions 
                WHERE created_at >= date('now', '-30 days')
            ''')
            active_groups = cursor.fetchone()['count']
            
            return jsonify({
                "total_vocabulary": total_words,
                "total_words_studied": studied_words,
                "success_rate": success_rate,
                "total_sessions": total_sessions,
                "active_groups": active_groups,
                "current_streak": 0  # Implement streak calculation if needed
            })
            
        except Exception as e:
            return jsonify({"error": str(e)}), 500
