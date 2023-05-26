import { json } from '@sveltejs/kit';

export async function POST(event) {
	const body = JSON.parse(await event.request.text());
	console.warn(body);

	return json(details);
}

let details: App.PackageDetails = {
	name: "Yuki's ECS",
	version: '0.2.3',
	description: 'A prototype ECS engine, written by Yuki',
	tags: ['ecs', 'engine'],
	versions: ['0.9.5', '0.5.3'],
	funding: ['patreon', 'github'],
	dependsOn: ['something'],
	usedBy: ['something_else'],
	requirements: { minCompilierVersion: 'dev-2023-05:118ab605' },
	links: { repo: 'https://github.com/NoahR02/odin-ecs', discord: '' },
	lastUpdated: '2/2/23',
	license: 'Unlicense',
	size: '123kb',
	kind: 'unstable',
	owners: [{ name: 'NoahR02', username: 'NoahR02' }],
	stats: { allTimeDownloads: 42 },
	readme: `<p align="center" style="width:"> 
 <img width="100%" height="250" src="https://raw.githubusercontent.com/NoahR02/odin-ecs/main/repo_images/ecs-readme.svg">
 </p>
<p> **Odin ECS** was built because I needed a way to dynamically add functionality to things in my game. I also just find entity component systems fun to work with and I couldn't find a general purpose one for Odin , so I made it myself.</p>
<p>**Features**:</p>
<ul>
 <li>Any type can be a Component</li>
 <li>Unlimited* Amount of Component Types</li>
 </ul>
<p>Example Usage:</p>
<code>
package main
import ecs "odin-ecs"
import "core:fmt"
// Context: Internal state that the ECS needs to manipulate.
ctx: ecs.Context
// You can add any type as a Component!
Name :: distinct string
main :: proc() {
  ctx = ecs.init_ecs()
  defer ecs.deinit_ecs(&ctx)
  player := ecs.create_entity(&ctx)
 
  // (Optional) Or you can let ecs.deinit_ecs()... clean this up.
  defer ecs.destroy_entity(&ctx, player)
  name_component, err := ecs.add_component(&ctx, player, Name("Yuki"))
  fmt.println(name_component^) // "Yuki"
 
  remove_err := ecs.remove_component(&ctx, player, Name)
  //(Optional) Or you can let ecs.destroy_entity()... clean this up.
}
</code>
`
};
