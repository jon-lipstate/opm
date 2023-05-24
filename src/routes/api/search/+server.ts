import { json } from '@sveltejs/kit';

export async function POST(event) {
  //@ts-ignore
  const body = await event.request.text()
  const query = JSON.parse(body).query
  console.log(query);
  return json(data) ;
}
  // Simulated data
  const data = [
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
  ];