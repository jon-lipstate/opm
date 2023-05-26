export async function getAuth(event): Promise<auth> {
	const session = await event.locals.getSession();
	//@ts-ignore
	const login = session?.user.name;
	const authHeader = {
		headers: {
			//@ts-ignore
			Authorization: `token ${session?.accessToken}`
		}
	};
	return {
		login,
		authHeader,
		session
	};
}
