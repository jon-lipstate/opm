import sql from '$lib/database';
import { json } from '@sveltejs/kit';

export async function POST(event) {
	const details = JSON.parse(await event.request.text());
	return json({ msg: 'Not Implemented' }, { status: 500 });
	// TODO: NAME MUST BE SLUGGABLE IMPLICITLY

	// await make_package_table();

	try {
		console.warn('INSERTING PKG');
		await insert_new_package({
			name: 'Test',
			updated_at: new Date(),
			created_at: new Date(),
			downloads: 42,
			description: 'at a time',
			readme: 'stuff and things',
			repository: 'https://google.com',
			size_kb: 42
		});

		return json(null, { status: 201 });
	} catch (error) {
		console.error('Error inserting package details', error);
		return json(null, { status: 500 });
	}
}

const insert_new_package = async function ({
	name,
	updated_at,
	created_at,
	downloads,
	description,
	readme,
	repository,
	size_kb
}) {
	console.warn(name, updated_at, created_at, downloads, description, readme, repository, size_kb);
	try {
		await sql`INSERT INTO packages (name, updated_at, created_at, downloads, description, readme, repository, size_kb)
				VALUES (${name}, ${updated_at}, ${created_at}, ${downloads}, ${description}, ${readme}, ${repository}, ${size_kb})`;

		console.log('Package inserted successfully');
	} catch (error) {
		console.error('Error inserting package: ', error);
	}
};
