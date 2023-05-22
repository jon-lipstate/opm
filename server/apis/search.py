from flask import Blueprint, request, jsonify

# Create a blueprint for the search api
search_bp = Blueprint('search', __name__)

# Create the search route in this blueprint
@search_bp.route('/search', methods=['POST'])
def search():
    query = request.args.get('query', '')
    # Simulated data for demonstration
    data = [
            {
                'name': "Yuki's ECS",
                'version': '1.2.3',
                'updated': '1/1/1901',
                'downloads': 100000,
                'tags': ['ecs', 'engine'],
                'kind': 'Curated Library'
            },
            {
                'name': 'pico editor',
                'version': '0.0.1',
                'updated': '2/2/2902',
                'downloads': 10,
                'tags': ['tui', 'text-editor'],
                'kind': 'demo'
            }
        ]
    return jsonify(data), 200
