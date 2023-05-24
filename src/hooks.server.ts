import { SvelteKitAuth } from "@auth/sveltekit"
import GitHub from "@auth/core/providers/github"
import { GITHUB_ID, GITHUB_SECRET, AUTH_SECRET } from "$env/static/private"

export const handle = SvelteKitAuth({
    //@ts-ignore
  providers: [GitHub({ clientId: GITHUB_ID, clientSecret: GITHUB_SECRET })],
  secret: AUTH_SECRET,
  trustHost:true,
})