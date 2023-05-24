import type { LayoutServerLoad } from "./$types"

export const load: LayoutServerLoad = async (event) => {
  const sesh = event.cookies.get('session')
  const token = event.cookies.get('next-auth.session-token')
  // console.warn(event.cookies.getAll());
// console.warn(event.request.headers);
  return {
    session: await event.locals.getSession(),
    // bearer:sesh,
    // token:token,
    userAgent: event.request.headers.get('user-agent')
  }
}
