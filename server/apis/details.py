from flask import Blueprint, request, jsonify

# Create a blueprint for the details api
details_bp = Blueprint('details', __name__)

# Create the details route in this blueprint
@details_bp.route('/details', methods=['GET'])
def details():
    package = request.args.get('package', '')
    version = request.args.get('version', '')

    data = {
        'name': 'Yuki\'s ECS',
        'version': '0.2.3',
        'description': 'A prototype ECS engine, written by Yuki',
        'tags': ["ecs","engine"],
        'versions': [{"version":"1.2.3","date":"may 4, 2020","changes":"stuff and things"}], #just get IDs ???
        'funding': ["patreon","github"],
        'dependsOn': ["something"],
        'usedBy': ["something_else"],
        'requirements': { 'minCompilierVersion': 'dev-2023-05:118ab605' },
        'links': {'repo':'https://github.com/NoahR02/odin-ecs','discord':''},
        'lastUpdated': '2/2/23',
        'license': 'Unlicense',
        'size': '123kb',
        'kind': 'unstable',
        'owners': [{'name':"NoahR02",'username':"NoahR02"}],
        'stats': { 'allTimeDownloads': 42 },
        'readme': '<p align="center" style="width:"> <img width="100%" height="250" src="https://raw.githubusercontent.com/NoahR02/odin-ecs/main/repo_images/ecs-readme.svg"></p><p> **Odin ECS** was built because I needed a way to dynamically add functionality to things in my game. I also just find entity component systems fun to work with and I couldn\'t find a general purpose one for Odin , so I made it myself.</p><p>**Features**:</p><ul><li>Any type can be a Component</li><li>Unlimited* Amount of Component Types</li></ul><p>Example Usage:</p><code>package main\n\nimport ecs "odin-ecs"\nimport "core:fmt"\n\n// Context: Internal state that the ECS needs to manipulate.\nctx: ecs.Context\n\n// You can add any type as a Component!\nName :: distinct string\n\nmain :: proc() {\n\n  ctx = ecs.init_ecs()\n  defer ecs.deinit_ecs(&ctx)\n\n  player := ecs.create_entity(&ctx)\n \n  // (Optional) Or you can let ecs.deinit_ecs()... clean this up.\n  defer ecs.destroy_entity(&ctx, player)\n\n  name_component, err := ecs.add_component(&ctx, player, Name("Yuki"))\n  fmt.println(name_component^) // "Yuki"\n \n  remove_err := ecs.remove_component(&ctx, player, Name)\n  //(Optional) Or you can let ecs.destroy_entity()... clean this up.\n}\n</code>'
    }
    return jsonify(data), 200
